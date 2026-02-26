package applicationstate

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) timerStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if m.sessionTimer.IsPaused() {
				m.sessionTimer.UnPause()
				return m, doTick()
			} else {
				m.sessionTimer.Pause()
				return m, nil // stop ticking
			}

		case "x":
			m.sessionTimer.EndEarly()
			m.programState = breakState
			m.breakTimer.Reset()
			m.breakTimer.Start()
			m.persistSession()
			return m, nil
		}

	case tickMsg:
		// Only process ticks if not paused
		if m.sessionTimer.IsPaused() {
			return m, nil
		}

		if m.sessionTimer.IsExpired() {
			m.programState = breakState
			m.breakTimer.Reset()
			m.breakTimer.Start()
			m.persistSession()
			return m, nil
		}

		return m, doTick()

	case spinner.TickMsg:
		if m.sessionTimer.IsPaused() {
			return m, nil
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
