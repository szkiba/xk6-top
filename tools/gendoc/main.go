// Package main contains CLI doc generator.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"
	"github.com/szkiba/xk6-top/internal/cmd"
)

//nolint:forbidigo
func checkerr(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "gendoc: error: %s\n", err)
	os.Exit(1)
}

func linkHandler(name string) string {
	link := strings.ReplaceAll(strings.TrimSuffix(name, ".md"), "_", "-")

	return "#" + link
}

func fprintf(out io.Writer, format string, args ...any) {
	_, err := fmt.Fprintf(out, format, args...)
	checkerr(err)
}

//nolint:forbidigo
func main() {
	if len(os.Args) != 2 { //nolint:gomnd
		fmt.Fprint(os.Stderr, "usage: gendoc filename")
		os.Exit(1)
	}

	root := cmd.RootCmd()

	var buff bytes.Buffer

	regions := map[string]string{}

	checkerr(doc.GenMarkdownCustom(root, &buff, linkHandler))

	for _, cmd := range root.Commands() {
		if strings.HasPrefix(cmd.Use, "help") || !cmd.Runnable() {
			continue
		}

		fprintf(&buff, "---\n")
		checkerr(doc.GenMarkdownCustom(cmd, &buff, linkHandler))
	}

	cli := buff.String()

	cli = strings.ReplaceAll(cli, "### Options inherited from parent commands", "### Global Flags")
	cli = strings.ReplaceAll(cli, "### Options", "### Flags")

	regions["cli"] = cli

	help, err := os.ReadFile(filepath.Join("..", "..", "docs", "help.md"))
	checkerr(err)

	regions["help"] = string(help)

	readme := filepath.Clean(os.Args[1])

	src, err := os.ReadFile(readme)
	checkerr(err)

	for name, value := range regions {
		res, found, err := replace(src, name, []byte(value))
		checkerr(err)

		if found {
			src = res
		}
	}

	checkerr(os.WriteFile(readme, src, 0o600)) //nolint:gomnd
}
