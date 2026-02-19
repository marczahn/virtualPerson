package bio

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
// Cascade effects must be applied via ApplyThresholdCascades after clampAll.
type ThresholdEvent struct {
	Variable    string   // bio variable name (e.g., "stress", "body_temp")
	Severity    Severity
	Description string   // human-readable condition description
	Cascade     []Delta  // bio effects triggered by this threshold crossing
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

// EvaluateThresholds evaluates state against all threshold conditions and returns events.
// Called AFTER clampAll (step 5 of tick pipeline). Cascade deltas are returned inside
// the events, not applied here — the engine applies them and re-clamps.
//
// For stepped tiers (e.g., BodyTemp <35, <34, <33), only the MOST SEVERE applicable
// event is emitted per variable group.
func EvaluateThresholds(s *State, cfg ThresholdConfig, dt float64) []ThresholdEvent {
	var events []ThresholdEvent

	// --- Body Temperature: hypothermia (stepped, most-severe-only) ---
	switch {
	case s.BodyTemp < 33.0:
		events = append(events, ThresholdEvent{
			Variable:    "body_temp",
			Severity:    Critical,
			Description: "Severe hypothermia, crisis",
			Cascade: []Delta{
				{Field: "stress", Amount: 0.3},
				{Field: "cognitive_capacity", Amount: -0.4},
			},
		})
	case s.BodyTemp < 34.0:
		events = append(events, ThresholdEvent{
			Variable:    "body_temp",
			Severity:    Warning,
			Description: "Moderate hypothermia, mental slowing",
			Cascade: []Delta{
				{Field: "physical_tension", Amount: 0.3},
				{Field: "cognitive_capacity", Amount: -0.2},
			},
		})
	case s.BodyTemp < 35.0:
		events = append(events, ThresholdEvent{
			Variable:    "body_temp",
			Severity:    Mild,
			Description: "Mild hypothermia, shivering",
			Cascade: []Delta{
				{Field: "physical_tension", Amount: 0.2},
			},
		})
	}

	// --- Body Temperature: hyperthermia (stepped, most-severe-only) ---
	switch {
	case s.BodyTemp > 40.5:
		events = append(events, ThresholdEvent{
			Variable:    "body_temp",
			Severity:    Critical,
			Description: "Dangerous hyperthermia",
			Cascade: []Delta{
				{Field: "stress", Amount: 0.4},
				{Field: "cognitive_capacity", Amount: -0.3},
			},
		})
	case s.BodyTemp > 39.5:
		events = append(events, ThresholdEvent{
			Variable:    "body_temp",
			Severity:    Warning,
			Description: "Fever, significant impairment",
			Cascade: []Delta{
				{Field: "stress", Amount: 0.2},
				{Field: "mood", Amount: -0.2},
			},
		})
	case s.BodyTemp > 38.5:
		events = append(events, ThresholdEvent{
			Variable:    "body_temp",
			Severity:    Mild,
			Description: "Elevated temperature, discomfort",
			Cascade: []Delta{
				{Field: "stress", Amount: 0.1},
				{Field: "cognitive_capacity", Amount: -0.1},
			},
		})
	}

	// --- Stress thresholds (dt-scaled cascades, most-severe-only) ---
	switch {
	case s.Stress > 0.95:
		events = append(events, ThresholdEvent{
			Variable:    "stress",
			Severity:    Critical,
			Description: "Crisis state",
			Cascade: []Delta{
				{Field: "mood", Amount: -0.03 * dt},
				{Field: "energy", Amount: -0.02 * dt},
			},
		})
	case s.Stress > 0.85:
		events = append(events, ThresholdEvent{
			Variable:    "stress",
			Severity:    Warning,
			Description: "High stress, impaired function",
			Cascade: []Delta{
				{Field: "cognitive_capacity", Amount: -0.02 * dt},
				{Field: "mood", Amount: -0.01 * dt},
			},
		})
	case s.Stress > 0.7:
		events = append(events, ThresholdEvent{
			Variable:    "stress",
			Severity:    Mild,
			Description: "Elevated stress",
			Cascade: []Delta{
				{Field: "physical_tension", Amount: 0.01 * dt},
			},
		})
	}

	// --- Energy thresholds (dt-scaled, most-severe-only) ---
	switch {
	case s.Energy < 0.05:
		events = append(events, ThresholdEvent{
			Variable:    "energy",
			Severity:    Critical,
			Description: "Near physical collapse",
			Cascade: []Delta{
				{Field: "stress", Amount: 0.03 * dt},
				{Field: "cognitive_capacity", Amount: -0.03 * dt},
			},
		})
	case s.Energy < 0.15:
		events = append(events, ThresholdEvent{
			Variable:    "energy",
			Severity:    Warning,
			Description: "Very low energy",
			Cascade: []Delta{
				{Field: "mood", Amount: -0.01 * dt},
				{Field: "stress", Amount: 0.01 * dt},
			},
		})
	case s.Energy < 0.3:
		events = append(events, ThresholdEvent{
			Variable:    "energy",
			Severity:    Mild,
			Description: "Low energy, effort costs more",
			Cascade: []Delta{
				{Field: "cognitive_capacity", Amount: -0.01 * dt},
			},
		})
	}

	// --- Hunger thresholds (dt-scaled, most-severe-only) ---
	switch {
	case s.Hunger > 0.95:
		events = append(events, ThresholdEvent{
			Variable:    "hunger",
			Severity:    Critical,
			Description: "Starving",
			Cascade: []Delta{
				{Field: "stress", Amount: 0.02 * dt},
				{Field: "energy", Amount: -0.01 * dt},
			},
		})
	case s.Hunger > 0.85:
		events = append(events, ThresholdEvent{
			Variable:    "hunger",
			Severity:    Warning,
			Description: "Very hungry, difficulty focusing",
			Cascade: []Delta{
				{Field: "cognitive_capacity", Amount: -0.01 * dt},
				{Field: "stress", Amount: 0.005 * dt},
			},
		})
	case s.Hunger > 0.7:
		events = append(events, ThresholdEvent{
			Variable:    "hunger",
			Severity:    Mild,
			Description: "Noticeably hungry",
			Cascade: []Delta{
				{Field: "mood", Amount: -0.005 * dt},
			},
		})
	}

	return events
}

// ApplyThresholdCascades applies the cascade Deltas from all threshold events to s.
// After calling this, clampAll must be called to keep values in range.
func ApplyThresholdCascades(s *State, events []ThresholdEvent) {
	for _, e := range events {
		for _, d := range e.Cascade {
			applyDelta(s, d)
		}
	}
}
