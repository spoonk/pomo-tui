package ui

import (
	"fmt"
	"strings"
	"time"
)

// SessionListEntry is the plain-data representation of a past session for list rendering.
type SessionListEntry struct {
	Completed   bool
	ProjectName string // empty = no project
	StartedAt   time.Time
	EndedAt     time.Time
	Duration    time.Duration // active work time only
}

// SessionList renders past sessions as a list with project name, time range, and duration.
func SessionList(entries []SessionListEntry) string {
	if len(entries) == 0 {
		return ""
	}

	const projectWidth = 16

	var sb strings.Builder
	for i, e := range entries {
		projectName := e.ProjectName
		if projectName == "" {
			projectName = "No project"
		}
		// Truncate long project names
		if len(projectName) > projectWidth {
			projectName = projectName[:projectWidth-1] + "…"
		}

		project := DimStyle.Render(fmt.Sprintf("%-*s", projectWidth, projectName))
		timeRange := VeryDimStyle.Render(fmt.Sprintf("%s > %s", e.StartedAt.Format("3:04 PM"), e.EndedAt.Format("3:04 PM")))
		duration := VeryDimStyle.Render(FormatDuration(e.Duration))

		sb.WriteString(project + "  " + timeRange + "  " + duration)
		if i < len(entries)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
