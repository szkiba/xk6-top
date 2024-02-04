// Package navbar contains navbar UI component.
package navbar

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

// NavChangedMsg message sent when active item changed.
type NavChangedMsg int

type disableMsg int

type enableMsg int

// Disable disables given navbar item.
func Disable(idx int) tea.Cmd {
	return func() tea.Msg {
		return disableMsg(idx)
	}
}

// Enable enables given navbar item.
func Enable(idx int) tea.Cmd {
	return func() tea.Msg {
		return enableMsg(idx)
	}
}

func navChangedCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return NavChangedMsg(idx)
	}
}

// Item describes a navbar item.
type Item struct {
	Label    string
	Disabled bool
}

// Model contains navbar UI component.
type Model struct {
	theme  *theme.Theme
	items  []*Item
	Active int
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	idx := m.Active

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "right", "n":
			for next := m.Active + 1; next < len(m.items); next++ {
				if !m.items[next].Disabled {
					m.Active = next
					break
				}
			}
		case "shift+tab", "left", "p":
			for prev := m.Active - 1; prev >= 0; prev-- {
				if !m.items[prev].Disabled {
					m.Active = prev
					break
				}
			}
		}
	case disableMsg:
		if len(m.items) > int(msg) {
			m.items[msg].Disabled = true
		}
	case enableMsg:
		if len(m.items) > int(msg) {
			m.items[msg].Disabled = false
		}
	default:
	}

	var cmd tea.Cmd

	if idx != m.Active {
		cmd = navChangedCmd(m.Active)
	}

	return m, cmd
}

// View implements tea.Model.
func (m Model) View() string {
	var doc strings.Builder

	doc.WriteString(m.theme.Primary.Render(" k6 "))
	for idx, item := range m.items {
		style := inactiveStyle
		if item.Disabled {
			style = m.theme.Disabled
		} else if idx == m.Active {
			style = activeStyle
		}

		style = style.Copy().Padding(0, 1)

		doc.WriteString(style.Render(item.Label))
	}

	return doc.String()
}

// New creates new navbar instance.
func New(theme *theme.Theme, items []*Item) Model {
	m := Model{
		theme: theme,
		items: items,
	}

	return m
}

//nolint:gochecknoglobals
var (
	inactiveStyle = lipgloss.NewStyle().Padding(0, 1)
	activeStyle   = inactiveStyle.Copy().Reverse(true)
)
