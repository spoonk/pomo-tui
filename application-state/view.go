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

	var timerText = ""
	timerText += d.String()
	timerText += " / "
	timerText += m.timeLimit.Round(time.Second).String()

	return timerStyle.Render(timerText)
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
	if m.sessionRunning {
		return m.runningSessionView()
	}

	return m.inputView()
}
