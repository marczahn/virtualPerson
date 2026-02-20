package biology_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestThresholds_SteppedTierOnlyMostSevere(t *testing.T) {
	cfg := biology.DefaultThresholdConfig()

	tests := []struct {
		name         string
		stress       float64
		wantSeverity biology.Severity
		wantEvent    bool
	}{
		{"Stress=0.96 -> Critical", 0.96, biology.Critical, true},
		{"Stress=0.90 -> Warning", 0.90, biology.Warning, true},
		{"Stress=0.72 -> Mild", 0.72, biology.Mild, true},
		{"Stress=0.50 -> no event", 0.50, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := biology.NewDefaultState()
			s.Stress = tt.stress
			events := biology.EvaluateThresholds(s, cfg, 1.0)

			stressEvents := filterEvents(events, "stress")

			if tt.wantEvent {
				if len(stressEvents) != 1 {
					t.Fatalf("expected exactly 1 stress event, got %d", len(stressEvents))
				}
				if stressEvents[0].Severity != tt.wantSeverity {
					t.Errorf("stress severity = %v, want %v", stressEvents[0].Severity, tt.wantSeverity)
				}
			} else {
				if len(stressEvents) != 0 {
					t.Errorf("expected no stress event at Stress=%v, got %d", tt.stress, len(stressEvents))
				}
			}
		})
	}
}

func TestThresholds_BodyTemp(t *testing.T) {
	cfg := biology.DefaultThresholdConfig()

	tests := []struct {
		name         string
		bodyTemp     float64
		wantSeverity biology.Severity
		wantEvent    bool
	}{
		{"BodyTemp=32.0 -> Critical hypothermia", 32.0, biology.Critical, true},
		{"BodyTemp=34.5 -> Mild hypothermia", 34.5, biology.Mild, true},
		{"BodyTemp=33.5 -> Warning hypothermia", 33.5, biology.Warning, true},
		{"BodyTemp=36.6 -> no event", 36.6, 0, false},
		{"BodyTemp=41.0 -> Critical hyperthermia", 41.0, biology.Critical, true},
		{"BodyTemp=40.0 -> Warning hyperthermia", 40.0, biology.Warning, true},
		{"BodyTemp=39.0 -> Mild hyperthermia", 39.0, biology.Mild, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := biology.NewDefaultState()
			s.BodyTemp = tt.bodyTemp
			events := biology.EvaluateThresholds(s, cfg, 1.0)

			tempEvents := filterEvents(events, "body_temp")

			if tt.wantEvent {
				if len(tempEvents) != 1 {
					t.Fatalf("expected exactly 1 body_temp event at BodyTemp=%v, got %d: %v",
						tt.bodyTemp, len(tempEvents), tempEvents)
				}
				if tempEvents[0].Severity != tt.wantSeverity {
					t.Errorf("body_temp severity at %v = %v, want %v",
						tt.bodyTemp, tempEvents[0].Severity, tt.wantSeverity)
				}
			} else {
				if len(tempEvents) != 0 {
					t.Errorf("expected no body_temp event at BodyTemp=%v, got %d", tt.bodyTemp, len(tempEvents))
				}
			}
		})
	}
}

func TestThresholds_EnergyCritical(t *testing.T) {
	cfg := biology.DefaultThresholdConfig()
	s := biology.NewDefaultState()
	s.Energy = 0.04

	events := biology.EvaluateThresholds(s, cfg, 1.0)
	energyEvents := filterEvents(events, "energy")

	if len(energyEvents) != 1 {
		t.Fatalf("expected 1 energy event at Energy=0.04, got %d", len(energyEvents))
	}
	if energyEvents[0].Severity != biology.Critical {
		t.Errorf("energy severity = %v, want Critical", energyEvents[0].Severity)
	}
}

func TestThresholds_HungerCritical(t *testing.T) {
	cfg := biology.DefaultThresholdConfig()
	s := biology.NewDefaultState()
	s.Hunger = 0.96

	events := biology.EvaluateThresholds(s, cfg, 1.0)
	hungerEvents := filterEvents(events, "hunger")

	if len(hungerEvents) != 1 {
		t.Fatalf("expected 1 hunger event at Hunger=0.96, got %d", len(hungerEvents))
	}
	if hungerEvents[0].Severity != biology.Critical {
		t.Errorf("hunger severity = %v, want Critical", hungerEvents[0].Severity)
	}
}

func TestThresholds_CascadeApplied(t *testing.T) {
	cfg := biology.DefaultThresholdConfig()
	s := biology.NewDefaultState()
	s.Stress = 0.96 // Critical stress
	origMood := s.Mood
	origEnergy := s.Energy

	events := biology.EvaluateThresholds(s, cfg, 1.0)
	biology.ApplyThresholdCascades(s, events)

	// Critical stress cascade: Mood -= 0.03*dt, Energy -= 0.02*dt
	if s.Mood >= origMood {
		t.Errorf("Mood should decrease from critical stress cascade, got %v (was %v)", s.Mood, origMood)
	}
	if s.Energy >= origEnergy {
		t.Errorf("Energy should decrease from critical stress cascade, got %v (was %v)", s.Energy, origEnergy)
	}
}

func TestThresholds_NoCascadeAtBaseline(t *testing.T) {
	cfg := biology.DefaultThresholdConfig()
	s := biology.NewDefaultState()
	before := *s

	events := biology.EvaluateThresholds(s, cfg, 1.0)
	biology.ApplyThresholdCascades(s, events)

	if s.Energy != before.Energy {
		t.Errorf("Energy changed at baseline: %v -> %v", before.Energy, s.Energy)
	}
	if s.Mood != before.Mood {
		t.Errorf("Mood changed at baseline: %v -> %v", before.Mood, s.Mood)
	}
	if s.Stress != before.Stress {
		t.Errorf("Stress changed at baseline: %v -> %v", before.Stress, s.Stress)
	}
}
