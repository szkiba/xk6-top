// Package status contains status bar UI component.
package status

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

// Model contains status bar UI component.
type Model struct {
	digest *digest.Digest
	theme  *theme.Theme
	ready  bool
	width  int

	message string
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case *digest.Digest:
		m.digest = msg
		m.message = ""
		m.update()

		if msg.EventType == digest.EventTypeCumulative {
			cmds = append(cmds, tickCmd())
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.update()

	case error:
		parts := strings.SplitN(msg.Error(), "\n", 2)
		m.message = m.theme.Error.Render(parts[0])

	case tickMsg:
		if m.digest != nil && !m.digest.Playback && m.digest.State == digest.StateRunning {
			cmds = append(cmds, tickCmd())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) update() {
	m.ready = m.width != 0
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return ""
	}

	var buff strings.Builder

	buff.WriteString(m.progress())
	buff.WriteRune('\n')

	buff.WriteString(m.state())
	buff.WriteString(m.message)

	return buff.String()
}

func (m *Model) progress() string {
	empty := m.width
	full := 0

	if m.digest != nil && !m.digest.Playback && m.digest.State == digest.StateRunning {
		percent := float64(time.Since(m.digest.Time())) / float64(m.digest.Period())

		full = int(float64(m.width) * percent)
		if full > m.width {
			full = m.width
		}
		if full < 0 {
			full = 0
		}

		empty = m.width - full
	}

	return m.theme.Primary.Render(strings.Repeat("─", full)) +
		m.theme.Divider.Render(strings.Repeat("─", empty))
}

// New creates new status instance.
func New(theme *theme.Theme) Model {
	m := Model{
		theme: theme,
	}

	return m
}

func (m *Model) state() string {
	if m.digest == nil {
		return m.theme.State.Render(digest.StateWaiting.String())
	}

	return m.theme.State.Render(m.digest.State.String())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
