package digest

import (
	"sync"
	"time"
)

// Digester is the metrics handling model.
type Digester struct {
	mu     sync.RWMutex
	digest *Digest

	series *series
	peaks  Peaks

	thresholdsEvaluator *thresholdsEvaluator
}

// NewDigester returns new Digester instance.
func NewDigester() *Digester {
	d := new(Digester)

	d.digest = newDigest(nil)

	d.series = newSeries()
	d.thresholdsEvaluator = newThresholdsEvaluator(nil)
	d.peaks = make(Peaks)

	return d
}

func cast[T ConfigData | *ParamData | Aggregates | Metrics](data interface{}) T {
	value, ok := data.(T)
	if !ok {
		panic("")
	}

	return value
}

// Digest returns the latest update result.
func (d *Digester) Digest() *Digest {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.digest
}

// Update process new event and returns a digest data.
func (d *Digester) Update(event *Event) *Digest {
	d.mu.Lock()
	defer d.mu.Unlock()

	dig := newDigest(d.digest)

	switch event.Type {
	case EventTypeConfig:
		dig = newDigest(nil)
		dig.State = StatePreparing
		dig.Config = cast[ConfigData](event.Data)
	case EventTypeParam:
		dig.State = StatePreparing
		dig.Param = *cast[*ParamData](event.Data)
		d.thresholdsEvaluator.setConfig(dig.Param.Thresholds)
	case EventTypeStart:
		dig.State = StateStarting
		dig.Start = cast[Aggregates](event.Data).Time()
	case EventTypeStop:
		dig.State = StateFinished
		dig.Stop = cast[Aggregates](event.Data).Time()
	case EventTypeMetric:
		dig.State = StateRunning
		for name, value := range cast[Metrics](event.Data) {
			dig.Metrics[name] = value
		}
	case EventTypeCumulative:
		dig.State = StateRunning
		aggs := cast[Aggregates](event.Data)
		d.peaks.inject(aggs)
		for name, value := range aggs {
			dig.Cumulative[name] = value
		}
		dig.Thresholds = d.thresholdsEvaluator.update(aggs)
	case EventTypeSnapshot:
		dig.State = StateRunning
		data := cast[Aggregates](event.Data)
		for name, value := range data {
			dig.Snapshot[name] = value
		}
		d.peaks.update(data)
		d.series.update(data)
	case EventTypeConnect:
		dig.State = StateConnected
	case EventTypeDisconnect:
		dig.State = StateDetached
	default:
	}

	dig.Playback = dig.State == StateRunning && time.Since(dig.Time()) > 3*dig.Period()

	dig.EventType = event.Type
	d.digest = dig

	return dig
}

// Serie returns data series for a given metric and aggregate name.
func (d *Digester) Serie(metric string, aggregate string) []float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.series.get(metric, aggregate)
}

// Collect registers metric name and aggregate name pair to collect series.
func (d *Digester) Collect(metric string, aggregate string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.series.collect(metric, aggregate)
}
