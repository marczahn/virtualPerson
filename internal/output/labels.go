package output

import (
	"time"

	"github.com/marczahn/person/internal/i18n"
)

// Source identifies which layer produced an output entry.
type Source int

const (
	Sense  Source = iota // sensory event parsing
	Bio                  // biological state changes
	Psych                // psychological state changes
	Mind                 // consciousness thoughts/emotions
	Review               // psychologist reviewer notes
)

func (s Source) String() string {
	tr := i18n.T()
	switch s {
	case Sense:
		return tr.Output.SourceLabels.Sense
	case Bio:
		return tr.Output.SourceLabels.Bio
	case Psych:
		return tr.Output.SourceLabels.Psych
	case Mind:
		return tr.Output.SourceLabels.Mind
	case Review:
		return tr.Output.SourceLabels.Review
	default:
		return tr.Output.Unknown
	}
}

// Entry is a single line of output from any layer of the simulation.
type Entry struct {
	Source    Source
	Message  string
	Timestamp time.Time

	// ThoughtType and Trigger are set only for Mind entries that represent
	// consciousness thoughts. They carry structured data so listeners can
	// avoid parsing the formatted Message string.
	ThoughtType string
	Trigger     string
}
