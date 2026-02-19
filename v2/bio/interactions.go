package bio

// Delta represents a change to a specific bio state field.
type Delta struct {
	Field  string  // field name for logging
	Amount float64 // the change amount (positive or negative)
}

// Rule is a data-driven interaction rule between bio variables.
type Rule struct {
	Name      string
	Condition func(s *State) bool
	Apply     func(s *State, dt float64) Delta
}

// motivationRules is the complete set of 22 interaction rules between bio variables.
// Rules are evaluated against a pre-tick snapshot to prevent single-tick feedback explosions.
var motivationRules = []Rule{
	// --- Stress interactions ---
	{
		Name:      "stress->physical_tension: high stress tightens muscles",
		Condition: func(s *State) bool { return s.Stress > 0.6 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"physical_tension", s.Stress * 0.3 * dt} },
	},
	{
		Name:      "stress->cognitive_capacity: stress depletes mental capacity",
		Condition: func(s *State) bool { return s.Stress > 0.5 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"cognitive_capacity", -s.Stress * 0.2 * dt} },
	},
	{
		Name:      "stress->mood: severe stress dampens mood",
		Condition: func(s *State) bool { return s.Stress > 0.7 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.002 * dt} },
	},

	// --- Hunger interactions ---
	{
		Name:      "hunger->stress: severe hunger raises stress",
		Condition: func(s *State) bool { return s.Hunger > 0.7 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", 0.001 * dt} },
	},
	{
		Name:      "hunger->cognitive_capacity: extreme hunger depletes cognition",
		Condition: func(s *State) bool { return s.Hunger > 0.8 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"cognitive_capacity", -0.002 * dt} },
	},

	// --- Energy interactions ---
	{
		Name:      "energy->mood: low energy worsens mood",
		Condition: func(s *State) bool { return s.Energy < 0.3 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.001 * dt} },
	},
	{
		Name:      "energy->stress: very low energy raises stress",
		Condition: func(s *State) bool { return s.Energy < 0.2 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", 0.002 * dt} },
	},
	{
		Name:      "energy->cognitive_capacity: very low energy depletes cognition",
		Condition: func(s *State) bool { return s.Energy < 0.2 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"cognitive_capacity", -0.002 * dt} },
	},

	// --- Physical tension interactions ---
	{
		Name:      "physical_tension->stress: high tension feeds back to stress",
		Condition: func(s *State) bool { return s.PhysicalTension > 0.7 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", 0.001 * dt} },
	},
	{
		Name:      "physical_tension->mood: elevated tension dampens mood",
		Condition: func(s *State) bool { return s.PhysicalTension > 0.6 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.001 * dt} },
	},

	// --- Cognitive capacity interactions (inverted: low capacity = high load) ---
	{
		Name:      "cognitive_capacity->stress: severe depletion raises stress",
		Condition: func(s *State) bool { return s.CognitiveCapacity < 0.2 }, // equiv. CognitiveLoad > 0.8
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", 0.002 * dt} },
	},
	{
		Name:      "cognitive_capacity->mood: depleted cognition lowers mood",
		Condition: func(s *State) bool { return s.CognitiveCapacity < 0.3 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.001 * dt} },
	},

	// --- Mood interactions ---
	{
		Name:      "mood->stress: dysphoria elevates stress",
		Condition: func(s *State) bool { return s.Mood < 0.2 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", 0.001 * dt} },
	},
	{
		Name:      "mood->social_deficit: dysphoria deepens isolation",
		Condition: func(s *State) bool { return s.Mood < 0.2 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"social_deficit", 0.001 * dt} },
	},

	// --- Social deficit interactions ---
	{
		Name:      "social_deficit->mood: high isolation lowers mood",
		Condition: func(s *State) bool { return s.SocialDeficit > 0.7 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.001 * dt} },
	},
	{
		Name:      "social_deficit->stress: extreme isolation raises stress",
		Condition: func(s *State) bool { return s.SocialDeficit > 0.8 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", 0.001 * dt} },
	},

	// --- Body temperature interactions (hypothermia) ---
	{
		Name:      "body_temp->stress: hypothermia raises stress",
		Condition: func(s *State) bool { return s.BodyTemp < 35.5 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", (35.5 - s.BodyTemp) * 0.01 * dt} },
	},
	{
		Name:      "body_temp->physical_tension: hypothermia causes muscle tension (shivering)",
		Condition: func(s *State) bool { return s.BodyTemp < 35.5 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"physical_tension", (35.5 - s.BodyTemp) * 0.05 * dt} },
	},

	// --- Body temperature interactions (hyperthermia) ---
	{
		Name:      "body_temp->stress: hyperthermia raises stress",
		Condition: func(s *State) bool { return s.BodyTemp > 38.5 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"stress", (s.BodyTemp - 38.5) * 0.01 * dt} },
	},
	{
		Name:      "body_temp->cognitive_capacity: hyperthermia depletes cognition",
		Condition: func(s *State) bool { return s.BodyTemp > 38.5 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"cognitive_capacity", -(s.BodyTemp - 38.5) * 0.03 * dt} },
	},

	// --- Compound spiral rules ---
	{
		Name:      "energy+hunger->mood: low energy AND high hunger collapses mood faster",
		Condition: func(s *State) bool { return s.Energy < 0.4 && s.Hunger > 0.6 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.002 * dt} },
	},
	{
		Name:      "stress+cognitive_capacity->mood: overwhelmed-depleted spiral crushes mood",
		Condition: func(s *State) bool { return s.Stress > 0.8 && s.CognitiveCapacity < 0.3 },
		Apply:     func(s *State, dt float64) Delta { return Delta{"mood", -0.003 * dt} },
	},
}

// ApplyInteractions evaluates all motivation rules against the pre-tick state snapshot
// and applies the combined deltas to s. Clamp is not called here.
// Single-pass evaluation: all conditions read from snapshot, not from post-rule state.
// This prevents single-tick feedback explosions.
func ApplyInteractions(s *State, dt float64) []Delta {
	snap := *s // snapshot pre-tick state for condition evaluation
	var deltas []Delta
	for _, rule := range motivationRules {
		if rule.Condition(&snap) {
			d := rule.Apply(&snap, dt)
			deltas = append(deltas, d)
			// Apply to real state immediately (but conditions already evaluated from snap)
			applyDelta(s, d)
		}
	}
	return deltas
}

// applyDelta applies a single Delta to the State by field name.
func applyDelta(s *State, d Delta) {
	switch d.Field {
	case "energy":
		s.Energy += d.Amount
	case "stress":
		s.Stress += d.Amount
	case "cognitive_capacity":
		s.CognitiveCapacity += d.Amount
	case "mood":
		s.Mood += d.Amount
	case "physical_tension":
		s.PhysicalTension += d.Amount
	case "hunger":
		s.Hunger += d.Amount
	case "social_deficit":
		s.SocialDeficit += d.Amount
	case "body_temp":
		s.BodyTemp += d.Amount
	}
}
