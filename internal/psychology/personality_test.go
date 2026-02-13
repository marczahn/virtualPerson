package psychology

import (
	"math"
	"testing"
)

func TestNegativeEmotionMultiplier_HighNeuroticism(t *testing.T) {
	p := Personality{Neuroticism: 0.8}
	got := NegativeEmotionMultiplier(p)
	want := 1.36

	if math.Abs(got-want) > 0.01 {
		t.Errorf("NegativeEmotionMultiplier(N=0.8) = %f, want ~%f", got, want)
	}
}

func TestNegativeEmotionMultiplier_LowNeuroticism(t *testing.T) {
	p := Personality{Neuroticism: 0.2}
	got := NegativeEmotionMultiplier(p)
	want := 0.64

	if math.Abs(got-want) > 0.01 {
		t.Errorf("NegativeEmotionMultiplier(N=0.2) = %f, want ~%f", got, want)
	}
}

func TestIsolationDistressRate_HighExtraversion(t *testing.T) {
	p := Personality{Extraversion: 0.85}
	got := IsolationDistressRate(p)

	// 1.0 + (0.85 - 0.5) * 1.5 = 1.525
	if got < 1.4 || got > 1.6 {
		t.Errorf("IsolationDistressRate(E=0.85) = %f, expected 1.4-1.6", got)
	}
}

func TestIsolationDistressRate_LowExtraversion(t *testing.T) {
	p := Personality{Extraversion: 0.2}
	got := IsolationDistressRate(p)

	// 1.0 + (0.2 - 0.5) * 1.5 = 0.55
	if got < 0.4 || got > 0.7 {
		t.Errorf("IsolationDistressRate(E=0.2) = %f, expected 0.4-0.7", got)
	}
}

func TestBaselineRegulation_HighResources(t *testing.T) {
	p := Personality{
		Conscientiousness: 0.9,
		Openness:          0.8,
		Neuroticism:       0.2,
	}
	got := BaselineRegulation(p)

	// 0.3 + 0.9*0.2 + 0.8*0.15 + 0.8*0.2 = 0.3 + 0.18 + 0.12 + 0.16 = 0.76
	if got < 0.7 || got > 0.85 {
		t.Errorf("BaselineRegulation(high resources) = %f, expected 0.7-0.85", got)
	}
}

func TestBaselineRegulation_LowResources(t *testing.T) {
	p := Personality{
		Conscientiousness: 0.1,
		Openness:          0.2,
		Neuroticism:       0.9,
	}
	got := BaselineRegulation(p)

	// 0.3 + 0.1*0.2 + 0.2*0.15 + 0.1*0.2 = 0.3 + 0.02 + 0.03 + 0.02 = 0.37
	if got < 0.3 || got > 0.45 {
		t.Errorf("BaselineRegulation(low resources) = %f, expected 0.3-0.45", got)
	}
}

func TestReappraisalAbility_HighOpenness(t *testing.T) {
	p := Personality{Openness: 0.9}
	got := ReappraisalAbility(p)

	// 0.2 + 0.9 * 0.5 = 0.65
	if math.Abs(got-0.65) > 0.01 {
		t.Errorf("ReappraisalAbility(O=0.9) = %f, want ~0.65", got)
	}
}

func TestIsolationResilience_IntrovertLowNeuroticism(t *testing.T) {
	p := Personality{
		Extraversion:      0.1,
		Neuroticism:       0.1,
		Conscientiousness: 0.8,
	}
	got := IsolationResilience(p)

	// (0.9)*0.5 + (0.9)*0.3 + 0.8*0.2 = 0.45 + 0.27 + 0.16 = 0.88
	if got < 0.8 {
		t.Errorf("IsolationResilience(introvert, low N) = %f, expected > 0.8", got)
	}
}

func TestIsolationResilience_ExtrovertHighNeuroticism(t *testing.T) {
	p := Personality{
		Extraversion:      0.9,
		Neuroticism:       0.9,
		Conscientiousness: 0.2,
	}
	got := IsolationResilience(p)

	// (0.1)*0.5 + (0.1)*0.3 + 0.2*0.2 = 0.05 + 0.03 + 0.04 = 0.12
	if got > 0.2 {
		t.Errorf("IsolationResilience(extrovert, high N) = %f, expected < 0.2", got)
	}
}

func TestDestabilizationThresholdHours(t *testing.T) {
	resilient := Personality{Extraversion: 0.1, Neuroticism: 0.1, Conscientiousness: 0.8}
	vulnerable := Personality{Extraversion: 0.9, Neuroticism: 0.9, Conscientiousness: 0.2}

	rHours := DestabilizationThresholdHours(resilient)
	vHours := DestabilizationThresholdHours(vulnerable)

	if rHours <= vHours {
		t.Errorf("resilient threshold (%f) should be > vulnerable (%f)", rHours, vHours)
	}
	if rHours < 100 {
		t.Errorf("resilient threshold = %f, expected > 100 hours", rHours)
	}
	if vHours > 80 {
		t.Errorf("vulnerable threshold = %f, expected < 80 hours", vHours)
	}
}
