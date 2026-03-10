package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	applicationstate "pomo-tui/application-state"
	"pomo-tui/storage"
)

func main() {
	debugStateFlag := flag.String("debug-state", "", "boot into state: edit|running|complete|stopped|break")
	sessionMinsFlag := flag.Int("session-duration", 25, "session timer limit in minutes")
	breakMinsFlag := flag.Int("break-duration", 5, "break timer limit in minutes")
	elapsedFlag := flag.String("elapsed", "5m", "elapsed time in the active timer phase (e.g. 5m, 90s)")
	flag.Parse()

	store, err := storage.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not open database: %v\n", err)
	} else {
		defer store.Close()
	}

	var model tea.Model
	if *debugStateFlag != "" {
		elapsed, err := time.ParseDuration(*elapsedFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --elapsed value %q: %v\n", *elapsedFlag, err)
			os.Exit(1)
		}
		model = applicationstate.DebugModel(applicationstate.DebugConfig{
			State:           applicationstate.DebugState(*debugStateFlag),
			SessionDuration: time.Duration(*sessionMinsFlag) * time.Minute,
			BreakDuration:   time.Duration(*breakMinsFlag) * time.Minute,
			Elapsed:         elapsed,
		}, store)
	} else {
		model = applicationstate.InitialModel(store)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
