// Package divider contains divider bar UI component.
package divider

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

// Model contains divider bar UI component.
type Model struct {
	theme  *theme.Theme
	digest *digest.Digest
	ready  bool
	width  int
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *digest.Digest:
		m.digest = msg
		m.update()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.update()
	}

	return m, nil
}

func (m *Model) update() {
	m.ready = m.width != 0
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return ""
	}

	return m.progress()
}

func (m *Model) progress() string {
	empty := m.width
	full := 0

	if m.digest != nil {
		percent := m.digest.ProgressPercent()

		full = int(float64(m.width) * percent)
		if full > m.width {
			full = m.width
		}
		if full < 0 {
			full = 0
		}

		empty = m.width - full
	}

	return m.theme.Primary.Render(strings.Repeat("━", full)) +
		m.theme.Divider.Render(strings.Repeat("━", empty))
}

// New creates new divider instance.
func New(theme *theme.Theme) Model {
	return Model{theme: theme}
}
