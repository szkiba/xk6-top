package digest

// A MetricType specifies the type of a metric.
type MetricType int

//go:generate go run github.com/dmarkham/enumer@latest -text -json -transform lower -trimprefix MetricType -type MetricType

// Possible values for MetricType.
const (
	MetricTypeCounter MetricType = iota // A counter that sums its data points
	MetricTypeGauge                     // A gauge that displays the latest value
	MetricTypeTrend                     // A trend, min/max/avg/med are interesting
	MetricTypeRate                      // A rate, displays % of values that aren't 0
)

//nolint:gochecknoglobals
var (
	trendKeys   = []string{"avg", "max", "med", "min", "p(90)", "p(95)", "p(99)"}
	rateKeys    = []string{"rate", "peak"}
	counterKeys = []string{"count", "rate", "peak"}
	gaugeKeys   = []string{"value"}
)

// Aggregates returns aggregate names for a given metric type.
func (md MetricType) Aggregates() []string {
	switch md {
	case MetricTypeTrend:
		return trendKeys
	case MetricTypeCounter:
		return counterKeys
	case MetricTypeRate:
		return rateKeys
	case MetricTypeGauge:
		return gaugeKeys
	default:
		panic("unreachable")
	}
}

// ValueType holds the type of values a metric contains.
type ValueType int

//go:generate go run github.com/dmarkham/enumer@latest -text -json -transform lower -trimprefix ValueType -type ValueType

// Possible values for ValueType.
const (
	ValueTypeDefault ValueType = iota // Values are presented as-is
	ValueTypeTime                     // Values are time durations (milliseconds)
	ValueTypeData                     // Values are data amounts (bytes)
)

// Format formats value for a given value type.
func (vt ValueType) Format(value float64) string {
	switch vt {
	case ValueTypeTime:
		return formatDuration(value)
	case ValueTypeData:
		return formatData(value)
	case ValueTypeDefault:
		fallthrough
	default:
		return formatDefault(value)
	}
}

// Metric holds metric metadata.
type Metric struct {
	Name     string     `json:"name,omitempty"`
	Type     MetricType `json:"type,omitempty"`
	Contains ValueType  `json:"contains,omitempty"`
	Tainted  bool       `json:"tainted,omitempty"`
}
