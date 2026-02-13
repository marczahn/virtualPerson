package sense

import (
	"strings"
	"time"
)

// Parser converts raw text input into sensory events.
type Parser interface {
	Parse(input string) []Event
}

// KeywordParser classifies text into sensory events using keyword matching.
// Deterministic and fast — no external API calls.
type KeywordParser struct{}

// NewKeywordParser creates a keyword-based sensory parser.
func NewKeywordParser() *KeywordParser {
	return &KeywordParser{}
}

// Parse analyzes raw text and returns zero or more sensory events.
// A single input can produce multiple events (e.g., "a freezing, painful wind"
// produces both thermal and pain events).
func (p *KeywordParser) Parse(input string) []Event {
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "" {
		return nil
	}

	now := time.Now()
	var events []Event

	for _, rule := range keywordRules {
		if matched, intensity := rule.match(lower); matched {
			events = append(events, Event{
				Channel:   rule.channel,
				Intensity: intensity,
				RawInput:  input,
				Parsed:    rule.description(lower),
				Timestamp: now,
			})
		}
	}

	// If nothing matched, treat as a generic auditory event (someone speaking).
	if len(events) == 0 {
		events = append(events, Event{
			Channel:   Auditory,
			Intensity: 0.3,
			RawInput:  input,
			Parsed:    "hearing speech or ambient sound",
			Timestamp: now,
		})
	}

	return events
}

type keywordRule struct {
	channel     Channel
	keywords    []keywordEntry
	description func(string) string
}

type keywordEntry struct {
	phrase    string
	intensity float64
}

func (r *keywordRule) match(lower string) (bool, float64) {
	var maxIntensity float64
	matched := false
	for _, kw := range r.keywords {
		if strings.Contains(lower, kw.phrase) {
			matched = true
			if kw.intensity > maxIntensity {
				maxIntensity = kw.intensity
			}
		}
	}
	return matched, maxIntensity
}

