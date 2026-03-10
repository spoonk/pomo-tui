package ui

import (
	"fmt"
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
