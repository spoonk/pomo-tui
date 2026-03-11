package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectEntry represents a project item in the selector.
type ProjectEntry struct {
	ID          uint
	Name        string
	Description string
}

// ProjectSelectorNavigateUpMsg is emitted when the user presses up at the top
// of the list, signalling the parent to move focus to the previous field.
type ProjectSelectorNavigateUpMsg struct{}

// ProjectCreateRequestMsg is emitted when the user confirms a name that does
// not exist yet, requesting the parent to persist it and call AddProject.
type ProjectCreateRequestMsg struct{ Name string }

// ProjectSelectorEnterMsg is emitted when the user confirms their selection.
// Selected is nil when the selector has no items and no query.
type ProjectSelectorEnterMsg struct{ Selected *ProjectEntry }

// ProjectSelector is a searchable dropdown component for picking a project.
type ProjectSelector struct {
	textInput  textinput.Model
	all        []ProjectEntry
	filtered   []ProjectEntry
	cursor     int
	viewOffset int
	selected   *ProjectEntry
	focused    bool
}

// NewProjectSelector creates a ProjectSelector pre-populated with projects.
func NewProjectSelector(projects []ProjectEntry) ProjectSelector {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "none"
	ti.CharLimit = 64
	ti.Width = 15

	filtered := make([]ProjectEntry, len(projects))
	copy(filtered, projects)

	return ProjectSelector{
		textInput: ti,
		all:       projects,
		filtered:  filtered,
	}
}

// Focus gives focus to the selector and returns a Cmd to enable cursor blink.
func (ps ProjectSelector) Focus() (ProjectSelector, tea.Cmd) {
	ps.focused = true
	cmd := ps.textInput.Focus()
	return ps, cmd
}

// Blur removes focus from the selector.
func (ps ProjectSelector) Blur() ProjectSelector {
	ps.focused = false
	ps.textInput.Blur()
	return ps
}

// Focused reports whether the selector currently has focus.
func (ps ProjectSelector) Focused() bool { return ps.focused }

// Selected returns the currently selected project entry, or nil if none.
func (ps ProjectSelector) Selected() *ProjectEntry { return ps.selected }

// AddProject appends a new entry to the list and selects it. Use this after a
// successful project creation to update the available options.
func (ps ProjectSelector) AddProject(entry ProjectEntry) ProjectSelector {
	ps.all = append(ps.all, entry)
	ps.refilter()
	for i, p := range ps.filtered {
		if p.Name == entry.Name {
			e := ps.filtered[i]
			ps.selected = &e
			ps.cursor = i
			break
		}
	}
	ps.adjustViewport()
	ps.textInput.SetValue(entry.Name)
	return ps
}

// adjustViewport keeps cursor inside [viewOffset, viewOffset+10).
func (ps *ProjectSelector) adjustViewport() {
	if ps.cursor < ps.viewOffset {
		ps.viewOffset = ps.cursor
	} else if ps.cursor >= ps.viewOffset+10 {
		ps.viewOffset = ps.cursor - 9
	}
}

// refilter updates ps.filtered from the current text input value and clamps
// ps.cursor to the valid range (including the optional "create" row).
func (ps *ProjectSelector) refilter() {
	query := strings.ToLower(ps.textInput.Value())
	if query == "" {
		ps.filtered = make([]ProjectEntry, len(ps.all))
		copy(ps.filtered, ps.all)
	} else {
		ps.filtered = nil
		for _, p := range ps.all {
			if strings.Contains(strings.ToLower(p.Name), query) {
				ps.filtered = append(ps.filtered, p)
			}
		}
	}
	total := len(ps.filtered)
	if ps.showCreateOption() {
		total++
	}
	if total == 0 {
		ps.cursor = 0
	} else if ps.cursor >= total {
		ps.cursor = total - 1
	}
	ps.viewOffset = 0
}

