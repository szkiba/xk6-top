// Package statbar contains stat panel bar UI component.
package statbar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

// PanelChangedMsg sent when active panel changed.
type PanelChangedMsg int

func panelChangedCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return PanelChangedMsg(idx)
	}
}

// Panel contains stat panel definition.
type Panel struct {
	Label     string
	Caption   string
	Metric    string
	Aggregate string
}

// Model contains stat panel bar UI component.
type Model struct {
	theme  *theme.Theme
	digest *digest.Digest
	panels []*Panel
	ready  bool
	width  int
	Active int

	captionsWidth int
	minPanelWidth int
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
		case "shift+right", "ctrl+right":
			m.Active = min(m.Active+1, len(m.panels)-1)
		case "shift+left", "ctrl+left":
			m.Active = max(m.Active-1, 0)
		}

	case *digest.Digest:
		m.digest = msg
		m.update()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.update()
	}

	var cmd tea.Cmd

	if idx != m.Active {
		cmd = panelChangedCmd(m.Active)
	}

	return m, cmd
}

func (m *Model) update() {
	m.ready = m.width != 0 && m.digest != nil && !m.digest.Start.IsZero()
}

func (m Model) getLabels() []string {
	labels := make([]string, 0, len(m.panels))

	useCaptions := m.captionsWidth <= m.width

	for _, panel := range m.panels {
		if useCaptions {
			labels = append(labels, panel.Caption)
		} else {
			labels = append(labels, panel.Label)
		}
	}

	return labels
}

func (m Model) getPanelWidth() int {
	width := m.width / len(m.panels)
	if width < m.minPanelWidth {
		width = m.minPanelWidth
	}

	return width
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.ready {
		return ""
	}

	pan := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.getPanelWidth())

	snapshots := m.digest.Snapshot

	all := make([]string, 0, len(m.panels))

	labels := m.getLabels()

	for idx, panel := range m.panels {
		sty := pan.Copy().Reverse(idx == m.Active)

		met, found := m.digest.FindMetric(panel.Metric)
		if !found {
			continue
		}

		val := met.Contains.Format(snapshots[panel.Metric][panel.Aggregate])
		if panel.Aggregate == "rate" || panel.Aggregate == "peak" {
			val += "/s"
		}

		str := fmt.Sprintf("%s\n%s", labels[idx], val)
		all = append(all, sty.Render(str))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, all...)
}

// New creates new statbar instance.
func New(theme *theme.Theme, panels []*Panel) Model {
	m := Model{
		theme:  theme,
		panels: panels,
	}

	maxCaptionLen := 0
	maxLabelLen := 0

	for _, panel := range m.panels {
		if len(panel.Caption) > maxCaptionLen {
			maxCaptionLen = len(panel.Caption)
		}
		if len(panel.Label) > maxLabelLen {
			maxLabelLen = len(panel.Label)
		}
	}

	m.captionsWidth = len(panels) * maxCaptionLen
	m.minPanelWidth = maxLabelLen

	return m
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
