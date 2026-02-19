package applicationstate

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type tickMsg time.Time

type Model struct {
	sessionStartSeconds     time.Time
	sessionTimeLimitSeconds float64
	sessionBreakTimeSeconds float64

	sessionRunning bool
}

func InitialModel() Model {
	return Model{
		sessionStartSeconds:     time.Now(),
		sessionTimeLimitSeconds: 25 * 60,
		sessionBreakTimeSeconds: 5 * 60,
		sessionRunning:          false,
	}
}

func doTick() tea.Cmd {
	return tea.Every(20*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return doTick()
}
