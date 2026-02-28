package ui

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
		}
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", seconds)
}
