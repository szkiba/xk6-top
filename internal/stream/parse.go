package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/r3labs/sse/v2"
	"github.com/szkiba/xk6-top/internal/digest"
)

type parser struct {
	metrics digest.Metrics
	names   []string
}

func newParser() *parser {
	return &parser{metrics: make(digest.Metrics)}
}

func (p *parser) parse(msg *sse.Event) (*digest.Event, error) {
	var (
		etype digest.EventType
		edata interface{}
		err   error
	)

	if err = etype.UnmarshalText(msg.Event); err != nil {
		return nil, err
	}

	edata, err = p.unmarshalData(etype, msg.Data)
	if err != nil {
		return nil, err
	}

	return &digest.Event{Type: etype, Data: edata}, nil
}

func (p *parser) unmarshalData(etype digest.EventType, data []byte) (interface{}, error) {
	switch etype {
	case digest.EventTypeMetric:
		return p.parseMetric(data)

	case digest.EventTypeParam:
		return p.parseParam(data)

	case digest.EventTypeConfig:
		return p.parseConfig(data)

	case digest.EventTypeStart,
		digest.EventTypeStop,
		digest.EventTypeSnapshot,
		digest.EventTypeCumulative:
		if len(data) > 0 && data[0] == '{' {
			return p.parseAggregatesLegacy(data)
		}

		return p.parseAggregates(data)

	default:
		return nil, nil //nolint:nilnil
	}
}

func (p *parser) parseMetric(data []byte) (interface{}, error) {
	target := make(digest.Metrics)

	if err := json.Unmarshal(data, &target); err != nil {
		return nil, err
	}

	for k, v := range target {
		v.Name = k

		p.metrics[k] = v
	}

	names := make([]string, 0, len(p.metrics))

	for name := range p.metrics {
		names = append(names, name)
	}

	sort.Strings(names)

	p.names = names

	return target, nil
}

func (p *parser) parseParam(data []byte) (interface{}, error) {
	target := new(digest.ParamData)

	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}

	return target, nil
}

func (p *parser) parseConfig(data []byte) (interface{}, error) {
	target := make(digest.ConfigData)

	if err := json.Unmarshal(data, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (p *parser) parseAggregatesLegacy(data []byte) (interface{}, error) {
	target := make(digest.Aggregates)

	if err := json.Unmarshal(data, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (p *parser) parseAggregates(data []byte) (interface{}, error) {
	var samples [][]float64

	if err := json.Unmarshal(data, &samples); err != nil {
		return nil, err
	}

	target := make(digest.Aggregates)

	for metricIdx := range samples {
		metric, err := p.getMetric(metricIdx)
		if err != nil {
			return nil, err
		}

		agg, err := p.parseAggregate(samples[metricIdx], metric.Type)
		if err != nil {
			return nil, err
		}

		target[metric.Name] = agg
	}

	return target, nil
}

func (p *parser) getMetric(idx int) (*digest.Metric, error) {
	if idx >= len(p.names) {
		return nil, fmt.Errorf("%w: metric index out of range %d", errData, idx)
	}

	name := p.names[idx]

	metric, found := p.metrics[name]
	if !found {
		return nil, fmt.Errorf("%w: unknown metric name %s", errData, name)
	}

	return metric, nil
}

func (p *parser) parseAggregate(data []float64, mt digest.MetricType) (digest.Aggregate, error) {
	names := aggregateNames(mt)
	if len(names) == 0 {
		return nil, fmt.Errorf("%w: no metric names for type %s", errData, mt.String())
	}

	if len(data) != len(names) {
		return nil, fmt.Errorf(
			"%w: metric definition mismatch %d - %d - %s",
			errData,
			len(data),
			len(names),
			names,
		)

		//		return nil, fmt.Errorf("%w: metric definition mismatch %s", errData, mt.String())
	}

	agg := make(digest.Aggregate, len(names))

	for idx := range names {
		agg[names[idx]] = data[idx]
	}

	return agg, nil
}

func aggregateNames(mtype digest.MetricType) []string {
	switch mtype {
	case digest.MetricTypeGauge:
		return gaugeAggregateNames
	case digest.MetricTypeRate:
		return rateAggregateNames
	case digest.MetricTypeCounter:
		return counterAggregateNames
	case digest.MetricTypeTrend:
		return trendAggregateNames
	default:
		return nil
	}
}

//nolint:gochecknoglobals
var (
	gaugeAggregateNames   = []string{"value"}
	rateAggregateNames    = []string{"rate"}
	counterAggregateNames = []string{"count", "rate"}
	trendAggregateNames   = []string{"avg", "max", "med", "min", "p(90)", "p(95)", "p(99)"}
)

var errData = errors.New("invalid data")
