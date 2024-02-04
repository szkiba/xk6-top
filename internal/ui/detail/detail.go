// Package detail contains detail tab UI component.
package detail

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/chart"
	"github.com/szkiba/xk6-top/internal/ui/statbar"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

// Model contains detail tab UI component.
type Model struct {
	theme    *theme.Theme
	digest   *digest.Digest
	statbar  statbar.Model
	charts   []tea.Model
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
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case statbar.PanelChangedMsg:
		wsmsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
		cmds = append(cmds, func() tea.Msg { return wsmsg })
		m.update()
	case *digest.Digest:
		m.digest = msg
		m.update()

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
	}

	m.statbar, cmd = m.statbar.Update(msg)
	cmds = append(cmds, cmd)

	m.charts[m.statbar.Active], cmd = m.charts[m.statbar.Active].Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.update()

	return m, tea.Batch(cmds...)
}

func (m *Model) update() {
	graph := m.charts[m.statbar.Active].View()

	div := m.theme.Divider.
		Render(strings.Repeat("─", m.width))

	m.viewport.SetContent(m.statbar.View() + "\n" + div + "\n" + graph)
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return m.viewport.View()
}

// New creates new detail instance.
func New(theme *theme.Theme, digest *digest.Digester, panels []*statbar.Panel) Model {
	m := Model{theme: theme}

	m.statbar = statbar.New(theme, panels)

	for _, panel := range panels {
		digest.Collect(panel.Metric, panel.Aggregate)

		series := []*chart.Serie{{Metric: panel.Metric, Aggregate: panel.Aggregate}}

		if panel.Aggregate == "avg" {
			digest.Collect(panel.Metric, "p(90)")
			digest.Collect(panel.Metric, "p(95)")
			digest.Collect(panel.Metric, "p(99)")
			series = append(series, &chart.Serie{Metric: panel.Metric, Aggregate: "p(90)"})
			series = append(series, &chart.Serie{Metric: panel.Metric, Aggregate: "p(95)"})
			series = append(series, &chart.Serie{Metric: panel.Metric, Aggregate: "p(99)"})
		}

		chart := chart.New(digest, series)

		m.charts = append(m.charts, chart)
	}

	return m
}
