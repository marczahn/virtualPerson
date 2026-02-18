package server

import (
	"fmt"
	"time"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/psychology"
)

// ClientMessage is sent from the TUI client to the server.
type ClientMessage struct {
	Type    string `json:"type"`    // "speech", "action", "environment"
	Content string `json:"content"`
}

// Validate checks that the message has a known type and non-empty content.
func (m ClientMessage) Validate() error {
	switch m.Type {
	case "speech", "action", "environment":
	default:
		return fmt.Errorf("unknown message type: %q", m.Type)
	}
	if m.Content == "" {
		return fmt.Errorf("content must not be empty")
	}
	return nil
}

// ToInputLine converts a client message to the simulation's input format:
// speech → plain text, action → *text*, environment → ~text.
func (m ClientMessage) ToInputLine() string {
	switch m.Type {
	case "action":
		return "*" + m.Content + "*"
	case "environment":
		return "~" + m.Content
	default:
		return m.Content
	}
}

// ServerMessage is sent from the server to connected clients.
// The Type field determines which payload fields are populated.
// "thought" messages use Content, ThoughtType, and Trigger.
// "bio_state" messages use BioState.
// "psych_state" messages use PsychState.
type ServerMessage struct {
	Type        string    `json:"type"`
	Content     string    `json:"content,omitempty"`
	ThoughtType string    `json:"thought_type,omitempty"`
	Trigger     string    `json:"trigger,omitempty"`
	Timestamp   time.Time `json:"timestamp"`

	BioState   *BioStatePayload   `json:"bio_state,omitempty"`
	PsychState *PsychStatePayload `json:"psych_state,omitempty"`
}

// BioStatePayload carries a snapshot of all 20 biological variables
// plus any currently active threshold breaches.
type BioStatePayload struct {
	BodyTemp        float64 `json:"body_temp"`
	HeartRate       float64 `json:"heart_rate"`
	BloodPressure   float64 `json:"blood_pressure"`
	RespiratoryRate float64 `json:"respiratory_rate"`
	Hunger          float64 `json:"hunger"`
	Thirst          float64 `json:"thirst"`
	Fatigue         float64 `json:"fatigue"`
	Pain            float64 `json:"pain"`
	MuscleTension   float64 `json:"muscle_tension"`
	BloodSugar      float64 `json:"blood_sugar"`
	Cortisol        float64 `json:"cortisol"`
	Adrenaline      float64 `json:"adrenaline"`
	Serotonin       float64 `json:"serotonin"`
	Dopamine        float64 `json:"dopamine"`
	ImmuneResponse  float64 `json:"immune_response"`
	CircadianPhase  float64 `json:"circadian_phase"`
	SpO2            float64 `json:"spo2"`
	Hydration       float64 `json:"hydration"`
	Glycogen        float64 `json:"glycogen"`
	Endorphins      float64 `json:"endorphins"`

	// Active threshold breaches; empty slice (not nil) when all normal.
	Thresholds []ThresholdPayload `json:"thresholds"`
}

// ThresholdPayload describes one active threshold breach.
type ThresholdPayload struct {
	Condition   string `json:"condition"`
	System      string `json:"system"`
	Description string `json:"description"`
}

// PsychStatePayload carries the current psychological state and stable personality.
type PsychStatePayload struct {
	Arousal            float64 `json:"arousal"`
	Valence            float64 `json:"valence"`
	Energy             float64 `json:"energy"`
	CognitiveLoad      float64 `json:"cognitive_load"`
	RegulationCapacity float64 `json:"regulation_capacity"`

	ActiveDistortions []string `json:"active_distortions"`
	ActiveCoping      []string `json:"active_coping"`
	IsolationPhase    string   `json:"isolation_phase"`
	LonelinessLevel   float64  `json:"loneliness_level"`

	Personality PersonalityPayload `json:"personality"`
}

// PersonalityPayload carries the Big Five trait scores.
type PersonalityPayload struct {
	Openness          float64 `json:"openness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Extraversion      float64 `json:"extraversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Neuroticism       float64 `json:"neuroticism"`
}

// BioStatePayloadFromState converts a biology.State into a BioStatePayload,
// evaluating current thresholds. Returns a fully populated payload with a
// non-nil Thresholds slice even when no thresholds are active.
func BioStatePayloadFromState(s *biology.State) BioStatePayload {
	thresholds := biology.EvaluateThresholds(s)
	tp := make([]ThresholdPayload, 0, len(thresholds))
	for _, t := range thresholds {
		tp = append(tp, ThresholdPayload{
			Condition:   t.Condition.String(),
			System:      t.System,
			Description: t.Description,
		})
	}
	return BioStatePayload{
		BodyTemp:        s.BodyTemp,
		HeartRate:       s.HeartRate,
		BloodPressure:   s.BloodPressure,
		RespiratoryRate: s.RespiratoryRate,
		Hunger:          s.Hunger,
		Thirst:          s.Thirst,
		Fatigue:         s.Fatigue,
		Pain:            s.Pain,
		MuscleTension:   s.MuscleTension,
		BloodSugar:      s.BloodSugar,
		Cortisol:        s.Cortisol,
		Adrenaline:      s.Adrenaline,
		Serotonin:       s.Serotonin,
		Dopamine:        s.Dopamine,
		ImmuneResponse:  s.ImmuneResponse,
		CircadianPhase:  s.CircadianPhase,
		SpO2:            s.SpO2,
		Hydration:       s.Hydration,
		Glycogen:        s.Glycogen,
		Endorphins:      s.Endorphins,
		Thresholds:      tp,
	}
}

// PsychStatePayloadFromState converts a psychology.State and Personality
// into a PsychStatePayload with string-serialized distortions and coping strategies.
func PsychStatePayloadFromState(s psychology.State, p psychology.Personality) PsychStatePayload {
	distortions := make([]string, len(s.ActiveDistortions))
	for i, d := range s.ActiveDistortions {
		distortions[i] = d.String()
	}
	coping := make([]string, len(s.ActiveCoping))
	for i, c := range s.ActiveCoping {
		coping[i] = c.String()
	}
	return PsychStatePayload{
		Arousal:            s.Arousal,
		Valence:            s.Valence,
		Energy:             s.Energy,
		CognitiveLoad:      s.CognitiveLoad,
		RegulationCapacity: s.RegulationCapacity,
		ActiveDistortions:  distortions,
		ActiveCoping:       coping,
		IsolationPhase:     s.Isolation.Phase.String(),
		LonelinessLevel:    s.Isolation.LonelinessLevel,
		Personality: PersonalityPayload{
			Openness:          p.Openness,
			Conscientiousness: p.Conscientiousness,
			Extraversion:      p.Extraversion,
			Agreeableness:     p.Agreeableness,
			Neuroticism:       p.Neuroticism,
		},
	}
}
