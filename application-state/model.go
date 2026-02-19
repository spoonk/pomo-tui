package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type tickMsg time.Time

type Model struct {
	sessionStart            time.Time
	sessionTimeLimitSeconds time.Duration
	sessionBreakTimeSeconds time.Duration

	sessionRunning bool
}

func InitialModel() Model {
	return Model{
		sessionStart:            time.Now(),
		sessionTimeLimitSeconds: 25 * time.Minute,
		sessionBreakTimeSeconds: 5 * time.Minute,
		sessionRunning:          false,
	}
}

func doTick() tea.Cmd {
	return tea.Every(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(doTick(), tea.ClearScreen)
}
