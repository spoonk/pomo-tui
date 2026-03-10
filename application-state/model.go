package applicationstate

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pomo-tui/storage"
	"pomo-tui/ui"
	"pomo-tui/utils"
)

type appState int

const (
	editTimerState appState = iota
	sessionRunningState
	sessionCompleteState
	sessionEndedEarlyState
	sessionEndedCanceledState
	breakState
)

type Model struct {
	programState appState

	sessionTimer utils.PausableTimer
	breakTimer   utils.PausableTimer

	width  int
	height int

	timeLimitInput  textinput.Model
	breakLimitInput textinput.Model
	spinner         spinner.Model

	projectSelector ui.ProjectSelector
	selectedProject *ui.ProjectEntry

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

	sessionTimer := utils.NewPausableTimer()
	sessionTimer.SetDuration(25 * time.Minute)
	breakTimer := utils.NewPausableTimer()
	breakTimer.SetDuration(5 * time.Minute)

	var projectEntries []ui.ProjectEntry
	if store != nil {
		for _, p := range store.GetProjects() {
			projectEntries = append(projectEntries, ui.ProjectEntry{
				ID:   p.ID,
				Name: p.Name,
			})
		}
	}

	return Model{
		programState:    editTimerState,
		sessionTimer:    sessionTimer,
		breakTimer:      breakTimer,
		timeLimitInput:  sessionInput,
		breakLimitInput: breakInput,
		spinner:         s,
		projectSelector: ui.NewProjectSelector(projectEntries),
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
	switch m.programState {
	case sessionRunningState, breakState:
		return tea.Batch(tea.ClearScreen, doTick())
	default:
		return tea.Batch(tea.ClearScreen, textinput.Blink)
	}
}
