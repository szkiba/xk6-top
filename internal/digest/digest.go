// Package digest contains metrics handling model.
package digest

import (
	"maps"
	"time"
)

// Digest is the metrics digest data model.
type Digest struct {
	Config     map[string]interface{}
	Param      ParamData
	Metrics    Metrics
	Cumulative Aggregates
	Snapshot   Aggregates
	Start      time.Time
	Stop       time.Time
	Thresholds *Thresholds
	EventType  EventType
	State      State
	Playback   bool
}

func newDigest(from *Digest) *Digest {
	d := new(Digest)

	if from == nil {
		d.Cumulative = make(Aggregates)
		d.Snapshot = make(Aggregates)
		d.Metrics = make(Metrics)
		d.Config = make(map[string]interface{})
	} else {
		d.Start = from.Start
		d.Stop = from.Stop
		d.Param = from.Param
		d.Cumulative = maps.Clone(from.Cumulative)
		d.Snapshot = maps.Clone(from.Snapshot)
		d.Metrics = maps.Clone(from.Metrics)
		d.State = from.State
		if d.Config == nil {
			d.Config = make(map[string]interface{})
		}
		for key, value := range from.Config {
			d.Config[key] = value
		}
		d.Thresholds = from.Thresholds
		d.Playback = from.Playback
	}

	return d
}

// GetState returns state or StateWaiting if digest is nil.
func (d *Digest) GetState() State {
	if d == nil {
		return StateWaiting
	}

	return d.State
}

// FindMetric finds parent metric's metadata by metric name.
func (d *Digest) FindMetric(metric string) (*Metric, bool) {
	return d.Metrics.Find(metric)
}

// GetMetric returns metric metadata by metric name (exact match).
func (d *Digest) GetMetric(metric string) (*Metric, bool) {
	value, found := d.Metrics[metric]

	return value, found
}

// Period returns period parameter as duration.
func (d *Digest) Period() time.Duration {
	return time.Millisecond * time.Duration(d.Param.Period)
}

// Duration returns test run time intervall as duration.
func (d *Digest) Duration() time.Duration {
	return time.Millisecond * time.Duration(d.Param.EndOffset)
}

// Time returns the latest "time" metric value from event stream.
func (d *Digest) Time() time.Time {
	if tim := d.Cumulative.Time(); !tim.IsZero() {
		return tim
	}

	return d.Start
}

// TimeLeft returns the remaining time as duration.
func (d *Digest) TimeLeft() time.Duration {
	total := d.Duration()

	now := d.Time()
	if now.IsZero() {
		return total
	}

	left := d.Start.Add(total).Sub(now) + time.Second

	return left.Truncate(d.Period())
}

// TimePassed returns elapsed time as duration.
func (d *Digest) TimePassed() time.Duration {
	total := d.Duration()
	now := d.Time()

	var elapsed time.Duration

	if !d.Stop.IsZero() {
		elapsed = total
	} else if !now.IsZero() {
		elapsed = now.Sub(d.Start) + time.Second
		elapsed = elapsed.Truncate(d.Period())
	}

	return elapsed
}

// ProgressPercent return percent value of elapsed time.
func (d *Digest) ProgressPercent() float64 {
	if d.State == StateFinished {
		return 100
	}

	total := d.Duration()
	left := d.TimeLeft()

	if left <= 0 || total == 0 {
		return 0
	}

	return 1.0 - float64(left)/float64(total)
}

// ProgressLabel return elapsed / total value as string.
func (d *Digest) ProgressLabel() string {
	total := d.Duration()
	passed := d.TimePassed()

	return passed.String() + "/" + total.String()
}
