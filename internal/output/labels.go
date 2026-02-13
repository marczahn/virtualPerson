package output

import "time"

// Source identifies which layer produced an output entry.
type Source int

const (
	Sense  Source = iota // sensory event parsing
	Bio                  // biological state changes
	Psych                // psychological state changes
	Mind                 // consciousness thoughts/emotions
	Review               // psychologist reviewer notes
)

var sourceLabels = [...]string{
	"SENSE",
	"BIO",
	"PSYCH",
	"MIND",
	"REVIEW",
}

func (s Source) String() string {
	if int(s) < len(sourceLabels) {
		return sourceLabels[s]
	}
	return "UNKNOWN"
}

// Entry is a single line of output from any layer of the simulation.
type Entry struct {
	Source    Source
	Message  string
	Timestamp time.Time
}
