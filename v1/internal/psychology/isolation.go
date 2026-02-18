package psychology

import "time"

// UpdateIsolation advances the isolation state by dt seconds.
// The person's personality modulates how quickly isolation becomes distressing.
func UpdateIsolation(current IsolationState, p Personality, dt float64) IsolationState {
	current.Duration += time.Duration(dt * float64(time.Second))

	hours := current.Duration.Hours()
	extraFactor := IsolationDistressRate(p)

	current.LonelinessLevel = computeLoneliness(hours, extraFactor)
	current.Phase = computeIsolationPhase(hours)

	return current
}

// computeLoneliness derives loneliness level from isolation duration.
// Uses the timeline from the psychologist advisory.
func computeLoneliness(hours, extraversionFactor float64) float64 {
	var loneliness float64

	switch {
	case hours <= 2:
		loneliness = 0.05 * extraversionFactor
	case hours <= 8:
		loneliness = (0.1 + (hours-2)/6*0.15) * extraversionFactor
	case hours <= 24:
		loneliness = (0.25 + (hours-8)/16*0.25) * extraversionFactor
	case hours <= 72:
		loneliness = (0.5 + (hours-24)/48*0.2) * extraversionFactor
	case hours <= 168:
		loneliness = 0.7 + (hours-72)/96*0.15
	default:
		loneliness = 0.85 + clamp((hours-168)/336*0.1, 0, 0.15)
	}

	return clamp(loneliness, 0, 1)
}

// computeIsolationPhase maps hours to the isolation phase.
func computeIsolationPhase(hours float64) IsolationPhase {
	switch {
	case hours <= 2:
		return IsolationNone
	case hours <= 8:
		return IsolationBoredom
	case hours <= 24:
		return IsolationLoneliness
	case hours <= 72:
		return IsolationSignificant
	case hours <= 168:
		return IsolationDestabilizing
	default:
		return IsolationSevere
	}
}
