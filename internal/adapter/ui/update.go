package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isSearching {
			switch msg.String() {
			case "enter", "esc":
				m.isSearching = false
				m.searchInput.Blur()
				return m, nil
			}
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "/":
			m.isSearching = true
			m.searchInput.Focus()
			return m, textinput.Blink

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			if m.cursor < m.offset {
				m.offset = m.cursor
			}

		case "down":
			if m.cursor < len(m.processes)-1 {
				m.cursor++
			}
			if m.cursor >= m.offset+viewLimit {
				m.offset = m.cursor - viewLimit + 1
			}

		case "m":
			m.sortBy = "mem"
			m.cursor, m.offset = 0, 0
		case "c":
			m.sortBy = "cpu"
			m.cursor, m.offset = 0, 0
		case "p":
			m.sortBy = "pid"
			m.cursor, m.offset = 0, 0
		case "t":
			m.sortBy = "tree"
			m.cursor, m.offset = 0, 0

		case "k", "delete":
			if len(m.processes) > 0 {
				pidToKill := m.processes[m.cursor].PID
				if err := m.killUseCase.Execute(pidToKill); err != nil {
					m.message = fmt.Sprintf("Erro: %v", err)
				} else {
					m.message = fmt.Sprintf("Sucesso: Sinal enviado para PID %d", pidToKill)
				}
			}
		}

	case tickMsg:
		processes, stats, err := m.listUseCase.Execute(m.sortBy, m.searchInput.Value())
		if err != nil {
			m.err = err
		} else {
			m.processes = processes
			m.sysStats = stats

			if m.cursor >= len(m.processes) {
				m.cursor = len(m.processes) - 1
				if m.cursor < 0 {
					m.cursor = 0
				}
			}
		}
		return m, tickCmd()
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}
