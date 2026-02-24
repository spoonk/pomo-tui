package applicationstate

import (
	"strconv"
	"time"

	"pomo-tui/storage"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

func (m Model) inputStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.breakLimitInput.Focused() {
				m.timeLimit = parseMins(m.timeLimitInput.Value(), 25)
				m.breakLimit = parseMins(m.breakLimitInput.Value(), 5)
				m.sessionStart = time.Now()
				m.programState = sessionRunningState
				return m, tea.Batch(doTick(), m.spinner.Tick)
			} else if m.timeLimitInput.Focused() {
				m.timeLimitInput.Blur()
				return m, m.breakLimitInput.Focus()
			}

		case "up", "k":
			m.breakLimitInput.Blur()
			return m, m.timeLimitInput.Focus()

		case "down", "j":
			m.timeLimitInput.Blur()
			return m, m.breakLimitInput.Focus()
		}
	}

	if m.timeLimitInput.Focused() {
		m.timeLimitInput, cmd = m.timeLimitInput.Update(msg)
	} else {
		m.breakLimitInput, cmd = m.breakLimitInput.Update(msg)
	}

	return m, cmd
}

func (m Model) timerStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if m.isPaused {
				// Resume: accumulate the paused time
				m.pausedDuration += time.Since(m.pausedAt)
				m.isPaused = false
				return m, doTick()
			} else {
				// Pause: record when we paused
				m.pausedAt = time.Now()
				m.isPaused = true
				return m, nil // stop ticking
			}

		case "x":
			m.sessionEnd = time.Now()
			m.programState = sessionEndedEarlyState
			m.persistSession()
			return m, nil
		}

	case tickMsg:
		// Only process ticks if not paused
		if m.isPaused {
			return m, nil
		}

		d := time.Since(m.sessionStart) - m.pausedDuration

		if d >= m.timeLimit {
			m.sessionEnd = time.Now()
			m.programState = sessionCompleteState
			m.persistSession()
			return m, nil
		}

		return m, doTick()

	case spinner.TickMsg:
		if m.isPaused {
			return m, nil
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) endStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "r":
			m.programState = sessionRunningState
			m.sessionStart = time.Now()
			m.isPaused = false
			m.pausedDuration = 0 * time.Second
			return m, doTick()
		case "e":
			m.programState = editTimerState
			m.isPaused = false
			m.pausedDuration = 0 * time.Second
			return m, textinput.Blink
		}
	}

	return m, cmd
}

func (m Model) persistSession() {
	if m.store == nil {
		return
	}
	activeDuration := m.sessionEnd.Sub(m.sessionStart) - m.pausedDuration
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
	default:
		return m.endStateUpdate(msg)
	}
}
