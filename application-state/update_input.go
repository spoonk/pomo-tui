package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) inputStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.breakLimitInput.Focused() {
				m.sessionTimer.SetDuration(parseMins(m.timeLimitInput.Value(), 25))
				m.breakTimer.SetDuration(parseMins(m.breakLimitInput.Value(), 5))
				m.sessionTimer.Start()
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
