package applicationstate

import (
	"fmt"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var timerStyle = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(4).
		Faint(true)
	var inputStyle = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(4)

	d := time.Since(m.sessionStart)
	d = d.Round(time.Second)

	var timerText = ""
	timerText += d.String()
	timerText += " / "
	timerText += m.sessionTimeLimitSeconds.Round(time.Second).String()

	var rendered = inputStyle.Render(fmt.Sprintf("session length %s", m.sessionTimeLimitInput.View())) + "\n" + timerStyle.Render(timerText)
	return rendered

	// return lipgloss.Place(100, 22, lipgloss.Center, lipgloss.Center, rendered)
}
