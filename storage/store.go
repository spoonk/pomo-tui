package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store interface {
	Close() error
	SaveSession(session *Session) error
	GetSessions() []Session
	SaveProject(project *Project) error
	GetProjects() []Project
}

type sqliteStore struct {
	db *gorm.DB
}

func Open() (Store, error) {
	path, err := dbPath()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		// Silence GORM's query log so it doesn't bleed into the TUI output.
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.AutoMigrate(&Session{}, &Project{}); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &sqliteStore{db: db}, nil
}

func (s *sqliteStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func dbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".local", "share", "pomo-tui", "sessions.db"), nil
}
