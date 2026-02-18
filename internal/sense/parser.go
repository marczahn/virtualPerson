package sense

import (
	"strings"
	"time"

	"github.com/marczahn/person/internal/i18n"
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
	tr := i18n.T()
	var events []Event

	rules := buildRules(tr)
	for _, rule := range rules {
		if matched, intensity := rule.match(lower); matched {
			events = append(events, Event{
				Channel:   rule.channel,
				Intensity: intensity,
				RawInput:  input,
				Parsed:    rule.describe(lower, tr),
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
			Parsed:    tr.Sense.Fallback,
			Timestamp: now,
		})
	}

	return events
}

type keywordEntry struct {
	phrase    string
	intensity float64
}

type keywordRule struct {
	channel  Channel
	keywords []keywordEntry
	describe func(string, *i18n.Translations) string
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

func buildRules(tr *i18n.Translations) []keywordRule {
	return []keywordRule{
		{
			channel:  Thermal,
			keywords: toKeywordEntries(tr.Sense.Keywords.Thermal),
			describe: describeThermal,
		},
		{
			channel:  Pain,
			keywords: toKeywordEntries(tr.Sense.Keywords.Pain),
			describe: func(_ string, tr *i18n.Translations) string {
				return tr.Sense.Descriptions.ExperiencingPain
			},
		},
		{
			channel:  Auditory,
			keywords: toKeywordEntries(tr.Sense.Keywords.Auditory),
			describe: describeAuditory,
		},
		{
			channel:  Visual,
			keywords: toKeywordEntries(tr.Sense.Keywords.Visual),
			describe: describeVisual,
		},
		{
			channel:  Tactile,
			keywords: toKeywordEntries(tr.Sense.Keywords.Tactile),
			describe: describeTactile,
		},
		{
			channel:  Olfactory,
			keywords: toKeywordEntries(tr.Sense.Keywords.Olfactory),
			describe: describeOlfactory,
		},
		{
			channel:  Gustatory,
			keywords: toKeywordEntries(tr.Sense.Keywords.Gustatory),
			describe: describeGustatory,
		},
		{
			channel:  Interoceptive,
			keywords: toKeywordEntries(tr.Sense.Keywords.Interoceptive),
			describe: describeInteroceptive,
		},
		{
			channel:  Vestibular,
			keywords: toKeywordEntries(tr.Sense.Keywords.Vestibular),
			describe: describeVestibular,
		},
	}
}

func toKeywordEntries(entries []i18n.KeywordEntry) []keywordEntry {
	out := make([]keywordEntry, len(entries))
	for i, e := range entries {
		out[i] = keywordEntry{phrase: e.Phrase, intensity: e.Intensity}
	}
	return out
}

// Description functions use explicit keyword groups from the translations
// to select the most specific description for each sensory channel.

func describeThermal(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.ThermalCold...) {
		return tr.Sense.Descriptions.FeelingCold
	}
	return tr.Sense.Descriptions.FeelingHeat
}

func describeAuditory(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.AuditoryStartling...) {
		return tr.Sense.Descriptions.HearingStartling
	}
	if containsAny(s, tr.Sense.DescriptionGroups.AuditoryQuiet...) {
		return tr.Sense.Descriptions.HearingQuiet
	}
	return tr.Sense.Descriptions.HearingSound
}

func describeVisual(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.VisualDark...) {
		return tr.Sense.Descriptions.VisualDarkness
	}
	if containsAny(s, tr.Sense.DescriptionGroups.VisualThreat...) {
		return tr.Sense.Descriptions.SeeingThreatening
	}
	return tr.Sense.Descriptions.VisualPerception
}

func describeTactile(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.TactileViolent...) {
		return tr.Sense.Descriptions.PhysicallyStruck
	}
	if containsAny(s, tr.Sense.DescriptionGroups.TactileGentle...) {
		return tr.Sense.Descriptions.GentleContact
	}
	return tr.Sense.Descriptions.TactileSensation
}

func describeOlfactory(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.OlfactoryUnpleasant...) {
		return tr.Sense.Descriptions.SmellingUnpleasant
	}
	return tr.Sense.Descriptions.DetectingSmell
}

func describeGustatory(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.GustatoryEating...) {
		return tr.Sense.Descriptions.TastingEating
	}
	if containsAny(s, tr.Sense.DescriptionGroups.GustatoryDrinking...) {
		return tr.Sense.Descriptions.Drinking
	}
	return tr.Sense.Descriptions.GustatorySensation
}

func describeInteroceptive(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.InteroBreathing...) {
		return tr.Sense.Descriptions.DifficultyBreathing
	}
	if containsAny(s, tr.Sense.DescriptionGroups.InteroDizzy...) {
		return tr.Sense.Descriptions.FeelingDizzy
	}
	if containsAny(s, tr.Sense.DescriptionGroups.InteroGastro...) {
		return tr.Sense.Descriptions.GastrointestinalDistress
	}
	return tr.Sense.Descriptions.InternalBodySensation
}

func describeVestibular(s string, tr *i18n.Translations) string {
	if containsAny(s, tr.Sense.DescriptionGroups.VestibularFalling...) {
		return tr.Sense.Descriptions.SensationFalling
	}
	if containsAny(s, tr.Sense.DescriptionGroups.VestibularSpinning...) {
		return tr.Sense.Descriptions.RotationalDisorientation
	}
	return tr.Sense.Descriptions.BalanceDisruption
}

func containsAny(s string, patterns ...string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
