package psychology

import (
	"math"
	"time"

	"github.com/marczahn/person/internal/biology"
)

// EmotionalMemory represents a stored emotional association.
type EmotionalMemory struct {
	ID        string
	Stimulus  string  // what triggered it (e.g., "cold", "darkness")
	Valence   float64 // -1 to 1
	Intensity float64 // 0 to 1
	CreatedAt time.Time
	Traumatic bool // if true, decays slower and intrudes
}

// EmotionalMemoryStore holds accumulated emotional memories.
type EmotionalMemoryStore struct {
	memories []EmotionalMemory
}

// NewEmotionalMemoryStore creates an empty memory store.
func NewEmotionalMemoryStore() *EmotionalMemoryStore {
	return &EmotionalMemoryStore{}
}

// Add stores a new emotional memory.
func (s *EmotionalMemoryStore) Add(mem EmotionalMemory) {
	s.memories = append(s.memories, mem)
}

// Memories returns all stored memories.
func (s *EmotionalMemoryStore) Memories() []EmotionalMemory {
	return s.memories
}

// QueryActivations returns emotional memories activated by similarity to
// the current biological state. This is a simplified model that checks
// for stimulus-relevant biological conditions.
func (s *EmotionalMemoryStore) QueryActivations(bio *biology.State) []EmotionalMemoryActivation {
	now := time.Now()
	var activations []EmotionalMemoryActivation

	for _, mem := range s.memories {
		daysSince := now.Sub(mem.CreatedAt).Hours() / 24
		if daysSince < 0 {
			daysSince = 0
		}

		decay := recencyDecay(daysSince, mem.Traumatic)
		if decay*mem.Intensity < 0.01 {
			continue // too weak to activate
		}

		similarity := stimulusSimilarity(mem.Stimulus, bio)
		if similarity < 0.1 {
			continue // not relevant to current state
		}

		activations = append(activations, EmotionalMemoryActivation{
			MemoryID:    mem.ID,
			Similarity:  similarity,
			ValenceSign: sign(mem.Valence),
			Intensity:   mem.Intensity * decay * similarity,
		})
	}

	return activations
}

// MemoryModifier computes the emotional intensity modifier from active memories.
// Negative memories are weighted 1.5x more heavily (negativity bias).
// Returns a value in [-0.5, 1.0] to multiply emotional intensity.
func MemoryModifier(activations []EmotionalMemoryActivation, p Personality) float64 {
	var posWeight, negWeight float64

	for _, a := range activations {
		salience := emotionalSalience(a.ValenceSign) * a.Intensity
		if a.ValenceSign < 0 {
			negWeight += salience
		} else {
			posWeight += salience
		}
	}

	neuroticismScaling := 0.7 + p.Neuroticism*0.6
	modifier := (negWeight - posWeight) * neuroticismScaling

	return clamp(modifier, -0.5, 1.0)
}

// recencyDecay follows a power law: 1 / (1 + days^exponent).
// Normal memories use exponent 0.5, traumatic memories use 0.3 (slower decay).
func recencyDecay(daysSince float64, traumatic bool) float64 {
	exp := 0.5
	if traumatic {
		exp = 0.3
	}
	return 1.0 / (1.0 + math.Pow(daysSince, exp))
}

// emotionalSalience weights negative memories more heavily.
// Negativity bias: negative memories are ~1.5x more salient.
func emotionalSalience(valenceSign float64) float64 {
	if valenceSign < 0 {
		return 1.5
	}
	return 1.0
}

// stimulusSimilarity checks if the current biological state matches the
// conditions that originally created the memory. Returns 0-1.
func stimulusSimilarity(stimulus string, bio *biology.State) float64 {
	switch stimulus {
	case "cold":
		if bio.BodyTemp < 35.5 {
			return clamp((35.5-bio.BodyTemp)/3.0, 0, 1)
		}
	case "heat":
		if bio.BodyTemp > 38.0 {
			return clamp((bio.BodyTemp-38.0)/4.0, 0, 1)
		}
	case "pain":
		if bio.Pain > 0.3 {
			return bio.Pain
		}
	case "hunger":
		if bio.Hunger > 0.3 {
			return bio.Hunger
		}
	case "darkness":
		// Activated during circadian night phases.
		if bio.CircadianPhase < 6 || bio.CircadianPhase > 21 {
			return 0.5
		}
	case "exertion":
		if bio.HeartRate > 100 {
			return clamp((bio.HeartRate-100)/100, 0, 1)
		}
	}
	return 0
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}
