package biology

import (
	"math/rand"
	"time"
)

// Config holds all engine configuration.
type Config struct {
	Decay      DecayConfig
	Noise      NoiseConfig
	Thresholds ThresholdConfig
}

// DefaultConfig returns development-friendly defaults.
func DefaultConfig() Config {
	return Config{
		Decay:      DefaultDecayConfig(),
		Noise:      DefaultNoiseConfig(),
		Thresholds: DefaultThresholdConfig(),
	}
}

// TickResult holds the output of one Engine.Tick call.
type TickResult struct {
	Deltas     []Delta
	Thresholds []ThresholdEvent
}

// Engine runs the bio simulation tick pipeline.
type Engine struct {
	config Config
	rng    *rand.Rand
}

// NewEngine creates an Engine with the given config and a random seed.
func NewEngine(cfg Config) *Engine {
	return &Engine{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewEngineWithSeed creates a deterministic Engine for testing.
func NewEngineWithSeed(cfg Config, seed int64) *Engine {
	return &Engine{
		config: cfg,
		rng:    rand.New(rand.NewSource(seed)),
	}
}

// Tick advances the bio state by dt seconds.
func (e *Engine) Tick(s *State, dt float64) TickResult {
	var result TickResult

	ApplyDecay(s, e.config.Decay, dt)
	result.Deltas = ApplyInteractions(s, dt)
	ApplyNoise(s, e.rng, e.config.Noise, dt)
	ClampAll(s)

	events := EvaluateThresholds(s, e.config.Thresholds, dt)
	result.Thresholds = events
	ApplyThresholdCascades(s, events)
	ClampAll(s)

	s.UpdatedAt = time.Now()
	return result
}
