package ui

import lipgloss "charm.land/lipgloss/v2"

// Renders overlay over base at x, y
func CompositeOver(base, overlay string, x, y int) string {
	baseLayer := lipgloss.NewLayer(base) // positioned at 0, 0 initially
	overlayLayer := lipgloss.NewLayer(overlay).X(x).Y(y)

	compositor := lipgloss.NewCompositor(baseLayer, overlayLayer)

	return compositor.Render()
}
