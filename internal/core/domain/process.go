package domain

type Process struct {
	PID        int
	PPID       int
	Name       string
	State      string
	User       string
	Memory     uint64
	CPU        float64
	TreePrefix string
}

type SystemStats struct {
	TotalRAM uint64
	UsedRAM  uint64
	LoadAvg  string
}

type ProcessRepository interface {
	GetAll() ([]Process, error)
	Terminate(pid int) error
	GetSystemStats() (SystemStats, error)
}
