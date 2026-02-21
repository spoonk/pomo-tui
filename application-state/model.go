package applicationstate

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type appState int

const (
	initialState appState = iota
	sessionRunningState
	sessionEndedState
)

type Model struct {
	programState appState
	sessionStart time.Time
	timeLimit    time.Duration
	breakLimit   time.Duration

	timeLimitInput  textinput.Model
	breakLimitInput textinput.Model
}

func InitialModel() Model {
	sessionInput := textinput.New()
	sessionInput.Focus()
	sessionInput.Placeholder = "25"
	sessionInput.CharLimit = 3
	sessionInput.Width = 5
	sessionInput.Validate = validateMinutes

	breakInput := textinput.New()
	breakInput.Placeholder = "5"
	breakInput.CharLimit = 3
	breakInput.Width = 5
	breakInput.Validate = validateMinutes

	return Model{
		programState:    initialState,
		sessionStart:    time.Now(),
		timeLimit:       25 * time.Minute,
		breakLimit:      5 * time.Minute,
		timeLimitInput:  sessionInput,
		breakLimitInput: breakInput,
	}
}

func validateMinutes(s string) error {
	for _, c := range s {
		if c < '0' || c > '9' {
			return fmt.Errorf("invalid character: %c", c)
		}
	}
	return nil
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Every(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(doTick(), tea.ClearScreen, textinput.Blink)
}
