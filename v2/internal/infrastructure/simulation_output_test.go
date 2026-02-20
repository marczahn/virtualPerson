package infrastructure_test

import (
	"reflect"
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/consciousness"
	"github.com/marczahn/person/v2/internal/infrastructure"
	"github.com/marczahn/person/v2/internal/motivation"
)

func TestBuildTaggedOutputLines_EmitsBIOAndMINDAndSignificantDRIVES(t *testing.T) {
	previous := motivation.MotivationState{
		EnergyUrgency:      0.40,
		SocialUrgency:      0.30,
		StimulationUrgency: 0.30,
		SafetyUrgency:      0.30,
		IdentityUrgency:    0.30,
	}
	result := infrastructure.TickResult{
		Bio: biology.TickResult{},
		Motivation: motivation.MotivationState{
			EnergyUrgency:      0.58,
			SocialUrgency:      0.35,
			StimulationUrgency: 0.31,
			SafetyUrgency:      0.31,
			IdentityUrgency:    0.30,
		},
		Parsed: consciousness.ParsedResponse{
			Narrative: "I should eat now.",
		},
	}

	got := infrastructure.BuildTaggedOutputLines(result, previous, 0.15)
	want := []string{
		"[BIO] no significant biological deltas",
		"[DRIVES] energy: 0.40 -> 0.58 (+0.18)",
		"[MIND] I should eat now.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected tagged output lines: got=%v want=%v", got, want)
	}
}
