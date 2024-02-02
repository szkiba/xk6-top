// Package app contains application model.
package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/szkiba/xk6-top/internal/digest"
	"github.com/szkiba/xk6-top/internal/stream"
	"github.com/szkiba/xk6-top/internal/ui/detail"
	"github.com/szkiba/xk6-top/internal/ui/divider"
	"github.com/szkiba/xk6-top/internal/ui/help"
	"github.com/szkiba/xk6-top/internal/ui/navbar"
	"github.com/szkiba/xk6-top/internal/ui/overview"
	"github.com/szkiba/xk6-top/internal/ui/statbar"
	"github.com/szkiba/xk6-top/internal/ui/status"
	"github.com/szkiba/xk6-top/internal/ui/summary"
	"github.com/szkiba/xk6-top/internal/ui/theme"
)

// Model contains application model.
type Model struct {
	stream   chan tea.Msg
	digester *digest.Digester
	digest   *digest.Digest

	theme *theme.Theme

	ready bool

	width  int
	height int

	navbar      navbar.Model
	tabContents []tea.Model
	divider     tea.Model
	status      tea.Model

	sseEndpoint string

	sseContext context.Context
	sseCancel  context.CancelFunc
}

type quitMsg struct{}

// StopMsg sent from outside of UI to stop the UI.
type StopMsg struct{}

// New creates new model instance.
func New(sseEndpoint string) *Model {
	m := new(Model)

	m.theme = theme.Detect()

	m.sseContext, m.sseCancel = context.WithCancel(context.TODO())
	m.sseEndpoint = sseEndpoint
	m.navbar = navbar.New(m.theme, navItems)
	m.stream = make(chan tea.Msg)
	m.digester = digest.NewDigester()
	m.tabContents = []tea.Model{
		overview.New(m.theme),
		summary.New(digest.MetricTypeTrend),
		summary.New(digest.MetricTypeCounter),
		summary.New(digest.MetricTypeRate),
		summary.New(digest.MetricTypeGauge),
		detail.New(m.theme, m.digester, httpPanels),
		detail.New(m.theme, m.digester, gRPCPanels),
		detail.New(m.theme, m.digester, wsPanels),
		detail.New(m.theme, m.digester, browserPanels),
		help.New(m.theme),
	}
	m.divider = divider.New(m.theme)
	m.status = status.New(m.theme)

	return m
}

func (m *Model) readEvent() tea.Msg {
	for {
		select {
		case event := <-m.stream:
			return event
		case <-m.sseContext.Done():
			return quitMsg{}
		}
	}
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		stream.Subscribe(m.sseContext, m.sseEndpoint, m.stream),
		m.readEvent,
	)
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.sseCancel()
			if m.digest.GetState() == digest.StateDetached {
				return m, tea.Quit
			}

			return m, nil
		default:
		}
	case StopMsg:
		m.sseCancel()
		return m, tea.Quit
	case quitMsg:
		return m, tea.Quit
	case *digest.Event:
		cmds = append(cmds, m.readEvent)
		m.digest = m.digester.Update(msg)
		cmds = append(cmds, m.enableNav(m.digest)...)
		if !m.digest.Playback {
			cmds = append(cmds, digestCmd(m.digest))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case navbar.NavChangedMsg:
		wsmsg := tea.WindowSizeMsg{Width: m.width, Height: m.height}
		cmds = append(cmds, func() tea.Msg { return wsmsg })

		if m.digest != nil {
			cmds = append(cmds, digestCmd(m.digest))
		}
	}

	m.navbar, cmd = m.navbar.Update(msg)
	cmds = append(cmds, cmd)

	m.tabContents[m.navbar.Active], cmd = m.tabContents[m.navbar.Active].Update(msg)
	cmds = append(cmds, cmd)

	m.divider, cmd = m.divider.Update(msg)
	cmds = append(cmds, cmd)

	m.status, cmd = m.status.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) enableNav(dig *digest.Digest) []tea.Cmd {
	var cmds []tea.Cmd

	for idx, item := range navItems {
		if !item.Disabled {
			continue
		}

		metric, hasEnabler := navItemEnabler[item.Label]
		if !hasEnabler {
			continue
		}

		_, found := dig.FindMetric(metric)
		if found {
			cmds = append(cmds, navbar.Enable(idx))
		}
	}

	return cmds
}

