package digest

import (
	"fmt"
	"time"
)

func trunc(dur time.Duration) time.Duration {
	if dur < time.Nanosecond {
		return time.Nanosecond / 2
	}

	if dur < time.Microsecond {
		return time.Nanosecond
	}

	if dur < time.Millisecond {
		return time.Microsecond
	}

	if dur < time.Second {
		return time.Millisecond
	}

	if dur < time.Minute {
		return time.Second / 10
	}

	if dur < time.Hour {
		return time.Minute / 10
	}

	if dur < time.Hour*24 {
		return time.Hour / 10
	}

	return time.Hour
}

func formatDuration(value float64) string {
	dur := time.Duration(float64(time.Millisecond) * value)

	return dur.Truncate(trunc(dur)).String()
}

func formatData(value float64) string {
	unitI, prec := conv(value, dataPrec)
	unit := dataUnit(unitI)

	suffix := unit.String()

	return fmt.Sprintf("%.*f%s", prec, value/float64(unit), suffix)
}

// Unit returns unit and precision for displaying a given value.
func Unit(value float64) (int, int) {
	return conv(value, dataPrec)
}

//nolint:gochecknoglobals
var (
	countPrec = [][]int{
		{10 * int(countUnitP), int(countUnitP), 0},
		{int(countUnitP), int(countUnitP), 1},

		{10 * int(countUnitT), int(countUnitT), 0},
		{int(countUnitT), int(countUnitT), 1},

		{10 * int(countUnitG), int(countUnitG), 0},
		{int(countUnitG), int(countUnitG), 1},

		{10 * int(countUnitM), int(countUnitM), 0},
		{int(countUnitM), int(countUnitM), 1},

		{10 * int(countUnitK), int(countUnitK), 0},
		{int(countUnitK), int(countUnitK), 1},

		{100 * int(countUnitOne), int(countUnitOne), 0},
		{10 * int(countUnitOne), int(countUnitOne), 1},
		{int(countUnitOne), int(countUnitOne), 2},
	}

	dataPrec = [][]int{
		{10 * int(dataUnitPB), int(dataUnitPB), 0},
		{int(dataUnitPB), int(dataUnitPB), 1},

		{10 * int(dataUnitTB), int(dataUnitTB), 0},
		{int(dataUnitTB), int(dataUnitTB), 1},

		{10 * int(dataUnitGB), int(dataUnitGB), 0},
		{int(dataUnitGB), int(dataUnitGB), 1},

		{10 * int(dataUnitMB), int(dataUnitMB), 0},
		{int(dataUnitMB), int(dataUnitMB), 1},

		{10 * int(dataUnitKB), int(dataUnitKB), 0},
		{int(dataUnitKB), int(dataUnitKB), 1},

		{100 * int(dataUnitB), int(dataUnitB), 0},
		{10 * int(dataUnitB), int(dataUnitB), 1},
		{int(dataUnitB), int(dataUnitB), 2},
	}
)

func conv(value float64, prec [][]int) (int, int) {
	for i := 0; i < len(prec); i++ {
		if value > float64(prec[i][0]) {
			return prec[i][1], prec[i][2]
		}
	}

	return prec[len(prec)-1][1], prec[len(prec)-1][2]
}

func formatDefault(value float64) string {
	unitI, prec := conv(value, countPrec)
	unit := countUnit(unitI)

	var suffix string

	if unit != countUnitOne {
		suffix = unit.String()
	}

	return fmt.Sprintf("%.*f%s", prec, value/float64(unit), suffix)
}

type countUnit int

const (
	countUnitOne countUnit = 1
	countUnitK   countUnit = countUnitOne * 1000
	countUnitM   countUnit = countUnitK * 1000
	countUnitG   countUnit = countUnitM * 1000
	countUnitT   countUnit = countUnitG * 1000
	countUnitP   countUnit = countUnitT * 1000
	countUnitE   countUnit = countUnitP * 1000
)

//go:generate go run github.com/dmarkham/enumer@latest -text -transform lower -trimprefix countUnit -type countUnit

type dataUnit int

const (
	dataUnitB  dataUnit = 1
	dataUnitKB dataUnit = dataUnitB * 1000
	dataUnitMB dataUnit = dataUnitKB * 1000
	dataUnitGB dataUnit = dataUnitMB * 1000
	dataUnitTB dataUnit = dataUnitGB * 1000
	dataUnitPB dataUnit = dataUnitTB * 1000
	dataUnitEB dataUnit = dataUnitPB * 1000
)

//go:generate go run github.com/dmarkham/enumer@latest -text -transform lower -trimprefix dataUnit -type dataUnit
