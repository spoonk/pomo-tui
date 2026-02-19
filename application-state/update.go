package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

	m.sessionTimeLimitInput, cmd = m.sessionTimeLimitInput.Update(msg)
	return m, cmd

}
