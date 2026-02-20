package biology

// Severity represents the urgency level of a threshold crossing.
type Severity int

const (
	Mild     Severity = iota // noticeable but manageable
	Warning                  // impairs function, attention required
	Critical                 // dangerous, demands immediate response
)

func (s Severity) String() string {
	switch s {
	case Mild:
		return "mild"
	case Warning:
		return "warning"
	case Critical:
		return "critical"
	default:
		return "unknown"
	}
}

// ThresholdEvent describes a bio variable that has crossed a threshold tier.
// Cascade contains the bio effects triggered by this crossing.
// Cascade effects must be applied via ApplyThresholdCascades after ClampAll.
type ThresholdEvent struct {
	Variable    string // bio variable name (e.g., "stress", "body_temp")
	Severity    Severity
	Description string  // human-readable condition description
	Cascade     []Delta // bio effects triggered by this threshold crossing
}

// ThresholdConfig controls threshold behavior.
// TerminalStatesEnabled=false (development default) does not suppress detection —
// it only prevents cascade effects from driving variables toward lethal extremes
// by capping cascade magnitudes. Reserved for future use.
type ThresholdConfig struct {
	TerminalStatesEnabled bool // false = development default
}

// DefaultThresholdConfig returns the development-default threshold configuration.
func DefaultThresholdConfig() ThresholdConfig {
	return ThresholdConfig{TerminalStatesEnabled: false}
}
