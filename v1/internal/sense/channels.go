package sense

import "time"

// Channel represents a sensory modality through which input is perceived.
type Channel int

const (
	Visual Channel = iota
	Auditory
	Tactile
	Thermal
	Pain
	Olfactory
	Gustatory
	Vestibular  // balance, spatial orientation
	Interoceptive // internal body signals (nausea, breathlessness, etc.)
)

var channelNames = [...]string{
	"visual",
	"auditory",
	"tactile",
	"thermal",
	"pain",
	"olfactory",
	"gustatory",
	"vestibular",
	"interoceptive",
}

func (c Channel) String() string {
	if int(c) < len(channelNames) {
		return channelNames[c]
	}
	return "unknown"
}

// Event represents a parsed sensory input with its characteristics.
type Event struct {
	Channel   Channel
	Intensity float64   // 0-1, how strong the stimulus is
	RawInput  string    // the original text that triggered this event
	Parsed    string    // LLM-parsed description of what the person perceives
	Timestamp time.Time
}
