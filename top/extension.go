// Package top contains top k6 extension.
package top

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/szkiba/xk6-top/internal/ui/app"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

const (
	// OutputName contains extensions's output name.
	OutputName = "top"

	envDashboardPort = "K6_WEB_DASHBOARD_PORT"
	envDashboardHost = "K6_WEB_DASHBOARD_HOST"
	envDashboard     = "K6_WEB_DASHBOARD"

	defaultPort = "5665"
	defaultHost = "127.0.0.1"
)

type extension struct {
	description  string
	url          string
	prog         *tea.Program
	stopCallback func(error)

	normalStop atomic.Bool
}

// New returns new top extension instance.
func New(params output.Params) (output.Output, error) { //nolint:ireturn
	log.SetOutput(params.StdErr)

	port, hasPort := params.Environment[envDashboardPort]
	if !hasPort {
		port = defaultPort
	}

	host, hasHost := params.Environment[envDashboardHost]
	if !hasHost {
		host = defaultHost
	}

	addr := net.JoinHostPort(host, port)

	return &extension{
		description: params.OutputType,
		url:         fmt.Sprintf("http://%s/events", addr),
	}, nil
}

func (ext *extension) Description() string {
	return ext.description
}

func (ext *extension) Start() error {
	ext.prog = tea.NewProgram(app.New(ext.url), tea.WithAltScreen())

	go func() {
		if _, err := ext.prog.Run(); err != nil {
			ext.stopCallback(err)
		}

		if !ext.normalStop.Load() {
			ext.stopCallback(errAbort)
		}
	}()

	return nil
}

// SetTestRunStopCallback accepts a callback which can stop the test run mid-way through.
func (ext *extension) SetTestRunStopCallback(callback func(error)) {
	ext.stopCallback = callback
}

// Stop stops the extension.
func (ext *extension) Stop() error {
	ext.normalStop.Store(true)
	ext.prog.Send(app.StopMsg{})

	return captureEnd()
}

func (*extension) AddMetricSamples(_ []metrics.SampleContainer) {}

var errAbort = errext.WithExitCodeIfNone(
	errext.WithAbortReasonIfNone(
		tea.ErrProgramKilled,
		errext.AbortedByUser,
	),
	exitcodes.ExternalAbort,
)
