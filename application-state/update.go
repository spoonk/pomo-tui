package applicationstate

import (
	"strconv"
	"time"

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
