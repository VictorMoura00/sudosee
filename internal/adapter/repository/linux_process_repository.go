package repository

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/VictorMoura00/sudosee/internal/core/domain"
)

const linuxClockTicks float64 = 100.0

type LinuxProcessRepository struct {
	pageSize uint64
	users    map[string]string
}

func NewLinuxProcessRepository() *LinuxProcessRepository {
	repo := &LinuxProcessRepository{
		pageSize: uint64(os.Getpagesize()),
		users:    make(map[string]string),
	}
	repo.loadUsers()
	return repo
}

func (r *LinuxProcessRepository) loadUsers() {
	data, err := os.ReadFile("/etc/passwd")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				r.users[parts[2]] = parts[0]
			}
		}
	}
}

func (r *LinuxProcessRepository) getPIDs() ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	pids := make([]int, 0, len(entries)/2)
	for _, entry := range entries {
		if entry.IsDir() {
			if pid, err := strconv.Atoi(entry.Name()); err == nil {
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}

func (r *LinuxProcessRepository) getSystemUptime() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	parts := strings.SplitN(string(data), " ", 2)
	return strconv.ParseFloat(parts[0], 64)
}

func (r *LinuxProcessRepository) parseStat(pid int, sysUptime float64) (domain.Process, error) {
	path := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Process{}, err
	}

	line := string(data)
	startIdx := strings.IndexByte(line, '(')
	endIdx := strings.LastIndexByte(line, ')')

	if startIdx == -1 || endIdx == -1 {
		return domain.Process{}, fmt.Errorf("formato inválido")
	}

	name := line[startIdx+1 : endIdx]
	parts := strings.Split(line[endIdx+2:], " ")

	state := "?"
	if len(parts) > 0 {
		state = parts[0]
	}

	var ppid int
	if len(parts) > 1 {
		ppid, _ = strconv.Atoi(parts[1])
	}

	var memoryBytes uint64
	if len(parts) > 21 {
		if rssPages, err := strconv.ParseUint(parts[21], 10, 64); err == nil {
			memoryBytes = rssPages * r.pageSize
		}
	}

	var cpuUsage float64
	if len(parts) > 19 {
		utime, _ := strconv.ParseFloat(parts[11], 64)
		stime, _ := strconv.ParseFloat(parts[12], 64)
		starttime, _ := strconv.ParseFloat(parts[19], 64)

		totalTimeTicks := utime + stime
		secondsActive := sysUptime - (starttime / linuxClockTicks)

		if secondsActive > 0 {
			cpuUsage = 100 * ((totalTimeTicks / linuxClockTicks) / secondsActive)
		}
	}

	userName := "unknown"
	statusData, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err == nil {
		statusLines := strings.Split(string(statusData), "\n")
		for _, sLine := range statusLines {
			if strings.HasPrefix(sLine, "Uid:") {
				uidParts := strings.Fields(sLine)
				if len(uidParts) > 1 {
					uidStr := uidParts[1]
					if name, exists := r.users[uidStr]; exists {
						userName = name
					} else {
						userName = uidStr
					}
				}
				break
			}
		}
	}

	return domain.Process{
		PID:    pid,
		PPID:   ppid,
		Name:   name,
		State:  state,
		User:   userName,
		Memory: memoryBytes,
		CPU:    cpuUsage,
	}, nil
}

func (r *LinuxProcessRepository) GetAll() ([]domain.Process, error) {
	sysUptime, err := r.getSystemUptime()
	if err != nil {
		return nil, err
	}

	pids, err := r.getPIDs()
	if err != nil {
		return nil, err
	}

	processes := make([]domain.Process, 0, len(pids))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, pid := range pids {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			if proc, err := r.parseStat(id, sysUptime); err == nil {
				mu.Lock()
				processes = append(processes, proc)
				mu.Unlock()
			}
		}(pid)
	}

	wg.Wait()

	return processes, nil
}

func (r *LinuxProcessRepository) Terminate(pid int) error {
	if pid == 1 {
		return fmt.Errorf("operação negada")
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}

func (r *LinuxProcessRepository) GetSystemStats() (domain.SystemStats, error) {
	var stats domain.SystemStats

	memData, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		var available uint64
		for _, line := range strings.Split(string(memData), "\n") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				if fields[0] == "MemTotal:" {
					stats.TotalRAM = val * 1024
				} else if fields[0] == "MemAvailable:" {
					available = val * 1024
				}
			}
		}
		stats.UsedRAM = stats.TotalRAM - available
	}

	loadData, err := os.ReadFile("/proc/loadavg")
	if err == nil {
		stats.LoadAvg = strings.Fields(string(loadData))[0]
	}

	return stats, nil
}
