package digest

// Level defines log and alert levels.
type Level int

const (
	// None means regular.
	None Level = iota
	// Info means informational.
	Info
	// Ready means everything is ok.
	Ready
	// Notice means something happened.
	Notice
	// Warning means possible problem.
	Warning
	// Error means error happened.
	Error
)

//go:generate go run github.com/dmarkham/enumer@latest -text -type Level
