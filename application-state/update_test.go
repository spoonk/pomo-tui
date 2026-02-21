package applicationstate

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestParseMins tests the minute parsing helper function
func TestParseMins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fallback int
		expected time.Duration
	}{
		{
			name:     "empty string uses fallback",
			input:    "",
			fallback: 25,
			expected: 25 * time.Minute,
		},
		{
			name:     "valid number",
			input:    "30",
			fallback: 25,
			expected: 30 * time.Minute,
		},
		{
			name:     "single digit",
			input:    "5",
			fallback: 10,
			expected: 5 * time.Minute,
		},
		{
			name:     "zero uses fallback",
			input:    "0",
			fallback: 25,
			expected: 25 * time.Minute,
		},
		{
			name:     "negative uses fallback",
			input:    "-5",
			fallback: 25,
			expected: 25 * time.Minute,
		},
		{
			name:     "invalid string uses fallback",
			input:    "abc",
			fallback: 25,
			expected: 25 * time.Minute,
		},
		{
			name:     "mixed invalid uses fallback",
			input:    "25abc",
			fallback: 10,
			expected: 10 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMins(tt.input, tt.fallback)
			assert.Equal(t, tt.expected, result,
				"parseMins(%q, %d) should return %v", tt.input, tt.fallback, tt.expected)
		})
	}
}

// TestInputStateUpdate_Quit tests that quit keys work in input state
func TestInputStateUpdate_Quit(t *testing.T) {
	m := InitialModel()

	tests := []struct {
		name string
		key  string
	}{
		{"ctrl+c quits", "ctrl+c"},
		{"q quits", "q"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			_, cmd := m.inputStateUpdate(msg)
			assert.NotNil(t, cmd, "quit command should be returned")
		})
	}
}

// TestInputStateUpdate_Navigation tests focus navigation between inputs
func TestInputStateUpdate_Navigation(t *testing.T) {
	tests := []struct {
		name            string
		key             string
		initialFocus    bool // true = timeLimit focused, false = breakLimit focused
		expectedFocused string
	}{
		{
			name:            "down moves from time to break",
			key:             "down",
			initialFocus:    true,
			expectedFocused: "break",
		},
		{
			name:            "j moves from time to break",
			key:             "j",
			initialFocus:    true,
			expectedFocused: "break",
		},
		{
			name:            "up moves from break to time",
			key:             "up",
			initialFocus:    false,
			expectedFocused: "time",
		},
		{
			name:            "k moves from break to time",
			key:             "k",
			initialFocus:    false,
			expectedFocused: "time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := InitialModel()

			// Set initial focus
			if tt.initialFocus {
				m.timeLimitInput.Focus()
				m.breakLimitInput.Blur()
			} else {
				m.timeLimitInput.Blur()
				m.breakLimitInput.Focus()
			}

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			switch tt.key {
			case "up":
				msg = tea.KeyMsg{Type: tea.KeyUp}
			case "down":
				msg = tea.KeyMsg{Type: tea.KeyDown}
			}

			updatedModel, _ := m.inputStateUpdate(msg)
			m = updatedModel.(Model)

			if tt.expectedFocused == "time" {
				assert.True(t, m.timeLimitInput.Focused(), "time limit input should be focused")
				assert.False(t, m.breakLimitInput.Focused(), "break limit input should not be focused")
			} else {
				assert.False(t, m.timeLimitInput.Focused(), "time limit input should not be focused")
				assert.True(t, m.breakLimitInput.Focused(), "break limit input should be focused")
			}
		})
	}
}

