// Package cmd contains mdcode CLI interface.
package cmd

import (
	"context"
	_ "embed"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/szkiba/xk6-top/internal/ui/app"
)

//go:generate go run ../../tools/gendoc ../../README.md

//nolint:gochecknoglobals
var (
	version = "dev"
	appname = "k6top"
)

// Execute executes root command.
func Execute(args []string) {
	cmd := RootCmd()

	cmd.SetArgs(args)
	cobra.CheckErr(RootCmd().Execute())
}

//go:embed root.md
var rootHelp string

// RootCmd returns a roo cobra command. Exported for documentation generation purpose only.
func RootCmd() *cobra.Command {
	var baseURL string

	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     appname + " [flags]",
		Short:   "Terminal based metrics dashboard viewer for k6",
		Long:    rootHelp,
		Version: version,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			endpoint, err := url.JoinPath(baseURL, "/events")
			if err != nil {
				return err
			}

			return runRoot(endpoint)
		},

		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	flags := cmd.Flags()

	flags.StringVarP(&baseURL, "url", "u", "http://127.0.0.1:5665", "k6 web dashboard URL")

	cmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s" .}}{{end}}{{printf " version %s\n" .Version}}`,
	)

	cmd.AddCommand(runCmd())

	return cmd
}

func runRoot(url string) error {
	ctx := context.TODO()

	prog := tea.NewProgram(
		app.New(url),
		tea.WithContext(ctx),
		tea.WithAltScreen(),
	)

	prog.SetWindowTitle("k6 dashboard")

	_, err := prog.Run()

	return err
}
