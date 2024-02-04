// Package theme contains lipgloss styles for dark and light theme.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme holds lipgloss styles for dark or light theme.
type Theme struct {
	PrimaryColor string
	DividerColor string
	Primary      lipgloss.Style
	Heading      lipgloss.Style
	Disabled     lipgloss.Style
	Divider      lipgloss.Style
	State        lipgloss.Style
	Error        lipgloss.Style
}

// Detect returns terminal theme based on background color.
func Detect() *Theme {
	if lipgloss.HasDarkBackground() {
		return dark()
	}

	return light()
}

func dark() *Theme {
	return &Theme{
		PrimaryColor: "#6e59de",
		Primary:      lipgloss.NewStyle().Foreground(lipgloss.Color("#6e59de")),
		Heading:      lipgloss.NewStyle().Bold(true),
		Disabled:     lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#808080")),
		DividerColor: "#394160",
		Divider:      lipgloss.NewStyle().Foreground(lipgloss.Color("#394160")),
		State: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#808080")).
			Width(11).
			Padding(0, 1, 0, 1),
		Error: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#a00000")).
			Padding(0, 1, 0, 1),
	}
}

func light() *Theme {
	return &Theme{
		PrimaryColor: "#6e59de",
		Primary:      lipgloss.NewStyle().Foreground(lipgloss.Color("#6e59de")),
		Heading:      lipgloss.NewStyle().Bold(true),
		Disabled:     lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#a0a0a0")),
		DividerColor: "#e0e0e0",
		Divider:      lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0")),
		State: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#a0a0a0")).
			Width(11).
			Padding(0, 1, 0, 1),
		Error: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#c00000")).
			Padding(0, 1, 0, 1),
	}
}
