package biology_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestInteractions_StressCausesPhysicalTension(t *testing.T) {
	s := biology.NewDefaultState()
	s.Stress = 0.7 // above 0.6 threshold for rule 1
	origTension := s.PhysicalTension

	deltas := biology.ApplyInteractions(s, 1.0)

	if s.PhysicalTension <= origTension {
		t.Errorf("PhysicalTension should increase when Stress=0.7, got %v (was %v)",
			s.PhysicalTension, origTension)
	}

	// Verify the delta was returned
	found := false
	for _, d := range deltas {
		if d.Field == "physical_tension" && d.Amount > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected a positive physical_tension delta in returned deltas, got %v", deltas)
	}
}

func TestInteractions_StressDepletesCognitiveCapacity(t *testing.T) {
	s := biology.NewDefaultState()
	s.Stress = 0.6 // above 0.5 threshold for rule 2
	origCog := s.CognitiveCapacity

	biology.ApplyInteractions(s, 1.0)

	if s.CognitiveCapacity >= origCog {
		t.Errorf("CognitiveCapacity should decrease when Stress=0.6, got %v (was %v)",
			s.CognitiveCapacity, origCog)
	}
}

func TestInteractions_LowEnergyDepressesMood(t *testing.T) {
	s := biology.NewDefaultState()
	s.Energy = 0.2 // below 0.3 threshold for rule 6
	origMood := s.Mood

	biology.ApplyInteractions(s, 1.0)

	if s.Mood >= origMood {
		t.Errorf("Mood should decrease when Energy=0.2, got %v (was %v)",
			s.Mood, origMood)
	}
}

func TestInteractions_CompoundEnergyHungerMoodDrop(t *testing.T) {
	// Rule 21: Energy<0.4 AND Hunger>0.6 causes additional Mood drop
	s := biology.NewDefaultState()
	s.Energy = 0.35
	s.Hunger = 0.65
	origMood := s.Mood

	biology.ApplyInteractions(s, 1.0)

	// Both rule 6 (Energy<0.3 — NOT triggered at 0.35) and rule 21 (compound) should apply.
	// At Energy=0.35, rule 6 (Energy<0.3) does NOT fire, but rule 21 does.
	moodDeltas := 0.0
	for _, d := range biology.ApplyInteractions(biology.NewDefaultState(), 0.0) {
		_ = d // just to use deltas
	}

	// Reset and test properly
	s2 := biology.NewDefaultState()
	s2.Energy = 0.35
	s2.Hunger = 0.65
	deltas := biology.ApplyInteractions(s2, 1.0)

	for _, d := range deltas {
		if d.Field == "mood" {
			moodDeltas += d.Amount
		}
	}

	_ = origMood
	if moodDeltas >= 0 {
		t.Errorf("Expected negative mood delta from compound rule, total mood delta = %v", moodDeltas)
	}
	if s2.Mood >= 0.5 {
		t.Errorf("Mood should have decreased from 0.5, got %v", s2.Mood)
	}
}

func TestInteractions_SnapshotPreventsExplosion(t *testing.T) {
	// Rule 1: Stress>0.6 -> PhysicalTension increases
	// Rule 9: PhysicalTension>0.7 -> Stress increases
	// If Stress=0.65 and PhysicalTension=0.65, rule 1 fires and adds tension.
	// But rule 9 should NOT fire because it evaluates against the SNAPSHOT where
	// PhysicalTension was 0.65 (below 0.7), even though rule 1 already increased it.
	s := biology.NewDefaultState()
	s.Stress = 0.65
	s.PhysicalTension = 0.65
	origStress := s.Stress

	deltas := biology.ApplyInteractions(s, 1.0)

	// PhysicalTension should have increased (rule 1 fires)
	if s.PhysicalTension <= 0.65 {
		t.Errorf("PhysicalTension should increase from rule 1, got %v", s.PhysicalTension)
	}

	// But the Stress feedback from PhysicalTension>0.7 (rule 9) should NOT have fired
	// because snapshot had PhysicalTension=0.65
	stressDeltaFromTension := false
	for _, d := range deltas {
		if d.Field == "stress" {
			stressDeltaFromTension = true
		}
	}

	// At Stress=0.65, several stress-adding rules could fire:
	// rule 2 (Stress>0.5) fires -> cognitive_capacity delta (not stress)
	// rule 9 (PhysicalTension>0.7) should NOT fire from snapshot
	// No other rule should add stress at these values
	if stressDeltaFromTension {
		// Check if stress actually changed — it shouldn't from tension feedback
		// (mood->stress and other stress sources shouldn't be active here either)
		t.Errorf("Stress should not change from PhysicalTension feedback when snapshot PhysicalTension=0.65, "+
			"but stress deltas found. Stress went from %v to %v", origStress, s.Stress)
	}
}

func TestInteractions_DtZero_AllDeltasZero(t *testing.T) {
	// With dt=0, all rule apply functions multiply by dt, so all deltas should be 0
	s := biology.NewDefaultState()
	s.Stress = 0.9    // triggers multiple rules
	s.Energy = 0.1    // triggers multiple rules
	s.Hunger = 0.9    // triggers multiple rules
	s.Mood = 0.1      // triggers mood rules
	s.BodyTemp = 30.0 // triggers body temp rules
	before := *s

	deltas := biology.ApplyInteractions(s, 0.0)

	for _, d := range deltas {
		if d.Amount != 0.0 {
			t.Errorf("dt=0: expected zero delta, got %v for field %s", d.Amount, d.Field)
		}
	}

	// State should be unchanged
	if s.Energy != before.Energy {
		t.Errorf("dt=0: Energy changed from %v to %v", before.Energy, s.Energy)
	}
	if s.Stress != before.Stress {
		t.Errorf("dt=0: Stress changed from %v to %v", before.Stress, s.Stress)
	}
	if s.Mood != before.Mood {
		t.Errorf("dt=0: Mood changed from %v to %v", before.Mood, s.Mood)
	}
}
