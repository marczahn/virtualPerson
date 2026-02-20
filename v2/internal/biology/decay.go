package biology

// Decay rates: autonomous degradation per second at DecayMultiplier=1.0.
// Only 5 variables degrade autonomously (BIO-03).
// Calibrated for ~20% change from baseline within 4 minutes at 1x speed.
const (
	energyDecayRate = 0.00067 // Energy: 0.8→0.6 in ~300s at 1x
	hungerDecayRate = 0.00083 // Hunger: 0.1→0.3 in ~240s at 1x
	cogCapDecayRate = 0.00050 // CognitiveCapacity: 1.0→0.8 in ~400s at 1x
	moodDecayRate   = 0.00033 // Mood: 0.5→0.3 in ~600s at 1x (slower drift)
	socialDecayRate = 0.00033 // SocialDeficit: 0.0→0.2 in ~600s at 1x (slow isolation)
)

// DecayConfig holds the multiplier for autonomous decay rates.
// DecayMultiplier=1.0 is normal speed; 5.0 is fast development mode (5x faster degradation).
// HomeostasisEnabled=false is the V2 default — no variable auto-returns to baseline.
type DecayConfig struct {
	DecayMultiplier    float64
	HomeostasisEnabled bool // reserved for future use; always false in V2
}

// DefaultDecayConfig returns the development-friendly default:
// 5x speed so degradation is visible within ~1 minute.
func DefaultDecayConfig() DecayConfig {
	return DecayConfig{
		DecayMultiplier:    5.0,
		HomeostasisEnabled: false,
	}
}

// ApplyDecay applies autonomous linear decay to s in-place for elapsed dt seconds.
// Only 5 variables decay autonomously (BIO-03):
//   - Energy drifts toward 0 (exhaustion)
//   - Hunger rises toward 1 (starvation)
//   - CognitiveCapacity drifts toward 0 (mental depletion)
//   - Mood drifts toward 0 (dysphoria)
//   - SocialDeficit rises toward 1 (isolation)
//
// Stress, PhysicalTension, and BodyTemp are NOT touched — they only change
// from explicit causes (interactions, thresholds, external feedback).
//
// Clamp is NOT called here — caller must call ClampAll after all mutations.
func ApplyDecay(s *State, cfg DecayConfig, dt float64) {
	// Cap dt at 60s to prevent pause-recovery explosions.
	if dt > 60.0 {
		dt = 60.0
	}
	rate := cfg.DecayMultiplier * dt
	s.Energy -= energyDecayRate * rate
	s.Hunger += hungerDecayRate * rate
	s.CognitiveCapacity -= cogCapDecayRate * rate
	s.Mood -= moodDecayRate * rate
	s.SocialDeficit += socialDecayRate * rate
	// Stress, PhysicalTension, BodyTemp: no autonomous decay
}
