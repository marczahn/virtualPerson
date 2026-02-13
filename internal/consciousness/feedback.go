package consciousness

import (
	"strings"

	"github.com/marczahn/person/internal/biology"
)

// ParseFeedback analyzes a consciousness output (the LLM's response) and
// extracts coping strategies and distortions that should feed back into
// the biological state.
//
// This is a keyword-based heuristic. The LLM's output is scanned for
// indicators of specific psychological patterns.
func ParseFeedback(output string) ThoughtFeedback {
	lower := strings.ToLower(output)

	var fb ThoughtFeedback

	// Detect coping strategies.
	if containsAny(lower, "what if everything", "worst case", "disaster", "catastroph", "terrible", "going to die", "never going to") {
		fb.ActiveDistortions = append(fb.ActiveDistortions, "catastrophizing")
	}
	if containsAny(lower, "this always happens", "it always", "i always fail", "nothing ever works", "every time this", "it never gets") {
		fb.ActiveDistortions = append(fb.ActiveDistortions, "overgeneralization")
	}
	if containsAny(lower, "i feel like it must be", "feels true", "because i feel") {
		fb.ActiveDistortions = append(fb.ActiveDistortions, "emotional_reasoning")
	}
	if containsAny(lower, "my fault", "i caused", "blame myself", "i should have") {
		fb.ActiveDistortions = append(fb.ActiveDistortions, "personalization")
	}
	if containsAny(lower, "they must think", "they probably think", "judging me") {
		fb.ActiveDistortions = append(fb.ActiveDistortions, "mind_reading")
	}

	// Detect coping strategies.
	if containsAny(lower, "keep thinking about", "can't stop thinking", "going over and over", "why did", "replaying") {
		fb.ActiveCoping = append(fb.ActiveCoping, "rumination")
	}
	if containsAny(lower, "it's okay", "accept", "let it go", "can't change", "it is what it is") {
		fb.ActiveCoping = append(fb.ActiveCoping, "acceptance")
	}
	if containsAny(lower, "maybe it's not that bad", "another way to look", "reframe", "on the other hand", "could also mean") {
		fb.ActiveCoping = append(fb.ActiveCoping, "reappraisal")
	}
	if containsAny(lower, "think about something else", "distract myself", "focus on something else", "try not to think about") {
		fb.ActiveCoping = append(fb.ActiveCoping, "distraction")
	}
	if containsAny(lower, "figure this out", "need a plan", "solve this", "what can i do", "steps to", "how do i fix") {
		fb.ActiveCoping = append(fb.ActiveCoping, "problem_solving")
	}
	if containsAny(lower, "push it down", "don't feel", "ignore it", "bury") {
		fb.ActiveCoping = append(fb.ActiveCoping, "suppression")
	}
	if containsAny(lower, "it's fine", "nothing's wrong", "not happening", "refuse to believe") {
		fb.ActiveCoping = append(fb.ActiveCoping, "denial")
	}

	return fb
}

// FeedbackToChanges converts parsed thought feedback into biological state
// changes. These are per-second rates to be multiplied by the cycle duration.
func FeedbackToChanges(fb ThoughtFeedback) []biology.StateChange {
	var changes []biology.StateChange

	for _, coping := range fb.ActiveCoping {
		switch coping {
		case "rumination":
			changes = append(changes,
				biology.StateChange{Variable: biology.VarCortisol, Delta: 0.02, Source: "consciousness_rumination"},
				biology.StateChange{Variable: biology.VarSerotonin, Delta: -0.01, Source: "consciousness_rumination"},
			)
		case "acceptance", "reappraisal":
			changes = append(changes,
				biology.StateChange{Variable: biology.VarCortisol, Delta: -0.01, Source: "consciousness_acceptance"},
				biology.StateChange{Variable: biology.VarSerotonin, Delta: 0.005, Source: "consciousness_acceptance"},
			)
		case "suppression":
			changes = append(changes,
				biology.StateChange{Variable: biology.VarCortisol, Delta: 0.01, Source: "consciousness_suppression"},
				biology.StateChange{Variable: biology.VarAdrenaline, Delta: 0.005, Source: "consciousness_suppression"},
			)
		}
	}

	for _, distortion := range fb.ActiveDistortions {
		switch distortion {
		case "catastrophizing":
			changes = append(changes,
				biology.StateChange{Variable: biology.VarAdrenaline, Delta: 0.03, Source: "consciousness_catastrophizing"},
				biology.StateChange{Variable: biology.VarCortisol, Delta: 0.02, Source: "consciousness_catastrophizing"},
			)
		}
	}

	// Multiple distortions compound stress.
	if len(fb.ActiveDistortions) > 2 {
		changes = append(changes,
			biology.StateChange{Variable: biology.VarCortisol, Delta: 0.01, Source: "consciousness_distortion_load"},
		)
	}

	return changes
}

func containsAny(s string, patterns ...string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
