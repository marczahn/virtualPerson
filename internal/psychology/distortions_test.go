package psychology

import "testing"

func TestStressMultiplier_LowStress(t *testing.T) {
	got := stressMultiplier(0.2)
	if got != 1.0 {
		t.Errorf("stressMultiplier(0.2) = %f, want 1.0", got)
	}
}

func TestStressMultiplier_ModerateStress(t *testing.T) {
	got := stressMultiplier(0.5)
	// 1.5 + (0.5-0.5)*3.0 = 1.5
	if got < 1.4 || got > 1.6 {
		t.Errorf("stressMultiplier(0.5) = %f, want ~1.5", got)
	}
}

func TestStressMultiplier_HighStress(t *testing.T) {
	got := stressMultiplier(0.9)
	// 2.1 + (0.9-0.7)*4.0 = 2.1 + 0.8 = 2.9
	if got < 2.8 || got > 3.0 {
		t.Errorf("stressMultiplier(0.9) = %f, want ~2.9", got)
	}
}

func TestActivateDistortionsDeterministic_NoDistortionsAtLowStress(t *testing.T) {
	p := Personality{Neuroticism: 0.3}
	// Use a rand function that always returns 0.5 — should be above all
	// base_rate * 1.0 (stress mult) * ~1.3 (trait mult) * ~0.7 (reg) ≈ 0.03-0.05.
	always := func() float64 { return 0.5 }

	active := ActivateDistortionsDeterministic(0.1, 0.5, p, always)

	if len(active) != 0 {
		t.Errorf("expected no distortions at low stress with 0.5 rand, got %v", active)
	}
}

func TestActivateDistortionsDeterministic_AllDistortionsAtExtremeStress(t *testing.T) {
	p := Personality{Neuroticism: 0.9, Openness: 0.1, Conscientiousness: 0.1, Agreeableness: 0.9, Extraversion: 0.8}
	// Use a rand function that always returns 0.0 — every distortion fires.
	always := func() float64 { return 0.0 }

	active := ActivateDistortionsDeterministic(0.95, 0.0, p, always)

	if len(active) != len(distortionSpecs) {
		t.Errorf("expected all %d distortions at extreme stress with 0.0 rand, got %d", len(distortionSpecs), len(active))
	}
}

func TestActivateDistortionsDeterministic_RegulationReducesDistortions(t *testing.T) {
	p := Personality{Neuroticism: 0.8}

	// Fixed random at a moderate threshold.
	threshold := 0.15
	randFn := func() float64 { return threshold }

	lowReg := ActivateDistortionsDeterministic(0.7, 0.1, p, randFn)
	highReg := ActivateDistortionsDeterministic(0.7, 0.8, p, randFn)

	if len(highReg) > len(lowReg) {
		t.Errorf("high regulation (%d distortions) should produce <= low regulation (%d)", len(highReg), len(lowReg))
	}
}

func TestActivateDistortionsDeterministic_NeuroticismIncreasesDistortions(t *testing.T) {
	threshold := 0.1
	randFn := func() float64 { return threshold }

	lowN := ActivateDistortionsDeterministic(0.6, 0.3, Personality{Neuroticism: 0.2}, randFn)
	highN := ActivateDistortionsDeterministic(0.6, 0.3, Personality{Neuroticism: 0.9}, randFn)

	if len(highN) < len(lowN) {
		t.Errorf("high N (%d distortions) should produce >= low N (%d)", len(highN), len(lowN))
	}
}
