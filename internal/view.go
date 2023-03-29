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
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	v1 "go.k6.io/k6/api/v1"
	"go.k6.io/k6/metrics"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type TableView struct {
	*tview.Table
	propertyNames []string
	metricTypes   []metrics.MetricType
}

func NewTableView(title string, propertyNames []string, metricTypes []metrics.MetricType) *TableView {
	view := &TableView{
		Table:         tview.NewTable(),
		propertyNames: propertyNames,
		metricTypes:   metricTypes,
	}

	view.SetBorder(true).SetBorderColor(tcell.ColorMediumPurple).SetTitle(title).SetTitleColor(tcell.ColorMediumPurple)

	return view
}

func (view *TableView) reset() {
	view.Clear()
	view.addHeader()
}

func (view *TableView) addHeader() {
	view.InsertRow(0)
	view.InsertColumn(0)

	cell := tview.NewTableCell(propName)

	cell.SetMaxWidth(nameMaxWidth)
	cell.SetAlign(tview.AlignLeft)
	cell.SetTextColor(tcell.ColorYellow)
	view.SetCell(0, 0, cell)

	for i, prop := range view.propertyNames {
		idx := i + 1

		view.InsertColumn(idx)

		cell := tview.NewTableCell(fmt.Sprintf("%*s", valueMaxWidth, prop))

		cell.SetAlign(tview.AlignRight)
		cell.SetTextColor(tcell.ColorYellow)
		view.SetCell(0, idx, cell)
	}
}

func (view *TableView) renderValues(samples map[string]v1.Metric, metricNames []string) {
	for idx, key := range metricNames {
		row := idx + 1
		sample := samples[key]

		color := stripeColor[idx%2]

		cell := tview.NewTableCell(key)
		cell.SetTextColor(color)

		view.SetCell(row, 0, cell)

		col := 1

		for _, prop := range view.propertyNames {
			var cell *tview.TableCell

			if val, ok := sample.Sample[prop]; ok {
				cell = tview.NewTableCell(fmt.Sprintf("%*s", valueMaxWidth, formatNumber(val)))
			} else {
				cell = tview.NewTableCell("")
			}

			cell.SetAlign(tview.AlignRight)
			cell.SetTextColor(color)

			view.SetCell(row, col, cell)
			col++
		}

		view.SetCell(row, 0, cell)
	}
}

func (view *TableView) inMetricTypes(typ metrics.MetricType) bool {
	for _, t := range view.metricTypes {
		if t == typ {
			return true
		}
	}

	return false
}

func (view *TableView) extractMetricNames(samples map[string]v1.Metric) []string {
	names := []string{}

	for key, sample := range samples {
		if view.inMetricTypes(sample.Type.Type) {
			names = append(names, key)
		}
	}

	sort.Strings(names)

	return names
}

func (view *TableView) Update(samples map[string]v1.Metric) {
	names := view.extractMetricNames(samples)

	if rows := view.GetRowCount(); rows == 0 || len(names) < rows-1 {
		view.reset()
	}

	view.renderValues(samples, names)
}

func NewTrendView() *TableView {
	return NewTableView(" trend ", trendProps, []metrics.MetricType{metrics.Trend})
}

func NewRateView() *TableView {
	return NewTableView(" rate ", rateProps, []metrics.MetricType{metrics.Rate})
}

func NewCounterView() *TableView {
	return NewTableView(" counter ", counterProps, []metrics.MetricType{metrics.Counter})
}

func NewGaugeView() *TableView {
	return NewTableView(" gauge ", gaugeProps, []metrics.MetricType{metrics.Gauge})
}

const precision = 4

var printer = message.NewPrinter(language.English)

func formatNumber(num float64) string {
	return printer.Sprint(number.Decimal(num, number.Precision(precision)))
}

const (
	nameMaxWidth  = 40
	valueMaxWidth = 8
	propName      = "name"
)

var (
	trendProps   = []string{"avg", "min", "med", "max", "p(90)", "p(95)"}
	counterProps = []string{"rate", "count"}
	rateProps    = []string{"rate"}
	gaugeProps   = []string{"value"}
)

var stripeColor = []tcell.Color{tcell.ColorLightGray, tcell.ColorLightSlateGray}
