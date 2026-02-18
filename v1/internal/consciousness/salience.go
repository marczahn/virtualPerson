package consciousness

import (
	"math"

	"github.com/marczahn/person/internal/psychology"
)

// SalienceCalculator determines when biological/psychological changes
// are significant enough to break into conscious awareness.
type SalienceCalculator struct {
	// Previous state snapshot for computing rate of change.
	prevArousal float64
	prevValence float64
	prevEnergy  float64
	prevPain    float64
	initialized bool // false until first Compute establishes baseline

	// Tracks how recently each dimension was in conscious awareness.
	// Higher = longer since last awareness (makes it more novel).
	noveltyClock map[string]float64

	// Current attention direction: positive = inward (body-focused),
	// negative = outward (engaged with external). Range [-1, 1].
	AttentionDirection float64
}

// NewSalienceCalculator creates a calculator with default state.
func NewSalienceCalculator() *SalienceCalculator {
	return &SalienceCalculator{
		noveltyClock: map[string]float64{
			"arousal": 0,
			"valence": 0,
			"energy":  0,
			"pain":    0,
		},
	}
}

// SalienceResult holds the computed salience for a psychological state.
type SalienceResult struct {
	Score     float64 // overall salience score
	Trigger   string  // which dimension triggered the highest salience
	Threshold float64 // the dynamic threshold that was used
	Exceeded  bool    // whether salience exceeded the threshold
}

// Compute calculates the salience of the current psychological state change.
// dt is the elapsed time in seconds.
// The first call establishes baseline and never triggers salience.
func (sc *SalienceCalculator) Compute(ps *psychology.State, dt float64) SalienceResult {
	if dt <= 0 {
		return SalienceResult{}
	}

	// First call establishes baseline without triggering.
	if !sc.initialized {
		sc.prevArousal = ps.Arousal
		sc.prevValence = ps.Valence
		sc.prevEnergy = ps.Energy
		sc.prevPain = painFromValence(ps)
		sc.initialized = true
		return SalienceResult{}
	}

	// Advance novelty clocks for all dimensions.
	for k := range sc.noveltyClock {
		sc.noveltyClock[k] += dt
	}

	// Compute per-dimension salience.
	dimensions := []struct {
		name     string
		current  float64
		previous float64
	}{
		{"arousal", ps.Arousal, sc.prevArousal},
		{"valence", ps.Valence, sc.prevValence},
		{"energy", ps.Energy, sc.prevEnergy},
		{"pain", painFromValence(ps), sc.prevPain},
	}

	var maxSalience float64
	var maxTrigger string

	for _, dim := range dimensions {
		rateOfChange := math.Abs(dim.current-dim.previous) / dt
		novelty := noveltyWeight(sc.noveltyClock[dim.name])
		attention := attentionModifier(sc.AttentionDirection)
		breach := thresholdBreachBonus(dim.name, dim.current)

		salience := rateOfChange*novelty*attention + breach

		if salience > maxSalience {
			maxSalience = salience
			maxTrigger = dim.name
		}
	}

	// Dynamic threshold.
	threshold := dynamicThreshold(ps)

	// Update previous state for next cycle.
	sc.prevArousal = ps.Arousal
	sc.prevValence = ps.Valence
	sc.prevEnergy = ps.Energy
	sc.prevPain = painFromValence(ps)

	// Reset novelty clock for the triggering dimension if it exceeded threshold.
	exceeded := maxSalience > threshold
	if exceeded && maxTrigger != "" {
		sc.noveltyClock[maxTrigger] = 0
	}

	return SalienceResult{
		Score:     maxSalience,
		Trigger:   maxTrigger,
		Threshold: threshold,
		Exceeded:  exceeded,
	}
}

// ResetNovelty resets the novelty clock for a dimension (called when
// consciousness becomes aware of something).
func (sc *SalienceCalculator) ResetNovelty(dimension string) {
	sc.noveltyClock[dimension] = 0
}

// noveltyWeight returns a weight that increases the longer a dimension
// has been out of conscious awareness. Saturates around 2.0.
func noveltyWeight(secondsSinceAwareness float64) float64 {
	// Logarithmic growth: quickly rises then saturates.
	return 1.0 + math.Log1p(secondsSinceAwareness/60.0)*0.3
}

// attentionModifier amplifies body signals when attention is inward,
// suppresses them when outward.
func attentionModifier(direction float64) float64 {
	// direction: -1 (outward/engaged) to +1 (inward/introspective)
	return 1.0 + direction*0.5
}

// thresholdBreachBonus returns a large bonus when a dimension enters
// an extreme range.
func thresholdBreachBonus(dimension string, value float64) float64 {
	switch dimension {
	case "arousal":
		if value > 0.8 {
			return 0.5
		}
	case "valence":
		if value < -0.6 {
			return 0.4
		}
	case "energy":
		if value < 0.15 {
			return 0.3
		}
	case "pain":
		if value > 0.6 {
			return 0.6
		}
	}
	return 0
}

// dynamicThreshold adjusts the salience threshold based on psychological state.
// Lower when idle/anxious (more body-aware), higher when engaged/absorbed.
func dynamicThreshold(ps *psychology.State) float64 {
	base := 0.3

	// Anxiety (high arousal + negative valence) lowers threshold.
	if ps.Arousal > 0.5 && ps.Valence < -0.2 {
		base -= 0.1
	}

	// Low energy / boredom lowers threshold (more body-aware).
	if ps.Energy < 0.3 {
		base -= 0.05
	}

	// High cognitive load (absorption) raises threshold.
	if ps.CognitiveLoad > 0.6 {
		base += 0.1
	}

	if base < 0.1 {
		base = 0.1
	}
	if base > 0.6 {
		base = 0.6
	}

	return base
}

// painFromValence extracts a pain-like signal from the psychological state.
// When valence is very negative and arousal is high, it maps to pain awareness.
func painFromValence(ps *psychology.State) float64 {
	if ps.Valence >= 0 {
		return 0
	}
	return math.Min(1.0, -ps.Valence*ps.Arousal)
}
