# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**pomo-tui** is a Pomodoro timer TUI (Terminal User Interface) application built with Go and Bubble Tea. It follows the Elm Architecture pattern (Model-Update-View).

## Development Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests (uses gotestsum with testdox format)
make test

# Run a specific test
go test -run TestName ./application-state

# Run a specific test with verbose output
go test -v -run TestName ./application-state
```

## Architecture

The application uses the **Elm Architecture** via Bubble Tea:

### State Machine

Three distinct program states (appState enum):
- **initialState**: User inputs session and break lengths
- **sessionRunningState**: Timer actively counting up
- **sessionEndedState**: Session complete, showing summary

### Package Structure

**application-state package** contains all business logic:

- **model.go**:
  - Defines the `Model` struct (program state, timers, UI components)
  - `InitialModel()` constructor
  - `validateMinutes()` input validation
  - `doTick()` creates timer tick commands
  - `Init()` initializes the Bubble Tea program

- **update.go**:
  - Main `Update()` dispatcher handles all messages
  - `inputStateUpdate()` handles user input for session configuration
  - `timerStateUpdate()` processes tick messages and completion logic
  - `endStateUpdate()` handles end screen interactions
  - `parseMins()` helper for parsing minute inputs with fallback defaults

- **view.go**:
  - Main `View()` dispatcher renders appropriate state
  - `inputView()` renders session/break input screen with focus indicators
  - `runningSessionView()` renders active timer display
  - `endView()` renders completion screen with time summary
  - `formatDuration()` formats time.Duration into readable strings

**main.go**: Entry point that instantiates and runs the Bubble Tea program

### Key Patterns

- State transitions are unidirectional: initial → running → ended
- All state-specific logic is separated into dedicated update/view functions
- Window size (width/height) is tracked universally across all states
- Timer uses `tea.Every()` to generate periodic tick messages (200ms intervals)
- Input validation happens at the textinput level via `validateMinutes()`

### UI Components

- **textinput**: Two inputs for session and break duration (max 3 digits)
- **spinner**: Animated indicator (currently unused in view)
- **lipgloss**: All styling and layout centering

### Navigation

- `j`/`down`, `k`/`up`: Navigate between inputs (vim-style)
- `enter`: Advance to next input or start session
- `q`/`ctrl+c`: Quit at any time

## Testing

Uses testify for assertions and gotestsum for pretty output. Tests are comprehensive and cover:
- Model initialization
- Input validation
- State transitions
- Navigation between inputs
- Timer completion logic
- Window size handling across all states
