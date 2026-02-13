package consciousness

import (
	"testing"

	"github.com/marczahn/person/internal/psychology"
)

func TestSalienceCalculator_NoChangeNoSalience(t *testing.T) {
	sc := NewSalienceCalculator()
	ps := &psychology.State{Arousal: 0.2, Valence: 0.3, Energy: 0.5}

	// First call establishes baseline.
	sc.Compute(ps, 1.0)

	// Same state, no change.
	result := sc.Compute(ps, 1.0)

	if result.Score > 0.1 {
		t.Errorf("score = %f for unchanged state, expected near 0", result.Score)
	}
}

func TestSalienceCalculator_SuddenArousalSpike(t *testing.T) {
	sc := NewSalienceCalculator()

	calm := &psychology.State{Arousal: 0.1, Valence: 0.3, Energy: 0.5}
	sc.Compute(calm, 1.0)

	spike := &psychology.State{Arousal: 0.9, Valence: 0.3, Energy: 0.5}
	result := sc.Compute(spike, 1.0)

	if !result.Exceeded {
		t.Errorf("expected salience to exceed threshold on sudden arousal spike, score=%f thresh=%f",
			result.Score, result.Threshold)
	}
	if result.Trigger != "arousal" {
		t.Errorf("trigger = %q, expected arousal", result.Trigger)
	}
}

func TestSalienceCalculator_GradualChangeNotSalient(t *testing.T) {
	sc := NewSalienceCalculator()

	ps := &psychology.State{Arousal: 0.2, Valence: 0.3, Energy: 0.5}
	sc.Compute(ps, 1.0)

	// Very small change over a long period.
	ps2 := &psychology.State{Arousal: 0.22, Valence: 0.3, Energy: 0.5}
	result := sc.Compute(ps2, 60.0)

	if result.Exceeded {
		t.Errorf("gradual change should not be salient, score=%f", result.Score)
	}
}

func TestSalienceCalculator_ThresholdBreachBonus(t *testing.T) {
	sc := NewSalienceCalculator()

	// Establish baseline at moderate pain (negative valence + arousal).
	ps := &psychology.State{Arousal: 0.5, Valence: -0.3, Energy: 0.5}
	sc.Compute(ps, 1.0)

	// Move into extreme range.
	extreme := &psychology.State{Arousal: 0.85, Valence: -0.7, Energy: 0.5}
	result := sc.Compute(extreme, 1.0)

	// Should have threshold breach bonus.
	if result.Score < 0.4 {
		t.Errorf("extreme state should have high salience, got %f", result.Score)
	}
}

func TestDynamicThreshold_AnxietyLowers(t *testing.T) {
	calm := &psychology.State{Arousal: 0.2, Valence: 0.3, Energy: 0.5}
	anxious := &psychology.State{Arousal: 0.7, Valence: -0.4, Energy: 0.5}

	calmThresh := dynamicThreshold(calm)
	anxThresh := dynamicThreshold(anxious)

	if anxThresh >= calmThresh {
		t.Errorf("anxious threshold (%f) should be lower than calm (%f)", anxThresh, calmThresh)
	}
}

func TestDynamicThreshold_HighCogLoadRaises(t *testing.T) {
	idle := &psychology.State{CognitiveLoad: 0.2}
	absorbed := &psychology.State{CognitiveLoad: 0.8}

	idleThresh := dynamicThreshold(idle)
	absorbedThresh := dynamicThreshold(absorbed)

	if absorbedThresh <= idleThresh {
		t.Errorf("absorbed threshold (%f) should be higher than idle (%f)", absorbedThresh, idleThresh)
	}
}

func TestNoveltyWeight_IncreasesOverTime(t *testing.T) {
	w0 := noveltyWeight(0)
	w60 := noveltyWeight(60)
	w300 := noveltyWeight(300)

	if w60 <= w0 {
		t.Errorf("novelty at 60s (%f) should be > at 0s (%f)", w60, w0)
	}
	if w300 <= w60 {
		t.Errorf("novelty at 300s (%f) should be > at 60s (%f)", w300, w60)
	}
}

func TestAttentionModifier_InwardAmplifies(t *testing.T) {
	inward := attentionModifier(1.0)
	outward := attentionModifier(-1.0)

	if inward <= outward {
		t.Errorf("inward modifier (%f) should be > outward (%f)", inward, outward)
	}
}

func TestSalienceCalculator_ZeroDt(t *testing.T) {
	sc := NewSalienceCalculator()
	ps := &psychology.State{Arousal: 0.5}

	result := sc.Compute(ps, 0)
	if result.Score != 0 {
		t.Errorf("score = %f for dt=0, expected 0", result.Score)
	}
}

func TestPainFromValence_PositiveValence_NoPain(t *testing.T) {
	ps := &psychology.State{Valence: 0.5, Arousal: 0.8}
	got := painFromValence(ps)
	if got != 0 {
		t.Errorf("pain from positive valence = %f, want 0", got)
	}
}

func TestPainFromValence_NegativeValence_WithArousal(t *testing.T) {
	ps := &psychology.State{Valence: -0.8, Arousal: 0.9}
	got := painFromValence(ps)
	if got < 0.5 {
		t.Errorf("pain from very negative valence + high arousal = %f, expected > 0.5", got)
	}
}
