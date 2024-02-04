package digest

// Peaks contains peak rate values by metric name.
type Peaks map[string]float64

func (p Peaks) update(aggs Aggregates) {
	for name, agg := range aggs {
		rate, hasRate := agg["rate"]
		if hasRate {
			prev, hasPrev := p[name]
			if !hasPrev || rate > prev {
				p[name] = rate
			}
		}
	}
}

func (p Peaks) inject(aggs Aggregates) {
	for name, agg := range aggs {
		_, hasRate := agg["rate"]
		if hasRate {
			peak, hasPeak := p[name]
			if hasPeak {
				agg["peak"] = peak
			}
		}
	}
}
