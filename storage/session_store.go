package storage

func (s *sqliteStore) SaveSession(session *Session) error {
	return s.db.Create(session).Error
}

func (s *sqliteStore) GetSessions() []Session {
	var sessions []Session
	result := s.db.Preload("Project").Find(&sessions)
	if result.Error != nil {
		panic(result.Error)
	}

	return sessions
}
