package applicationstate

import (
	"fmt"
	"strings"
	"time"

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

func (m Model) runningSessionView() string {
	var timerStyle = lipgloss.NewStyle().
		Bold(true).
		Faint(true).Align(lipgloss.Center)

	var keybindStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	var breakTextStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("242")).
		Align(lipgloss.Center)

	var pauseTextStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("#A7C080")).
		Align(lipgloss.Center)

	d := time.Since(m.sessionStart) - m.pausedDuration
	d = d.Round(time.Second)

	pauseText := ""
	if m.isPaused {
		pauseText += " ⏸︎   "
	}
	var timerText = ""
	timerText += d.String()
	timerText += " / "
	timerText += formatDuration(m.timeLimit.Round(time.Second))

	keybinds := "p: pause · x: stop early · q: quit"
	if m.isPaused {
		keybinds = "p: resume · q: quit"
	}

	breakText := fmt.Sprintf(" - (break: %s)", formatDuration(m.breakLimit))

	pause := pauseTextStyle.Render(pauseText)
	timer := timerStyle.Render(timerText)
	hints := keybindStyle.Render(keybinds)
	breakUI := breakTextStyle.Render(breakText)

	content := pause + timer + breakUI + "\n" + hints
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) completedView() string {
	checkStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A7C080"))

	checkLine := checkStyle.Render("✓ Session complete")
	infoLine := elapsedTimeUI(m)
	keybinds := endStateKeybindsUI()
	list := sessionListUI(m)

	content := checkLine + "\n" + infoLine + "\n" + keybinds + "\n\n" + list
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) stoppedEarlyView() string {
	checkStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#E69875"))

	checkLine := checkStyle.Render("⏺  Session stopped early")
	infoLine := elapsedTimeUI(m)
	keybinds := endStateKeybindsUI()
	list := sessionListUI(m)

	content := checkLine + "\n" + infoLine + "\n" + keybinds + "\n\n" + list
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func elapsedTimeUI(m Model) string {
	infoStyle := lipgloss.NewStyle().
		Faint(true).
		PaddingLeft(2)

	startTime := m.sessionStart.Format("3:04 PM")
	endTime := m.sessionEnd.Format("3:04 PM")

	duration := m.sessionEnd.Sub(m.sessionStart).Round(time.Second)
	durationStr := formatDuration(duration)

	return infoStyle.Render(fmt.Sprintf("%s  >  %s (%s)", startTime, endTime, durationStr))
}

// Currently not feeling the spark of beauty from this :(
func sessionListUI(m Model) string {
	if m.store == nil {
		return ""
	}
	sessions := m.store.GetSessions()
	if len(sessions) == 0 {
		return ""
	}

	blockStyle := lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241"))

	completedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A7C080"))

	stoppedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E69875"))

	var rows []string
	for _, s := range sessions {
		duration := time.Duration(s.DurationSeconds) * time.Second
		planned := time.Duration(s.PlannedSeconds) * time.Second
		completed := s.DurationSeconds >= s.PlannedSeconds

		var icon, durationStr string
		if completed {
			icon = completedStyle.Render("✓")
			durationStr = formatDuration(planned)
		} else {
			icon = stoppedStyle.Render("⏺")
			durationStr = formatDuration(duration)
		}

		timeStr := s.StartedAt.Local().Format("3:04 PM")
		row := fmt.Sprintf("%s  %s  %s", icon, timeStr, blockStyle.Render(fmt.Sprintf("(%s)", durationStr)))
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func endStateKeybindsUI() string {
	var keybindStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	keybinds := "r: restart · e: edit duration · q: quit"

	return keybindStyle.Render(keybinds)
}

func formatDuration(d time.Duration) string {
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

func (m Model) View() string {
	switch m.programState {
	case editTimerState:
		return m.inputView()

	case sessionRunningState:
		return m.runningSessionView()

	case sessionCompleteState:
		return m.completedView()

	case sessionEndedEarlyState:
		return m.stoppedEarlyView()

	default:
		return m.completedView()
	}
}
