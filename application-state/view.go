package applicationstate

import (
	"fmt"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) runningSessionView() string {
	var timerStyle = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(4).
		Faint(true)

	d := time.Since(m.sessionStart)
	d = d.Round(time.Second)

	var timerText = m.spinner.View() + " "
	timerText += d.String()
	timerText += " / "
	timerText += m.timeLimit.Round(time.Second).String()

	return timerStyle.Render(timerText)
}

func (m Model) endView() string {
	checkStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		PaddingLeft(2)

	infoStyle := lipgloss.NewStyle().
		Faint(true).
		PaddingLeft(4)

	startTime := m.sessionStart.Format("3:04 PM")
	endTime := m.sessionEnd.Format("3:04 PM")

	duration := m.sessionEnd.Sub(m.sessionStart).Round(time.Second)
	durationStr := formatDuration(duration)

	checkLine := checkStyle.Render("✓ Session completed")
	timeLine := infoStyle.Render(fmt.Sprintf("Started: %s  →  Ended: %s", startTime, endTime))
	durationLine := infoStyle.Render(fmt.Sprintf("Duration: %s", durationStr))

	return checkLine + "\n" + timeLine + "\n" + durationLine
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

func (m Model) inputView() string {
	var inputStyleFocused = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(4)

	var inputStyleUnFocused = lipgloss.NewStyle().Bold(false).Faint(true).PaddingLeft(2)

	timeLimit := ""
	breakLimit := ""

	if m.timeLimitInput.Focused() {
		timeLimit = inputStyleFocused.Render(fmt.Sprintf("session length %s min", m.timeLimitInput.View()))
		breakLimit = inputStyleUnFocused.Render(fmt.Sprintf("break length %s min", m.breakLimitInput.View()))
	} else {
		timeLimit = inputStyleUnFocused.Render(fmt.Sprintf("session length %s min", m.timeLimitInput.View()))
		breakLimit = inputStyleFocused.Render(fmt.Sprintf("break length %s min", m.breakLimitInput.View()))
	}

	return timeLimit + "\n" + breakLimit

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
