// Package main contains help ansi text generator.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

const fileMode = 0o600

func render(src string, style string, dir string) error {
	input, err := os.ReadFile(filepath.Clean(src))
	if err != nil {
		return err
	}

	r, err := glamour.NewTermRenderer(glamour.WithStylePath(style), glamour.WithWordWrap(-1))
	if err != nil {
		return err
	}

	out, err := r.RenderBytes(input)
	if err != nil {
		return err
	}

	out = bytes.Trim(out, "\n\r")

	base := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src)) + "-" + style + ".ansi"

	return os.WriteFile(filepath.Join(dir, base), out, fileMode)
}

func main() {
	const argcount = 3

	if len(os.Args) != argcount {
		panic("Usage genhelp input-file output-dir")
	}

	var margin uint

	styles.DarkStyleConfig.Document.Margin = &margin
	styles.LightStyleConfig.Document.Margin = &margin

	styles.DarkStyleConfig.H1.Prefix = "# "
	styles.DarkStyleConfig.H1.Suffix = ""
	styles.DarkStyleConfig.H1.Color = nil
	styles.DarkStyleConfig.H1.BackgroundColor = nil
	styles.DarkStyleConfig.Code.BackgroundColor = nil

	styles.LightStyleConfig.H1.Prefix = "# "
	styles.LightStyleConfig.H1.Suffix = ""
	styles.LightStyleConfig.H1.Color = nil
	styles.LightStyleConfig.H1.BackgroundColor = nil
	styles.LightStyleConfig.Code.BackgroundColor = nil

	if err := render(os.Args[1], "dark", os.Args[2]); err != nil {
		panic(err)
	}

	if err := render(os.Args[1], "light", os.Args[2]); err != nil {
		panic(err)
	}
}
