// Package chart contains chart UI component.
package chart

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/statbar"
)

// Serie contains serie reference.
type Serie struct {
	Metric    string
	Aggregate string
}

// Model contains chart UI component.
type Model struct {
	digest *digest.Digester
	series []*Serie
	data   [][]float64
	width  int
	height int

	details int
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "+", "shift+down":
			if m.details < len(m.series)-1 {
				m.details++
			}
			m.update()
		case "-", "shift+up":
			if m.details > 0 {
				m.details--
			}
			m.update()
		default:
		}

	case statbar.PanelChangedMsg:
		m.update()
	case *digest.Digest:
		m.update()

	case tea.WindowSizeMsg:
		m.width = msg.Width - 8
		m.height = msg.Height - 9

		m.update()
	default:
	}

	return m, nil
}

func (m *Model) update() {
	dig := m.digest.Digest()

	scale := 1.0

	if agg, ok := dig.Snapshot[m.series[0].Metric]; ok {
		if val, hasAvg := agg[m.series[0].Aggregate]; hasAvg {
			unit, _ := digest.Unit(val)
			scale = float64(unit)
		}
	}

	data := make([][]float64, 0, len(m.series))

	for _, serie := range m.series {
		values := m.digest.Serie(serie.Metric, serie.Aggregate)

		for idx, value := range values {
			values[idx] = value / scale
		}

		data = append(data, values)
	}

	m.data = data
}

// View implements tea.Model.
func (m Model) View() string {
	if len(m.data) < 1 || len(m.data[0]) < 1 {
		return ""
	}

	graph := asciigraph.PlotMany(
		m.data[:m.details+1],
		asciigraph.Height(m.height),
		asciigraph.Width(m.width),
		asciigraph.SeriesColors(seriesColors...),
		asciigraph.AxisColor(asciigraph.Gray),
		asciigraph.LabelColor(asciigraph.Gray),
	)

	var legend strings.Builder

	for idx, serie := range m.series[:m.details+1] {
		style := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(seriesColors[idx]))
		legend.WriteString(style.Render("━━"))
		legend.WriteString(" ")
		legend.WriteString(serie.Aggregate)
		legend.WriteString("  ")
	}

	style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(m.width)

	graph += "\n" + style.Render(legend.String())

	return lipgloss.NewStyle().AlignVertical(lipgloss.Bottom).Height(m.height).Render(graph)
}

// New creates new chart instance.
func New(digest *digest.Digester, series []*Serie) Model {
	m := Model{
		digest:  digest,
		series:  series,
		details: len(series) - 1,
	}

	return m
}

//nolint:gochecknoglobals
var (
	seriesColors = []asciigraph.AnsiColor{
		asciigraph.Green,
		asciigraph.Yellow,
		asciigraph.IndianRed,
		asciigraph.Violet,
	}
)
