package digest

// Level defines log and alert levels.
type Level int

const (
	None    Level = iota // None means regular.
	Info                 // Info means informational.
	Ready                // Ready means everything is ok.
	Notice               // Notice means something happened.
	Warning              // Warning means possible problem.
	Error                // Error means error happened.
)

//go:generate go run github.com/dmarkham/enumer@latest -text -type Level
