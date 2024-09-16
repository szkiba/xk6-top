package digest

// EventType defines the type of the SSE event.
type EventType int

const (
	EventTypeConfig     EventType = iota // EventTypeConfig mean "config" SSE event.
	EventTypeParam                       // EventTypeParam mean "param" SSE event.
	EventTypeMetric                      // EventTypeMetric mean "metric" SSE event.
	EventTypeSnapshot                    // EventTypeSnapshot mean "snapshot" SSE event.
	EventTypeCumulative                  // EventTypeCumulative mean "cumulative" SSE event.
	EventTypeStart                       // EventTypeStart mean "start" SSE event.
	EventTypeStop                        // EventTypeStop mean "stop" SSE event.
	EventTypeThreshold                   // EventTypeThreshold mean "threshold" SSE event.
	EventTypeConnect                     // EventTypeConnect mean SSE channel connected.
	EventTypeDisconnect                  // EventTypeDisconnect mean SSE channel disconnected.
)

//go:generate go run github.com/dmarkham/enumer@latest -text -json -transform lower -trimprefix EventType -type EventType

// Event describes an SSE event.
type Event struct {
	Type EventType   `json:"event,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// ConfigData holds "config" event data.
type ConfigData map[string]interface{}

// ParamData holds "param" event data.
type ParamData struct {
	Thresholds map[string][]string `json:"thresholds,omitempty"`
	Scenarios  []string            `json:"scenarios,omitempty"`
	EndOffset  int64               `json:"endOffset,omitempty"`
	Period     int64               `json:"period,omitempty"`
	Tags       []string            `json:"tags,omitempty"`
}
