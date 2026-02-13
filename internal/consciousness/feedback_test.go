package consciousness

import (
	"testing"

	"github.com/marczahn/person/internal/biology"
)

func TestParseFeedback_DetectsRumination(t *testing.T) {
	output := "I keep thinking about what happened. I can't stop thinking about it, replaying it over and over."

	fb := ParseFeedback(output)

	found := false
	for _, c := range fb.ActiveCoping {
		if c == "rumination" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected rumination detection, got coping=%v", fb.ActiveCoping)
	}
}

func TestParseFeedback_DetectsAcceptance(t *testing.T) {
	output := "I guess I need to accept this. It is what it is. I can't change what happened."

	fb := ParseFeedback(output)

	found := false
	for _, c := range fb.ActiveCoping {
		if c == "acceptance" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected acceptance detection, got coping=%v", fb.ActiveCoping)
	}
}

func TestParseFeedback_DetectsCatastrophizing(t *testing.T) {
	output := "This is a disaster. What if everything falls apart? I'm going to die out here."

	fb := ParseFeedback(output)

	found := false
	for _, d := range fb.ActiveDistortions {
		if d == "catastrophizing" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected catastrophizing detection, got distortions=%v", fb.ActiveDistortions)
	}
}

func TestParseFeedback_DetectsPersonalization(t *testing.T) {
	output := "It's my fault. I caused this. I should have been more careful."

	fb := ParseFeedback(output)

	found := false
	for _, d := range fb.ActiveDistortions {
		if d == "personalization" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected personalization detection, got distortions=%v", fb.ActiveDistortions)
	}
}

func TestParseFeedback_DetectsReappraisal(t *testing.T) {
	output := "Maybe it's not that bad. On the other hand, this could also mean something positive."

	fb := ParseFeedback(output)

	found := false
	for _, c := range fb.ActiveCoping {
		if c == "reappraisal" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected reappraisal detection, got coping=%v", fb.ActiveCoping)
	}
}

func TestParseFeedback_DetectsProblemSolving(t *testing.T) {
	output := "What can I do about this? Let me figure this out. I need a plan."

	fb := ParseFeedback(output)

	found := false
	for _, c := range fb.ActiveCoping {
		if c == "problem_solving" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected problem_solving detection, got coping=%v", fb.ActiveCoping)
	}
}

func TestParseFeedback_NeutralOutput_NoDetection(t *testing.T) {
	output := "The air feels cool. I notice the quiet around me."

	fb := ParseFeedback(output)

	if len(fb.ActiveCoping) != 0 {
		t.Errorf("expected no coping for neutral output, got %v", fb.ActiveCoping)
	}
	if len(fb.ActiveDistortions) != 0 {
		t.Errorf("expected no distortions for neutral output, got %v", fb.ActiveDistortions)
	}
}

func TestFeedbackToChanges_Rumination(t *testing.T) {
	fb := ThoughtFeedback{ActiveCoping: []string{"rumination"}}
	changes := FeedbackToChanges(fb)

	var hasCortisol, hasSerotonin bool
	for _, c := range changes {
		if c.Variable == biology.VarCortisol && c.Delta > 0 {
			hasCortisol = true
		}
		if c.Variable == biology.VarSerotonin && c.Delta < 0 {
			hasSerotonin = true
		}
	}
	if !hasCortisol {
		t.Error("rumination should increase cortisol")
	}
	if !hasSerotonin {
		t.Error("rumination should decrease serotonin")
	}
}

func TestFeedbackToChanges_Acceptance(t *testing.T) {
	fb := ThoughtFeedback{ActiveCoping: []string{"acceptance"}}
	changes := FeedbackToChanges(fb)

	var hasCortisol, hasSerotonin bool
	for _, c := range changes {
		if c.Variable == biology.VarCortisol && c.Delta < 0 {
			hasCortisol = true
		}
		if c.Variable == biology.VarSerotonin && c.Delta > 0 {
			hasSerotonin = true
		}
	}
	if !hasCortisol {
		t.Error("acceptance should decrease cortisol")
	}
	if !hasSerotonin {
		t.Error("acceptance should increase serotonin")
	}
}

func TestFeedbackToChanges_Catastrophizing(t *testing.T) {
	fb := ThoughtFeedback{ActiveDistortions: []string{"catastrophizing"}}
	changes := FeedbackToChanges(fb)

	var hasAdrenaline, hasCortisol bool
	for _, c := range changes {
		if c.Variable == biology.VarAdrenaline && c.Delta > 0 {
			hasAdrenaline = true
		}
		if c.Variable == biology.VarCortisol && c.Delta > 0 {
			hasCortisol = true
		}
	}
	if !hasAdrenaline {
		t.Error("catastrophizing should increase adrenaline")
	}
	if !hasCortisol {
		t.Error("catastrophizing should increase cortisol")
	}
}

func TestFeedbackToChanges_Suppression(t *testing.T) {
	fb := ThoughtFeedback{ActiveCoping: []string{"suppression"}}
	changes := FeedbackToChanges(fb)

	var hasCortisol bool
	for _, c := range changes {
		if c.Variable == biology.VarCortisol && c.Delta > 0 {
			hasCortisol = true
		}
	}
	if !hasCortisol {
		t.Error("suppression should increase cortisol")
	}
}

func TestFeedbackToChanges_MultipleDistortions_CompoundStress(t *testing.T) {
	fb := ThoughtFeedback{ActiveDistortions: []string{"catastrophizing", "overgeneralization", "personalization"}}
	changes := FeedbackToChanges(fb)

	found := false
	for _, c := range changes {
		if c.Source == "consciousness_distortion_load" {
			found = true
		}
	}
	if !found {
		t.Error("3+ distortions should add compound cortisol")
	}
}

func TestFeedbackToChanges_Empty(t *testing.T) {
	fb := ThoughtFeedback{}
	changes := FeedbackToChanges(fb)

	if len(changes) != 0 {
		t.Errorf("expected no changes for empty feedback, got %d", len(changes))
	}
}
