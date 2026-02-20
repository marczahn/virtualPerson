package biology_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestEngineTick_ReturnsDeltasAndThresholds(t *testing.T) {
	s := biology.NewDefaultState()
	s.Stress = 0.9

	cfg := biology.Config{
		Decay:      biology.DecayConfig{DecayMultiplier: 0},
		Noise:      biology.NoiseConfig{Sigma: 0},
		Thresholds: biology.DefaultThresholdConfig(),
	}
	e := biology.NewEngineWithSeed(cfg, 1)

	result := e.Tick(s, 1.0)

	if len(result.Deltas) == 0 {
		t.Fatalf("expected interaction deltas, got none")
	}
	if len(result.Thresholds) == 0 {
		t.Fatalf("expected threshold events, got none")
	}
	if len(filterEvents(result.Thresholds, "stress")) != 1 {
		t.Fatalf("expected one stress threshold event, got %v", result.Thresholds)
	}
}

func TestEngineTick_SecondClampAfterCascades(t *testing.T) {
	s := biology.NewDefaultState()
	s.Stress = 0.96
	s.BodyTemp = 41.0

	cfg := biology.Config{
		Decay:      biology.DecayConfig{DecayMultiplier: 0},
		Noise:      biology.NoiseConfig{Sigma: 0},
		Thresholds: biology.DefaultThresholdConfig(),
	}
	e := biology.NewEngineWithSeed(cfg, 1)

	e.Tick(s, 1.0)

	if s.Stress != 1.0 {
		t.Fatalf("stress should be clamped to 1 after cascades, got %v", s.Stress)
	}
}

func TestEngineWithSeed_IsDeterministic(t *testing.T) {
	s1 := biology.NewDefaultState()
	s2 := biology.NewDefaultState()
	cfg := biology.Config{
		Decay:      biology.DecayConfig{DecayMultiplier: 0},
		Noise:      biology.NoiseConfig{Sigma: 0.005},
		Thresholds: biology.DefaultThresholdConfig(),
	}

	e1 := biology.NewEngineWithSeed(cfg, 99)
	e2 := biology.NewEngineWithSeed(cfg, 99)

	e1.Tick(s1, 2.0)
	e2.Tick(s2, 2.0)

	checks := []struct {
		name string
		a    float64
		b    float64
	}{
		{"Energy", s1.Energy, s2.Energy},
		{"Stress", s1.Stress, s2.Stress},
		{"CognitiveCapacity", s1.CognitiveCapacity, s2.CognitiveCapacity},
		{"Mood", s1.Mood, s2.Mood},
		{"PhysicalTension", s1.PhysicalTension, s2.PhysicalTension},
		{"Hunger", s1.Hunger, s2.Hunger},
		{"SocialDeficit", s1.SocialDeficit, s2.SocialDeficit},
		{"BodyTemp", s1.BodyTemp, s2.BodyTemp},
	}

	for _, c := range checks {
		if c.a != c.b {
			t.Fatalf("%s differs with same seed: %v vs %v", c.name, c.a, c.b)
		}
	}
}

func TestEngine_TenMinutesUnattended(t *testing.T) {
	cfg := biology.DefaultConfig()
	cfg.Decay.DecayMultiplier = 1.0
	s := biology.NewDefaultState()
	initialEnergy := s.Energy
	initialHunger := s.Hunger
	initialMood := s.Mood

	e := biology.NewEngineWithSeed(cfg, 42)
	for i := 0; i < 600; i++ {
		e.Tick(s, 1.0)
	}

	t.Logf("final state after 600 ticks: energy=%.4f hunger=%.4f mood=%.4f", s.Energy, s.Hunger, s.Mood)
	if s.Energy > initialEnergy-0.2 {
		t.Fatalf("expected energy drop >=0.2, got initial=%v final=%v", initialEnergy, s.Energy)
	}
	if s.Hunger < initialHunger+0.2 {
		t.Fatalf("expected hunger rise >=0.2, got initial=%v final=%v", initialHunger, s.Hunger)
	}
	if s.CognitiveCapacity > 1.0 {
		t.Fatalf("cognitive capacity exceeded max: %v", s.CognitiveCapacity)
	}
	if s.Mood > initialMood {
		t.Fatalf("expected mood to drift down or stay flat, got initial=%v final=%v", initialMood, s.Mood)
	}
}

func TestEngine_FastMode_TwoMinutes(t *testing.T) {
	cfg := biology.DefaultConfig()
	s := biology.NewDefaultState()
	initialEnergy := s.Energy
	initialHunger := s.Hunger

	e := biology.NewEngineWithSeed(cfg, 42)
	for i := 0; i < 120; i++ {
		e.Tick(s, 1.0)
	}

	t.Logf("final state after 120 ticks fast mode: energy=%.4f hunger=%.4f", s.Energy, s.Hunger)
	if s.Energy > initialEnergy-0.2 {
		t.Fatalf("expected energy drop >=0.2 in fast mode, got initial=%v final=%v", initialEnergy, s.Energy)
	}
	if s.Hunger < initialHunger+0.2 {
		t.Fatalf("expected hunger rise >=0.2 in fast mode, got initial=%v final=%v", initialHunger, s.Hunger)
	}
}

