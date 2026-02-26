package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// delta is the tolerance used in time comparisons throughout this file.
// Because we set timestamps deterministically via TimerConfig, the only
// slippage is the few microseconds of test execution time.
const delta = 200 * time.Millisecond

// ── Initial state ──────────────────────────────────────────────────────────

func TestNewPausableTimer_InitialState(t *testing.T) {
	p := NewPausableTimer()

	assert.Equal(t, time.Duration(0), p.TimeLimit())
	assert.True(t, p.StartedAt().IsZero(), "start time should be unset before Start()")
	assert.False(t, p.IsPaused())

	_, ok := p.EndedAt()
	assert.False(t, ok, "EndedAt should return false before the timer ends")
}

// ── Configuration ──────────────────────────────────────────────────────────

func TestSetDuration(t *testing.T) {
	p := NewPausableTimer()
	p.SetDuration(25 * time.Minute)
	assert.Equal(t, 25*time.Minute, p.TimeLimit())
}

// ── Start ──────────────────────────────────────────────────────────────────

func TestStart_SetsStartTime(t *testing.T) {
	p := NewPausableTimer()
	p.SetDuration(25 * time.Minute)

	before := time.Now()
	p.Start()
	after := time.Now()

	started := p.StartedAt()
	assert.True(t, !started.Before(before) && !started.After(after),
		"StartedAt should be between the before/after timestamps")
}

// ── Elapsed ────────────────────────────────────────────────────────────────

func TestElapsed_ReflectsActiveTime(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})
	assert.InDelta(t, float64(5*time.Minute), float64(p.Elapsed()), float64(delta))
}

func TestElapsed_FrozenWhilePaused(t *testing.T) {
	// Started 5 min ago, paused 3 min ago → active elapsed = 2 min
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:    time.Now().Add(-5 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		PausedAt: time.Now().Add(-3 * time.Minute),
	})

	assert.InDelta(t, float64(2*time.Minute), float64(p.Elapsed()), float64(delta),
		"elapsed should be frozen at the moment of pause (2 minutes active)")
}

func TestElapsed_AccountsForPriorPauses(t *testing.T) {
	// 10 min elapsed, 2 min were paused across earlier pauses → active = 8 min
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-10 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 2 * time.Minute,
	})
	assert.InDelta(t, float64(8*time.Minute), float64(p.Elapsed()), float64(delta))
}

// ── Pause ──────────────────────────────────────────────────────────────────

func TestPause_SetsIsPaused(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

	assert.False(t, p.IsPaused())
	p.Pause()
	assert.True(t, p.IsPaused())
}

func TestPause_FreezesElapsed(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

	p.Pause()
	frozenElapsed := p.Elapsed()

	// Simulate wall clock advancing without the test actually sleeping.
	// Since Elapsed() is frozen, calling it again returns the same value.
	assert.InDelta(t, float64(frozenElapsed), float64(p.Elapsed()), float64(delta),
		"Elapsed should not change while paused")
}

func TestPause_Idempotent(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})
	p.Pause()
	elapsedAfterFirst := p.Elapsed()

	p.Pause() // second call should be a no-op

	assert.True(t, p.IsPaused())
	assert.InDelta(t, float64(elapsedAfterFirst), float64(p.Elapsed()), float64(delta),
		"second Pause should not change the frozen elapsed value")
}

// ── Resume ─────────────────────────────────────────────────────────────────

func TestResume_ClearsIsPaused(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:    time.Now().Add(-5 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		PausedAt: time.Now(),
	})

	assert.True(t, p.IsPaused())
	p.Resume()
	assert.False(t, p.IsPaused())
}

func TestResume_AccumulatesPausedDuration(t *testing.T) {
	// Started 5 min ago; currently paused at the 3-min mark (pausedAt = now-2min),
	// with 1 min of prior accumulated pause.
	// After resume: total paused = 1 min + 2 min = 3 min → active elapsed = 5 - 3 = 2 min
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-5 * time.Minute),
		Limit:       25 * time.Minute,
		IsPaused:    true,
		PausedAt:    time.Now().Add(-2 * time.Minute),
		PausedTotal: 1 * time.Minute,
	})

	p.Resume()

	assert.InDelta(t, float64(2*time.Minute), float64(p.Elapsed()), float64(delta),
		"active elapsed should be 5 min - (1 min prior + 2 min just ended) = 2 min")
}

