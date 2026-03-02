package ui

import (
	"time"

	"github.com/VictorMoura00/sudosee/internal/core/domain"
	"github.com/VictorMoura00/sudosee/internal/core/usecase"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

const viewLimit = 15

type Model struct {
	listUseCase *usecase.ListProcessesUseCase
	killUseCase *usecase.KillProcessUseCase
	processes   []domain.Process
	sysStats    domain.SystemStats
	err         error
	sortBy      string
	cursor      int
	offset      int
	message     string
	searchInput textinput.Model
	isSearching bool
}

func NewModel(listUc *usecase.ListProcessesUseCase, killUc *usecase.KillProcessUseCase) Model {
	ti := textinput.New()
	ti.Placeholder = "Digite para filtrar..."
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		listUseCase: listUc,
		killUseCase: killUc,
		sortBy:      "mem",
		cursor:      0,
		offset:      0,
		searchInput: ti,
		isSearching: false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), textinput.Blink)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
