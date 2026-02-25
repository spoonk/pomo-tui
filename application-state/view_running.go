package applicationstate

import (
	"fmt"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) runningSessionView() string {
	var timerStyle = lipgloss.NewStyle().
		Bold(true).
		Faint(true).Align(lipgloss.Center)

	var keybindStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	var breakTextStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("242")).
		Align(lipgloss.Center)

	var pauseTextStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("#A7C080")).
		Align(lipgloss.Center)

	d := time.Since(m.sessionStart) - m.sessionPauseTimer.PausedDuration()
	d = d.Round(time.Second)

	pauseText := ""
	if m.sessionPauseTimer.IsPaused() {
		pauseText += " ⏸︎   "
	}
	var timerText = ""
	timerText += d.String()
	timerText += " / "
	timerText += formatDuration(m.timeLimit.Round(time.Second))

	keybinds := "p: pause · x: stop early · q: quit"
	if m.sessionPauseTimer.IsPaused() {
		keybinds = "p: resume · x: stop early · q: quit"
	}

	breakText := fmt.Sprintf(" - (break: %s)", formatDuration(m.breakLimit))

	pause := pauseTextStyle.Render(pauseText)
	timer := timerStyle.Render(timerText)
	hints := keybindStyle.Render(keybinds)
	breakUI := breakTextStyle.Render(breakText)

	content := pause + timer + breakUI + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
