package applicationstate

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type Model struct {
	sessionStart time.Time
	timeLimit    time.Duration
	breakLimit   time.Duration

	timeLimitInput  textinput.Model
	breakLimitInput textinput.Model

	sessionRunning bool
}

func InitialModel() Model {
	sessionInput := textinput.New()
	sessionInput.Focus()

	breakInput := textinput.New()

	return Model{
		sessionStart:    time.Now(),
		timeLimit:       25 * time.Minute,
		breakLimit:      5 * time.Minute,
		timeLimitInput:  sessionInput,
		breakLimitInput: breakInput,
		sessionRunning:  false,
	}
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
