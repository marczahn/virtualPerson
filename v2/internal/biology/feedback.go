package biology

// BioRate is a dt-scaled biological effect.
// PerSecond is multiplied by dt at application time.
type BioRate struct {
	Field     string
	PerSecond float64
}

// BioPulse is a one-shot biological effect.
// Amount is applied as-is and must not be dt-scaled.
type BioPulse struct {
	Field  string
	Amount float64
}

// FeedbackEnvelope carries feedback effects accumulated during a tick.
// Contract: rates are dt-scaled at apply time; pulses are absolute.
type FeedbackEnvelope struct {
	Rates  []BioRate
	Pulses []BioPulse
}

// TickFeedbackBuffer accumulates feedback during a tick and applies it once at tick end.
// Contract: collecting feedback never mutates State until ApplyAtTickEnd is called.
type TickFeedbackBuffer struct {
	rates  []BioRate
	pulses []BioPulse
}

func (b *TickFeedbackBuffer) AddRates(rates []BioRate) {
	if len(rates) == 0 {
		return
	}
	b.rates = append(b.rates, rates...)
}

func (b *TickFeedbackBuffer) AddPulses(pulses []BioPulse) {
	if len(pulses) == 0 {
		return
	}
	b.pulses = append(b.pulses, pulses...)
}

func (b *TickFeedbackBuffer) ApplyAtTickEnd(s *State, dt float64) {
	ApplyFeedbackAtTickEnd(s, dt, FeedbackEnvelope{
		Rates:  b.rates,
		Pulses: b.pulses,
	})
	b.rates = nil
	b.pulses = nil
}

// ApplyFeedbackAtTickEnd applies accumulated feedback effects to state in one commit.
// This is the only point where buffered feedback mutates state.
func ApplyFeedbackAtTickEnd(s *State, dt float64, feedback FeedbackEnvelope) {
	for _, rate := range feedback.Rates {
		applyDelta(s, Delta{
			Field:  rate.Field,
			Amount: rate.PerSecond * dt,
		})
	}
	for _, pulse := range feedback.Pulses {
		applyDelta(s, Delta{
			Field:  pulse.Field,
			Amount: pulse.Amount,
		})
	}
	ClampAll(s)
}