// showCreateOption returns true when the input text is non-empty and has no
// exact match in the filtered list, meaning a new project can be created.
func (ps ProjectSelector) showCreateOption() bool {
	q := ps.textInput.Value()
	if q == "" {
		return false
	}
	for _, p := range ps.filtered {
		if strings.EqualFold(p.Name, q) {
			return false
		}
	}
	return true
}

// Update handles Bubble Tea messages for the selector.
func (ps ProjectSelector) Update(msg tea.Msg) (ProjectSelector, tea.Cmd) {
	if !ps.focused {
		return ps, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "down":
			total := len(ps.filtered)
			if ps.showCreateOption() {
				total++
			}
			if ps.cursor < total-1 {
				ps.cursor++
				ps.adjustViewport()
			}
			return ps, nil

		case "up":
			if ps.cursor > 0 {
				ps.cursor--
				ps.adjustViewport()
				return ps, nil
			}
			// At the top — signal the parent to move focus up.
			return ps, func() tea.Msg { return ProjectSelectorNavigateUpMsg{} }

		case "enter":
			hasCreate := ps.showCreateOption()
			if hasCreate && ps.cursor == len(ps.filtered) {
				name := ps.textInput.Value()
				return ps, func() tea.Msg { return ProjectCreateRequestMsg{Name: name} }
			}
			if len(ps.filtered) > 0 {
				e := ps.filtered[ps.cursor]
				ps.selected = &e
				ps.textInput.SetValue(e.Name)
				selected := ps.selected
				return ps, func() tea.Msg { return ProjectSelectorEnterMsg{Selected: selected} }
			}
			// No items and no create option (empty query, no projects).
			return ps, func() tea.Msg { return ProjectSelectorEnterMsg{Selected: nil} }
		}
	}

	// Forward everything else to the embedded text input (cursor blink, etc.).
	var cmd tea.Cmd
	prevVal := ps.textInput.Value()
	ps.textInput, cmd = ps.textInput.Update(msg)
	if ps.textInput.Value() != prevVal {
		ps.refilter()
		ps.cursor = 0
	}
	return ps, cmd
}

// View renders the value portion of the selector (without its label prefix).
// When focused it shows the live text input; when unfocused it shows the
// selected name or a faint "none" placeholder.
func (ps ProjectSelector) View() string {
	if ps.selected != nil && !ps.focused {
		return ps.selected.Name
	}
	return ps.textInput.View()
}

// DropdownView renders the dropdown list indented by indent spaces, or an
// empty string when the selector is not focused or the list is empty.
func (ps ProjectSelector) DropdownView() string {
	if !ps.focused {
		return ""
	}
	hasCreate := ps.showCreateOption()
	if len(ps.filtered) == 0 && !hasCreate {
		return ""
	}

	selectedStyle := lipgloss.NewStyle().Bold(true).Faint(true)
	dimStyle := lipgloss.NewStyle().Faint(true)

	end := min(ps.viewOffset+10, len(ps.filtered))

	var lines []string

	if ps.viewOffset > 0 {
		lines = append(lines, dimStyle.Render(fmt.Sprintf("  \u2191 %d more", ps.viewOffset)))
	}

	for i, entry := range ps.filtered[ps.viewOffset:end] {
		idx := ps.viewOffset + i
		var line string
		if idx == ps.cursor {
			line = selectedStyle.Render(" " + entry.Name)
		} else {
			line = dimStyle.Render("  " + entry.Name)
		}
		lines = append(lines, line)
	}

	if hasCreate && end == len(ps.filtered) {
		label := fmt.Sprintf("+ create \"%s\"", ps.textInput.Value())
		var line string
		if ps.cursor == len(ps.filtered) {
			line = selectedStyle.Render("\u25b6 " + label)
		} else {
			line = dimStyle.Render("  " + label)
		}
		lines = append(lines, line)
	}

	below := len(ps.filtered) - end
	if hasCreate && end < len(ps.filtered) {
		below++ // create row also hidden
	}
	if below > 0 {
		lines = append(lines, dimStyle.Render(fmt.Sprintf("  \u2193 %d more", below)))
	}

	return strings.Join(lines, "\n")
}
