package tui

import "github.com/charmbracelet/lipgloss"

func centerContent(width, height int, content string) string {
	if width <= 0 || height <= 0 {
		return content
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
