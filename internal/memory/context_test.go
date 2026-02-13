package memory

import (
	"math"
	"testing"
	"time"
)

func TestSomaticSimilarity_IdenticalStates(t *testing.T) {
	a := BioSnapshot{Arousal: 0.5, Valence: 0.3, BodyTemp: 36.6, Pain: 0.1, Fatigue: 0.2, Hunger: 0.1}
	got := somaticSimilarity(a, a)

	if got != 1.0 {
		t.Errorf("similarity of identical states = %f, want 1.0", got)
	}
}

func TestSomaticSimilarity_OppositeStates(t *testing.T) {
	a := BioSnapshot{Arousal: 0.0, Valence: -1.0, BodyTemp: 34.0, Pain: 0.0, Fatigue: 0.0, Hunger: 0.0}
	b := BioSnapshot{Arousal: 1.0, Valence: 1.0, BodyTemp: 42.0, Pain: 1.0, Fatigue: 1.0, Hunger: 1.0}

	got := somaticSimilarity(a, b)

	if got > 0.3 {
		t.Errorf("similarity of opposite states = %f, expected < 0.3", got)
	}
}

func TestSomaticSimilarity_SimilarStates(t *testing.T) {
	a := BioSnapshot{Arousal: 0.5, Valence: 0.3, BodyTemp: 36.6, Pain: 0.1, Fatigue: 0.2, Hunger: 0.1}
	b := BioSnapshot{Arousal: 0.55, Valence: 0.25, BodyTemp: 36.7, Pain: 0.15, Fatigue: 0.25, Hunger: 0.12}

	got := somaticSimilarity(a, b)

	if got < 0.9 {
		t.Errorf("similarity of close states = %f, expected > 0.9", got)
	}
}

func TestContextSelector_SelectsAllWhenUnderLimit(t *testing.T) {
	cs := NewContextSelector(10)
	memories := []EpisodicMemory{
		{ID: "a", Importance: 0.5},
		{ID: "b", Importance: 0.8},
	}

	selected := cs.Select(memories, BioSnapshot{})
	if len(selected) != 2 {
		t.Errorf("expected 2 memories, got %d", len(selected))
	}
}

func TestContextSelector_LimitsToMax(t *testing.T) {
	cs := NewContextSelector(2)
	memories := make([]EpisodicMemory, 10)
	for i := range memories {
		memories[i] = EpisodicMemory{
			ID:               string(rune('a' + i)),
			Importance:       float64(i) * 0.1,
			EmotionalValence: 0.1,
		}
	}

	selected := cs.Select(memories, BioSnapshot{})
	if len(selected) != 2 {
		t.Errorf("expected 2 memories, got %d", len(selected))
	}
}

func TestContextSelector_PrefersImportantMemories(t *testing.T) {
	cs := NewContextSelector(1)
	memories := []EpisodicMemory{
		{ID: "low", Importance: 0.1, EmotionalValence: 0.1},
		{ID: "high", Importance: 0.9, EmotionalValence: 0.1},
	}

	selected := cs.Select(memories, BioSnapshot{})
	if len(selected) != 1 {
		t.Fatalf("expected 1 memory, got %d", len(selected))
	}
	if selected[0].ID != "high" {
		t.Errorf("expected high-importance memory, got %q", selected[0].ID)
	}
}

func TestContextSelector_PrefersSomaticallySimilar(t *testing.T) {
	cs := NewContextSelector(1)

	current := BioSnapshot{Arousal: 0.8, Valence: -0.5, BodyTemp: 34.0, Pain: 0.5, Fatigue: 0.1, Hunger: 0}

	memories := []EpisodicMemory{
		{
			ID:               "matching",
			Importance:       0.5,
			EmotionalValence: -0.3,
			BioSnapshot:      BioSnapshot{Arousal: 0.75, Valence: -0.4, BodyTemp: 34.2, Pain: 0.4, Fatigue: 0.1, Hunger: 0},
		},
		{
			ID:               "different",
			Importance:       0.5,
			EmotionalValence: -0.3,
			BioSnapshot:      BioSnapshot{Arousal: 0.1, Valence: 0.8, BodyTemp: 37.0, Pain: 0, Fatigue: 0.8, Hunger: 0.9},
		},
	}

	selected := cs.Select(memories, current)
	if selected[0].ID != "matching" {
		t.Errorf("expected somatically similar memory, got %q", selected[0].ID)
	}
}

func TestContextSelector_PrefersEmotionallyIntense(t *testing.T) {
	cs := NewContextSelector(1)

	memories := []EpisodicMemory{
		{ID: "mild", Importance: 0.5, EmotionalValence: 0.1},
		{ID: "intense", Importance: 0.5, EmotionalValence: -0.9},
	}

	selected := cs.Select(memories, BioSnapshot{})
	if selected[0].ID != "intense" {
		t.Errorf("expected emotionally intense memory, got %q", selected[0].ID)
	}
}

func TestContextSelector_IntegrationWithStore(t *testing.T) {
	store := tempDB(t)

	now := time.Now().Truncate(time.Microsecond)
	for i := 0; i < 5; i++ {
		m := EpisodicMemory{
			ID:               string(rune('a' + i)),
			Content:          "memory " + string(rune('a'+i)),
			Timestamp:        now.Add(time.Duration(i) * time.Hour),
			EmotionalValence: float64(i) * 0.2,
			Importance:       float64(i) * 0.2,
			BioSnapshot:      BioSnapshot{Arousal: float64(i) * 0.1},
		}
		store.SaveMemory(&m)
	}

	memories, _ := store.LoadMemories()
	if len(memories) != 5 {
		t.Fatalf("expected 5 memories, got %d", len(memories))
	}

	cs := NewContextSelector(2)
	selected := cs.Select(memories, BioSnapshot{})
	if len(selected) != 2 {
		t.Errorf("expected 2 selected memories, got %d", len(selected))
	}
}

func TestSomaticSimilarity_Symmetry(t *testing.T) {
	a := BioSnapshot{Arousal: 0.3, Valence: 0.1, BodyTemp: 35.0, Pain: 0.2, Fatigue: 0.4, Hunger: 0.3}
	b := BioSnapshot{Arousal: 0.7, Valence: -0.3, BodyTemp: 37.0, Pain: 0.5, Fatigue: 0.1, Hunger: 0.6}

	ab := somaticSimilarity(a, b)
	ba := somaticSimilarity(b, a)

	if math.Abs(ab-ba) > 1e-10 {
		t.Errorf("similarity is not symmetric: %f vs %f", ab, ba)
	}
}
