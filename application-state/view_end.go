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

func sessionListUI(m Model) string {
	if m.store == nil {
		return ""
	}
	sessions := m.store.GetSessions()
	if len(sessions) == 0 {
		return ""
	}

	dimStyle := lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241"))

	completedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A7C080")).Faint(true)

	stoppedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E69875")).Faint(true)

	type rowData struct {
		completed bool
		timeStr   string
		durStr    string
	}

	rows := make([]rowData, len(sessions))
	maxTimeWidth := 0
	maxDurWidth := 0

	for i, s := range sessions {
		duration := time.Duration(s.DurationSeconds) * time.Second
		planned := time.Duration(s.PlannedSeconds) * time.Second
		completed := s.DurationSeconds >= s.PlannedSeconds

		var durStr string
		if completed {
			durStr = formatDuration(planned)
		} else {
			durStr = formatDuration(duration)
		}

		timeStr := s.StartedAt.Local().Format("3:04 PM")

		if len(timeStr) > maxTimeWidth {
			maxTimeWidth = len(timeStr)
		}
		if len(durStr) > maxDurWidth {
			maxDurWidth = len(durStr)
		}

		rows[i] = rowData{completed: completed, timeStr: timeStr, durStr: durStr}
	}

	var lines []string
	for _, r := range rows {
		var icon string
		if r.completed {
			icon = completedStyle.Render("✓")
		} else {
			icon = stoppedStyle.Render("⏺")
		}

		iconCell := lipgloss.NewStyle().Width(4).Render(icon)
		timeCell := dimStyle.Width(maxTimeWidth + 3).Render(r.timeStr)
		durCell := dimStyle.Render(fmt.Sprintf("(%s)", r.durStr))

		line := lipgloss.JoinHorizontal(lipgloss.Top, iconCell, timeCell, durCell)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func endStateKeybindsUI() string {
	var keybindStyle = lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	keybinds := "r: restart · e: edit duration · q: quit"

	return keybindStyle.Render(keybinds)
}
