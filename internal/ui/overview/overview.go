// Package overview contains overview tab UI component.
package overview

import (
	_ "embed"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/ui/badge"
	"github.com/szkiba/xk6-top/internal/ui/theme"
	"github.com/szkiba/xk6-top/internal/ui/throughput"
)

// Model contains overview tab UI component.
type Model struct {
	digest     *digest.Digest
	theme      *theme.Theme
	viewport   viewport.Model
	throughput tea.Model
	ready      bool
	width      int
	height     int

	started time.Time
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *digest.Digest:
		m.digest = msg
		if msg.State == digest.StateStarting {
			cmds = append(cmds, tickCmd())
			m.started = msg.Time()
		}
		m.update()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-4)
			m.viewport.YPosition = 2
			m.viewport.Style = m.viewport.Style.Padding(0, 1, 0, 1)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 4
		}

		m.update()

	case tickMsg:
		if m.digest.State == digest.StateStarting {
			cmds = append(cmds, tickCmd())
		}

		m.update()
	}

	m.throughput, cmd = m.throughput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

//go:embed nodata.txt
var nodata string

func (m *Model) nodata() {
	style := m.theme.Primary.Copy().Align(lipgloss.Center, lipgloss.Center).
		Width(m.viewport.Width).
		Height(m.viewport.Height)

	m.viewport.SetContent(style.Render(nodata))
}

func (m *Model) update() {
	if m.digest.GetState() == digest.StateWaiting {
		m.nodata()
		return
	}

	thresholds := m.digest.Thresholds

	style := m.theme.Heading.Copy().
		Width(m.viewport.Width).
		AlignHorizontal(lipgloss.Center)

	title := style.Render("k6 Terminal Dashboard")

	style = lipgloss.NewStyle().Width(m.viewport.Width).AlignHorizontal(lipgloss.Center)

	infos := m.infoLine(thresholds)
	indicators := m.indicatorLine()

	var buff strings.Builder

	buff.WriteString(title)
	buff.WriteString("\n\n")

	buff.WriteString(style.Render(infos))
	buff.WriteString("\n\n")

	if len(indicators) > 0 {
		buff.WriteString(style.Render(indicators))
		buff.WriteString("\n\n")
	}

	if m.digest.GetState() == digest.StateStarting {
		buff.WriteString(m.countdown())
	} else {
		thresholdsSection := m.thresholdsSection(thresholds)
		if len(thresholdsSection) > 0 {
			buff.WriteString(thresholdsSection)
		}

		if throughput := m.throughput.View(); len(throughput) > 0 {
			buff.WriteString("\n")
			buff.WriteString(m.theme.Heading.Render("Throughputs"))
			buff.WriteString("\n\n")
			buff.WriteString(throughput)
		}
	}

	m.viewport.SetContent(buff.String())
}

func (m *Model) infoLine(tresholds *digest.Thresholds) string {
	start := m.digest.Start
	if start.IsZero() {
		return ""
	}

	var items []string

	items = append(
		items,
		badge.Render("start", start.Format("2006-01-02 15:04:05"), digest.Info),
	)

	items = append(items, badge.Render("duration", m.digest.Duration().String(), digest.Info))

	stop := m.digest.Stop
	if stop.IsZero() {
		items = append(
			items,
			badge.Render("remaining", m.digest.TimeLeft().String(), digest.Info),
		)
	} else {
		var lvl digest.Level
		if tresholds == nil {
			lvl = digest.None
		} else {
			lvl = tresholds.Result
		}

		var result string
		switch lvl {
		case digest.Error:
			result = "failed"
		case digest.Ready:
			result = "passed"
		default:
			result = "unknown"
		}

		items = append(items,
			badge.Render("result", result, lvl),
		)
	}

	return strings.Join(items, " ")
}

func (m *Model) indicatorLine() string {
	aggs := m.digest.Cumulative

	ratebadge := func(metric string) string {
		value, found := aggs[metric]
		if !found {
			return ""
		}

		rate, ok := value["rate"]
		if !ok {
			return ""
		}

		lev := digest.Ready

		if rate > 0 {
			lev = digest.Warning
		}

		return badge.Render(metric, digest.ValueTypeDefault.Format(rate)+"/s", lev)
	}

	var items []string

	for _, metric := range []string{"checks", "http_req_failed", "browser_http_req_failed"} {
		if item := ratebadge(metric); len(item) != 0 {
			items = append(items, item)
		}
	}

	return strings.Join(items, " ")
}

func (m *Model) thresholdsSection(thresholds *digest.Thresholds) string {
	if thresholds == nil || len(thresholds.Brief) == 0 {
		return ""
	}

	var buff strings.Builder

	buff.WriteString(m.theme.Heading.Render("Thresholds"))
	buff.WriteString("\n\n")

	width := 0
	names := make([]string, 0, len(thresholds.Brief))

	for metric := range thresholds.Brief {
		if l := len(metric); l > width {
			width = l
		}
		names = append(names, metric)
	}

	sort.Strings(names)

	padd := lipgloss.NewStyle().
		Width(width).
		AlignHorizontal(lipgloss.Left).
		MarginRight(2).
		MarginLeft(1)

	for _, metric := range names {
		results := thresholds.Details[metric]
		sources := thresholds.Source[metric]

		buff.WriteString(padd.Render(metric))

		last := len(sources) - 1
		for idx, src := range sources {
			lvl := results[src]
			buff.WriteString(badge.Value(src, lvl))
			if idx < last {
				buff.WriteString(", ")
			}
		}
		buff.WriteString("\n")
	}

	return buff.String()
}

func (m *Model) countdown() string {
	if m.started.IsZero() {
		return ""
	}

	now := time.Now()
	elapsed := now.Sub(m.started)
	period := time.Millisecond * time.Duration(m.digest.Param.Period)

	if elapsed > period {
		return ""
	}

	var buff strings.Builder

	buff.WriteString(m.theme.Heading.Render("Waiting for the first aggregates..."))
	buff.WriteString("\n\n")

	percent := float64(period-elapsed) / float64(period)

	width := m.viewport.Width - 2

	empty := int(float64(width) * percent)
	full := width - empty

	buff.WriteRune(' ')
	buff.WriteString(m.theme.Primary.Render(strings.Repeat("█", full)))
	buff.WriteString(m.theme.Divider.Render(strings.Repeat("█", empty)))

	return buff.String()
}

// View implements tea.Model.
func (m Model) View() string {
	return m.viewport.View()
}

// New creates new overview instance.
func New(theme *theme.Theme) Model {
	return Model{
		theme:      theme,
		throughput: throughput.New(),
	}
}
