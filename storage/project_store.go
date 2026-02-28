package storage

func (s *sqliteStore) SaveProject(project *Project) error {
	return s.db.Create(project).Error
}

func (s *sqliteStore) GetProjects() []Project {
	var projects []Project
	result := s.db.Find(&projects)
	if result.Error != nil {
		panic(result.Error)
	}

	return projects
}
