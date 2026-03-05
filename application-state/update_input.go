package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
	"pomo-tui/storage"
	"pomo-tui/ui"
)

func (m Model) inputStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ui.ProjectSelectorNavigateUpMsg:
		m.projectSelector = m.projectSelector.Blur()
		return m, m.breakLimitInput.Focus()

	case ui.ProjectCreateRequestMsg:
		entry := ui.ProjectEntry{Name: msg.Name}
		if m.store != nil {
			proj := &storage.Project{Name: msg.Name}
			_ = m.store.SaveProject(proj)
			entry.ID = proj.ID
		}
		m.projectSelector = m.projectSelector.AddProject(entry)
		m.selectedProject = m.projectSelector.Selected()
		return m.startSession()

	case ui.ProjectSelectorEnterMsg:
		m.selectedProject = msg.Selected
		return m.startSession()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// When the project selector is focused, route all key messages to it.
		if m.projectSelector.Focused() {
			m.projectSelector, cmd = m.projectSelector.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "enter":
			if m.breakLimitInput.Focused() {
				m.breakLimitInput.Blur()
				m.projectSelector, cmd = m.projectSelector.Focus()
				return m, cmd
			} else if m.timeLimitInput.Focused() {
				m.timeLimitInput.Blur()
				return m, m.breakLimitInput.Focus()
			} else {
				return m.startSession()
			}

		case "up", "k":
			if m.breakLimitInput.Focused() {
				m.breakLimitInput.Blur()
				return m, m.timeLimitInput.Focus()
			}

		case "down", "j":
			if m.timeLimitInput.Focused() {
				m.timeLimitInput.Blur()
				return m, m.breakLimitInput.Focus()
			} else if m.breakLimitInput.Focused() {
				m.breakLimitInput.Blur()
				m.projectSelector, cmd = m.projectSelector.Focus()
				return m, cmd
			}
		}

		// Pass remaining key messages to the focused text input.
		if m.timeLimitInput.Focused() {
			m.timeLimitInput, cmd = m.timeLimitInput.Update(msg)
		} else if m.breakLimitInput.Focused() {
			m.breakLimitInput, cmd = m.breakLimitInput.Update(msg)
		}

	default:
		// Non-key messages (cursor blink, etc.): route to the focused component.
		if m.projectSelector.Focused() {
			m.projectSelector, cmd = m.projectSelector.Update(msg)
		} else if m.timeLimitInput.Focused() {
			m.timeLimitInput, cmd = m.timeLimitInput.Update(msg)
		} else if m.breakLimitInput.Focused() {
			m.breakLimitInput, cmd = m.breakLimitInput.Update(msg)
		}
	}

	return m, cmd
}

// startSession transitions the model into sessionRunningState using the
// current input values and resets focus on the input fields.
func (m Model) startSession() (tea.Model, tea.Cmd) {
	m.timeLimitInput.Blur()
	m.breakLimitInput.Blur()
	m.projectSelector = m.projectSelector.Blur()
	m.sessionTimer.SetDuration(parseMins(m.timeLimitInput.Value(), 25))
	m.breakTimer.SetDuration(parseMins(m.breakLimitInput.Value(), 5))
	m.sessionTimer.Start()
	m.programState = sessionRunningState
	return m, tea.Batch(doTick(), m.spinner.Tick)
}
