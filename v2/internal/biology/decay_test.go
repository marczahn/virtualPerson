package biology_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestApplyDecay_EnergyDrains(t *testing.T) {
	s := biology.NewDefaultState()
	cfg := biology.DecayConfig{DecayMultiplier: 1.0}
	origEnergy := s.Energy
	origHunger := s.Hunger

	biology.ApplyDecay(s, cfg, 1.0)

	// Energy should decrease: 0.8 - 0.00067*1.0*1.0 = 0.79933
	if s.Energy >= origEnergy {
		t.Errorf("Energy should decrease after decay, got %v (was %v)", s.Energy, origEnergy)
	}
	// Hunger should increase: 0.1 + 0.00083*1.0*1.0 = 0.10083
	if s.Hunger <= origHunger {
		t.Errorf("Hunger should increase after decay, got %v (was %v)", s.Hunger, origHunger)
	}
}

func TestApplyDecay_DtZero_NoChange(t *testing.T) {
	s := biology.NewDefaultState()
	cfg := biology.DecayConfig{DecayMultiplier: 5.0}
	before := *s

	biology.ApplyDecay(s, cfg, 0.0)

	checks := []struct {
		name string
		got  float64
		want float64
	}{
		{"Energy", s.Energy, before.Energy},
		{"Hunger", s.Hunger, before.Hunger},
		{"CognitiveCapacity", s.CognitiveCapacity, before.CognitiveCapacity},
		{"Mood", s.Mood, before.Mood},
		{"SocialDeficit", s.SocialDeficit, before.SocialDeficit},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("dt=0: %s changed from %v to %v", c.name, c.want, c.got)
		}
	}
}

func TestApplyDecay_NoTouchStressOrTension(t *testing.T) {
	s := biology.NewDefaultState()
	cfg := biology.DecayConfig{DecayMultiplier: 5.0}
	origStress := s.Stress
	origTension := s.PhysicalTension
	origTemp := s.BodyTemp

	biology.ApplyDecay(s, cfg, 1.0)

	if s.Stress != origStress {
		t.Errorf("Stress changed from %v to %v (should not decay autonomously)", origStress, s.Stress)
	}
	if s.PhysicalTension != origTension {
		t.Errorf("PhysicalTension changed from %v to %v (should not decay autonomously)", origTension, s.PhysicalTension)
	}
	if s.BodyTemp != origTemp {
		t.Errorf("BodyTemp changed from %v to %v (should not decay autonomously)", origTemp, s.BodyTemp)
	}
}

func TestApplyDecay_DtCap(t *testing.T) {
	// dt=100 should produce same result as dt=60 (cap enforced)
	s1 := biology.NewDefaultState()
	s2 := biology.NewDefaultState()
	cfg := biology.DecayConfig{DecayMultiplier: 1.0}

	biology.ApplyDecay(s1, cfg, 100.0)
	biology.ApplyDecay(s2, cfg, 60.0)

	checks := []struct {
		name string
		v1   float64
		v2   float64
	}{
		{"Energy", s1.Energy, s2.Energy},
		{"Hunger", s1.Hunger, s2.Hunger},
		{"CognitiveCapacity", s1.CognitiveCapacity, s2.CognitiveCapacity},
		{"Mood", s1.Mood, s2.Mood},
		{"SocialDeficit", s1.SocialDeficit, s2.SocialDeficit},
	}
	for _, c := range checks {
		if c.v1 != c.v2 {
			t.Errorf("dt=100 vs dt=60: %s differs: %v vs %v", c.name, c.v1, c.v2)
		}
	}
}

func TestApplyDecay_RateCalibration(t *testing.T) {
	s := biology.NewDefaultState()
	cfg := biology.DecayConfig{DecayMultiplier: 1.0}
	origEnergy := s.Energy
	origHunger := s.Hunger

	biology.ApplyDecay(s, cfg, 60.0)

	if s.Energy >= origEnergy {
		t.Errorf("Energy should decrease with decay: got %v want less than %v", s.Energy, origEnergy)
	}
	if s.Hunger <= origHunger {
		t.Errorf("Hunger should increase with decay: got %v want greater than %v", s.Hunger, origHunger)
	}
}

func TestApplyDecay_AllFiveVariablesMove(t *testing.T) {
	s := biology.NewDefaultState()
	cfg := biology.DecayConfig{DecayMultiplier: 5.0}

	biology.ApplyDecay(s, cfg, 10.0)

	// Energy decreases: 0.8 - 0.00067*5*10 = 0.8 - 0.0335 = 0.7665
	if s.Energy >= 0.80 {
		t.Errorf("Energy should decrease from baseline, got %v", s.Energy)
	}

	// Hunger increases: 0.1 + 0.00083*5*10 = 0.1 + 0.0415 = 0.1415
	if s.Hunger <= 0.10 {
		t.Errorf("Hunger should increase from baseline, got %v", s.Hunger)
	}

	// CognitiveCapacity decreases: 1.0 - 0.00050*5*10 = 1.0 - 0.025 = 0.975
	if s.CognitiveCapacity >= 1.00 {
		t.Errorf("CognitiveCapacity should decrease from baseline, got %v", s.CognitiveCapacity)
	}

	// Mood decreases: 0.5 - 0.00033*5*10 = 0.5 - 0.0165 = 0.4835
	if s.Mood >= 0.50 {
		t.Errorf("Mood should decrease from baseline, got %v", s.Mood)
	}

	// SocialDeficit increases: 0.0 + 0.00033*5*10 = 0.0 + 0.0165 = 0.0165
	if s.SocialDeficit <= 0.00 {
		t.Errorf("SocialDeficit should increase from baseline, got %v", s.SocialDeficit)
	}
}