// TestInputStateUpdate_StartSession tests transition from input to running state
func TestInputStateUpdate_StartSession(t *testing.T) {
	m := InitialModel()
	m.timeLimitInput.Blur()
	m.breakLimitInput.Focus()

	// Set some values
	m.timeLimitInput.SetValue("30")
	m.breakLimitInput.SetValue("10")

	// Press enter on break input (should start session)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	startTime := time.Now()

	updatedModel, cmd := m.inputStateUpdate(msg)
	m = updatedModel.(Model)

	// Check state transition
	assert.Equal(t, sessionRunningState, m.programState, "should transition to sessionRunningState")

	// Check time limits were parsed
	assert.Equal(t, 30*time.Minute, m.timeLimit, "should parse time limit from input")
	assert.Equal(t, 10*time.Minute, m.breakLimit, "should parse break limit from input")

	// Check session start time
	assert.WithinDuration(t, startTime, m.sessionStart, 1*time.Second,
		"session start should be set to approximately now")

	// Check that commands were batched
	assert.NotNil(t, cmd, "should return batched commands for tick and spinner")
}

// TestInputStateUpdate_EnterOnTimeInput tests that enter on time input moves to break input
func TestInputStateUpdate_EnterOnTimeInput(t *testing.T) {
	m := InitialModel()
	assert.True(t, m.timeLimitInput.Focused(), "time input should start focused")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.inputStateUpdate(msg)
	m = updatedModel.(Model)

	// Should still be in initial state
	assert.Equal(t, initialState, m.programState, "should remain in initialState")

	// Should have moved focus to break input
	assert.False(t, m.timeLimitInput.Focused(), "time input should no longer be focused")
	assert.True(t, m.breakLimitInput.Focused(), "break input should now be focused")
}

// TestTimerStateUpdate_Completion tests timer completion transition
func TestTimerStateUpdate_Completion(t *testing.T) {
	m := InitialModel()
	m.programState = sessionRunningState
	m.sessionStart = time.Now().Add(-26 * time.Minute) // Started 26 minutes ago
	m.timeLimit = 25 * time.Minute

	// Send a tick message
	msg := tickMsg("tick")
	endTime := time.Now()

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Should transition to ended state
	assert.Equal(t, sessionEndedState, m.programState, "should transition to sessionEndedState")

	// Should set end time
	assert.WithinDuration(t, endTime, m.sessionEnd, 1*time.Second,
		"session end should be set to approximately now")

	// Should not return a command (timer stops)
	assert.Nil(t, cmd, "should not return a command after completion")
}

// TestTimerStateUpdate_StillRunning tests that timer continues when time remains
func TestTimerStateUpdate_StillRunning(t *testing.T) {
	m := InitialModel()
	m.programState = sessionRunningState
	m.sessionStart = time.Now().Add(-5 * time.Minute) // Started 5 minutes ago
	m.timeLimit = 25 * time.Minute

	// Send a tick message
	msg := tickMsg("tick")

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Should still be running
	assert.Equal(t, sessionRunningState, m.programState, "should remain in sessionRunningState")

	// Should return a tick command to continue
	assert.NotNil(t, cmd, "should return a tick command to continue timer")
}

// TestUpdate_WindowSize tests that window size is captured across all states
func TestUpdate_WindowSize(t *testing.T) {
	states := []struct {
		name  string
		state appState
	}{
		{"initial state", initialState},
		{"running state", sessionRunningState},
		{"ended state", sessionEndedState},
	}

	for _, tt := range states {
		t.Run(tt.name, func(t *testing.T) {
			m := InitialModel()
			m.programState = tt.state

			msg := tea.WindowSizeMsg{Width: 120, Height: 40}

			updatedModel, _ := m.Update(msg)
			m = updatedModel.(Model)

			assert.Equal(t, 120, m.width, "width should be updated")
			assert.Equal(t, 40, m.height, "height should be updated")
		})
	}
}

// TestEndStateUpdate_Quit tests that quit keys work in end state
func TestEndStateUpdate_Quit(t *testing.T) {
	m := InitialModel()
	m.programState = sessionEndedState

	tests := []struct {
		name string
		key  string
	}{
		{"ctrl+c quits", "ctrl+c"},
		{"q quits", "q"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			_, cmd := m.endStateUpdate(msg)
			assert.NotNil(t, cmd, "quit command should be returned")
		})
	}
}
