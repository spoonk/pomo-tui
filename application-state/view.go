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
	var inputStyle = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(4)

	return inputStyle.Render(
		fmt.Sprintf(
			"session length %s", m.timeLimitInput.View()) +
			"\n" +
			fmt.Sprintf("break length % s", m.breakLimitInput.View()))
}

func (m Model) View() string {
	if m.sessionRunning {
		return m.runningSessionView()
	}

	return m.inputView()
}
