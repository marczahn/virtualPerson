package motivation_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/motivation"
)

func baselineBio() biology.State {
	return biology.State{
		Energy:            0.8,
		Stress:            0.1,
		CognitiveCapacity: 1.0,
		Mood:              0.5,
		PhysicalTension:   0.05,
		Hunger:            0.1,
		SocialDeficit:     0.0,
		BodyTemp:          36.6,
	}
}

func baselinePersonality() motivation.Personality {
	return motivation.Personality{
		StressSensitivity:    0.5,
		EnergyResilience:     0.5,
		Curiosity:            0.5,
		SelfObservation:      0.5,
		FrustrationTolerance: 0.5,
		RiskAversion:         0.5,
		SocialFactor:         0.5,
	}
}

func TestCompute_IsDeterministic(t *testing.T) {
	b := baselineBio()
	p := baselinePersonality()
	c := motivation.ChronicState{}

	first := motivation.Compute(b, p, c)
	second := motivation.Compute(b, p, c)

	if first != second {
		t.Fatalf("compute must be deterministic: first=%+v second=%+v", first, second)
	}
}

func TestCompute_DrivesAreClamped(t *testing.T) {
	b := biology.State{
		Energy:            -10,
		Stress:            10,
		CognitiveCapacity: -10,
		Mood:              -10,
		PhysicalTension:   10,
		Hunger:            10,
		SocialDeficit:     10,
		BodyTemp:          100,
	}
	p := motivation.Personality{
		StressSensitivity:    1,
		EnergyResilience:     0,
		Curiosity:            1,
		SelfObservation:      1,
		FrustrationTolerance: 0,
		RiskAversion:         1,
		SocialFactor:         1,
	}
	m := motivation.Compute(b, p, motivation.ChronicState{})

	for name, v := range map[string]float64{
		"energy":    m.EnergyUrgency,
		"social":    m.SocialUrgency,
		"stim":      m.StimulationUrgency,
		"safety":    m.SafetyUrgency,
		"identity":  m.IdentityUrgency,
		"goalValue": m.ActiveGoalUrgency,
	} {
		if v < 0 || v > 1 {
			t.Fatalf("%s urgency out of range: %f", name, v)
		}
	}
}

func TestCompute_MonotonicEnergyDrive(t *testing.T) {
	p := baselinePersonality()
	c := motivation.ChronicState{}

	b1 := baselineBio()
	b2 := baselineBio()
	b1.Energy = 0.8
	b2.Energy = 0.2

	m1 := motivation.Compute(b1, p, c)
	m2 := motivation.Compute(b2, p, c)
	if m2.EnergyUrgency <= m1.EnergyUrgency {
		t.Fatalf("lower energy should increase urgency: high=%f low=%f", m1.EnergyUrgency, m2.EnergyUrgency)
	}
}

func TestCompute_MonotonicSocialDrive(t *testing.T) {
	p := baselinePersonality()
	c := motivation.ChronicState{}

	b1 := baselineBio()
	b2 := baselineBio()
	b1.SocialDeficit = 0.1
	b2.SocialDeficit = 0.9

	m1 := motivation.Compute(b1, p, c)
	m2 := motivation.Compute(b2, p, c)
	if m2.SocialUrgency <= m1.SocialUrgency {
		t.Fatalf("higher social deficit should increase urgency: low=%f high=%f", m1.SocialUrgency, m2.SocialUrgency)
	}
}

func TestCompute_PersonalityModulatesTargetDrive(t *testing.T) {
	b := baselineBio()
	c := motivation.ChronicState{}

	lowCuriosity := baselinePersonality()
	highCuriosity := baselinePersonality()
	lowCuriosity.Curiosity = 0.0
	highCuriosity.Curiosity = 1.0

	mLow := motivation.Compute(b, lowCuriosity, c)
	mHigh := motivation.Compute(b, highCuriosity, c)

	if mHigh.StimulationUrgency <= mLow.StimulationUrgency {
		t.Fatalf("higher curiosity should raise stimulation urgency: low=%f high=%f", mLow.StimulationUrgency, mHigh.StimulationUrgency)
	}
	if mHigh.EnergyUrgency != mLow.EnergyUrgency {
		t.Fatalf("curiosity should not modulate energy urgency: low=%f high=%f", mLow.EnergyUrgency, mHigh.EnergyUrgency)
	}
}

func TestCompute_TieBreakOrderIsDeterministic(t *testing.T) {
	b := baselineBio()
	p := baselinePersonality()
	c := motivation.ChronicState{}

	// Construct a state that makes energy and safety equal and maximal.
	b.Energy = 0.0
	b.Hunger = 1.0
	b.Stress = 1.0
	b.PhysicalTension = 1.0
	b.BodyTemp = 42.6

	m := motivation.Compute(b, p, c)
	if m.ActiveGoalDrive != motivation.DriveSafety {
		t.Fatalf("expected safety to win deterministic tie-break, got %s", m.ActiveGoalDrive)
	}
}

func TestActionCandidates_RespectConstraints(t *testing.T) {
	c := motivation.ActionConstraints{
		HasFood:         false,
		HasPeopleNearby: true,
		CanRest:         true,
		CanExplore:      false,
		HasQuietSpace:   true,
	}

	actions := motivation.ActionCandidatesFor(motivation.DriveEnergy, c)
	if len(actions) == 0 {
		t.Fatal("expected at least one action for energy drive")
	}
	for _, a := range actions {
		if a == motivation.ActionEat {
			t.Fatal("eat should not be emitted when HasFood is false")
		}
	}
}
