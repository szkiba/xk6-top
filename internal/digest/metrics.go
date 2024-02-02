package digest

import "strings"

// Metrics contains metric by name.
type Metrics map[string]*Metric

// Find finds metric (or parent metric in case of submetric) by name.
func (mets Metrics) Find(name string) (*Metric, bool) {
	idx := strings.IndexRune(name, '{')
	if idx > 0 {
		name = name[:idx]
	}

	met, found := mets[name]

	return met, found
}
