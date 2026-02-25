package utils

import "time"

type PauseTimer struct {
	isPaused       bool
	pausedAt       time.Time
	pausedDuration time.Duration
}

func NewPauseTimer() PauseTimer {
	return PauseTimer{isPaused: false, pausedDuration: 0}
}

func (p *PauseTimer) IsPaused() bool {
	return p.isPaused
}

func (p *PauseTimer) PausedAt() time.Time {
	return p.pausedAt
}

func (p *PauseTimer) PausedDuration() time.Duration {
	return p.pausedDuration
}

func (p *PauseTimer) OverridePausedAt(pausedAtOverride time.Time) *PauseTimer {
	p.pausedAt = pausedAtOverride
	return p
}

func (p *PauseTimer) OverridePauseDuration(durationOverride time.Duration) *PauseTimer {
	p.pausedDuration = durationOverride
	return p
}

func (p *PauseTimer) Pause() {
	if p.isPaused {
		return
	}

	p.isPaused = true
	p.pausedAt = time.Now()
}

func (p *PauseTimer) UnPause() {
	if !p.isPaused {
		return
	}

	p.isPaused = false
	p.pausedDuration += time.Since(p.pausedAt)
}

func (p *PauseTimer) Reset() {
	p.isPaused = false
	p.pausedDuration = 0
}
