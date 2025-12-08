package digest

// State defines the SSE stream/connection states.
type State int //nolint:recvcheck

const (
	// StateWaiting means waiting for data.
	StateWaiting State = iota
	// StateConnected means SSE stream connected.
	StateConnected
	// StatePreparing means parems event received but test not stasrted yet.
	StatePreparing
	// StateStarting means test started but no data available yet.
	StateStarting
	// StateDetached means SSE stream disconnected.
	StateDetached
	// StateRunning means test is running.
	StateRunning
	// StateFinished means test execution is finished.
	StateFinished
)

//go:generate go run github.com/dmarkham/enumer@latest -text -json -trimprefix State -type State
