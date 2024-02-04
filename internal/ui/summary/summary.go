// Package summary contains summary tab UI component.
package summary

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/navbar"
)

// Model contains summary tab UI component.
type Model struct {
	mtype  digest.MetricType
	digest *digest.Digest
	table  table.Model
	width  int
	height int

	showTags bool
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
			m.showTags = true
			m.update()
		case "-", "shift+up":
			m.showTags = false
			m.update()
		default:
		}

	case navbar.NavChangedMsg:
		m.update()
	case *digest.Digest:
		m.digest = msg
		m.update()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(m.height - 5)
		m.table.SetWidth(m.width)

		m.update()
	}

	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *Model) update() {
	if m.digest == nil {
		return
	}

	var names []string

	cumulative := m.digest.Cumulative

	for name := range cumulative {
		metric, hasMetric := m.digest.FindMetric(name)
		if hasMetric && metric.Type == m.mtype && name != "time" {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	rows := make([]table.Row, 0, len(names))

	for _, name := range names {
		met, found := m.digest.FindMetric(name)
		if !found {
			continue
		}

		dig := cumulative[name]
		var row []string

		if strings.ContainsRune(name, '{') {
			if !m.showTags {
				continue
			}
			start := strings.IndexRune(name, '{')
			end := strings.LastIndexByte(name, '}')
			name = " { " + name[start+1:end] + " }"
		}

		row = append(row, name)
		for _, agg := range m.mtype.Aggregates() {
			str := met.Contains.Format(dig[agg])
			if agg == "rate" || agg == "peak" {
				str += "/s"
			}

			row = append(row, str)
		}

		rows = append(rows, row)
	}

	m.table.SetColumns(columnsFor(m.mtype, m.width))
	m.table.SetRows(rows)
}

// View implements tea.Model.
func (m Model) View() string {
	return m.table.View()
}

// New creates new summary instance.
func New(mtype digest.MetricType) Model {
	styles := table.Styles{
		Selected: lipgloss.NewStyle().Reverse(true),
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#6e59de")).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}

	m := Model{
		mtype: mtype,
		table: table.New(
			table.WithKeyMap(table.DefaultKeyMap()),
			table.WithFocused(true),
			table.WithStyles(styles),
		),
		showTags: true,
	}

	return m
}

func columnsFor(mtype digest.MetricType, width int) []table.Column {
	if width == 0 {
		width = defaultWidth
	}

	titles := mtype.Aggregates()

	cols := make([]table.Column, 0, len(titles)+1)

	cols = append(cols, table.Column{Title: "metric", Width: min(max(4*width/12, 26), 64)})

	minWidth := 7
	if mtype != digest.MetricTypeTrend {
		minWidth = 10
	}

	for _, title := range mtype.Aggregates() {
		cols = append(cols, table.Column{Title: title, Width: min(max(width/12, minWidth), 10)})
	}

	return cols
}

const defaultWidth = 80

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
