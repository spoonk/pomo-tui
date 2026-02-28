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

	pause := ui.PauseIndicator(m.breakTimer.IsPaused())
	breakIndicator := ui.BreakIndicator(true)
	timer := ui.TimerDisplay(m.breakTimer.Elapsed(), m.breakTimer.TimeLimit())
	hints := ui.Keybinds(keybinds)
	sessionLabel := ui.DimLabel(fmt.Sprintf(" - (session: %s)", ui.FormatDuration(m.sessionTimer.TimeLimit())))

	content := pause + " " + breakIndicator + timer + sessionLabel + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
