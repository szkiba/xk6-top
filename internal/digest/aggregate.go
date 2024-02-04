package digest

import (
	"time"
)

// Aggregate contains aggregated values by aggregation name.
type Aggregate map[string]float64

func (agg Aggregate) clone() Aggregate {
	other := make(Aggregate, len(agg))

	for key, value := range agg {
		other[key] = value
	}

	return other
}

// Aggregates contains aggregetes by metric name.
type Aggregates map[string]Aggregate

// Time returns "value" aggregate for "time" metric.
func (a Aggregates) Time() time.Time {
	agg, hasTime := a["time"]
	if !hasTime {
		return time.Time{}
	}

	val, hasValue := agg["value"]
	if !hasValue {
		return time.Time{}
	}

	return time.UnixMilli(int64(val))
}

func (a Aggregates) clone() Aggregates {
	other := make(Aggregates, len(a))

	for key, value := range a {
		other[key] = value.clone()
	}

	return other
}
