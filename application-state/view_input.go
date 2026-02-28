package applicationstate

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) inputView() string {
	var inputStyleFocused = lipgloss.NewStyle().Bold(true)
	var inputStyleUnFocused = lipgloss.NewStyle().Faint(true)

	sessionLengthPrompt := "session length"
	breakLengthPrompt := "break length"
	projectPrompt := "project"

	maxPromptWidth := lipgloss.Width(sessionLengthPrompt)
	justifiedSLPrompt := lipgloss.NewStyle().Width(maxPromptWidth).Render(sessionLengthPrompt)
	justifiedBLPrompt := lipgloss.NewStyle().Width(maxPromptWidth).Render(breakLengthPrompt)
	justifiedProjectPrompt := lipgloss.NewStyle().Width(maxPromptWidth).Render(projectPrompt)

	prefix := func(focused bool) string {
		if focused {
			return "> "
		}
		return "  "
	}
	style := func(focused bool) lipgloss.Style {
		if focused {
			return inputStyleFocused
		}
		return inputStyleUnFocused
	}

	timeFocused := m.timeLimitInput.Focused()
	breakFocused := m.breakLimitInput.Focused()
	projectFocused := m.projectSelector.Focused()

	timeLimit := style(timeFocused).Render(fmt.Sprintf("%s%s   %s min", prefix(timeFocused), justifiedSLPrompt, m.timeLimitInput.View()))
	breakLimit := style(breakFocused).Render(fmt.Sprintf("%s%s   %s min", prefix(breakFocused), justifiedBLPrompt, m.breakLimitInput.View()))
	project := style(projectFocused).Render(fmt.Sprintf("%s%s   %s", prefix(projectFocused), justifiedProjectPrompt, m.projectSelector.View()))

	maxInputLength := max(lipgloss.Width(timeLimit), lipgloss.Width(breakLimit), lipgloss.Width(project))

	// Indent dropdown items to align with the value column.
	// Layout: prefix (2) + prompt (maxPromptWidth) + separator (3)
	dropdown := m.projectSelector.DropdownView()

	content := applyWidth(timeLimit, maxInputLength) + "\n" + applyWidth(breakLimit, maxInputLength) + "\n" + applyWidth(project, maxInputLength)
	if dropdown != "" {
		content += "\n" + dropdown
	}

	borderedStyle := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, borderedStyle.Render(content))
}

func applyWidth(input string, width int) string {
	return lipgloss.NewStyle().Width(width).Render(input)
}
