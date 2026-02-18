package psychology

// RegulationState tracks the current depletable emotional regulation resource.
type RegulationState struct {
	Capacity float64 // 0-1, current available regulation capacity
}

// UpdateRegulation computes the new regulation state after dt seconds,
// given the current stress level and fatigue.
//
// Regulation depletes under stress and recovers during low-stress periods.
// The depletion curve is non-linear: people maintain regulation for a while,
// then it collapses relatively quickly.
func UpdateRegulation(current RegulationState, p Personality, stress, fatigue, dt float64) RegulationState {
	baseline := BaselineRegulation(p)
	dtHours := dt / 3600.0

	if stress > 0.3 {
		// Depletion: stress_level * time * rate
		depletion := stress * dtHours * 0.08
		current.Capacity -= depletion

		// Accelerating collapse after 60% depletion of baseline.
		if current.Capacity < baseline*0.4 {
			overshoot := baseline*0.4 - current.Capacity
			current.Capacity -= overshoot * overshoot * 0.5
		}
	} else {
		// Recovery during low-stress waking time: 0.1 per hour.
		recovery := 0.1 * dtHours
		current.Capacity += recovery
	}

	// Fatigue penalty: tired people regulate worse.
	fatiguePenalty := fatigue * 0.3
	effective := current.Capacity - fatiguePenalty

	// Don't let the stored capacity exceed baseline.
	if current.Capacity > baseline {
		current.Capacity = baseline
	}
	if current.Capacity < 0 {
		current.Capacity = 0
	}

	// The effective capacity (used for modulation) includes fatigue,
	// but the stored capacity doesn't drain from fatigue — it's a temporary penalty.
	_ = effective

	return current
}

// EffectiveCapacity returns the regulation capacity after accounting for
// the temporary fatigue penalty. This is what affects arousal/valence dampening.
func EffectiveCapacity(reg RegulationState, fatigue float64) float64 {
	effective := reg.Capacity - fatigue*0.3
	return clamp(effective, 0, 1)
}
