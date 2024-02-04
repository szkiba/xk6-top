// Package main contains help ansi text generator.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
)

//nolint:forbidigo
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

	return os.WriteFile(filepath.Join(dir, base), out, 0o600)
}

//nolint:forbidigo
func main() {
	if len(os.Args) != 3 {
		panic("Usage genhelp input-file output-dir")
	}

	var margin uint

	glamour.DarkStyleConfig.Document.Margin = &margin
	glamour.LightStyleConfig.Document.Margin = &margin

	glamour.DarkStyleConfig.H1.Prefix = "# "
	glamour.DarkStyleConfig.H1.Suffix = ""
	glamour.DarkStyleConfig.H1.Color = nil
	glamour.DarkStyleConfig.H1.BackgroundColor = nil
	glamour.DarkStyleConfig.Code.BackgroundColor = nil

	glamour.LightStyleConfig.H1.Prefix = "# "
	glamour.LightStyleConfig.H1.Suffix = ""
	glamour.LightStyleConfig.H1.Color = nil
	glamour.LightStyleConfig.H1.BackgroundColor = nil
	glamour.LightStyleConfig.Code.BackgroundColor = nil

	if err := render(os.Args[1], "dark", os.Args[2]); err != nil {
		panic(err)
	}

	if err := render(os.Args[1], "light", os.Args[2]); err != nil {
		panic(err)
	}
}