// View implements tea.Model.
func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.navbar.View(),
		m.divider.View(),
		m.tabContents[m.navbar.Active].View(),
		m.status.View(),
	)
}

func digestCmd(digest *digest.Digest) tea.Cmd {
	return func() tea.Msg {
		return digest
	}
}

//nolint:gochecknoglobals
var (
	navItems = []*navbar.Item{
		{Label: "Overview"},
		{Label: "Trends", Disabled: true},
		{Label: "Counters", Disabled: true},
		{Label: "Rates", Disabled: true},
		{Label: "Gauges", Disabled: true},
		{Label: "HTTP", Disabled: true},
		{Label: "gRPC", Disabled: true},
		{Label: "WS", Disabled: true},
		{Label: "Browser", Disabled: true},
		{Label: "Help"},
	}

	navItemEnabler = map[string]string{
		"Trends":   "iteration_duration",
		"Counters": "iterations",
		"Rates":    "vus_max",
		"Gauges":   "vus_max",
		"HTTP":     "http_reqs",
		"gRPC":     "grpc_streams_msgs_sent",
		"WS":       "ws_msgs_sent",
		"Browser":  "browser_http_req_duration",
	}

	httpPanels = []*statbar.Panel{
		{
			Label:     "Req Rate",
			Caption:   "HTTP Request Rate",
			Metric:    "http_reqs",
			Aggregate: "rate",
		},
		{
			Label:     "Req Duration",
			Caption:   "HTTP Request Duration",
			Metric:    "http_req_duration",
			Aggregate: "avg",
		},
		{
			Label:     "Req Failed",
			Caption:   "HTTP Failed Rate",
			Metric:    "http_req_failed",
			Aggregate: "rate",
		},
		{
			Label:     "Data Received",
			Caption:   "Data Received Rate",
			Metric:    "data_received",
			Aggregate: "rate",
		},
		{
			Label:     "Data Sent",
			Caption:   "Data Sent Rate",
			Metric:    "data_sent",
			Aggregate: "rate",
		},
	}

	gRPCPanels = []*statbar.Panel{
		{
			Label:     "Msgs Sent",
			Caption:   "Messages Sent Rate",
			Metric:    "grpc_streams_msgs_sent",
			Aggregate: "rate",
		},
		{
			Label:     "Msgs Received",
			Caption:   "Messages Received Rate",
			Metric:    "grpc_streams_msgs_received",
			Aggregate: "rate",
		},
		{
			Label:     "Req Duration",
			Caption:   "Request Duration",
			Metric:    "grpc_req_duration",
			Aggregate: "avg",
		},
		{
			Label:     "Streams",
			Caption:   "Streams Rate",
			Metric:    "grpc_streams",
			Aggregate: "rate",
		},
	}

	wsPanels = []*statbar.Panel{
		{
			Label:     "Msgs Sent",
			Caption:   "Messages Sent Rate",
			Metric:    "ws_msgs_sent",
			Aggregate: "rate",
		},
		{
			Label:     "Msgs Received",
			Caption:   "Messages Received Rate",
			Metric:    "ws_msgs_received",
			Aggregate: "rate",
		},
		{
			Label:     "Ping Duration",
			Caption:   "Ping Duration",
			Metric:    "ws_ping",
			Aggregate: "avg",
		},
		{
			Label:     "Session Duration",
			Caption:   "Session Duration",
			Metric:    "ws_session_duration",
			Aggregate: "avg",
		},
	}

	browserPanels = []*statbar.Panel{
		{
			Label:     "Req Duration",
			Caption:   "Request Duration",
			Metric:    "browser_http_req_duration",
			Aggregate: "avg",
		},
		{
			Label:     "Req Failed",
			Caption:   "Request Failed Rate",
			Metric:    "browser_http_req_failed",
			Aggregate: "rate",
		},
		{
			Label:     "Data Received",
			Caption:   "Data Received Rate",
			Metric:    "browser_data_received",
			Aggregate: "rate",
		},
		{
			Label:     "Data Sent",
			Caption:   "Data Sent Rate",
			Metric:    "browser_data_sent",
			Aggregate: "rate",
		},
	}
)
