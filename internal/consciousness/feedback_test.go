package consciousness

import (
	"strings"
	"testing"

	"github.com/marczahn/person/internal/biology"
)

func TestParseFeedback_DetectsRumination(t *testing.T) {
	output := "I keep thinking about what happened. I can't stop thinking about it, replaying it over and over."

	fb, _ := ParseFeedback(output)

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

	fb, _ := ParseFeedback(output)

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

	fb, _ := ParseFeedback(output)

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

	fb, _ := ParseFeedback(output)

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

	fb, _ := ParseFeedback(output)

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

	fb, _ := ParseFeedback(output)

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

	fb, _ := ParseFeedback(output)

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

// --- ParseEmotionalTag tests ---

func TestParseEmotionalTag_Valid(t *testing.T) {
	input := "I feel the heat rising in me.\n[STATE: arousal=0.8, valence=-0.7]"
	tag, content := ParseEmotionalTag(input)

	if tag.Arousal != 0.8 {
		t.Errorf("arousal = %v, expected 0.8", tag.Arousal)
	}
	if tag.Valence != -0.7 {
		t.Errorf("valence = %v, expected -0.7", tag.Valence)
	}
	if content != "I feel the heat rising in me." {
		t.Errorf("content = %q, expected stripped content", content)
	}
}

func TestParseEmotionalTag_NoTag(t *testing.T) {
	input := "I feel calm."
	tag, content := ParseEmotionalTag(input)

	if tag.Arousal != 0.0 || tag.Valence != 0.0 {
		t.Errorf("expected zero tag, got {%v, %v}", tag.Arousal, tag.Valence)
	}
	if content != "I feel calm." {
		t.Errorf("content = %q, expected original unchanged", content)
	}
}

func TestParseEmotionalTag_MalformedTag(t *testing.T) {
	input := "Something terrible.\n[STATE: arousal=abc]"
	tag, _ := ParseEmotionalTag(input)

	if tag.Arousal != 0.0 || tag.Valence != 0.0 {
		t.Errorf("malformed tag should produce zero values, got {%v, %v}", tag.Arousal, tag.Valence)
	}
}

func TestParseEmotionalTag_TagStrippedFromContent(t *testing.T) {
	input := "Anger floods through me.\n[STATE: arousal=0.9, valence=-0.8]"
	_, content := ParseEmotionalTag(input)

	if strings.Contains(content, "[STATE:") {
		t.Errorf("content still contains [STATE: tag: %q", content)
	}
}

// --- EmotionalPulses tests ---

func TestEmotionalPulses_AngryThought(t *testing.T) {
	tag := EmotionalTag{Arousal: 0.8, Valence: -0.7}
	changes := EmotionalPulses(tag)

	var cortisol, tension, hr, serotonin, adrenaline float64
	for _, c := range changes {
		switch c.Variable {
		case biology.VarCortisol:
			cortisol += c.Delta
		case biology.VarMuscleTension:
			tension += c.Delta
		case biology.VarHeartRate:
			hr += c.Delta
		case biology.VarSerotonin:
			serotonin += c.Delta
		case biology.VarAdrenaline:
			adrenaline += c.Delta
		}
	}

	if cortisol <= 0 {
		t.Errorf("angry thought should increase cortisol, got %v", cortisol)
	}
	if tension <= 0 {
		t.Errorf("angry thought should increase muscle tension, got %v", tension)
	}
	if hr <= 0 {
		t.Errorf("angry thought should increase heart rate, got %v", hr)
	}
	if serotonin >= 0 {
		t.Errorf("angry thought should decrease serotonin, got %v", serotonin)
	}
	if adrenaline != 0 {
		t.Errorf("non-extreme angry thought should have zero adrenaline, got %v", adrenaline)
	}
}

func TestEmotionalPulses_ExtremeFearthreshold(t *testing.T) {
	tag := EmotionalTag{Arousal: 0.9, Valence: -0.8}
	changes := EmotionalPulses(tag)

	var adrenaline float64
	for _, c := range changes {
		if c.Variable == biology.VarAdrenaline {
			adrenaline += c.Delta
		}
	}

	if adrenaline <= 0 {
		t.Errorf("extreme fear should produce adrenaline, got %v", adrenaline)
	}
}

func TestEmotionalPulses_PositiveCalm(t *testing.T) {
	tag := EmotionalTag{Arousal: 0.2, Valence: 0.6}
	changes := EmotionalPulses(tag)

	var cortisol, hr, serotonin, dopamine float64
	for _, c := range changes {
		switch c.Variable {
		case biology.VarCortisol:
			cortisol += c.Delta
		case biology.VarHeartRate:
			hr += c.Delta
		case biology.VarSerotonin:
			serotonin += c.Delta
		case biology.VarDopamine:
			dopamine += c.Delta
		}
	}

	if cortisol != 0 {
		t.Errorf("positive calm should have zero cortisol delta, got %v", cortisol)
	}
	if hr != 0 {
		t.Errorf("positive calm (arousal=0.2) should have zero HR delta, got %v", hr)
	}
	if serotonin <= 0 {
		t.Errorf("positive valence should increase serotonin, got %v", serotonin)
	}
}

func TestEmotionalPulses_ZeroTag_ReturnsNil(t *testing.T) {
	tag := EmotionalTag{Arousal: 0, Valence: 0}
	changes := EmotionalPulses(tag)

	if changes != nil {
		t.Errorf("zero tag should return nil, got %v", changes)
	}
}

func TestEmotionalPulses_AdrenalineSafeForNonExtreme(t *testing.T) {
	// arousal=0.8, valence=-0.7: does NOT cross extreme threshold (arousal>0.8 && negVal>0.5)
	tag := EmotionalTag{Arousal: 0.8, Valence: -0.7}
	changes := EmotionalPulses(tag)

	for _, c := range changes {
		if c.Variable == biology.VarAdrenaline && c.Delta > 0 {
			t.Errorf("non-extreme state (arousal=0.8) should not trigger adrenaline, got delta=%v", c.Delta)
		}
	}
}