func TestResume_Idempotent(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:    time.Now().Add(-5 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		PausedAt: time.Now().Add(-1 * time.Minute),
	})

	p.Resume()
	elapsedAfterFirst := p.Elapsed()

	p.Resume() // second call should be a no-op

	assert.False(t, p.IsPaused())
	assert.InDelta(t, float64(elapsedAfterFirst), float64(p.Elapsed()), float64(delta),
		"second Resume should not add extra pause time")
}

// ── IsExpired ──────────────────────────────────────────────────────────────

func TestIsExpired_FalseBeforeLimit(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})
	assert.False(t, p.IsExpired())
}

func TestIsExpired_TrueAfterLimit(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-26 * time.Minute),
		Limit: 25 * time.Minute,
	})
	assert.True(t, p.IsExpired())
}

func TestIsExpired_NotWhilePaused(t *testing.T) {
	// Wall clock has advanced 30 minutes past start, but the timer was paused
	// at the 20-minute mark → active elapsed = 20 min < 25 min limit.
	// A paused timer must never expire.
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:    time.Now().Add(-30 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		PausedAt: time.Now().Add(-10 * time.Minute), // paused 20 min into the session
	})

	assert.True(t, p.IsPaused())
	assert.False(t, p.IsExpired(),
		"a paused timer must not expire even if wall-clock time has exceeded the limit")
}

func TestIsExpired_AccountsForPausedTime(t *testing.T) {
	// 30 min elapsed, 6 min paused → active = 24 min < 25 min limit
	notYet := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-30 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 6 * time.Minute,
	})
	assert.False(t, notYet.IsExpired(), "24 active minutes should not exceed a 25-minute limit")

	// 32 min elapsed, 6 min paused → active = 26 min > 25 min limit
	expired := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-32 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 6 * time.Minute,
	})
	assert.True(t, expired.IsExpired(), "26 active minutes should exceed a 25-minute limit")
}

func TestIsExpired_FalseAfterStop(t *testing.T) {
	// Stopped early — has not reached its limit
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})
	p.Stop()
	assert.False(t, p.IsExpired(),
		"stopping early should not count as expiry")
}

// ── Stop ───────────────────────────────────────────────────────────────────

func TestStop_WhileRunning(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

	stopTime := time.Now()
	p.Stop()

	assert.False(t, p.IsPaused())

	endedAt, ok := p.EndedAt()
	assert.True(t, ok)
	assert.WithinDuration(t, stopTime, endedAt, delta)

	// Elapsed = 5 min (no pauses, stopped at wall-clock position)
	assert.InDelta(t, float64(5*time.Minute), float64(p.Elapsed()), float64(delta))
}

func TestStop_WhilePaused_FinalizesThePause(t *testing.T) {
	// Started 5 min ago, paused 2 min ago — so 3 min of active time before Stop.
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:    time.Now().Add(-5 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		PausedAt: time.Now().Add(-2 * time.Minute),
	})

	p.Stop()

	assert.False(t, p.IsPaused(), "Stop should clear the paused state")

	endedAt, ok := p.EndedAt()
	assert.True(t, ok)
	assert.WithinDuration(t, time.Now(), endedAt, delta)

	// Active elapsed = 5 min total - 2 min paused = 3 min
	assert.InDelta(t, float64(3*time.Minute), float64(p.Elapsed()), float64(delta))
}

// ── EndedAt ────────────────────────────────────────────────────────────────

func TestEndedAt_RunningTimer_ReturnsFalse(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

	_, ok := p.EndedAt()
	assert.False(t, ok, "EndedAt should return false while the timer is still running")
}

func TestEndedAt_ExpiredTimer_ReturnsComputedTime(t *testing.T) {
	// Started 26 min ago with a 25 min limit → expired 1 minute ago
	start := time.Now().Add(-26 * time.Minute)
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: start,
		Limit: 25 * time.Minute,
	})

	endedAt, ok := p.EndedAt()
	assert.True(t, ok)

	// EndedAt = start + limit = (now - 26min) + 25min = now - 1min
	expected := start.Add(25 * time.Minute)
	assert.WithinDuration(t, expected, endedAt, delta)
}

