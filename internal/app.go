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
	"github.com/rivo/tview"
	v1 "go.k6.io/k6/api/v1"
)

type App struct {
	*tview.Application
	trend   *TableView
	counter *TableView
	rate    *TableView
	gauge   *TableView
}

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		trend:       NewTrendView(),
		counter:     NewCounterView(),
		rate:        NewRateView(),
		gauge:       NewGaugeView(),
	}

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(app.trend, 0, 3, false) // nolint:gomnd

	col := tview.NewFlex()
	col.SetDirection(tview.FlexRow)
	col.AddItem(app.rate, 0, 1, false).AddItem(app.gauge, 0, 1, false)

	row := tview.NewFlex()
	row.AddItem(app.counter, 0, 7, false).AddItem(col, 0, 5, false) // nolint:gomnd

	flex.AddItem(row, 0, 2, false) // nolint:gomnd
	app.SetRoot(flex, true)

	return app
}

func (app *App) Update(samples map[string]v1.Metric) {
	app.QueueUpdateDraw(func() {
		app.trend.Update(samples)
		app.counter.Update(samples)
		app.rate.Update(samples)
		app.gauge.Update(samples)
	})
}
