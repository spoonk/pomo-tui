package applicationstate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitialModel verifies that InitialModel(nil) creates a properly initialized model
func TestInitialModel(t *testing.T) {
	m := InitialModel(nil)

	// Check initial state
	assert.Equal(t, initialState, m.programState, "should start in initialState")

	// Check default time limits
	assert.Equal(t, 25*time.Minute, m.timeLimit, "default session should be 25 minutes")
	assert.Equal(t, 5*time.Minute, m.breakLimit, "default break should be 5 minutes")

	// Check that sessionStart is set (should be recent)
	assert.WithinDuration(t, time.Now(), m.sessionStart, 1*time.Second,
		"sessionStart should be initialized to approximately now")

	// Check input configuration
	assert.True(t, m.timeLimitInput.Focused(), "time limit input should be focused initially")
	assert.False(t, m.breakLimitInput.Focused(), "break limit input should not be focused initially")

	// Check input placeholders
	assert.Equal(t, "25", m.timeLimitInput.Placeholder, "time limit placeholder should be '25'")
	assert.Equal(t, "5", m.breakLimitInput.Placeholder, "break limit placeholder should be '5'")

	// Check input character limits
	assert.Equal(t, 3, m.timeLimitInput.CharLimit, "should limit input to 3 characters")
	assert.Equal(t, 3, m.breakLimitInput.CharLimit, "should limit input to 3 characters")
}

// TestValidateMinutes tests the input validation function
func TestValidateMinutes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid single digit",
			input:   "5",
			wantErr: false,
		},
		{
			name:    "valid double digit",
			input:   "25",
			wantErr: false,
		},
		{
			name:    "valid triple digit",
			input:   "999",
			wantErr: false,
		},
		{
			name:    "empty string is valid",
			input:   "",
			wantErr: false,
		},
		{
			name:    "invalid with letter",
			input:   "5a",
			wantErr: true,
		},
		{
			name:    "invalid with special char",
			input:   "5-",
			wantErr: true,
		},
		{
			name:    "invalid with space",
			input:   "5 ",
			wantErr: true,
		},
		{
			name:    "invalid with decimal",
			input:   "5.5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMinutes(tt.input)
			if tt.wantErr {
				assert.Error(t, err, "expected validation error for input: %q", tt.input)
			} else {
				assert.NoError(t, err, "expected no validation error for input: %q", tt.input)
			}
		})
	}
}

// TestDoTick verifies that doTick returns a valid command
func TestDoTick(t *testing.T) {
	cmd := doTick()
	require.NotNil(t, cmd, "doTick should return a non-nil command")

	// Execute the command and verify it returns a tickMsg
	msg := cmd()
	_, ok := msg.(tickMsg)
	assert.True(t, ok, "doTick command should produce a tickMsg")
}
