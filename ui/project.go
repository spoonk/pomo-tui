package ui

import (
	"fmt"

	lipgloss "charm.land/lipgloss/v2"
)

var ProjectStyle = lipgloss.NewStyle().Bold(true)

func ProjectTitle(title string) string {
	return ProjectStyle.Render(fmt.Sprintf("  %s", title))
}
