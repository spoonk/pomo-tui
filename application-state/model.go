package applicationstate

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pomo-tui/storage"
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
	sessionEnd   time.Time
	timeLimit    time.Duration
	breakLimit   time.Duration

	isPaused       bool
	pausedAt       time.Time
	pausedDuration time.Duration

	width  int
	height int

	timeLimitInput  textinput.Model
	breakLimitInput textinput.Model
	spinner         spinner.Model

	// store is injected at construction time. It may be nil if the database
	// could not be opened
	store storage.Store
}

func InitialModel(store storage.Store) Model {
	sessionInput := textinput.New()
	sessionInput.Focus()
	sessionInput.Prompt = ""
	sessionInput.Placeholder = "25"
	sessionInput.CharLimit = 3
	sessionInput.Width = 5
	sessionInput.Validate = validateMinutes

	breakInput := textinput.New()
	breakInput.Prompt = ""
	breakInput.Placeholder = "5"
	breakInput.CharLimit = 3
	breakInput.Width = 5
	breakInput.Validate = validateMinutes

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		programState:    initialState,
		sessionStart:    time.Now(),
		timeLimit:       25 * time.Minute,
		breakLimit:      5 * time.Minute,
		timeLimitInput:  sessionInput,
		breakLimitInput: breakInput,
		spinner:         s,
		store:           store,
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

type tickMsg string

func doTick() tea.Cmd {
	return tea.Every(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg("whatever")
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textinput.Blink)
}
