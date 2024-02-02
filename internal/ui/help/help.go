// Package help contains help tab UI component.
package help

import (
	_ "embed"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

//go:generate go run ../../../tools/genhelp ../../../docs/help.md .

//go:embed help-dark.ansi
var darkText string

//go:embed help-light.ansi
var lightText string

// Model contains help tab UI component.
type Model struct {
	theme    *theme.Theme
	viewport viewport.Model
	ready    bool
	width    int
	height   int
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-4)
			m.viewport.YPosition = 2
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 4
		}

		m.update()

	default:
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m *Model) update() {
	var buff strings.Builder

	if lipgloss.HasDarkBackground() {
		buff.WriteString(darkText)
	} else {
		buff.WriteString(lightText)
	}

	buff.WriteString(
		m.theme.Disabled.Render(
			"\n\nFor more information visit: https://github.com/szkiba/xk6-top",
		),
	)

	m.viewport.SetContent(
		lipgloss.NewStyle().Padding(0, 1, 0, 1).Width(m.width - 2).Render(buff.String()),
	)
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return ""
	}

	return m.viewport.View()
}

// New creates new help instance.
func New(theme *theme.Theme) Model {
	m := Model{theme: theme}

	return m
}
