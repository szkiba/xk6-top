package cmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/szkiba/xk6-top/internal/ui/app"
)

//go:embed run.md
var runHelp string

func runCmd() *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "run",
		Short:   "k6 test runner and terminal-based metrics dashboard viewer",
		Long:    runHelp,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(args)
		},

		SilenceUsage:       true,
		SilenceErrors:      true,
		DisableAutoGenTag:  true,
		DisableFlagParsing: true,
	}

	return cmd
}

const (
	envDashboardPort = "K6_WEB_DASHBOARD_PORT"
	envDashboardHost = "K6_WEB_DASHBOARD_HOST"
	envDashboard     = "K6_WEB_DASHBOARD"

	defaultPort = "5665"
	defaultHost = "127.0.0.1"
)

//nolint:forbidigo
func endpoint() string {
	host := os.Getenv(envDashboardHost)
	if len(host) == 0 {
		host = defaultHost
	}

	port := os.Getenv(envDashboardPort)
	if len(port) == 0 {
		port = defaultPort
	}

	return fmt.Sprintf("http://%s/events", net.JoinHostPort(host, port))
}

//nolint:forbidigo
func showAndRemove(output *os.File) {
	if err := output.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	saved, err := os.Open(output.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	_, err = io.Copy(os.Stdout, saved)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if err = saved.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if err = os.Remove(output.Name()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

//nolint:forbidigo
func runRun(args []string) error {
	if os.Getenv(envDashboard) != "true" {
		if err := os.Setenv(envDashboard, "true"); err != nil {
			return err
		}
	}

	k6args := make([]string, len(args)+2)

	k6args[0] = "run"
	k6args[1] = "-q"
	copy(k6args[2:], args)

	output, err := os.CreateTemp("", appname)
	if err != nil {
		return err
	}

	defer showAndRemove(output)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	cmd := exec.CommandContext(ctx, "k6", k6args...) //nolint:gosec

	cmd.Stdout = output
	cmd.Stderr = output

	err = cmd.Start()
	if err != nil {
		return err
	}

	prog := tea.NewProgram(app.New(endpoint()), tea.WithAltScreen())

	var k6err error

	go func() {
		k6err = cmd.Wait()

		var e *exec.ExitError

		if !errors.As(k6err, &e) || e.ExitCode() != -1 {
			prog.Quit()
		}
	}()

	_, teaerr := prog.Run()
	cancel()

	var e *exec.ExitError

	if !errors.As(k6err, &e) || e.ExitCode() != -1 {
		return k6err
	}

	return teaerr
}
