package utils

import "time"

// PausableTimer is a countdown timer that supports pausing and early termination.
// Call SetDuration and Start to begin timing. Elapsed time excludes any paused time.
//
// Expected lifecycle: SetDuration → Start → [Pause/Resume]* → [Stop | expires]
// Reset returns the timer to a fresh state (preserving the configured duration).
type PausableTimer interface {
	// Configuration
	SetDuration(d time.Duration)

	// Lifecycle
	Start()
	Pause()
	Resume()
	Stop()  // ends the timer before it expires naturally
	Reset() // resets timing state; preserves the configured duration

	// State
	IsPaused() bool
	IsExpired() bool

	// Time observations
	Elapsed() time.Duration       // active (non-paused) time elapsed so far
	TotalDuration() time.Duration // total wall-clock time elapsed, including pauses
	TimeLimit() time.Duration
	StartedAt() time.Time
	EndedAt() (time.Time, bool) // second return is false if the timer hasn't ended yet
}

type pausableTimer struct {
	start       time.Time
	limit       time.Duration
	isPaused    bool
	pausedAt    time.Time
	pausedTotal time.Duration
	stopped     bool
	stoppedAt   time.Time
}

// NewPausableTimer creates a timer ready for configuration.
func NewPausableTimer() PausableTimer {
	return &pausableTimer{}
}

// TimerConfig holds explicit state for constructing a timer in a specific condition.
// This is primarily useful in tests to avoid manipulating time via sleeps.
type TimerConfig struct {
	Start       time.Time
	Limit       time.Duration
	IsPaused    bool
	PausedAt    time.Time     // only meaningful when IsPaused is true
	PausedTotal time.Duration // accumulated paused time from previous pauses
}

// NewPausableTimerFromConfig creates a timer with pre-populated state.
func NewPausableTimerFromConfig(cfg TimerConfig) PausableTimer {
	t := &pausableTimer{
		start:       cfg.Start,
		limit:       cfg.Limit,
		isPaused:    cfg.IsPaused,
		pausedAt:    cfg.PausedAt,
		pausedTotal: cfg.PausedTotal,
	}
	if t.isPaused && t.pausedAt.IsZero() {
		t.pausedAt = time.Now()
	}
	return t
}

func (p *pausableTimer) SetDuration(d time.Duration) {
	p.limit = d
}

func (p *pausableTimer) Start() {
	p.start = time.Now()
}

func (p *pausableTimer) Pause() {
	if p.isPaused {
		return
	}
	p.isPaused = true
	p.pausedAt = time.Now()
}

func (p *pausableTimer) Resume() {
	if !p.isPaused {
		return
	}
	p.isPaused = false
	p.pausedTotal += time.Since(p.pausedAt)
}

// Stop ends the timer early. If the timer is currently paused, the ongoing pause
// is finalized before stopping so Elapsed() reflects an accurate active duration.
func (p *pausableTimer) Stop() {
	now := time.Now()
	if p.isPaused {
		p.pausedTotal += now.Sub(p.pausedAt)
		p.isPaused = false
	}
	p.stopped = true
	p.stoppedAt = now
}

// Reset returns the timer to its initial state, preserving the configured duration.
func (p *pausableTimer) Reset() {
	limit := p.limit
	*p = pausableTimer{limit: limit}
}

func (p *pausableTimer) IsPaused() bool {
	return p.isPaused
}

func (p *pausableTimer) IsExpired() bool {
	return p.Elapsed() >= p.limit
}

// Elapsed returns how much active (non-paused) time has passed.
// While paused, this value is frozen at the moment the pause began.
func (p *pausableTimer) Elapsed() time.Duration {
	if p.stopped {
		return p.stoppedAt.Sub(p.start) - p.pausedTotal
	}
	if p.isPaused {
		// Freeze elapsed at the moment the pause started.
		return p.pausedAt.Sub(p.start) - p.pausedTotal
	}
	return time.Since(p.start) - p.pausedTotal
}

// TotalDuration returns wall-clock time elapsed, including any paused periods.
func (p *pausableTimer) TotalDuration() time.Duration {
	if p.stopped {
		return p.stoppedAt.Sub(p.start)
	}
	if p.IsExpired() {
		return p.limit + p.pausedTotal
	}
	return time.Since(p.start)
}

func (p *pausableTimer) TimeLimit() time.Duration {
	return p.limit
}

func (p *pausableTimer) StartedAt() time.Time {
	return p.start
}

// EndedAt returns the time the timer ended and true, or the zero time and false
// if the timer is still running.
func (p *pausableTimer) EndedAt() (time.Time, bool) {
	if p.stopped {
		return p.stoppedAt, true
	}
	if p.IsExpired() {
		// The timer expired at the wall-clock moment when Elapsed() first hit the limit:
		// start + limit + total paused time = start + (limit + pausedTotal)
		return p.start.Add(p.limit + p.pausedTotal), true
	}
	return time.Time{}, false
}
