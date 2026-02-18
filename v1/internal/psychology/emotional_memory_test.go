package psychology

import (
	"math"
	"testing"
	"time"

	"github.com/marczahn/person/internal/biology"
)

func TestRecencyDecay_RecentMemory(t *testing.T) {
	got := recencyDecay(0, false)
	if got != 1.0 {
		t.Errorf("recencyDecay(0) = %f, want 1.0", got)
	}
}

func TestRecencyDecay_OneDayOld(t *testing.T) {
	got := recencyDecay(1, false)
	// 1 / (1 + 1^0.5) = 0.5
	if math.Abs(got-0.5) > 0.01 {
		t.Errorf("recencyDecay(1) = %f, want ~0.5", got)
	}
}

func TestRecencyDecay_TraumaticSlower(t *testing.T) {
	normal := recencyDecay(10, false)
	traumatic := recencyDecay(10, true)

	if traumatic <= normal {
		t.Errorf("traumatic decay (%f) should be > normal (%f) at 10 days", traumatic, normal)
	}
}

func TestEmotionalSalience_NegativeBias(t *testing.T) {
	neg := emotionalSalience(-1)
	pos := emotionalSalience(1)

	if neg <= pos {
		t.Errorf("negative salience (%f) should be > positive (%f)", neg, pos)
	}
	if neg != 1.5 {
		t.Errorf("negative salience = %f, want 1.5", neg)
	}
}

func TestStimulusSimilarity_Cold(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.BodyTemp = 34.0

	sim := stimulusSimilarity("cold", &bio)
	if sim <= 0 {
		t.Errorf("cold stimulus similarity at 34°C = %f, expected > 0", sim)
	}

	bio.BodyTemp = 36.6
	sim = stimulusSimilarity("cold", &bio)
	if sim != 0 {
		t.Errorf("cold stimulus similarity at 36.6°C = %f, expected 0", sim)
	}
}

func TestStimulusSimilarity_Pain(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Pain = 0.6

	sim := stimulusSimilarity("pain", &bio)
	if sim != 0.6 {
		t.Errorf("pain stimulus similarity = %f, want 0.6", sim)
	}
}

func TestEmotionalMemoryStore_AddAndQuery(t *testing.T) {
	store := NewEmotionalMemoryStore()
	store.Add(EmotionalMemory{
		ID:        "cold1",
		Stimulus:  "cold",
		Valence:   -0.7,
		Intensity: 0.8,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	})

	bio := biology.NewDefaultState()
	bio.BodyTemp = 34.0

	activations := store.QueryActivations(&bio)

	if len(activations) == 0 {
		t.Fatal("expected cold memory to activate at 34°C")
	}
	if activations[0].ValenceSign != -1 {
		t.Errorf("expected negative valence sign, got %f", activations[0].ValenceSign)
	}
	if activations[0].Intensity <= 0 {
		t.Errorf("expected positive activation intensity, got %f", activations[0].Intensity)
	}
}

func TestEmotionalMemoryStore_NoActivation_IrrelevantStimulus(t *testing.T) {
	store := NewEmotionalMemoryStore()
	store.Add(EmotionalMemory{
		ID:        "cold1",
		Stimulus:  "cold",
		Valence:   -0.7,
		Intensity: 0.8,
		CreatedAt: time.Now(),
	})

	bio := biology.NewDefaultState() // warm, no cold stimulus

	activations := store.QueryActivations(&bio)
	if len(activations) != 0 {
		t.Errorf("expected no activation at normal temp, got %v", activations)
	}
}

func TestMemoryModifier_NegativeBiasAmplified(t *testing.T) {
	activations := []EmotionalMemoryActivation{
		{ValenceSign: -1, Intensity: 0.5},
		{ValenceSign: 1, Intensity: 0.5},
	}

	p := Personality{Neuroticism: 0.5}
	mod := MemoryModifier(activations, p)

	// Negative gets 1.5x salience, positive gets 1.0x.
	// negWeight = 1.5*0.5 = 0.75, posWeight = 1.0*0.5 = 0.5
	// modifier = (0.75 - 0.5) * 1.0 = 0.25
	if mod <= 0 {
		t.Errorf("modifier with equal pos/neg = %f, expected positive (negativity bias)", mod)
	}
}

func TestMemoryModifier_HighNeuroticism_Amplified(t *testing.T) {
	activations := []EmotionalMemoryActivation{
		{ValenceSign: -1, Intensity: 0.6},
	}

	lowN := MemoryModifier(activations, Personality{Neuroticism: 0.2})
	highN := MemoryModifier(activations, Personality{Neuroticism: 0.9})

	if highN <= lowN {
		t.Errorf("high N modifier (%f) should be > low N (%f)", highN, lowN)
	}
}
