package applicationstate

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

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
			m.sessionPauseTimer.Reset()
			return m, doTick()
		case "e":
			m.programState = editTimerState
			m.sessionPauseTimer.Reset()
			return m, textinput.Blink
		}
	}

	return m, cmd
}
