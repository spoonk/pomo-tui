package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	applicationstate "pomo-tui/application-state"
	"pomo-tui/storage"
)

func main() {
	store, err := storage.Open()
	if err != nil {
		// Log the problem but don't exit: the TUI works fine without persistence.
		fmt.Fprintf(os.Stderr, "warning: could not open session database: %v\n", err)
	} else {
		defer store.Close()
	}

	p := tea.NewProgram(applicationstate.InitialModel(store))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
