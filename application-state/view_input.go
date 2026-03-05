package applicationstate

import (
	"fmt"
	"pomo-tui/ui"

	lipgloss "charm.land/lipgloss/v2"
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
	maxInputLength += 2 - maxInputLength%2 // insane hack to handle input width increasing by 1 when first character is entered

	dropdown := m.projectSelector.DropdownView()

	borderedStyle := lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
	content := applyWidth(timeLimit, maxInputLength) + "\n" + applyWidth(breakLimit, maxInputLength) + "\n" + applyWidth(project, maxInputLength)

	inputGroup := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, borderedStyle.Render(content))

	groupW := lipgloss.Width(content)
	groupH := lipgloss.Height(content)

	return ui.CompositeOver(inputGroup, dropdown, (m.width-groupW)/2, (m.height-groupH)/2+3)
}

func applyWidth(input string, width int) string {
	return lipgloss.NewStyle().Width(width).Render(input)
}
