package applicationstate

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"pomo-tui/utils"
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
	m := InitialModel(nil)

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
			m := InitialModel(nil)

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
	m := InitialModel(nil)
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
	assert.Equal(t, 30*time.Minute, m.sessionTimer.TimeLimit(), "should parse time limit from input")
	assert.Equal(t, 10*time.Minute, m.breakTimer.TimeLimit(), "should parse break limit from input")

	// Check session start time
	assert.WithinDuration(t, startTime, m.sessionTimer.StartedAt(), 1*time.Second,
		"session start should be set to approximately now")

	// Check that commands were batched
	assert.NotNil(t, cmd, "should return batched commands for tick and spinner")
}

// TestInputStateUpdate_EnterOnTimeInput tests that enter on time input moves to break input
func TestInputStateUpdate_EnterOnTimeInput(t *testing.T) {
	m := InitialModel(nil)
	assert.True(t, m.timeLimitInput.Focused(), "time input should start focused")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.inputStateUpdate(msg)
	m = updatedModel.(Model)

	// Should still be in initial state
	assert.Equal(t, editTimerState, m.programState, "should remain in initialState")

	// Should have moved focus to break input
	assert.False(t, m.timeLimitInput.Focused(), "time input should no longer be focused")
	assert.True(t, m.breakLimitInput.Focused(), "break input should now be focused")
}

// TestTimerStateUpdate_Completion tests timer completion transition
func TestTimerStateUpdate_Completion(t *testing.T) {
	m := InitialModel(nil)
	m.programState = sessionRunningState
	// Started 26 minutes ago with a 25-minute limit → already expired by 1 minute
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start: time.Now().Add(-26 * time.Minute),
		Limit: 25 * time.Minute,
	})

	msg := tickMsg("tick")

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Session expiry triggers a break, not the complete state directly
	assert.Equal(t, breakState, m.programState, "should transition to breakState on expiry")

	// EndedAt should be ~1 minute ago (start + 25min limit = now - 1min)
	endedAt, ok := m.sessionTimer.EndedAt()
	assert.True(t, ok, "timer should have an end time after expiry")
	assert.WithinDuration(t, time.Now().Add(-1*time.Minute), endedAt, 5*time.Second,
		"session end should be ~1 minute ago")

	// Should not return a command (timer stops ticking)
	assert.Nil(t, cmd, "should not return a command after completion")
}

// TestTimerStateUpdate_StillRunning tests that timer continues when time remains
func TestTimerStateUpdate_StillRunning(t *testing.T) {
	m := InitialModel(nil)
	m.programState = sessionRunningState
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

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
		{"initial state", editTimerState},
		{"running state", sessionRunningState},
		{"ended state", sessionCompleteState},
	}

	for _, tt := range states {
		t.Run(tt.name, func(t *testing.T) {
			m := InitialModel(nil)
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
	m := InitialModel(nil)
	m.programState = sessionCompleteState

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

// TestTimerStateUpdate_Pause tests that pressing 'p' pauses the timer
func TestTimerStateUpdate_Pause(t *testing.T) {
	m := InitialModel(nil)
	m.programState = sessionRunningState
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")}

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Should be paused
	assert.True(t, m.sessionTimer.IsPaused(), "should be paused after pressing 'p'")

	// Elapsed should be frozen at ~5 minutes
	assert.InDelta(t, float64(5*time.Minute), float64(m.sessionTimer.Elapsed()), float64(200*time.Millisecond),
		"elapsed should be frozen at ~5 minutes while paused")

	// Should not return a tick command
	assert.Nil(t, cmd, "should not return a tick command when paused")
}

// TestTimerStateUpdate_Resume tests that pressing 'p' while paused resumes the timer
func TestTimerStateUpdate_Resume(t *testing.T) {
	m := InitialModel(nil)
	m.programState = sessionRunningState
	// Started 5 min ago, paused at the 3-min mark (pausedAt = now-2min),
	// with 1 min of prior accumulated pause. Active elapsed should be 5-1-2 = 2 min.
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start:       time.Now().Add(-5 * time.Minute),
		Limit:       25 * time.Minute,
		IsPaused:    true,
		PausedAt:    time.Now().Add(-2 * time.Minute),
		PausedTotal: 1 * time.Minute,
	})

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")}

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Should no longer be paused
	assert.False(t, m.sessionTimer.IsPaused(), "should not be paused after pressing 'p' while paused")

	// After resuming: total paused = 1min prior + 2min just ended = 3min
	// Active elapsed = 5min elapsed - 3min paused = 2min
	assert.InDelta(t, float64(2*time.Minute), float64(m.sessionTimer.Elapsed()), float64(200*time.Millisecond),
		"elapsed should reflect total active time (~2 minutes)")

	// Should return a tick command to resume
	assert.NotNil(t, cmd, "should return a tick command when resuming")
}

// TestTimerStateUpdate_PausedTick tests that tick messages are ignored when paused
func TestTimerStateUpdate_PausedTick(t *testing.T) {
	m := InitialModel(nil)
	m.programState = sessionRunningState
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start:    time.Now().Add(-5 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		PausedAt: time.Now(),
	})

	msg := tickMsg("tick")

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Should still be in running state (not completed)
	assert.Equal(t, sessionRunningState, m.programState, "should remain in sessionRunningState")

	// Should still be paused
	assert.True(t, m.sessionTimer.IsPaused(), "should still be paused")

	// Should not return a command
	assert.Nil(t, cmd, "should not return a command while paused")
}

// TestTimerStateUpdate_CompletionWithPause tests timer completion accounting for paused time
func TestTimerStateUpdate_CompletionWithPause(t *testing.T) {
	m := InitialModel(nil)
	m.programState = sessionRunningState
	// Started 30 min ago, 6 min paused → active elapsed = 24 min (should NOT complete)
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start:       time.Now().Add(-30 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 6 * time.Minute,
	})

	msg := tickMsg("tick")

	updatedModel, cmd := m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Should still be running (active elapsed = 24 min < 25 min limit)
	assert.Equal(t, sessionRunningState, m.programState, "should still be running")
	assert.NotNil(t, cmd, "should return a tick command to continue")

	// Now with 32 min elapsed, 6 min paused → active = 26 min (should complete → break)
	m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
		Start:       time.Now().Add(-32 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 6 * time.Minute,
	})

	updatedModel, cmd = m.timerStateUpdate(msg)
	m = updatedModel.(Model)

	// Session expiry triggers a break
	assert.Equal(t, breakState, m.programState, "should transition to breakState on expiry")
	assert.Nil(t, cmd, "should not return a command after completion")
}
