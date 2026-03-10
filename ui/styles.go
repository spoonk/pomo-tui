package ui

import "github.com/charmbracelet/lipgloss"

var (
	TimerStyle   = lipgloss.NewStyle().Bold(true).Faint(true).Align(lipgloss.Center)
	VeryDimStyle = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
	DimStyle     = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("242")).Align(lipgloss.Center)
	PauseStyle   = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#A7C080")).Align(lipgloss.Center)
	BreakStyle   = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#D699B6")).Align(lipgloss.Center)

	SuccessColor = lipgloss.Color("#A7C080")
	WarnColor    = lipgloss.Color("#E69875")
	BreakColor   = lipgloss.Color("#7FBBB3")
)
