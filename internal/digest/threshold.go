package digest

import (
	"strings"

	"github.com/Knetic/govaluate"
)

// Failure holds thresholds failure parameters.
type Failure struct {
	Metric     string
	Thresholds []string
}

// Thresholds holds all the result of the thresholds checks.
type Thresholds struct {
	Result   Level
	Details  map[string]map[string]Level
	Brief    map[string]Level
	Source   map[string][]string
	Failures []*Failure
}

func newThresholds(
	source map[string][]string,
	failures []*Failure,
	details map[string]map[string]Level,
) *Thresholds {
	tre := &Thresholds{
		Details:  make(map[string]map[string]Level),
		Brief:    make(map[string]Level),
		Source:   source,
		Failures: failures,
	}

	var result Level

	if len(details) > 0 {
		result = Ready
	}

	for metric, results := range details {
		tre.Details[metric] = results

		lev := Ready

		for _, result := range results {
			if result > lev {
				lev = result
			}
		}

		tre.Brief[metric] = lev

		if lev > result {
			result = lev
		}
	}

	tre.Result = result

	return tre
}

type thresholdsEvaluator struct {
	results map[string]map[string]Level
	source  map[string][]string
	config  thresholdsConfig
}

func newThresholdsEvaluator(source map[string][]string) *thresholdsEvaluator {
	t := new(thresholdsEvaluator)

	t.setConfig(source)
	t.results = make(map[string]map[string]Level)

	return t
}

func (t *thresholdsEvaluator) setConfig(source map[string][]string) {
	t.source = source
	t.config = newThresholdsConfig(source)
}

func (t *thresholdsEvaluator) update(aggs Aggregates) *Thresholds {
	failures := make([]*Failure, 0)

	for metric, agg := range aggs {
		exprs, ok := t.config[metric]
		if !ok || len(exprs) == 0 {
			continue
		}

		srcs, ok := t.source[metric]
		if !ok || len(srcs) == 0 {
			continue
		}

		vars := make(govaluate.MapParameters, len(agg))
		for key, value := range agg {
			vars[fixParentheses(key)] = value
		}

		results, ok := t.results[metric]
		if !ok {
			results = make(map[string]Level)
			t.results[metric] = results
		}

		var failed []string

		for _, src := range srcs {
			lvl := None
			exp := exprs[src]

			res, err := exp.Eval(vars)
			if err == nil {
				if val, ok := res.(bool); ok && val {
					lvl = Ready
				} else {
					lvl = Error
					failed = append(failed, src)
				}
			}

			results[src] = lvl

			if len(failed) > 0 {
				failures = append(failures, &Failure{Metric: metric, Thresholds: failed})
			}
		}
	}

	return newThresholds(t.source, failures, t.results)
}

type thresholdsConfig map[string]map[string]*govaluate.EvaluableExpression

func newThresholdsConfig(source map[string][]string) thresholdsConfig {
	tc := make(thresholdsConfig)

	fail, _ := govaluate.NewEvaluableExpression("false")

	for metric, list := range source {
		emap := make(map[string]*govaluate.EvaluableExpression)

		for _, src := range list {
			expr, err := govaluate.NewEvaluableExpression(fixParentheses(src))
			if err != nil {
				expr = fail
			}

			emap[src] = expr
		}

		if len(emap) > 0 {
			tc[metric] = emap
		}
	}

	return tc
}

var fixParentheses = strings.NewReplacer( //nolint:gochecknoglobals
	"p(90)", "p90",
	"p(95)", "p95",
	"p(99)", "p99",
).Replace
