package psychology

import "math/rand"

// distortionSpec defines the activation parameters for a cognitive distortion.
type distortionSpec struct {
	distortion     Distortion
	baseRate       float64
	traitMultipler func(p Personality) float64
}

var distortionSpecs = []distortionSpec{
	{
		distortion: Catastrophizing,
		baseRate:   0.05,
		traitMultipler: func(p Personality) float64 {
			return 1.0 + p.Neuroticism*2.0 + (1-p.Openness)*0.5
		},
	},
	{
		distortion: AllOrNothing,
		baseRate:   0.04,
		traitMultipler: func(p Personality) float64 {
			return 1.0 + p.Neuroticism*1.5 + (1-p.Openness)*1.0
		},
	},
	{
		distortion: Personalization,
		baseRate:   0.03,
		traitMultipler: func(p Personality) float64 {
			return 1.0 + p.Neuroticism*1.5 + p.Agreeableness*0.8
		},
	},
	{
		distortion: EmotionalReasoning,
		baseRate:   0.06,
		traitMultipler: func(p Personality) float64 {
			return 1.0 + p.Neuroticism*1.0 + (1-p.Conscientiousness)*0.8
		},
	},
	{
		distortion: Overgeneralization,
		baseRate:   0.04,
		traitMultipler: func(p Personality) float64 {
			return 1.0 + p.Neuroticism*1.5
		},
	},
	{
		distortion: MindReading,
		baseRate:   0.03,
		traitMultipler: func(p Personality) float64 {
			return 1.0 + p.Neuroticism*1.5 + p.Extraversion*0.5
		},
	},
}

// stressMultiplier returns the distortion probability multiplier for a given stress level.
// At stress < 0.3, distortions are minimal. They accelerate sharply above 0.7.
func stressMultiplier(stress float64) float64 {
	switch {
	case stress <= 0.3:
		return 1.0
	case stress <= 0.5:
		return 1.0 + (stress-0.3)*2.5
	case stress <= 0.7:
		return 1.5 + (stress-0.5)*3.0
	default:
		return 2.1 + (stress-0.7)*4.0
	}
}

// ActivateDistortions returns the cognitive distortions that are active
// given the current stress level, regulation capacity, and personality.
// Regulation reduces distortion probability.
func ActivateDistortions(stress, regulationCapacity float64, p Personality) []Distortion {
	sm := stressMultiplier(stress)
	regReduction := 1.0 - regulationCapacity*0.6

	var active []Distortion
	for _, spec := range distortionSpecs {
		prob := spec.baseRate * sm * spec.traitMultipler(p) * regReduction
		if prob > 1.0 {
			prob = 1.0
		}
		if rand.Float64() < prob {
			active = append(active, spec.distortion)
		}
	}
	return active
}

// ActivateDistortionsDeterministic is like ActivateDistortions but uses
// a provided random value source for testability. Each call to randFn
// should return a value in [0,1).
func ActivateDistortionsDeterministic(stress, regulationCapacity float64, p Personality, randFn func() float64) []Distortion {
	sm := stressMultiplier(stress)
	regReduction := 1.0 - regulationCapacity*0.6

	var active []Distortion
	for _, spec := range distortionSpecs {
		prob := spec.baseRate * sm * spec.traitMultipler(p) * regReduction
		if prob > 1.0 {
			prob = 1.0
		}
		if randFn() < prob {
			active = append(active, spec.distortion)
		}
	}
	return active
}