// keywordRules defines the mapping from text patterns to sensory channels.
// Rules are evaluated independently — a single input can match multiple rules.
var keywordRules = []keywordRule{
	// Thermal: cold stimuli (intensity < 0.5 means cold in the biology processor)
	{
		channel: Thermal,
		keywords: []keywordEntry{
			{"freezing", 0.1},
			{"frozen", 0.1},
			{"ice cold", 0.1},
			{"bitter cold", 0.1},
			{"hypothermi", 0.05},
			{"frost", 0.15},
			{"cold", 0.25},
			{"chilly", 0.3},
			{"cool", 0.35},
			{"snow", 0.2},
			{"blizzard", 0.1},
			// Heat stimuli (intensity > 0.5 means hot)
			{"burning", 0.95},
			{"scalding", 0.95},
			{"on fire", 0.95},
			{"boiling", 0.9},
			{"scorching", 0.9},
			{"hot", 0.75},
			{"warm", 0.6},
			{"heat", 0.7},
			{"fever", 0.7},
			{"sweltering", 0.85},
		},
		description: func(s string) string {
			if containsAny(s, "freezing", "frozen", "ice cold", "bitter cold", "cold", "chilly", "cool", "snow", "blizzard", "frost", "hypothermi") {
				return "feeling cold"
			}
			return "feeling heat"
		},
	},
	// Pain
	{
		channel: Pain,
		keywords: []keywordEntry{
			{"agony", 0.95},
			{"excruciating", 0.95},
			{"stabbing pain", 0.9},
			{"sharp pain", 0.85},
			{"intense pain", 0.85},
			{"severe pain", 0.8},
			{"throbbing", 0.6},
			{"aching", 0.4},
			{"sore", 0.35},
			{"pain", 0.5},
			{"hurt", 0.5},
			{"sting", 0.4},
			{"burn", 0.5}, // pain sense of burning, distinct from thermal
			{"cramp", 0.5},
			{"headache", 0.5},
			{"migraine", 0.7},
			{"broken bone", 0.85},
			{"fracture", 0.8},
			{"wound", 0.6},
			{"bleeding", 0.6},
			{"cut", 0.4},
			{"bruise", 0.3},
			{"injury", 0.6},
			{"injured", 0.6},
		},
		description: func(string) string { return "experiencing pain" },
	},
	// Auditory: loud/startling
	{
		channel: Auditory,
		keywords: []keywordEntry{
			{"explosion", 0.95},
			{"gunshot", 0.95},
			{"thunder", 0.8},
			{"scream", 0.8},
			{"screaming", 0.8},
			{"siren", 0.75},
			{"alarm", 0.75},
			{"bang", 0.8},
			{"crash", 0.8},
			{"loud", 0.7},
			{"deafening", 0.9},
			{"roar", 0.7},
			{"shout", 0.6},
			{"yell", 0.6},
			{"whisper", 0.2},
			{"quiet", 0.15},
			{"silence", 0.1},
			{"music", 0.4},
			{"singing", 0.35},
			{"noise", 0.5},
			{"ring", 0.45},
			{"knock", 0.4},
		},
		description: func(s string) string {
			if containsAny(s, "explosion", "gunshot", "bang", "crash", "thunder") {
				return "hearing a startling sound"
			}
			if containsAny(s, "whisper", "quiet", "silence") {
				return "hearing quiet or silence"
			}
			return "hearing a sound"
		},
	},
	// Visual: threatening or notable
	{
		channel: Visual,
		keywords: []keywordEntry{
			{"dark", 0.3},
			{"darkness", 0.35},
			{"pitch black", 0.5},
			{"blind", 0.6},
			{"flash", 0.7},
			{"bright light", 0.6},
			{"blinding", 0.85},
			{"shadow", 0.4},
			{"figure approaching", 0.7},
			{"someone approaching", 0.6},
			{"blood", 0.75},
			{"fire", 0.8},
			{"flames", 0.8},
			{"weapon", 0.9},
			{"knife", 0.85},
			{"gun", 0.9},
			{"threat", 0.8},
			{"threatening", 0.8},
			{"danger", 0.8},
			{"beautiful", 0.3},
			{"sunrise", 0.25},
			{"sunset", 0.25},
		},
		description: func(s string) string {
			if containsAny(s, "dark", "darkness", "pitch black", "blind") {
				return "visual deprivation or darkness"
			}
			if containsAny(s, "weapon", "knife", "gun", "threat", "danger", "blood") {
				return "seeing something threatening"
			}
			return "visual perception"
		},
	},
	// Tactile
	{
		channel: Tactile,
		keywords: []keywordEntry{
			{"touch", 0.3},
			{"touching", 0.3},
			{"grabbed", 0.7},
			{"shoved", 0.7},
			{"pushed", 0.6},
			{"hit", 0.75},
			{"slap", 0.65},
			{"punch", 0.8},
			{"kick", 0.8},
			{"hug", 0.3},
			{"caress", 0.25},
			{"stroke", 0.25},
			{"rough", 0.5},
			{"soft", 0.2},
			{"pressure", 0.5},
			{"squeeze", 0.5},
			{"vibration", 0.4},
			{"itch", 0.3},
			{"tickle", 0.3},
			{"numb", 0.4},
		},
		description: func(s string) string {
			if containsAny(s, "hit", "punch", "kick", "slap", "grabbed", "shoved") {
				return "being physically struck or grabbed"
			}
			if containsAny(s, "hug", "caress", "stroke", "soft") {
				return "gentle physical contact"
			}
			return "tactile sensation"
		},
	},
	// Olfactory
	{
		channel: Olfactory,
		keywords: []keywordEntry{
			{"smell", 0.4},
			{"stench", 0.7},
			{"stink", 0.7},
			{"odor", 0.5},
			{"fragrance", 0.3},
			{"perfume", 0.3},
			{"smoke", 0.6},
			{"gas", 0.7},
			{"rotten", 0.7},
			{"fresh air", 0.2},
			{"aroma", 0.3},
		},
		description: func(s string) string {
			if containsAny(s, "stench", "stink", "rotten", "gas") {
				return "smelling something unpleasant"
			}
			return "detecting a smell"
		},
	},
	// Gustatory
	{
		channel: Gustatory,
		keywords: []keywordEntry{
			{"taste", 0.4},
			{"eating", 0.4},
			{"eat", 0.4},
			{"food", 0.4},
			{"drink", 0.35},
			{"drinking", 0.35},
			{"water", 0.3},
			{"bitter", 0.5},
			{"sour", 0.4},
			{"sweet", 0.3},
			{"salty", 0.3},
			{"delicious", 0.3},
			{"disgusting", 0.6},
			{"vomit", 0.7},
			{"nausea", 0.6},
			{"hungry", 0.3},
			{"thirsty", 0.3},
			{"starving", 0.5},
		},
		description: func(s string) string {
			if containsAny(s, "eat", "food", "taste", "delicious") {
				return "tasting or eating"
			}
			if containsAny(s, "drink", "water", "thirsty") {
				return "drinking"
			}
			return "gustatory sensation"
		},
	},
	// Interoceptive
	{
		channel: Interoceptive,
		keywords: []keywordEntry{
			{"nausea", 0.6},
			{"nauseous", 0.6},
			{"dizzy", 0.6},
			{"dizziness", 0.6},
			{"faint", 0.7},
			{"fainting", 0.7},
			{"breathless", 0.7},
			{"can't breathe", 0.85},
			{"suffocating", 0.9},
			{"choking", 0.85},
			{"heart racing", 0.6},
			{"heart pounding", 0.7},
			{"chest tight", 0.7},
			{"stomach", 0.4},
			{"gut", 0.4},
			{"shaking", 0.5},
			{"trembling", 0.5},
			{"sweating", 0.4},
			{"chills", 0.4},
		},
		description: func(s string) string {
			if containsAny(s, "breathless", "can't breathe", "suffocating", "choking") {
				return "difficulty breathing"
			}
			if containsAny(s, "dizzy", "faint") {
				return "feeling dizzy or faint"
			}
			if containsAny(s, "nausea", "nauseous", "stomach") {
				return "gastrointestinal distress"
			}
			return "internal body sensation"
		},
	},
	// Vestibular
	{
		channel: Vestibular,
		keywords: []keywordEntry{
			{"falling", 0.8},
			{"fall", 0.7},
			{"spinning", 0.7},
			{"vertigo", 0.8},
			{"tilt", 0.5},
			{"off balance", 0.6},
			{"stumble", 0.5},
			{"earthquake", 0.85},
			{"shaking ground", 0.7},
		},
		description: func(s string) string {
			if containsAny(s, "falling", "fall") {
				return "sensation of falling"
			}
			if containsAny(s, "spinning", "vertigo") {
				return "rotational disorientation"
			}
			return "balance disruption"
		},
	},
}

func containsAny(s string, patterns ...string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