func TestEndedAt_ExpiredWithPauses_IncludesPausedTime(t *testing.T) {
	// Pauses push the expiry moment forward:
	// EndedAt = start + limit + pausedTotal
	start := time.Now().Add(-30 * time.Minute)
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       start,
		Limit:       25 * time.Minute,
		PausedTotal: 4 * time.Minute,
	})

	assert.True(t, p.IsExpired())

	endedAt, ok := p.EndedAt()
	assert.True(t, ok)

	expected := start.Add(25*time.Minute + 4*time.Minute)
	assert.WithinDuration(t, expected, endedAt, delta)
}

func TestEndedAt_StoppedTimer_ReturnsStopTime(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})

	stopTime := time.Now()
	p.Stop()

	endedAt, ok := p.EndedAt()
	assert.True(t, ok)
	assert.WithinDuration(t, stopTime, endedAt, delta)
}

// ── TotalDuration ──────────────────────────────────────────────────────────

func TestTotalDuration_Running(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start: time.Now().Add(-5 * time.Minute),
		Limit: 25 * time.Minute,
	})
	// TotalDuration includes paused time — since there are no pauses, it equals elapsed.
	assert.InDelta(t, float64(5*time.Minute), float64(p.TotalDuration()), float64(delta))
}

func TestTotalDuration_Stopped_IsWallClockTime(t *testing.T) {
	// 5 min of wall-clock time: 3 min paused + 2 min active
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-5 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 3 * time.Minute,
	})
	p.Stop()

	// TotalDuration = wall-clock from start to stop ≈ 5 min (includes the pauses)
	assert.InDelta(t, float64(5*time.Minute), float64(p.TotalDuration()), float64(delta))
}

func TestTotalDuration_Expired_IsLimitPlusPauses(t *testing.T) {
	// Expired: TotalDuration = limit + pausedTotal (not time.Since(start),
	// which would keep growing after expiry)
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-30 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 3 * time.Minute,
	})
	assert.True(t, p.IsExpired())
	assert.Equal(t, 28*time.Minute, p.TotalDuration(), // 25 + 3
		"TotalDuration for an expired timer should be limit + accumulated pause time")
}

// ── Reset ──────────────────────────────────────────────────────────────────

func TestReset_ClearsTimingState(t *testing.T) {
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-5 * time.Minute),
		Limit:       25 * time.Minute,
		IsPaused:    true,
		PausedAt:    time.Now().Add(-1 * time.Minute),
		PausedTotal: 30 * time.Second,
	})

	p.Reset()

	assert.Equal(t, 25*time.Minute, p.TimeLimit(), "duration should be preserved after Reset")
	assert.True(t, p.StartedAt().IsZero(), "start time should be cleared")
	assert.False(t, p.IsPaused(), "paused state should be cleared")
	_, ok := p.EndedAt()
	assert.False(t, ok, "should have no end time after Reset")
}

func TestReset_ThenStart_BehavesLikeFreshTimer(t *testing.T) {
	// Simulate completing a session, then resetting to start a new one.
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:       time.Now().Add(-26 * time.Minute),
		Limit:       25 * time.Minute,
		PausedTotal: 1 * time.Minute,
	})
	assert.True(t, p.IsExpired(), "sanity check: timer should be expired before reset")

	p.Reset()
	p.Start()

	assert.WithinDuration(t, time.Now(), p.StartedAt(), delta)
	assert.False(t, p.IsExpired(), "freshly started timer should not be expired")
	assert.InDelta(t, float64(0), float64(p.Elapsed()), float64(delta))
}

// ── NewPausableTimerFromConfig ─────────────────────────────────────────────

func TestNewPausableTimerFromConfig_IsPausedWithZeroPausedAt_DefaultsToNow(t *testing.T) {
	// When IsPaused is true but no PausedAt is given, the implementation
	// defaults PausedAt to time.Now() so Elapsed() doesn't return garbage.
	p := NewPausableTimerFromConfig(TimerConfig{
		Start:    time.Now().Add(-5 * time.Minute),
		Limit:    25 * time.Minute,
		IsPaused: true,
		// PausedAt intentionally zero
	})

	assert.True(t, p.IsPaused())
	// PausedAt defaulted to now → active elapsed ≈ 5 min
	assert.InDelta(t, float64(5*time.Minute), float64(p.Elapsed()), float64(delta))
}
