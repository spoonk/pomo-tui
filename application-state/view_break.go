package applicationstate

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
	"pomo-tui/ui"
)

func (m Model) breakView() string {
	breakTitle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(ui.BreakColor).Render("[ break ]")

	keybinds := "p: pause · x: stop early · q: quit"
	if m.breakTimer.IsPaused() {
		keybinds = "p: resume · x: stop early · q: quit"
	}

	pause         := ui.PauseIndicator(m.breakTimer.IsPaused())
	timer         := ui.TimerDisplay(m.breakTimer.Elapsed(), m.breakTimer.TimeLimit())
	hints         := ui.Keybinds(keybinds)
	sessionLabel  := ui.DimLabel(fmt.Sprintf(" - (session: %s)", ui.FormatDuration(m.sessionTimer.TimeLimit())))

	content := breakTitle + "\n" + pause + timer + sessionLabel + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
