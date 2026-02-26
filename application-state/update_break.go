package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) breakStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if m.breakTimer.IsPaused() {
				m.breakTimer.Resume()
				return m, doTick()
			} else {
				m.breakTimer.Pause()
				return m, nil // stop ticking
			}

		case "x":
			m.breakTimer.Stop()
			m.programState = sessionEndedEarlyState
			return m, nil
		}

	case tickMsg:
		if m.breakTimer.IsPaused() {
			return m, nil
		}

		if m.breakTimer.IsExpired() {
			m.programState = sessionCompleteState // TODO: new session here
			// TODO: should I persist breaks or not?
			return m, nil
		}

		return m, doTick()
	}

	return m, nil
}
