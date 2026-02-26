package utils

import "time"

// this utility should let me represent an entire
// session (including pauses)
type PausableTimer interface {
	Pause()
	UnPause()
	Reset()
	Start()
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetTotalDuration() time.Duration
	GetUnpausedDuration() time.Duration
	SetDuration(newDuration time.Duration)
	EndEarly()
	IsExpired() bool
	IsPaused() bool
	OverrideStartTime(newStartTime time.Time)
	OverridePausedDuration(newDuration time.Duration)
	GetTimeLimit() time.Duration
	PausedAt() time.Time
	OverridePausedAt(time.Time)
	PausedDuration() time.Duration
}

type pausableTimer struct {
	pauseTimer       PauseTimer
	sessionStart     time.Time
	sessionTimeLimit time.Duration
	isTerminated     bool
	terminatedAt     time.Time
}

func (p *pausableTimer) Pause() {
	p.pauseTimer.Pause()
}

func (p *pausableTimer) EndEarly() {
	p.isTerminated = true
	p.terminatedAt = time.Now()
	p.pauseTimer.UnPause()
}

func (p *pausableTimer) GetStartTime() time.Time {
	return p.sessionStart
}

func (p *pausableTimer) UnPause() {
	p.pauseTimer.UnPause()
}

func (p *pausableTimer) Reset() {
	p.pauseTimer.Reset()
	p.isTerminated = false
}

func (p *pausableTimer) Start() {
	p.sessionStart = time.Now()
}

func (p *pausableTimer) GetTotalDuration() time.Duration {
	if p.IsExpired() {
		return p.sessionTimeLimit + p.pauseTimer.pausedDuration
	}

	if p.isTerminated {
		return p.terminatedAt.Sub(p.sessionStart)
	}

	return time.Since(p.sessionStart)
}

func (p *pausableTimer) GetEndTime() time.Time { // TODO: should return an error
	if !p.IsExpired() {
		return time.Now().Add(999 * time.Hour)
	}

	if p.isTerminated {
		return p.terminatedAt
	}

	return p.sessionStart.Add(p.GetTotalDuration())
}

func (p *pausableTimer) GetUnpausedDuration() time.Duration {
	if p.isTerminated {
		return p.terminatedAt.Sub(p.sessionStart) - p.pauseTimer.pausedDuration
	}

	return time.Since(p.sessionStart) - p.pauseTimer.PausedDuration()
}

func (p *pausableTimer) SetDuration(newDuration time.Duration) {
	p.sessionTimeLimit = newDuration
}

func (p *pausableTimer) IsExpired() bool {
	unpausedDuration := p.GetUnpausedDuration()
	return unpausedDuration >= p.sessionTimeLimit
}

func (p *pausableTimer) IsPaused() bool {
	return p.pauseTimer.isPaused
}

func (p *pausableTimer) OverrideStartTime(newStartTime time.Time) {
	p.sessionStart = newStartTime
}

func (p *pausableTimer) GetTimeLimit() time.Duration {
	return p.sessionTimeLimit
}

func (p *pausableTimer) OverridePausedDuration(newDuration time.Duration) {
	p.pauseTimer.pausedDuration = newDuration
}

func (p *pausableTimer) PausedAt() time.Time {
	return p.pauseTimer.pausedAt
}

func (p *pausableTimer) OverridePausedAt(newPausedAt time.Time) {
	p.pauseTimer.pausedAt = newPausedAt
}

func (p *pausableTimer) PausedDuration() time.Duration {
	return p.pauseTimer.pausedDuration
}

func NewPausableTimer() PausableTimer {
	return &pausableTimer{pauseTimer: NewPauseTimer(), isTerminated: false}
}
