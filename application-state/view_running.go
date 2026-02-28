package applicationstate

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
	"pomo-tui/ui"
)

func (m Model) runningSessionView() string {
	keybinds := "p: pause · x: terminate · q: quit"
	if m.sessionTimer.IsPaused() {
		keybinds = "p: resume · x: terminate · q: quit"
	}

	project := ""
	if m.selectedProject != nil {
		project = ui.ProjectTitle(m.selectedProject.Name)
	}

	pause := ui.PauseIndicator(m.sessionTimer.IsPaused())
	timer := ui.TimerDisplay(m.sessionTimer.Elapsed(), m.sessionTimer.TimeLimit())
	hints := ui.Keybinds(keybinds)
	breakLabel := ui.DimLabel(fmt.Sprintf(" - (break: %s)", ui.FormatDuration(m.breakTimer.TimeLimit())))

	content := lipgloss.JoinVertical(lipgloss.Center, project, pause+timer+breakLabel)
	content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	centeredHints := lipgloss.Place(m.width, 0, lipgloss.Center, 0, hints)

	return ui.CompositeOver(content, centeredHints, 0, m.height-2)
}
