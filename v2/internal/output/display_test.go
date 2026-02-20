package output

import (
	"reflect"
	"testing"

	"github.com/marczahn/person/v2/internal/motivation"
)

func TestSignificantDriveChanges_FiltersByThreshold(t *testing.T) {
	previous := motivation.MotivationState{
		EnergyUrgency:      0.40,
		SocialUrgency:      0.20,
		StimulationUrgency: 0.20,
		SafetyUrgency:      0.30,
		IdentityUrgency:    0.30,
	}
	current := motivation.MotivationState{
		EnergyUrgency:      0.57,
		SocialUrgency:      0.25,
		StimulationUrgency: 0.37,
		SafetyUrgency:      0.31,
		IdentityUrgency:    0.10,
	}

	changes := SignificantDriveChanges(previous, current, 0.15)

	want := []DriveChange{
		{Name: motivation.DriveEnergy, Previous: 0.40, Current: 0.57},
		{Name: motivation.DriveStimulation, Previous: 0.20, Current: 0.37},
		{Name: motivation.DriveIdentityCoherence, Previous: 0.30, Current: 0.10},
	}
	if !reflect.DeepEqual(changes, want) {
		t.Fatalf("unexpected significant changes: got=%v want=%v", changes, want)
	}
}

func TestFormatDriveChangeLines_DeterministicOutput(t *testing.T) {
	changes := []DriveChange{
		{Name: motivation.DriveEnergy, Previous: 0.40, Current: 0.57},
		{Name: motivation.DriveIdentityCoherence, Previous: 0.30, Current: 0.10},
	}

	got := FormatDriveChangeLines(changes)
	want := []string{
		"[DRIVES] energy: 0.40 -> 0.57 (+0.17)",
		"[DRIVES] identity_coherence: 0.30 -> 0.10 (-0.20)",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected formatted drive lines: got=%v want=%v", got, want)
	}
}
