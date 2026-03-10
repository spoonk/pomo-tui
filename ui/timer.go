package ui

import "time"

// TimerDisplay renders "elapsed / limit" in the shared timer style.
func TimerDisplay(elapsed, limit time.Duration) string {
	text := elapsed.Round(time.Second).String() + " / " + FormatDuration(limit.Round(time.Second))
	return TimerStyle.Render(text)
}

// PauseIndicator renders a pause icon when the timer is paused, or an empty
// string when it is running.
func PauseIndicator(isPaused bool) string {
	if !isPaused {
		return "  "
	}
	return PauseStyle.Render("  ")
}

func BreakIndicator(isBreak bool) string {
	if !isBreak {
		return "  "
	}

	return BreakStyle.Render(" 󰒲  ")
}

// DimLabel renders a secondary label in the muted dim style.
func DimLabel(s string) string {
	return DimStyle.Render(s)
}
