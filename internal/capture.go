// MIT License
//
// Copyright (c) 2023 Iván Szkiba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package internal

import (
	"io"
	"os"
	"strings"
)

var (
	extensionActive bool
	stderrOrig      *os.File
)

func createTemp(name string) *os.File {
	tmp, err := os.CreateTemp(os.TempDir(), name)
	if err != nil {
		panic(err)
	}

	return tmp
}

func isEnablerFlag(args []string, idx int) bool {
	return (args[idx] == "--out" || args[idx] == "-o") &&
		idx < len(args)-1 && isEnablerFlagValue(args[idx+1])
}

func isEnablerFlagValue(value string) bool {
	return value == "top" || strings.HasPrefix(value, "top=")
}

func checkArgs(args []string) (bool, bool, int) {
	argn := len(args)

	var runIndex, outIndex, quietIndex int

	for idx := 0; idx < argn; idx++ {
		arg := args[idx]
		if arg == "run" && runIndex == 0 {
			runIndex = idx

			continue
		}

		if runIndex != 0 && isEnablerFlag(args, idx) {
			outIndex = idx

			continue
		}

		if arg == "-q" || arg == "--quiet" {
			quietIndex = idx
		}
	}

	if outIndex == 0 && isEnablerFlagValue(os.Getenv("K6_OUT")) {
		outIndex = -1
	}

	return outIndex != 0, quietIndex != 0, runIndex
}

func addQuietFlag(args []string, runIndex int) []string {
	quietIndex := runIndex + 1

	if quietIndex >= len(args) {
		return args
	}

	args = append(args, "")
	copy(args[quietIndex+1:], args[quietIndex:])
	args[quietIndex] = "-q"

	return args
}

func fixArgs(args []string) ([]string, bool) {
	active, quiet, runIndex := checkArgs(os.Args)

	if !active || quiet {
		return args, active
	}

	return addQuietFlag(os.Args, runIndex), active
}

func CaptureStart() {
	os.Args, extensionActive = fixArgs(os.Args)

	if extensionActive {
		stderrOrig = os.Stderr
		os.Stderr = createTemp("stderr")
	}
}

func copyAndRemove(dst *os.File, src *os.File) error {
	if err := src.Sync(); err != nil {
		return err
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return os.Remove(src.Name())
}

func CaptureEnd() {
	if !extensionActive {
		return
	}

	if err := copyAndRemove(stderrOrig, os.Stderr); err != nil {
		panic(err)
	}

	os.Stderr = stderrOrig
}
