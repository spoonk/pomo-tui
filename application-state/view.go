package applicationstate

import (
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(4).
		Faint(true).Align(lipgloss.Right)

	d := time.Since(m.sessionStart)
	d = d.Round(time.Second)

	s := ""

	s += d.String()
	s += " / "
	s += m.sessionTimeLimitSeconds.Round(time.Second).String()
	return lipgloss.Place(100, 22, lipgloss.Center, lipgloss.Center, style.Render(s))
}