func TestEngine_AllVariablesInRange(t *testing.T) {
	cfg := biology.DefaultConfig()
	s := biology.NewDefaultState()
	e := biology.NewEngineWithSeed(cfg, 42)

	for i := 0; i < 3000; i++ {
		e.Tick(s, 1.0)

		checks := []struct {
			name string
			v    float64
			min  float64
			max  float64
		}{
			{"energy", s.Energy, 0, 1},
			{"stress", s.Stress, 0, 1},
			{"cognitive_capacity", s.CognitiveCapacity, 0, 1},
			{"mood", s.Mood, 0, 1},
			{"physical_tension", s.PhysicalTension, 0, 1},
			{"hunger", s.Hunger, 0, 1},
			{"social_deficit", s.SocialDeficit, 0, 1},
			{"body_temp", s.BodyTemp, 25, 43},
		}
		for _, c := range checks {
			if c.v < c.min || c.v > c.max {
				t.Fatalf("tick %d: %s out of range: %v not in [%v,%v]", i, c.name, c.v, c.min, c.max)
			}
		}
	}
}

func TestEngine_DtZero_NoDecay(t *testing.T) {
	cfg := biology.DefaultConfig()
	cfg.Decay.DecayMultiplier = 1.0
	s := biology.NewDefaultState()
	before := *s
	e := biology.NewEngineWithSeed(cfg, 42)

	for i := 0; i < 10; i++ {
		e.Tick(s, 0.0)
	}

	if s.Energy != before.Energy {
		t.Fatalf("energy changed at dt=0: %v -> %v", before.Energy, s.Energy)
	}
	if s.Stress != before.Stress {
		t.Fatalf("stress changed at dt=0: %v -> %v", before.Stress, s.Stress)
	}
	if s.CognitiveCapacity != before.CognitiveCapacity {
		t.Fatalf("cognitive capacity changed at dt=0: %v -> %v", before.CognitiveCapacity, s.CognitiveCapacity)
	}
	if s.Mood != before.Mood {
		t.Fatalf("mood changed at dt=0: %v -> %v", before.Mood, s.Mood)
	}
	if s.PhysicalTension != before.PhysicalTension {
		t.Fatalf("physical tension changed at dt=0: %v -> %v", before.PhysicalTension, s.PhysicalTension)
	}
	if s.Hunger != before.Hunger {
		t.Fatalf("hunger changed at dt=0: %v -> %v", before.Hunger, s.Hunger)
	}
	if s.SocialDeficit != before.SocialDeficit {
		t.Fatalf("social deficit changed at dt=0: %v -> %v", before.SocialDeficit, s.SocialDeficit)
	}
	if s.BodyTemp != before.BodyTemp {
		t.Fatalf("body temp changed at dt=0: %v -> %v", before.BodyTemp, s.BodyTemp)
	}
}

func TestEngine_ThresholdsDetectedDuringDegradation(t *testing.T) {
	cfg := biology.DefaultConfig()
	cfg.Decay.DecayMultiplier = 10.0
	s := biology.NewDefaultState()
	s.Energy = 0.04
	e := biology.NewEngineWithSeed(cfg, 42)

	foundCriticalEnergy := false
	for i := 0; i < 5; i++ {
		result := e.Tick(s, 1.0)
		for _, event := range result.Thresholds {
			if event.Variable == "energy" && event.Severity == biology.Critical {
				foundCriticalEnergy = true
				break
			}
		}
		if foundCriticalEnergy {
			break
		}
	}

	if !foundCriticalEnergy {
		t.Fatalf("expected at least one critical energy threshold event, got none")
	}
}

func TestEngine_SlowMode_SixtyTicksDiffersFromStart(t *testing.T) {
	cfg := biology.DefaultConfig()
	cfg.Decay.DecayMultiplier = 1.0
	s := biology.NewDefaultState()
	before := *s
	e := biology.NewEngineWithSeed(cfg, 42)

	for i := 0; i < 60; i++ {
		e.Tick(s, 1.0)
	}

	if s.Energy == before.Energy &&
		s.Hunger == before.Hunger &&
		s.CognitiveCapacity == before.CognitiveCapacity &&
		s.Mood == before.Mood &&
		s.SocialDeficit == before.SocialDeficit {
		t.Fatalf("expected slow-path degradation to produce state change after 60 ticks, state unchanged")
	}
}
