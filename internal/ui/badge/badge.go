// Package badge contains badger rendering functions.
package badge

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/digest"
)

// Render renders a badge.
func Render(label string, value string, level digest.Level) string {
	theme := darkTheme

	if !lipgloss.HasDarkBackground() {
		theme = lightTheme
	}

	return theme.label.Render(label) + theme.value[level].Render(value)
}

// Value renders a colored value. Color determined from level.
func Value(value string, level digest.Level) string {
	theme := darkTheme

	if !lipgloss.HasDarkBackground() {
		theme = lightTheme
	}

	style := lipgloss.NewStyle().Foreground(theme.value[level].GetBackground())

	return style.Render(value)
}

type theme struct {
	label lipgloss.Style
	value map[digest.Level]lipgloss.Style
}

func badgeStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color(color)).Padding(0, 1, 0, 1)
}

//nolint:gochecknoglobals
var (
	darkTheme = theme{
		label: badgeStyle("#404040"),
		value: map[digest.Level]lipgloss.Style{
			digest.None:    badgeStyle("#606060"),
			digest.Info:    badgeStyle("#005fd7"),
			digest.Ready:   badgeStyle("#008700"),
			digest.Notice:  badgeStyle("#5f5f00"),
			digest.Warning: badgeStyle("#ff8700").Foreground(lipgloss.Color("#404040")),
			digest.Error:   badgeStyle("#d70000"),
		},
	}
	lightTheme = theme{
		label: badgeStyle("#dadada"),
		value: map[digest.Level]lipgloss.Style{
			digest.None:    badgeStyle("#eeeeee"),
			digest.Info:    badgeStyle("#00afff"),
			digest.Ready:   badgeStyle("#5fd75f"),
			digest.Notice:  badgeStyle("#d7d75f"),
			digest.Warning: badgeStyle("#FFFF00"),
			digest.Error:   badgeStyle("#ff5f5f"),
		},
	}
)
