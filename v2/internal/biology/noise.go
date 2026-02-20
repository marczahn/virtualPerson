package biology

import (
	"math"
	"math/rand"
)

// NoiseConfig controls Gaussian noise parameters.
type NoiseConfig struct {
	// Sigma is the base noise standard deviation per second.
	// Noise per tick scales with sqrt(dt) for Brownian consistency.
	Sigma float64
}

func DefaultNoiseConfig() NoiseConfig {
	return NoiseConfig{Sigma: 0.002}
}

// ApplyNoise adds Gaussian noise to all bio variables, scaled by sqrt(dt) for
// Brownian consistency (same total noise variance regardless of tick rate).
// BodyTemp receives smaller noise (0.1x sigma) since it has a narrower functional range.
// Clamp is NOT called here — caller must call ClampAll after noise to absorb boundary violations.
// Noise changes are not individually tracked (too granular for delta logging).
func ApplyNoise(s *State, rng *rand.Rand, cfg NoiseConfig, dt float64) {
	if dt <= 0 {
		return
	}

	sigma := cfg.Sigma * math.Sqrt(dt)
	s.Energy += rng.NormFloat64() * sigma
	s.Stress += rng.NormFloat64() * sigma
	s.CognitiveCapacity += rng.NormFloat64() * sigma
	s.Mood += rng.NormFloat64() * sigma
	s.PhysicalTension += rng.NormFloat64() * sigma
	s.Hunger += rng.NormFloat64() * sigma
	s.SocialDeficit += rng.NormFloat64() * sigma
	// BodyTemp: 0.1x sigma because it's a 18C range, not 0-1.
	s.BodyTemp += rng.NormFloat64() * sigma * 0.1
}
