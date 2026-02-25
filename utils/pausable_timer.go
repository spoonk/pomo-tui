package utils

import "time"

// this utility should let me represent an entire
// session (including pauses)
type PausableTimer interface {
	Pause()
	UnPause()
	Reset()
	Start()
	GetTotalDuration() time.Duration
	GetUnpausedDuration() time.Duration
	SetDuration(newDuration time.Duration)
	IsExpired() bool
}

type pausableTimer struct {
	pauseTimer      PauseTimer
	sessionStart    time.Time
	sessionDuration time.Duration
}

func (p *pausableTimer) Pause() {
	p.pauseTimer.Pause()
}

func (p *pausableTimer) UnPause() {
	p.pauseTimer.UnPause()
}

func (p *pausableTimer) Reset() {
	p.pauseTimer.Reset()
	p.sessionDuration = 0 * time.Second
}

func (p *pausableTimer) Start() {
	p.sessionStart = time.Now()
}

func (p *pausableTimer) GetTotalDuration() time.Duration {
	return time.Since(p.sessionStart)
}

func (p *pausableTimer) GetUnpausedDuration() time.Duration {
	pausedDuration := p.pauseTimer.pausedDuration
	return time.Since(p.sessionStart) - pausedDuration
}

func (p *pausableTimer) SetDuration(newDuration time.Duration) {
	p.sessionDuration = newDuration
}

func (p *pausableTimer) IsExpired() bool {
	unpausedDuration := p.GetUnpausedDuration()
	return unpausedDuration >= p.sessionDuration

}

func NewPausableTiemr() PausableTimer {
	return &pausableTimer{pauseTimer: NewPauseTimer()}
}
