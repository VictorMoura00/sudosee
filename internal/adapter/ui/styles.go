package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle  = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Bold(true).Padding(0, 1)
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false)

	colPID    = lipgloss.NewStyle().Width(8)
	colName   = lipgloss.NewStyle().Width(25)
	colStatus = lipgloss.NewStyle().Width(10)
	colMem    = lipgloss.NewStyle().Width(12)
	colCPU    = lipgloss.NewStyle().Width(10)
	colUser   = lipgloss.NewStyle().Width(12)

	styleRunning  = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	styleSleeping = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleIdle     = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))

	selectedRowStyle = lipgloss.NewStyle().Background(lipgloss.Color("236"))
)
