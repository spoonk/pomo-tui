package applicationstate

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
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

func (m Model) endStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m Model) timerStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tickMsg:
		d := time.Since(m.sessionStart)
		d = d.Round(time.Second)

		if true || d.Minutes() >= m.timeLimit.Minutes() {
			m.sessionEnd = time.Now()
			m.programState = sessionEndedState
			return m, nil // do nothing
		}

		return m, doTick()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// universal messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.programState {
	case initialState:
		return m.inputStateUpdate(msg)
	case sessionRunningState:
		return m.timerStateUpdate(msg)
	default:
		return m.endStateUpdate(msg)
	}
}
