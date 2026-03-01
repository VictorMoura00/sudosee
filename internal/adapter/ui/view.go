package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Erro crítico: %v\n", m.err)
	}

	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString(m.renderSearchAndMessage())
	b.WriteString(m.renderTableHeader())
	b.WriteString(m.renderTableRows())

	return b.String()
}

func (m Model) renderHeader() string {
	header := titleStyle.Render("SudoSee") + " Monitor de Tarefas (Pressione 'q' para sair)\n"

	if m.sysStats.TotalRAM > 0 {
		percent := float64(m.sysStats.UsedRAM) / float64(m.sysStats.TotalRAM)
		filledWidth := int(percent * 30.0)
		emptyWidth := 30 - filledWidth

		bar := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(strings.Repeat("█", filledWidth))
		bar += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("░", emptyWidth))

		usedMB := float64(m.sysStats.UsedRAM) / 1024 / 1024
		totalMB := float64(m.sysStats.TotalRAM) / 1024 / 1024

		header += fmt.Sprintf("RAM: [%s] %.0f/%.0f MB   |   Load Avg: %s\n", bar, usedMB, totalMB, m.sysStats.LoadAvg)
	}

	header += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
		fmt.Sprintf("Atalhos: [↑/↓] Navegar | [m/c/p/t] Ordenar | [k] Matar | [/] Buscar ---> Ordenando por: %s", m.sortBy),
	) + "\n"
	return header
}

func (m Model) renderSearchAndMessage() string {
	var out string
	if m.isSearching || m.searchInput.Value() != "" {
		out += lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("🔍 Buscar: ") + m.searchInput.View() + "\n"
	} else {
		out += "\n"
	}

	if m.message != "" {
		out += lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render(m.message) + "\n"
	} else {
		out += "\n"
	}
	return out
}

func (m Model) renderTableHeader() string {
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		colPID.Render("PID"),
		colUser.Render("USER"),
		colName.Render("NOME"),
		colStatus.Render("STATUS"),
		colMem.Render("MEM (MB)"),
		colCPU.Render("CPU (%)"),
	)
	return headerStyle.Render(header) + "\n"
}
func (m Model) renderTableRows() string {
	end := m.offset + viewLimit
	if end > len(m.processes) {
		end = len(m.processes)
	}

	var rows string
	for i := m.offset; i < end; i++ {
		proc := m.processes[i]
		memMB := float64(proc.Memory) / 1024 / 1024
		fullDisplayName := proc.TreePrefix + proc.Name

		chars := []rune(fullDisplayName)
		if len(chars) > 22 {
			fullDisplayName = string(chars[:19]) + "..."
		}

		displayName := proc.Name
		if len(displayName) > 22 {
			displayName = displayName[:19] + "..."
		}

		statusText := proc.State
		switch proc.State {
		case "R":
			statusText = styleRunning.Render("R (Run)")
		case "S":
			statusText = styleSleeping.Render("S (Slp)")
		case "I":
			statusText = styleIdle.Render("I (Idl)")
		}

		row := lipgloss.JoinHorizontal(lipgloss.Left,
			colPID.Render(fmt.Sprintf("%d", proc.PID)),
			colUser.Render(proc.User),
			colName.Render(fullDisplayName),
			colStatus.Render(statusText),
			colMem.Render(fmt.Sprintf("%6.2f", memMB)),
			colCPU.Render(fmt.Sprintf("%5.1f", proc.CPU)),
		)

		if i == m.cursor {
			row = selectedRowStyle.Render(row)
			row = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(" > ") + row
		} else {
			row = "   " + row
		}

		rows += row + "\n"
	}
	return rows
}
