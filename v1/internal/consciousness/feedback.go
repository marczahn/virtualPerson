package consciousness

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/i18n"
)

// stateTagRe matches [STATE: arousal=X, valence=Y] at the end of LLM output.
// Both fields are required for a valid parse.
var stateTagRe = regexp.MustCompile(`\[STATE:\s*arousal=([-\d.]+),\s*valence=([-\d.]+)\]`)

// ParseEmotionalTag extracts a [STATE: arousal=X, valence=Y] annotation from
// the LLM output. It returns the parsed tag and the content with the tag
// stripped and trimmed. If the tag is absent or malformed, zero values are
// returned and the original content is unchanged.
func ParseEmotionalTag(output string) (EmotionalTag, string) {
	m := stateTagRe.FindStringIndex(output)
	if m == nil {
		return EmotionalTag{}, output
	}

	match := stateTagRe.FindStringSubmatch(output)
	arousal, err1 := strconv.ParseFloat(match[1], 64)
	valence, err2 := strconv.ParseFloat(match[2], 64)
	if err1 != nil || err2 != nil {
		return EmotionalTag{}, output
	}

	clean := strings.TrimSpace(output[:m[0]])
	return EmotionalTag{Arousal: arousal, Valence: valence}, clean
}

// ParseFeedback analyzes a consciousness output (the LLM's response) and
// extracts coping strategies and distortions that should feed back into
// the biological state. It also strips the [STATE:] annotation and returns
// the cleaned content as the second return value.
//
// The keyword detection is a heuristic. The EmotionalState field is populated
// from the machine-readable annotation.
func ParseFeedback(output string) (ThoughtFeedback, string) {
	tag, clean := ParseEmotionalTag(output)

	lower := strings.ToLower(clean)
	tr := i18n.T()

	var fb ThoughtFeedback
	fb.EmotionalState = tag

	// Detect distortions.
	for name, keywords := range tr.Feedback.Distortions {
		if containsAny(lower, keywords...) {
			fb.ActiveDistortions = append(fb.ActiveDistortions, name)
		}
	}

	// Detect coping strategies.
	for name, keywords := range tr.Feedback.Coping {
		if containsAny(lower, keywords...) {
			fb.ActiveCoping = append(fb.ActiveCoping, name)
		}
	}

	return fb, clean
}

// EmotionalPulses converts an EmotionalTag into calibrated one-time absolute
// biological state changes. Returns nil when the tag is zero (no annotation).
//
// Pulse magnitudes are chosen so that 5 consecutive angry thoughts bring
// cortisol above the 0.3 tachycardia threshold, but a single thought does not.
// Extreme fear (arousal>0.8 && negVal>0.5) triggers adrenaline — kept safely
// below the 12 bpm/s cascade threshold.
func EmotionalPulses(tag EmotionalTag) []biology.StateChange {
	if tag.Arousal == 0 && tag.Valence == 0 {
		return nil
	}

	negVal := max(0, -tag.Valence) // 0..1 for negative valence
	posVal := max(0, tag.Valence)  // 0..1 for positive valence

	var changes []biology.StateChange

	// Cortisol: gated on negative valence — pure positive arousal doesn't raise it.
	// 5 angry thoughts (arousal=0.8, negVal=0.7) raise cortisol by ~0.245 from baseline 0.1 → above 0.3 tachycardia threshold.
	cortisol := negVal * (tag.Arousal*0.05 + 0.03)
	if cortisol > 0 {
		changes = append(changes, biology.StateChange{
			Variable: biology.VarCortisol,
			Delta:    cortisol,
			Source:   "consciousness_emotional",
		})
	}

	// Muscle tension: scales with arousal.
	if tag.Arousal > 0 {
		changes = append(changes, biology.StateChange{
			Variable: biology.VarMuscleTension,
			Delta:    tag.Arousal * 0.15,
			Source:   "consciousness_emotional",
		})
	}

	// Heart rate: only above arousal threshold of 0.3.
	if tag.Arousal > 0.3 {
		changes = append(changes, biology.StateChange{
			Variable: biology.VarHeartRate,
			Delta:    tag.Arousal * 10 * (1 + negVal*0.5),
			Source:   "consciousness_emotional",
		})
	}

	// Serotonin: negative valence drains it; positive valence restores it (both terms additive).
	serotonin := tag.Valence*0.025 + posVal*0.04
	if serotonin != 0 {
		changes = append(changes, biology.StateChange{
			Variable: biology.VarSerotonin,
			Delta:    serotonin,
			Source:   "consciousness_emotional",
		})
	}

	// Dopamine: positive excited state produces a reward signal.
	if posVal > 0.2 {
		changes = append(changes, biology.StateChange{
			Variable: biology.VarDopamine,
			Delta:    posVal * tag.Arousal * 0.1,
			Source:   "consciousness_emotional",
		})
	}

	// Adrenaline: extreme fear only (arousal>0.8 AND negVal>0.5).
	if tag.Arousal > 0.8 && negVal > 0.5 {
		changes = append(changes, biology.StateChange{
			Variable: biology.VarAdrenaline,
			Delta:    (tag.Arousal - 0.8) * (negVal - 0.5) * 0.5,
			Source:   "consciousness_emotional",
		})
	}

	return changes
}

// FeedbackToChanges converts parsed thought feedback into biological state
// changes. Applied as absolute one-time pulses (not dt-scaled).
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
