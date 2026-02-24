package applicationstate

import (
	"fmt"

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
