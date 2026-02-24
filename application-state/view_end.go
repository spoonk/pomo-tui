package applicationstate

import (
	"fmt"
	"strings"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

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
		Foreground(lipgloss.Color("#A7C080")).Faint(true)

	stoppedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E69875")).Faint(true)

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
