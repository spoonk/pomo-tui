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
	sessionLengthPrompt := "session length"
	breakLengthPrompt := "break length"

	maxPromptWidth := lipgloss.Width(sessionLengthPrompt)
	justifiedSLPrompt := lipgloss.NewStyle().Width(maxPromptWidth).Render(sessionLengthPrompt)
	justifiedBLPrompt := lipgloss.NewStyle().Width(maxPromptWidth).Render(breakLengthPrompt)

	if m.timeLimitInput.Focused() {
		timeLimit = inputStyleFocused.Render(fmt.Sprintf("> %s   %s min", justifiedSLPrompt, m.timeLimitInput.View()))
		breakLimit = inputStyleUnFocused.Render(fmt.Sprintf("  %s   %s min", justifiedBLPrompt, m.breakLimitInput.View()))
	} else {
		timeLimit = inputStyleUnFocused.Render(fmt.Sprintf("  %s   %s min", justifiedSLPrompt, m.timeLimitInput.View()))
		breakLimit = inputStyleFocused.Render(fmt.Sprintf("> %s   %s min", justifiedBLPrompt, m.breakLimitInput.View()))
	}

	// maxWidth = max(lipgloss.Width(timeLimit), lipgloss.Width(breakLimit))
	content := timeLimit + "\n" + breakLimit
	// leftAlignedContent := lipgloss.NewStyle().Width(maxWidth).Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
