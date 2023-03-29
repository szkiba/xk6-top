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

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

type Top struct {
	buffer *output.SampleBuffer

	flusher *output.PeriodicFlusher
	logger  logrus.FieldLogger

	options *Options

	cumulative *Meter

	description string

	app *App
}

var _ output.Output = (*Top)(nil)

func NewTop(params output.Params) (*Top, error) { //nolint:ireturn
	opts, err := ParseOptions(params.ConfigArgument)
	if err != nil {
		return nil, err
	}

	top := &Top{
		logger:      params.Logger,
		options:     opts,
		description: fmt.Sprintf("%s (period=%s)", params.OutputType, opts.Period),
		buffer:      nil,
		flusher:     nil,
		cumulative:  nil,
		app:         nil,
	}

	return top, nil
}

func (top *Top) Description() string {
	return top.description
}

func (top *Top) Start() error {
	var err error

	top.cumulative = NewMeter(0)

	top.buffer = new(output.SampleBuffer)

	flusher, err := output.NewPeriodicFlusher(top.options.Period, top.flush)
	if err != nil {
		return err
	}

	top.flusher = flusher

	top.app = NewApp()

	top.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyETX {
			top.logger.Error("test run was aborted because k6 received a 'interrupt' signal")
			top.stopApp()
			top.logger.Fatal()
		}

		return event
	})

	go func() {
		if err := top.app.Run(); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (top *Top) stopApp() {
	top.app.Stop()
	CaptureEnd()
}

func (top *Top) Stop() error {
	top.flusher.Stop()
	top.stopApp()

	return nil
}

func (top *Top) AddMetricSamples(samples []metrics.SampleContainer) {
	top.buffer.AddMetricSamples(samples)
}

func (top *Top) flush() {
	samples := top.buffer.GetBufferedSamples()

	top.update(samples, top.cumulative)
}

func (top *Top) update(containers []metrics.SampleContainer, meter *Meter) {
	data, err := meter.Update(containers)
	if err != nil {
		top.logger.WithError(err).Warn("Error while processing samples")

		return
	}

	top.app.Update(data)
}
