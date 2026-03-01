package usecase

import (
	"sort"
	"strings"

	"github.com/VictorMoura00/sudosee/internal/core/domain"
)

type ListProcessesUseCase struct {
	repo domain.ProcessRepository
}

func NewListProcessesUseCase(repo domain.ProcessRepository) *ListProcessesUseCase {
	return &ListProcessesUseCase{repo: repo}
}

func (uc *ListProcessesUseCase) Execute(sortBy string, filter string) ([]domain.Process, domain.SystemStats, error) {
	allProcesses, err := uc.repo.GetAll()
	if err != nil {
		return nil, domain.SystemStats{}, err
	}

	stats, _ := uc.repo.GetSystemStats()

	var processes []domain.Process
	filterLower := strings.ToLower(filter)

	for _, p := range allProcesses {
		if filter == "" || strings.Contains(strings.ToLower(p.Name), filterLower) {
			processes = append(processes, p)
		}
	}

	if sortBy == "tree" {
		processes = buildTree(processes)
	} else {
		sort.Slice(processes, func(i, j int) bool {
			switch sortBy {
			case "mem":
				if processes[i].Memory == processes[j].Memory {
					return processes[i].PID < processes[j].PID
				}
				return processes[i].Memory > processes[j].Memory
			case "cpu":
				if processes[i].CPU == processes[j].CPU {
					return processes[i].PID < processes[j].PID
				}
				return processes[i].CPU > processes[j].CPU
			default: // "pid"
				return processes[i].PID < processes[j].PID
			}
		})
	}

	return processes, stats, nil
}

func buildTree(procs []domain.Process) []domain.Process {
	procMap := make(map[int]domain.Process)
	childrenMap := make(map[int][]domain.Process)

	for _, p := range procs {
		procMap[p.PID] = p
		childrenMap[p.PPID] = append(childrenMap[p.PPID], p)
	}

	for _, children := range childrenMap {
		sort.Slice(children, func(i, j int) bool {
			return children[i].PID < children[j].PID
		})
	}

	var roots []domain.Process
	for _, p := range procs {
		if _, exists := procMap[p.PPID]; !exists {
			roots = append(roots, p)
		}
	}
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].PID < roots[j].PID
	})

	var result []domain.Process

	var traverse func(p domain.Process, prefix string)
	traverse = func(p domain.Process, prefix string) {
		result = append(result, p)
		children := childrenMap[p.PID]

		for i, child := range children {
			isLast := i == len(children)-1
			branch := "├─ "
			if isLast {
				branch = "└─ "
			}

			child.TreePrefix = prefix + branch

			nextPrefix := prefix + "│  "
			if isLast {
				nextPrefix = prefix + "   "
			}
			traverse(child, nextPrefix)
		}
	}

	for _, root := range roots {
		traverse(root, "")
	}

	return result
}
