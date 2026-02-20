package biology_test

import (
	"math"
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestApplyFeedbackAtTickEnd_PulseNotDtScaledAndRateIsDtScaled(t *testing.T) {
	s := biology.NewDefaultState()
	s.Mood = 0.30

	feedback := biology.FeedbackEnvelope{
		Rates: []biology.BioRate{
			{Field: "mood", PerSecond: 0.01},
		},
		Pulses: []biology.BioPulse{
			{Field: "mood", Amount: 0.10},
		},
	}

	biology.ApplyFeedbackAtTickEnd(s, 10, feedback)

	// Mood starts at 0.30.
	// Rate contributes +0.10 (0.01 * 10), pulse contributes +0.10 (absolute).
	if s.Mood != 0.50 {
		t.Fatalf("expected mood=0.50, got %f", s.Mood)
	}
}

func TestApplyFeedbackAtTickEnd_PulseAppliesAtDtZeroRateDoesNot(t *testing.T) {
	s := biology.NewDefaultState()
	s.Stress = 0.20

	feedback := biology.FeedbackEnvelope{
		Rates: []biology.BioRate{
			{Field: "stress", PerSecond: 0.50},
		},
		Pulses: []biology.BioPulse{
			{Field: "stress", Amount: -0.10},
		},
	}

	biology.ApplyFeedbackAtTickEnd(s, 0, feedback)

	if s.Stress != 0.10 {
		t.Fatalf("expected stress=0.10 from pulse-only effect at dt=0, got %f", s.Stress)
	}
}

func TestTickFeedbackBuffer_AppliesOnlyAtEndOfTick(t *testing.T) {
	s := biology.NewDefaultState()
	s.Hunger = 0.70
	s.Stress = 0.10

	var buffer biology.TickFeedbackBuffer
	buffer.AddPulses([]biology.BioPulse{{Field: "hunger", Amount: -0.20}})
	buffer.AddRates([]biology.BioRate{{Field: "stress", PerSecond: 0.05}})

	// Contract: collecting feedback does not mutate state mid-tick.
	if s.Hunger != 0.70 || s.Stress != 0.10 {
		t.Fatalf("state changed mid-tick while buffering feedback: hunger=%f stress=%f", s.Hunger, s.Stress)
	}

	buffer.ApplyAtTickEnd(s, 2)

	if math.Abs(s.Hunger-0.50) > 0.000001 {
		t.Fatalf("expected end-of-tick hunger=0.50, got %f", s.Hunger)
	}
	if math.Abs(s.Stress-0.20) > 0.000001 {
		t.Fatalf("expected end-of-tick stress=0.20, got %f", s.Stress)
	}
}
