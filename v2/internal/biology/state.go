package biology

import "time"

// State holds all 8 motivation-shaped biological variables.
// These are proxies calibrated for drive pressure, not physiological measurements.
//
// Drive mapping (BIO-02 contract):
//
//	Energy drive:              Energy (primary), Hunger (secondary)
//	Social connection drive:   SocialDeficit (primary)
//	Stimulation/novelty drive: CognitiveCapacity (primary), Mood (secondary)
//	Safety drive:              Stress (primary), PhysicalTension (secondary), BodyTemp (deviation)
//	Identity coherence drive:  Mood (primary), CognitiveCapacity (secondary)
type State struct {
	Energy            float64 // 0-1: 1=fully rested, 0=exhausted. Decays toward 0.
	Stress            float64 // 0-1: 0=calm, 1=overwhelmed. No autonomous decay.
	CognitiveCapacity float64 // 0-1: 1=fresh, 0=depleted. Decays toward 0.
	Mood              float64 // 0-1: 0=dysphoric, 0.5=neutral, 1=euphoric. Decays toward 0.
	PhysicalTension   float64 // 0-1: 0=relaxed, 1=tense. No autonomous decay.
	Hunger            float64 // 0-1: 0=full/satiated, 1=starving. Decays toward 1.
	SocialDeficit     float64 // 0-1: 0=connected, 1=isolated. Decays toward 1.
	BodyTemp          float64 // Celsius 25-43: baseline 36.6. No autonomous decay.
	UpdatedAt         time.Time
}

// NewDefaultState returns a *State initialised to physiological baselines.
// These values represent a healthy, rested person at the start of a simulation.
func NewDefaultState() *State {
	return &State{
		Energy:            0.80,
		Stress:            0.10,
		CognitiveCapacity: 1.00,
		Mood:              0.50,
		PhysicalTension:   0.05,
		Hunger:            0.10,
		SocialDeficit:     0.00,
		BodyTemp:          36.6,
		UpdatedAt:         time.Now(),
	}
}

// VarRange defines the valid [min, max] range for a bio variable.
type VarRange struct {
	Min, Max float64
}

// Ranges holds the valid range for each bio variable.
// BodyTemp uses {25, 43} — wider than V1's {34,42} to cover
// physiologically meaningful thresholds (33°C hypothermia, 35°C mild).
var Ranges = struct {
	Energy, Stress, CognitiveCapacity, Mood, PhysicalTension,
	Hunger, SocialDeficit, BodyTemp VarRange
}{
	Energy:            VarRange{0, 1},
	Stress:            VarRange{0, 1},
	CognitiveCapacity: VarRange{0, 1},
	Mood:              VarRange{0, 1},
	PhysicalTension:   VarRange{0, 1},
	Hunger:            VarRange{0, 1},
	SocialDeficit:     VarRange{0, 1},
	BodyTemp:          VarRange{25, 43},
}

// Clamp constrains v to [lo, hi].
func Clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ClampAll enforces all variable ranges on s in-place (BIO-06).
func ClampAll(s *State) {
	s.Energy = Clamp(s.Energy, Ranges.Energy.Min, Ranges.Energy.Max)
	s.Stress = Clamp(s.Stress, Ranges.Stress.Min, Ranges.Stress.Max)
	s.CognitiveCapacity = Clamp(s.CognitiveCapacity, Ranges.CognitiveCapacity.Min, Ranges.CognitiveCapacity.Max)
	s.Mood = Clamp(s.Mood, Ranges.Mood.Min, Ranges.Mood.Max)
	s.PhysicalTension = Clamp(s.PhysicalTension, Ranges.PhysicalTension.Min, Ranges.PhysicalTension.Max)
	s.Hunger = Clamp(s.Hunger, Ranges.Hunger.Min, Ranges.Hunger.Max)
	s.SocialDeficit = Clamp(s.SocialDeficit, Ranges.SocialDeficit.Min, Ranges.SocialDeficit.Max)
	s.BodyTemp = Clamp(s.BodyTemp, Ranges.BodyTemp.Min, Ranges.BodyTemp.Max)
}
