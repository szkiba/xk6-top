// Package main contains CLI main entry point.
package main

import (
	"os"

	"github.com/szkiba/xk6-top/internal/cmd"
)

//go:generate go run ../../tools/gendoc README.md

func main() {
	cmd.Execute(os.Args[1:]) //nolint:forbidigo
}
