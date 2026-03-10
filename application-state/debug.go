// debug.go — debug boot scaffolding; not intended for production use.
package applicationstate

import (
	"time"

	"pomo-tui/storage"
	"pomo-tui/utils"
)

type DebugState string

const (
	DebugStateEdit     DebugState = "edit"
	DebugStateRunning  DebugState = "running"
	DebugStateComplete DebugState = "complete"
	DebugStateStopped  DebugState = "stopped"
	DebugStateBreak    DebugState = "break"
)

type DebugConfig struct {
	State           DebugState
	SessionDuration time.Duration
	BreakDuration   time.Duration
	Elapsed         time.Duration // time elapsed in the active phase (running/break/stopped)
}

func DebugModel(cfg DebugConfig, store storage.Store) Model {
	m := InitialModel(store)

	sessionDur := orDefault(cfg.SessionDuration, 25*time.Minute)
	breakDur := orDefault(cfg.BreakDuration, 5*time.Minute)
	elapsed := orDefault(cfg.Elapsed, 5*time.Minute)

	switch cfg.State {
	case DebugStateRunning:
		m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
			Start: time.Now().Add(-elapsed),
			Limit: sessionDur,
		})
		m.breakTimer.SetDuration(breakDur)
		m.programState = sessionRunningState

	case DebugStateBreak:
		// Session timer expired (session just ended).
		m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
			Start: time.Now().Add(-sessionDur),
			Limit: sessionDur,
		})
		// Break timer is actively running.
		m.breakTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
			Start: time.Now().Add(-elapsed),
			Limit: breakDur,
		})
		m.programState = breakState

	case DebugStateComplete:
		// Session ran to completion; ended just now.
		m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
			Start: time.Now().Add(-sessionDur),
			Limit: sessionDur,
		})
		m.programState = sessionCompleteState

	case DebugStateStopped:
		// Session ran for `elapsed` then was stopped early.
		m.sessionTimer = utils.NewPausableTimerFromConfig(utils.TimerConfig{
			Start: time.Now().Add(-elapsed),
			Limit: sessionDur,
		})
		m.sessionTimer.Stop()
		m.programState = sessionEndedEarlyState
	}

	return m
}

func orDefault(d, fallback time.Duration) time.Duration {
	if d == 0 {
		return fallback
	}
	return d
}
