package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/psychology"
)

func TestClientMessage_Validate_ValidTypes(t *testing.T) {
	for _, typ := range []string{"speech", "action", "environment"} {
		msg := ClientMessage{Type: typ, Content: "hello"}
		if err := msg.Validate(); err != nil {
			t.Errorf("Validate(%q) returned error: %v", typ, err)
		}
	}
}

func TestClientMessage_Validate_UnknownType(t *testing.T) {
	msg := ClientMessage{Type: "shout", Content: "hello"}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestClientMessage_Validate_EmptyContent(t *testing.T) {
	msg := ClientMessage{Type: "speech", Content: ""}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for empty content")
	}
}

func TestClientMessage_ToInputLine_Speech(t *testing.T) {
	msg := ClientMessage{Type: "speech", Content: "hello"}
	if got := msg.ToInputLine(); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestClientMessage_ToInputLine_Action(t *testing.T) {
	msg := ClientMessage{Type: "action", Content: "pushes you"}
	if got := msg.ToInputLine(); got != "*pushes you*" {
		t.Errorf("got %q, want %q", got, "*pushes you*")
	}
}

func TestClientMessage_ToInputLine_Environment(t *testing.T) {
	msg := ClientMessage{Type: "environment", Content: "cold wind"}
	if got := msg.ToInputLine(); got != "~cold wind" {
		t.Errorf("got %q, want %q", got, "~cold wind")
	}
}

func TestBioStatePayloadFromState_DefaultState(t *testing.T) {
	s := biology.NewDefaultState()
	payload := BioStatePayloadFromState(&s)

	if payload.BodyTemp != s.BodyTemp {
		t.Errorf("BodyTemp: got %v, want %v", payload.BodyTemp, s.BodyTemp)
	}
	if payload.HeartRate != s.HeartRate {
		t.Errorf("HeartRate: got %v, want %v", payload.HeartRate, s.HeartRate)
	}
	if payload.BloodSugar != s.BloodSugar {
		t.Errorf("BloodSugar: got %v, want %v", payload.BloodSugar, s.BloodSugar)
	}
	// Default state has no threshold breaches; must be empty slice, not nil.
	if payload.Thresholds == nil {
		t.Error("Thresholds should be an empty slice, not nil")
	}
	if len(payload.Thresholds) != 0 {
		t.Errorf("expected no thresholds for default state, got %d", len(payload.Thresholds))
	}
}

func TestBioStatePayloadFromState_ThresholdBreachIncluded(t *testing.T) {
	s := biology.NewDefaultState()
	s.BloodSugar = 30 // critically below any safe threshold
	payload := BioStatePayloadFromState(&s)

	if len(payload.Thresholds) == 0 {
		t.Fatal("expected threshold breaches for critically low blood sugar")
	}
	for _, th := range payload.Thresholds {
		if th.System == "" {
			t.Error("threshold System must not be empty")
		}
		if th.Condition == "" {
			t.Error("threshold Condition must not be empty")
		}
	}
}

func TestBioStatePayloadFromState_AllVarsMapped(t *testing.T) {
	s := biology.NewDefaultState()
	s.Cortisol = 0.77
	s.Serotonin = 0.33
	s.CircadianPhase = 14.5
	payload := BioStatePayloadFromState(&s)

	if payload.Cortisol != 0.77 {
		t.Errorf("Cortisol: got %v, want 0.77", payload.Cortisol)
	}
	if payload.Serotonin != 0.33 {
		t.Errorf("Serotonin: got %v, want 0.33", payload.Serotonin)
	}
	if payload.CircadianPhase != 14.5 {
		t.Errorf("CircadianPhase: got %v, want 14.5", payload.CircadianPhase)
	}
}

func TestPsychStatePayloadFromState_AffectDimensions(t *testing.T) {
	s := psychology.State{
		Arousal:            0.8,
		Valence:            -0.5,
		Energy:             0.3,
		CognitiveLoad:      0.7,
		RegulationCapacity: 0.4,
	}
	p := psychology.Personality{Openness: 0.6}
	payload := PsychStatePayloadFromState(s, p)

	if payload.Arousal != 0.8 {
		t.Errorf("Arousal: got %v, want 0.8", payload.Arousal)
	}
	if payload.Valence != -0.5 {
		t.Errorf("Valence: got %v, want -0.5", payload.Valence)
	}
	if payload.Energy != 0.3 {
		t.Errorf("Energy: got %v, want 0.3", payload.Energy)
	}
	if payload.CognitiveLoad != 0.7 {
		t.Errorf("CognitiveLoad: got %v, want 0.7", payload.CognitiveLoad)
	}
	if payload.RegulationCapacity != 0.4 {
		t.Errorf("RegulationCapacity: got %v, want 0.4", payload.RegulationCapacity)
	}
}

func TestPsychStatePayloadFromState_ActiveProcesses(t *testing.T) {
	s := psychology.State{
		ActiveDistortions: []psychology.Distortion{psychology.Catastrophizing, psychology.MindReading},
		ActiveCoping:      []psychology.CopingStrategy{psychology.Suppression},
		Isolation: psychology.IsolationState{
			Phase:           psychology.IsolationLoneliness,
			LonelinessLevel: 0.4,
		},
	}
	p := psychology.Personality{
		Openness: 0.6, Conscientiousness: 0.5,
		Extraversion: 0.4, Agreeableness: 0.6, Neuroticism: 0.5,
	}
	payload := PsychStatePayloadFromState(s, p)

	if len(payload.ActiveDistortions) != 2 {
		t.Fatalf("expected 2 distortions, got %d", len(payload.ActiveDistortions))
	}
	if payload.ActiveDistortions[0] != "catastrophizing" {
		t.Errorf("expected catastrophizing, got %q", payload.ActiveDistortions[0])
	}
	if payload.ActiveDistortions[1] != "mind_reading" {
		t.Errorf("expected mind_reading, got %q", payload.ActiveDistortions[1])
	}
	if len(payload.ActiveCoping) != 1 {
		t.Fatalf("expected 1 coping, got %d", len(payload.ActiveCoping))
	}
	if payload.ActiveCoping[0] != "suppression" {
		t.Errorf("expected suppression, got %q", payload.ActiveCoping[0])
	}
	if payload.IsolationPhase != "loneliness" {
		t.Errorf("expected loneliness phase, got %q", payload.IsolationPhase)
	}
	if payload.LonelinessLevel != 0.4 {
		t.Errorf("LonelinessLevel: got %v, want 0.4", payload.LonelinessLevel)
	}
	if payload.Personality.Openness != 0.6 {
		t.Errorf("personality openness: got %v, want 0.6", payload.Personality.Openness)
	}
}

func TestServerMessage_BioState_JSONOmitsContent(t *testing.T) {
	payload := BioStatePayload{BodyTemp: 36.6}
	msg := ServerMessage{
		Type:      "bio_state",
		Timestamp: time.Now(),
		BioState:  &payload,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if v, ok := m["content"]; ok && v != "" {
		t.Errorf("bio_state message should not include non-empty content field, got %v", v)
	}
	if _, ok := m["bio_state"]; !ok {
		t.Error("bio_state message must include bio_state field")
	}
	if _, ok := m["psych_state"]; ok {
		t.Error("bio_state message must not include psych_state field")
	}
}

func TestServerMessage_Thought_JSONOmitsStateFields(t *testing.T) {
	msg := ServerMessage{
		Type:      "thought",
		Content:   "I think therefore I am",
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(msg)
	var m map[string]any
	json.Unmarshal(data, &m)

	if _, ok := m["bio_state"]; ok {
		t.Error("thought message must not include bio_state field")
	}
	if _, ok := m["psych_state"]; ok {
		t.Error("thought message must not include psych_state field")
	}
}
