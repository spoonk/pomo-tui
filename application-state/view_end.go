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
	infoLine  := elapsedTimeUI(m)
	keybinds  := ui.Keybinds("r: restart · e: edit duration · q: quit")
	grid      := ui.SessionGrid(toSessionEntries(m.store))

	content := checkLine + "\n" + infoLine + "\n" + keybinds + "\n\n" + grid
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) stoppedEarlyView() string {
	checkStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.WarnColor)

	checkLine := checkStyle.Render("⏺  Session stopped early")
	infoLine  := elapsedTimeUI(m)
	keybinds  := ui.Keybinds("r: restart · e: edit duration · q: quit")
	grid      := ui.SessionGrid(toSessionEntries(m.store))

	content := checkLine + "\n" + infoLine + "\n" + keybinds + "\n\n" + grid
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func elapsedTimeUI(m Model) string {
	endedAt, _ := m.sessionTimer.EndedAt()
	return ui.ElapsedTime(m.sessionTimer.StartedAt(), endedAt, m.breakTimer.TotalDuration())
}

// toSessionEntries converts storage sessions into the plain ui.SessionEntry type,
// keeping storage types out of the ui package.
func toSessionEntries(store storage.Store) []ui.SessionEntry {
	if store == nil {
		return nil
	}
	sessions := store.GetSessions()
	entries := make([]ui.SessionEntry, len(sessions))
	for i, s := range sessions {
		entries[i] = ui.SessionEntry{
			Completed: s.DurationSeconds >= s.PlannedSeconds,
			StartedAt: s.StartedAt,
			Duration:  time.Duration(s.DurationSeconds) * time.Second,
			Planned:   time.Duration(s.PlannedSeconds) * time.Second,
		}
	}
	return entries
}
