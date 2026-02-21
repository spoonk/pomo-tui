package applicationstate

import (
	"fmt"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) inputView() string {
	var inputStyleFocused = lipgloss.NewStyle().
		Bold(true)

	var inputStyleUnFocused = lipgloss.NewStyle().Faint(true)

	timeLimit := ""
	breakLimit := ""

	if m.timeLimitInput.Focused() {
		timeLimit = inputStyleFocused.Render(fmt.Sprintf("> session length %s min", m.timeLimitInput.View()))
		breakLimit = inputStyleUnFocused.Render(fmt.Sprintf("  break length %s min", m.breakLimitInput.View()))
	} else {
		timeLimit = inputStyleUnFocused.Render(fmt.Sprintf("  session length %s min", m.timeLimitInput.View()))
		breakLimit = inputStyleFocused.Render(fmt.Sprintf("> break length %s min", m.breakLimitInput.View()))
	}

	content := timeLimit + "\n" + breakLimit
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

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
		Foreground(lipgloss.Color("129")).
		Align(lipgloss.Center)

	d := time.Since(m.sessionStart) - m.pausedDuration
	d = d.Round(time.Second)

	pauseText := ""
	if m.isPaused {
		pauseText += "[paused] "
	}
	var timerText = ""
	timerText += d.String()
	timerText += " / "
	timerText += formatDuration(m.timeLimit.Round(time.Second))

	keybinds := "p: pause · q: quit"
	if m.isPaused {
		keybinds = "p: resume · q: quit"
	}

	breakText := fmt.Sprintf(" - (break: %s)", formatDuration(m.breakLimit))

	pause := pauseTextStyle.Render(pauseText)
	timer := timerStyle.Render(timerText)
	hints := keybindStyle.Render(keybinds)
	breakUI := breakTextStyle.Render(breakText)

	content := pause + timer + breakUI + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) endView() string {
	checkStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A7C080"))

	infoStyle := lipgloss.NewStyle().
		Faint(true).
		PaddingLeft(2)

	startTime := m.sessionStart.Format("3:04 PM")
	endTime := m.sessionEnd.Format("3:04 PM")

	duration := m.sessionEnd.Sub(m.sessionStart).Round(time.Second)
	durationStr := formatDuration(duration)

	checkLine := checkStyle.Render("✓ Session complete")
	timeLine := infoStyle.Render(fmt.Sprintf("%s  →  %s (%s)", startTime, endTime, durationStr))
	// durationLine := infoStyle.Render(fmt.Sprintf("%s", durationStr))

	content := checkLine + "\n" + timeLine
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
		}
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", seconds)
}

func (m Model) View() string {
	switch m.programState {
	case initialState:
		return m.inputView()

	case sessionRunningState:
		return m.runningSessionView()

	default:
		return m.endView()
	}

}
