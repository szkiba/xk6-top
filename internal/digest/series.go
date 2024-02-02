package digest

type series struct {
	values   map[string][]float64
	retained map[string]struct{}
}

func newSeries() *series {
	s := new(series)

	s.values = make(map[string][]float64)
	s.retained = make(map[string]struct{})

	s.collect("time", "value")

	return s
}

func (s *series) collect(metric string, aggregate string) {
	s.retained[s.key(metric, aggregate)] = struct{}{}
}

func (s *series) update(aggs Aggregates) {
	var length int
	if times, hasTime := s.values[s.key("time", "value")]; hasTime {
		length = len(times)
	}

	for metric, agg := range aggs {
		for prop, value := range agg {
			key := s.key(metric, prop)
			if _, retained := s.retained[key]; !retained {
				continue
			}
			serie, found := s.values[key]
			if !found {
				serie = make([]float64, length, length+1)
			}

			serie = append(serie, value)
			s.values[key] = serie
		}
	}
}

func (s *series) key(metric string, agg string) string {
	return metric + "#" + agg
}

func (s *series) get(metric string, agg string) []float64 {
	serie, found := s.values[s.key(metric, agg)]
	if !found {
		return []float64{}
	}

	ret := make([]float64, len(serie))
	copy(ret, serie)

	return ret
}
