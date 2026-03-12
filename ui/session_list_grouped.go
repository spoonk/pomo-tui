package ui

import (
	"time"

	"charm.land/lipgloss/v2"
)

type GroupedSessionListEntry struct {
	Completed   bool
	ProjectName *string
	StartedAt   time.Time
	EndedAt     time.Time
	Duration    time.Duration // active work time only
}

func GroupedSessionList(entries []GroupedSessionListEntry) string {
	if len(entries) == 0 {
		return ""
	}

	const maxProjectWidth = 64

	// coalesce adjacent sessions by project id
	groupedSessions := make([][]GroupedSessionListEntry, 0)
	currGroup := []GroupedSessionListEntry{}
	currProject := entries[0].ProjectName

	for i, entry := range entries {
		if i == 0 {
			continue
		}

		if entry.ProjectName == nil && currProject == nil {
			currGroup = append(currGroup, entry)
		} else if entry.ProjectName == nil {
			groupedSessions = append(groupedSessions, currGroup)
			currGroup = []GroupedSessionListEntry{entry}
			currProject = entry.ProjectName
		} else if *entry.ProjectName != *currProject {
			groupedSessions = append(groupedSessions, currGroup)
			currGroup = []GroupedSessionListEntry{entry}
			currProject = entry.ProjectName
		} else {
			currGroup = append(currGroup, entry)
		}
	}

	// add final group
	groupedSessions = append(groupedSessions, currGroup)

	// now, display some UI for each session
	groupedSessionsUi := []string{}
	for _, groupedSession := range groupedSessions {
		groupedSessionsUi = append(groupedSessionsUi, GroupedSession(groupedSession))
	}

	return lipgloss.JoinVertical(lipgloss.Left, groupedSessionsUi...)
}

func GroupedSession(entriesForProject []GroupedSessionListEntry) string {
	if len(entriesForProject) == 0 {
		return ""
	}
	// all entries have the same project
	if entriesForProject[0].ProjectName == nil {
		return "huh"

	}
	project := entriesForProject[0].ProjectName

	var projectName string
	if project != nil {
		projectName = *project
	} else {
		projectName = "None"
	}

	projectName = lipgloss.NewStyle().Faint(false).Render(projectName)

	entryRows := []string{}

	for _, entry := range entriesForProject {
		entryDuration := ElapsedTime(entry.StartedAt, entry.EndedAt, entry.Duration)
		// entryListItem := fmt.Sprintf("%s", entryDuration)
		entryRows = append(entryRows, entryDuration)
	}

	allEntries := lipgloss.JoinVertical(lipgloss.Left, entryRows...)
	return lipgloss.JoinVertical(lipgloss.Left, projectName, allEntries)
}
