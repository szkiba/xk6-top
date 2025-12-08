package digest

// EventType defines the type of the SSE event.
type EventType int //nolint:recvcheck

const (
	// EventTypeConfig mean "config" SSE event.
	EventTypeConfig EventType = iota
	// EventTypeParam mean "param" SSE event.
	EventTypeParam
	// EventTypeMetric mean "metric" SSE event.
	EventTypeMetric
	// EventTypeSnapshot mean "snapshot" SSE event.
	EventTypeSnapshot
	// EventTypeCumulative mean "cumulative" SSE event.
	EventTypeCumulative
	// EventTypeStart mean "start" SSE event.
	EventTypeStart
	// EventTypeStop mean "stop" SSE event.
	EventTypeStop
	// EventTypeThreshold mean "threshold" SSE event.
	EventTypeThreshold
	// EventTypeConnect mean SSE channel connected.
	EventTypeConnect
	// EventTypeDisconnect mean SSE channel disconnected.
	EventTypeDisconnect
)

//go:generate go run github.com/dmarkham/enumer@latest -text -json -transform lower -trimprefix EventType -type EventType

// Event describes an SSE event.
type Event struct {
	Type EventType `json:"event,omitempty"`
	Data any       `json:"data,omitempty"`
}

// ConfigData holds "config" event data.
type ConfigData map[string]any

// ParamData holds "param" event data.
type ParamData struct {
	Thresholds map[string][]string `json:"thresholds,omitempty"`
	Scenarios  []string            `json:"scenarios,omitempty"`
	EndOffset  int64               `json:"endOffset,omitempty"`
	Period     int64               `json:"period,omitempty"`
	Tags       []string            `json:"tags,omitempty"`
}
