package applicationstate

import (
	"fmt"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) breakView() string {
	var breakTextStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#7FBBB3"))

	var timerStyle = lipgloss.NewStyle().
		Bold(true).
		Faint(true).Align(lipgloss.Center)

	var keybindStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	var sessionTextStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("242")).
		Align(lipgloss.Center)

	var pauseTextStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("#A7C080")).
		Align(lipgloss.Center)

	d := time.Since(m.breakStart) - m.breakPauseTimer.PausedDuration()
	d = d.Round(time.Second)

	pauseText := ""
	if m.breakPauseTimer.IsPaused() {
		pauseText += " ⏸︎   "
	}
	var timerText = ""
	timerText += d.String()
	timerText += " / "
	timerText += formatDuration(m.breakLimit.Round(time.Second))

	keybinds := "p: pause · x: stop early · q: quit"
	if m.breakPauseTimer.IsPaused() {
		keybinds = "p: resume · x: stop early · q: quit"
	}

	sessionText := fmt.Sprintf(" - (session: %s)", formatDuration(m.timeLimit))

	breakText := breakTextStyle.Render("[ break ]")
	pause := pauseTextStyle.Render(pauseText)
	timer := timerStyle.Render(timerText)
	hints := keybindStyle.Render(keybinds)
	breakUI := sessionTextStyle.Render(sessionText)

	content := breakText + "\n" + pause + timer + breakUI + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
