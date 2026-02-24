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

Five distinct program states (appState enum):
- **initialState**: User inputs session and break lengths
- **sessionRunningState**: Timer actively counting up (can be paused/resumed)
- **sessionCompleteState**: Timer reached the full planned duration
- **sessionEndedEarlyState**: User stopped the session early with `x`
- **sessionEndedCanceledState**: Reserved (unused)

The running state supports pause/resume functionality via the `isPaused` flag. When paused, tick messages are ignored and the `pausedDuration` accumulates time spent paused, which is subtracted from elapsed time calculations.

### Package Structure

**application-state package** contains all business logic:

- **model.go**:
  - Defines the `Model` struct (program state, timers, pause state, UI components)
  - Pause tracking fields: `isPaused`, `pausedAt`, `pausedDuration`
  - `InitialModel()` constructor
  - `validateMinutes()` input validation
  - `doTick()` creates timer tick commands
  - `Init()` initializes the Bubble Tea program

- **update.go**:
  - Main `Update()` dispatcher handles all messages
  - `inputStateUpdate()` handles user input for session configuration
  - `timerStateUpdate()` processes tick messages, pause/resume ('p' key), and completion logic
  - `endStateUpdate()` handles end screen interactions
  - `parseMins()` helper for parsing minute inputs with fallback defaults
  - Pause/resume logic: toggles `isPaused`, tracks `pausedDuration`, stops/resumes tick commands

- **view.go**:
  - Main `View()` dispatcher renders appropriate state
  - `inputView()` renders session/break input screen with focus indicators
  - `runningSessionView()` renders active timer display with pause indicator and keybind hints
  - `completedView()` / `stoppedEarlyView()` render the end screen with time summary
  - `elapsedTimeUI()` shows start/end times and total duration
  - `sessionListUI()` renders the historical session list from the store
  - `formatDuration()` formats time.Duration into readable strings
  - Elapsed time calculation accounts for `pausedDuration` to show accurate progress

**storage package** handles SQLite persistence via GORM:

- **session.go**: `Session` GORM model with fields `ID`, `StartedAt`, `EndedAt`, `DurationSeconds` (active time, pauses excluded), `PlannedSeconds`
- **store.go**:
  - `Store` interface with `SaveSession()`, `GetSessions()`, `Close()`
  - `sqliteStore` implementation backed by GORM/SQLite
  - `Open()` constructor — creates DB at `~/.local/share/pomo-tui/sessions.db` (XDG convention), runs `AutoMigrate`
  - GORM driver: `gorm.io/driver/sqlite` (requires CGO / mattn/go-sqlite3)

**main.go**: Entry point — calls `storage.Open()`, injects the store into `InitialModel(store)`, and runs the Bubble Tea program. If `Open()` fails, `store` is nil and the TUI runs without persistence.

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

- `j`/`down`, `k`/`up`: Navigate between inputs in initial state (vim-style)
- `enter`: Advance to next input or start session
- `p`: Pause/resume timer in running state
- `q`/`ctrl+c`: Quit at any time

## Testing

Uses testify for assertions and gotestsum for pretty output. Tests are comprehensive and cover:
- Model initialization
- Input validation
- State transitions
- Navigation between inputs
- Timer completion logic
- Window size handling across all states
