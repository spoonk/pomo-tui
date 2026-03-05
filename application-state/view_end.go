package applicationstate

import (
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
	"pomo-tui/storage"
	"pomo-tui/ui"
)

func (m Model) completedView() string {
	checkStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.SuccessColor)

	checkLine := checkStyle.Render("✓ Session complete")
	infoLine := elapsedTimeUI(m)
	keybinds := ui.Keybinds("r: restart · e: edit duration · q: quit")
	list := ui.SessionList(toSessionListEntries(m.store))

	content := checkLine + "\n" + infoLine + "\n" + keybinds + "\n\n" + list
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) stoppedEarlyView() string {
	checkStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.WarnColor)

	checkLine := checkStyle.Render("⏺  Session stopped early")
	infoLine := elapsedTimeUI(m)
	keybinds := ui.Keybinds("r: restart · e: edit duration · q: quit")
	list := ui.SessionList(toSessionListEntries(m.store))

	content := checkLine + "\n" + infoLine + "\n" + keybinds + "\n\n" + list
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func elapsedTimeUI(m Model) string {
	endedAt, _ := m.sessionTimer.EndedAt()
	return ui.ElapsedTime(m.sessionTimer.StartedAt(), endedAt, m.sessionTimer.TotalDuration())
}

// toSessionListEntries converts storage sessions into ui.SessionListEntry values,
// filtering to today's sessions only and reversing so newest is first.
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

	// Reverse so newest is first
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries
}
