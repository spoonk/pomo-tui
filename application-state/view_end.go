package applicationstate

import (
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
	"pomo-tui/storage"
	"pomo-tui/ui"
)

func (m Model) completedView() string {
	var endMessage string
	if m.programState == sessionCompleteState {
		endMessageStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.SuccessColor)
		endMessage = endMessageStyle.Render("✓ Session complete")
	} else {
		endMessageStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.WarnColor)
		endMessage = endMessageStyle.Render("⏺  Session stopped early")
	}

	infoLine := elapsedTimeUI(m)
	hints := ui.Keybinds("r: restart · e: edit duration · q: quit")

	centeredContent := lipgloss.JoinVertical(lipgloss.Center, endMessage, infoLine)
	contentH := lipgloss.Height(centeredContent)

	centeredContent = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, centeredContent)
	centeredHints := lipgloss.Place(m.width, 0, lipgloss.Center, 0, hints)
	// list := ui.SessionList(toSessionListEntries(m.store))
	list := ui.GroupedSessionList(toGroupedSessionListEntry(m.store))
	centeredList := lipgloss.Place(m.width, 0, lipgloss.Center, 0, list)

	withList := ui.CompositeOver(centeredContent, centeredList, 0, (m.height+contentH)/2+2)
	withHints := ui.CompositeOver(withList, centeredHints, 0, m.height-2)

	return withHints
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

func toGroupedSessionListEntry(store storage.Store) []ui.GroupedSessionListEntry {
	if store == nil {
		return nil
	}
	sessions := store.GetSessions()

	now := time.Now()
	todayYear, todayMonth, todayDay := now.Date()

	var entries []ui.GroupedSessionListEntry
	for _, s := range sessions {
		y, mo, d := s.StartedAt.In(time.Local).Date()
		if y != todayYear || mo != todayMonth || d != todayDay {
			continue
		}
		entries = append(entries, ui.GroupedSessionListEntry{
			Completed:   s.DurationSeconds >= s.PlannedSeconds,
			ProjectName: &s.Project.Name, // TODO: will be a bug when I let people not select a project
			StartedAt:   s.StartedAt.In(time.Local),
			EndedAt:     s.EndedAt.In(time.Local),
			Duration:    time.Duration(s.DurationSeconds) * time.Second,
		})
	}

	return entries
}
