package applicationstate

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func (m Model) timerStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if m.sessionPauseTimer.IsPaused() {
				// Resume: accumulate the paused time
				m.sessionPauseTimer.UnPause()
				return m, doTick()
			} else {
				m.sessionPauseTimer.Pause()
				return m, nil // stop ticking
			}

		case "x":
			m.sessionEnd = time.Now()
			m.breakStart = time.Now()
			m.breakPauseTimer.Pause()
			m.programState = breakState
			m.persistSession()
			return m, nil
		}

	case tickMsg:
		// Only process ticks if not paused
		if m.sessionPauseTimer.IsPaused() {
			return m, nil
		}

		d := time.Since(m.sessionStart) - m.sessionPauseTimer.PausedDuration()

		if d >= m.timeLimit {
			m.sessionEnd = time.Now()
			m.programState = sessionCompleteState
			m.persistSession()
			return m, nil
		}

		return m, doTick()

	case spinner.TickMsg:
		if m.sessionPauseTimer.IsPaused() {
			return m, nil
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
