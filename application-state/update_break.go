package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func (m Model) breakStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if m.breakPauseTimer.IsPaused() {
				m.breakPauseTimer.UnPause()
				return m, doTick()
			} else {
				m.breakPauseTimer.Pause()
				return m, nil // stop ticking
			}

		case "x":
			m.breakEnd = time.Now()
			m.programState = sessionEndedEarlyState
			return m, nil
		}

	case tickMsg:
		if m.breakPauseTimer.IsPaused() {
			return m, nil
		}

		d := time.Since(m.breakStart) - m.breakPauseTimer.PausedDuration()

		if d >= m.breakLimit {
			m.breakEnd = time.Now()
			m.programState = sessionCompleteState // TODO: new session here
			// TODO: should I persist breaks or not?
			return m, nil
		}

		return m, doTick()
	}

	return m, nil
}
