package psychology

import (
	"testing"
	"time"
)

func TestComputeIsolationPhase_Timeline(t *testing.T) {
	tests := []struct {
		hours float64
		want  IsolationPhase
	}{
		{0, IsolationNone},
		{1, IsolationNone},
		{3, IsolationBoredom},
		{7, IsolationBoredom},
		{12, IsolationLoneliness},
		{20, IsolationLoneliness},
		{48, IsolationSignificant},
		{70, IsolationSignificant},
		{100, IsolationDestabilizing},
		{160, IsolationDestabilizing},
		{200, IsolationSevere},
	}

	for _, tt := range tests {
		got := computeIsolationPhase(tt.hours)
		if got != tt.want {
			t.Errorf("computeIsolationPhase(%f) = %s, want %s", tt.hours, got, tt.want)
		}
	}
}

func TestComputeLoneliness_Progression(t *testing.T) {
	extFactor := 1.0 // average extraversion
	var prev float64

	for _, hours := range []float64{0, 2, 8, 24, 72, 168, 336} {
		loneliness := computeLoneliness(hours, extFactor)
		if loneliness < prev {
			t.Errorf("loneliness decreased: at %fh = %f, was %f", hours, loneliness, prev)
		}
		prev = loneliness
	}

	// At 336 hours (2 weeks), should be near maximum.
	if prev < 0.8 {
		t.Errorf("loneliness at 336h = %f, expected > 0.8", prev)
	}
}

func TestComputeLoneliness_HighExtraversionFaster(t *testing.T) {
	lowE := computeLoneliness(12, IsolationDistressRate(Personality{Extraversion: 0.2}))
	highE := computeLoneliness(12, IsolationDistressRate(Personality{Extraversion: 0.9}))

	if highE <= lowE {
		t.Errorf("high E loneliness (%f) should be > low E (%f) at 12h", highE, lowE)
	}
}

func TestUpdateIsolation_AdvancesDuration(t *testing.T) {
	p := Personality{Extraversion: 0.5, Neuroticism: 0.5}
	state := IsolationState{}

	state = UpdateIsolation(state, p, 7200) // 2 hours

	if state.Duration < 2*time.Hour {
		t.Errorf("duration = %v, expected >= 2h", state.Duration)
	}
	if state.Phase != IsolationNone {
		t.Errorf("phase = %s, expected none at 2h", state.Phase)
	}
}

func TestUpdateIsolation_PhaseProgression(t *testing.T) {
	p := Personality{Extraversion: 0.5, Neuroticism: 0.5}
	state := IsolationState{}

	// Advance to 5 hours.
	state = UpdateIsolation(state, p, 5*3600)
	if state.Phase != IsolationBoredom {
		t.Errorf("at 5h: phase = %s, expected boredom", state.Phase)
	}

	// Advance another 7 hours (total 12).
	state = UpdateIsolation(state, p, 7*3600)
	if state.Phase != IsolationLoneliness {
		t.Errorf("at 12h: phase = %s, expected loneliness", state.Phase)
	}
}

func TestUpdateIsolation_LonelinessBounded(t *testing.T) {
	p := Personality{Extraversion: 0.9}
	state := IsolationState{}

	// Advance to 1000 hours.
	state = UpdateIsolation(state, p, 1000*3600)

	if state.LonelinessLevel > 1.0 {
		t.Errorf("loneliness = %f, should not exceed 1.0", state.LonelinessLevel)
	}
}
