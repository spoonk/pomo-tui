package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
	"pomo-tui/storage"
	"strconv"
	"time"
)

type changeState string

func parseMins(s string, fallback int) time.Duration {
	if s == "" {
		return time.Duration(fallback) * time.Minute
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return time.Duration(fallback) * time.Minute
	}
	return time.Duration(n) * time.Minute
}

func (m Model) persistSession() {
	if m.store == nil {
		return
	}
	activeDuration := m.sessionEnd.Sub(m.sessionStart) - m.sessionPauseTimer.PausedDuration()
	session := &storage.Session{
		StartedAt:       m.sessionStart.UTC(),
		EndedAt:         m.sessionEnd.UTC(),
		DurationSeconds: int64(activeDuration.Seconds()),
		PlannedSeconds:  int64(m.timeLimit.Seconds()),
	}
	// Intentionally ignoring the error for now; the TUI stays usable even if
	// persistence fails. TODO: surface errors in endView if desired.
	_ = m.store.SaveSession(session)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// universal messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.programState {
	case editTimerState:
		return m.inputStateUpdate(msg)
	case sessionRunningState:
		return m.timerStateUpdate(msg)
	case breakState:
		return m.breakStateUpdate(msg)

	default:
		return m.endStateUpdate(msg)
	}
}
