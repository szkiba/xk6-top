package top

import (
	"io"
	"os"
	"strings"
)

//nolint:gochecknoglobals,forbidigo
var (
	extensionActive bool
	stderrOrig      *os.File
)

//nolint:forbidigo
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
	return value == OutputName || strings.HasPrefix(value, OutputName+"=")
}

//nolint:forbidigo
func checkArgs(args []string) (bool, bool, int) {
	argn := len(args)

	var runIndex, outIndex, quietIndex int

	for idx := 0; idx < argn; idx++ {
		arg := args[idx]
		if arg == "run" && runIndex == 0 {
			runIndex = idx

			continue
		}

		if isEnablerFlag(args, idx) {
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

//nolint:forbidigo
func fixArgs(args []string) ([]string, bool) {
	active, quiet, runIndex := checkArgs(os.Args)

	if !active || quiet {
		return args, active
	}

	return addQuietFlag(os.Args, runIndex), active
}

// CaptureStart starts capturing stderr.
//
//nolint:forbidigo
func CaptureStart() {
	os.Args, extensionActive = fixArgs(os.Args)

	if !extensionActive {
		return
	}

	stderrOrig = os.Stderr //nolint:forbidigo
	os.Stderr = createTemp("stderr")

	if os.Getenv(envDashboard) != "true" {
		os.Setenv(envDashboard, "true") //nolint:errcheck,gosec
	}
}

//nolint:forbidigo
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

//nolint:forbidigo
func captureEnd() error {
	if !extensionActive {
		return nil
	}

	if err := copyAndRemove(stderrOrig, os.Stderr); err != nil {
		return err
	}

	os.Stderr = stderrOrig

	return nil
}
