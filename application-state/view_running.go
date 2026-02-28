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

	pause := ui.PauseIndicator(m.sessionTimer.IsPaused())
	timer := ui.TimerDisplay(m.sessionTimer.Elapsed(), m.sessionTimer.TimeLimit())
	hints := ui.Keybinds(keybinds)
	breakLabel := ui.DimLabel(fmt.Sprintf(" - (break: %s)", ui.FormatDuration(m.breakTimer.TimeLimit())))

	content := pause + timer + breakLabel + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
