package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
)

type changeState string

func (m Model) inputStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.breakLimitInput.Focused() {
				// move to next state
				m.sessionRunning = true
				return m, doTick()
			} else if m.timeLimitInput.Focused() {
				m.timeLimitInput.Blur()
				return m, m.breakLimitInput.Focus()
			}
		}
	}

	if m.timeLimitInput.Focused() {
		m.timeLimitInput, cmd = m.timeLimitInput.Update(msg)
	} else {
		m.breakLimitInput, cmd = m.breakLimitInput.Update(msg)
	}

	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.sessionRunning {
		return m.inputStateUpdate(msg)
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tickMsg:
		return m, doTick()
	}

	return m, cmd
}
