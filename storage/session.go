package storage

import "time"

// Session represents a completed Pomodoro session persisted in the database.
type Session struct {
	ID uint `gorm:"primarykey;autoIncrement"`

	StartedAt time.Time // UTC wall-clock time the session began
	EndedAt   time.Time // UTC wall-clock time the session ended

	// DurationSeconds is the actual active time (wall-clock elapsed minus any
	// time spent paused). Stored as integer seconds
	DurationSeconds int64

	// PlannedSeconds is the session length the user configured
	PlannedSeconds int64

	Project Project
}
