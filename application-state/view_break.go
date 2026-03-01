package applicationstate

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
	"pomo-tui/ui"
)

func (m Model) breakView() string {
	// breakTitle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(ui.BreakColor).Render("[ break ]")

	keybinds := "p: pause · x: terminate · q: quit"
	if m.breakTimer.IsPaused() {
		keybinds = "p: resume · x: terminate · q: quit"
	}

	project := ""
	if m.selectedProject != nil {
		project = ui.ProjectTitle(m.selectedProject.Name)
	}
	pause := ui.PauseIndicator(m.breakTimer.IsPaused())
	breakIndicator := ui.BreakIndicator(true)
	timer := ui.TimerDisplay(m.breakTimer.Elapsed(), m.breakTimer.TimeLimit())
	hints := ui.Keybinds(keybinds)
	sessionLabel := ui.DimLabel(fmt.Sprintf(" - (session: %s)", ui.FormatDuration(m.sessionTimer.TimeLimit())))

	content := lipgloss.JoinVertical(lipgloss.Center, project, timer+sessionLabel)
	contentW := lipgloss.Width(content)
	content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	centeredHints := lipgloss.Place(m.width, 0, lipgloss.Center, 0, hints)

	contentWithHints := ui.CompositeOver(content, centeredHints, 0, m.height-2)

	return ui.CompositeOver(contentWithHints, pause+breakIndicator, (m.width-contentW)/2-6, m.height/2)

}
