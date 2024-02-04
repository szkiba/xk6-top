package digest

// State defines the SSE stream/connection states.
type State int

const (
	StateWaiting   State = iota // StateWaiting means waiting for data.
	StateConnected              // StateConnected means SSE stream connected.
	StatePreparing              // StatePreparing means parems event received but test not stasrted yet.
	StateStarting               // StateStarting means test started but no data available yet.
	StateDetached               // StateDetached means SSE stream disconnected.
	StateRunning                // StateRunning means test is running.
	StateFinished               // StateFinished means test execution is finished.
)

//go:generate go run github.com/dmarkham/enumer@latest -text -json -trimprefix State -type State
