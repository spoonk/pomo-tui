package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	applicationstate "pomo-tui/application-state"
)

func main() {
	p := tea.NewProgram(applicationstate.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
