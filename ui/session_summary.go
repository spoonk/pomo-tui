package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ElapsedTime renders a "start > end (duration)" summary line for a completed session.
func ElapsedTime(start, end time.Time, active time.Duration) string {
	infoStyle := lipgloss.NewStyle().Faint(true).PaddingLeft(2)
	return infoStyle.Render(fmt.Sprintf("%s  >  %s (%s)",
		start.Format("3:04 PM"),
		end.Format("3:04 PM"),
		FormatDuration(active),
	))
}

// SessionEntry is the plain-data representation of a past session used for rendering.
type SessionEntry struct {
	Completed bool
	StartedAt time.Time
	Duration  time.Duration
	Planned   time.Duration
}

// SessionGrid renders past sessions as a compact grid of status icons.
func SessionGrid(entries []SessionEntry) string {
	if len(entries) == 0 {
		return ""
	}

	completedStyle := lipgloss.NewStyle().Foreground(SuccessColor).Faint(true)
	stoppedStyle := lipgloss.NewStyle().Foreground(WarnColor).Faint(true)

	const columnsPerRow = 10

	var rowBuilder strings.Builder
	var outputBuilder strings.Builder

	pushed := 0
	var fullRowWidth int

	for _, e := range entries {
		var icon string
		if e.Completed {
			icon = completedStyle.Render("✓ ")
		} else {
			icon = stoppedStyle.Render("⏺ ")
		}
		rowBuilder.WriteString(icon)
		pushed++

		if pushed == columnsPerRow {
			row := rowBuilder.String()
			if fullRowWidth == 0 {
				fullRowWidth = lipgloss.Width(row)
			}
			outputBuilder.WriteString(row + "\n")
			rowBuilder.Reset()
			pushed = 0
		}
	}

	if pushed != 0 {
		outputBuilder.WriteString(rowBuilder.String())
	}

	if fullRowWidth == 0 {
		fullRowWidth = lipgloss.Width(outputBuilder.String())
	}

	return lipgloss.NewStyle().Width(fullRowWidth).Render(outputBuilder.String())
}
