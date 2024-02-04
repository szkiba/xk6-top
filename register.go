// Package top contains the assembly and registration of the output extension.
package top

import (
	"github.com/szkiba/xk6-top/top"
	"go.k6.io/k6/output"
)

func init() {
	top.CaptureStart()
	output.RegisterExtension(top.OutputName, top.New)
}
