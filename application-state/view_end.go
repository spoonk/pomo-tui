package applicationstate

import (
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
	"pomo-tui/storage"
	"pomo-tui/ui"
)

func (m Model) completedView() string {
	endMessageStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.WarnColor)
	var endMessage string
	if m.programState == sessionCompleteState {
		endMessage = endMessageStyle.Render("✓ Session complete")
	} else {
		endMessage = endMessageStyle.Render("⏺  Session stopped early")
	}

	infoLine := elapsedTimeUI(m)
	hints := ui.Keybinds("r: restart · e: edit duration · q: quit")
	list := ui.SessionList(toSessionListEntries(m.store))

	content := lipgloss.JoinVertical(lipgloss.Center, endMessage, infoLine, list)
	content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	centeredHints := lipgloss.Place(m.width, 0, lipgloss.Center, 0, hints)

	contentWithHints := ui.CompositeOver(content, centeredHints, 0, m.height-2)

	return contentWithHints
}

func elapsedTimeUI(m Model) string {
	endedAt, _ := m.sessionTimer.EndedAt()
	return ui.ElapsedTime(m.sessionTimer.StartedAt(), endedAt, m.sessionTimer.TotalDuration())
}

// toSessionListEntries converts storage sessions into ui.SessionListEntry values,
// filtering to today's sessions only.
func toSessionListEntries(store storage.Store) []ui.SessionListEntry {
	if store == nil {
		return nil
	}
	sessions := store.GetSessions()

	now := time.Now()
	todayYear, todayMonth, todayDay := now.Date()

	var entries []ui.SessionListEntry
	for _, s := range sessions {
		y, mo, d := s.StartedAt.In(time.Local).Date()
		if y != todayYear || mo != todayMonth || d != todayDay {
			continue
		}
		entries = append(entries, ui.SessionListEntry{
			Completed:   s.DurationSeconds >= s.PlannedSeconds,
			ProjectName: s.Project.Name,
			StartedAt:   s.StartedAt.In(time.Local),
			EndedAt:     s.EndedAt.In(time.Local),
			Duration:    time.Duration(s.DurationSeconds) * time.Second,
		})
	}

	return entries
}
