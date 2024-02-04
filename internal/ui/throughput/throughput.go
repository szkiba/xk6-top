// Package throughput contains throughput table UI component.
package throughput

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/navbar"
)

// Model contains throughput table UI component.
type Model struct {
	digest *digest.Digest
	table  table.Model
	width  int
	height int
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case navbar.NavChangedMsg:
		m.update()
	case *digest.Digest:
		m.digest = msg
		m.update()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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

	for name, agg := range cumulative {
		_, hasRate := agg["rate"]
		_, hasMetric := m.digest.FindMetric(name)
		if hasMetric && hasRate {
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
			continue
		}

		row = append(row, name)
		for _, agg := range aggregateNames {
			row = append(row, met.Contains.Format(dig[agg])+"/s")
		}

		rows = append(rows, row)
	}

	m.table.SetRows(rows)
}

// View implements tea.Model.
func (m Model) View() string {
	if len(m.table.Rows()) == 0 {
		return ""
	}

	return m.table.View()
}

// New creates new throughput instance.
func New() Model {
	styles := table.Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#6e59de")).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}

	m := Model{
		table: table.New(
			table.WithKeyMap(table.DefaultKeyMap()),
			table.WithFocused(false),
			table.WithStyles(styles),
			table.WithColumns(columns()),
		),
	}

	return m
}

func columns() []table.Column {
	cols := make([]table.Column, 0, len(aggregateNames)+1)

	cols = append(cols, table.Column{Title: "metric", Width: 32})

	for _, title := range aggregateNames {
		cols = append(cols, table.Column{Title: title, Width: 10})
	}

	return cols
}

var aggregateNames = []string{"rate", "peak"} //nolint:gochecknoglobals
