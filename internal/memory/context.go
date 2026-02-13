package memory

import (
	"math"
	"sort"
)

// ContextSelector picks the most relevant memories for inclusion in
// the consciousness prompt, subject to a token budget.
type ContextSelector struct {
	maxMemories int // maximum number of memories to return
}

// NewContextSelector creates a selector that returns at most maxMemories.
func NewContextSelector(maxMemories int) *ContextSelector {
	return &ContextSelector{maxMemories: maxMemories}
}

// Select returns the most relevant episodic memories for the current state.
// It scores memories by a combination of importance, recency, and somatic
// similarity to the current biological snapshot.
func (cs *ContextSelector) Select(memories []EpisodicMemory, current BioSnapshot) []EpisodicMemory {
	if len(memories) <= cs.maxMemories {
		return memories
	}

	type scored struct {
		memory EpisodicMemory
		score  float64
	}

	var scored_ []scored
	for _, m := range memories {
		s := cs.score(m, current)
		scored_ = append(scored_, scored{memory: m, score: s})
	}

	sort.Slice(scored_, func(i, j int) bool {
		return scored_[i].score > scored_[j].score
	})

	result := make([]EpisodicMemory, cs.maxMemories)
	for i := 0; i < cs.maxMemories; i++ {
		result[i] = scored_[i].memory
	}
	return result
}

// score computes the relevance score for a memory.
// Higher scores mean more relevant for inclusion in the prompt.
func (cs *ContextSelector) score(m EpisodicMemory, current BioSnapshot) float64 {
	// Importance contributes directly.
	importanceScore := m.Importance

	// Somatic similarity: memories formed in a similar body state are more
	// easily recalled (state-dependent memory / Bower's mood-congruent recall).
	somaticScore := somaticSimilarity(m.BioSnapshot, current)

	// Emotional intensity: stronger emotional memories are recalled more easily.
	emotionalScore := math.Abs(m.EmotionalValence)

	return 0.35*importanceScore + 0.35*somaticScore + 0.30*emotionalScore
}

// somaticSimilarity computes how similar two biological snapshots are.
// Returns 0-1 where 1 means identical states.
func somaticSimilarity(a, b BioSnapshot) float64 {
	// Compute Euclidean distance across normalized dimensions, then convert to similarity.
	dims := []struct {
		aVal, bVal, maxRange float64
	}{
		{a.Arousal, b.Arousal, 1.0},
		{a.Valence, b.Valence, 2.0}, // range is -1 to 1
		{a.BodyTemp, b.BodyTemp, 8.0}, // 34°C to 42°C range
		{a.Pain, b.Pain, 1.0},
		{a.Fatigue, b.Fatigue, 1.0},
		{a.Hunger, b.Hunger, 1.0},
	}

	var sumSqDiff float64
	for _, d := range dims {
		normalized := (d.aVal - d.bVal) / d.maxRange
		sumSqDiff += normalized * normalized
	}

	distance := math.Sqrt(sumSqDiff / float64(len(dims)))
	return 1.0 - math.Min(distance, 1.0)
}
