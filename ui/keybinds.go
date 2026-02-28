package ui

// Keybinds renders a keybind hint string in the shared muted style.
func Keybinds(hints string) string {
	return KeybindStyle.Render(hints)
}
