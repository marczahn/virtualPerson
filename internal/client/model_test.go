package client

import (
	"testing"
	"time"

	"github.com/marczahn/person/internal/server"
)

func TestModel_AddThought_IgnoresNonThoughtMessages(t *testing.T) {
	m := &Model{}

	m.addThought(server.ServerMessage{
		Type:      "bio_state",
		Timestamp: time.Now(),
		BioState:  &server.BioStatePayload{},
	})
	if len(m.thoughts) != 0 {
		t.Errorf("expected no thoughts after bio_state message, got %d", len(m.thoughts))
	}

	m.addThought(server.ServerMessage{
		Type:       "psych_state",
		Timestamp:  time.Now(),
		PsychState: &server.PsychStatePayload{},
	})
	if len(m.thoughts) != 0 {
		t.Errorf("expected no thoughts after psych_state message, got %d", len(m.thoughts))
	}
}

func TestModel_AddThought_ProcessesThoughtMessages(t *testing.T) {
	m := &Model{}

	m.addThought(server.ServerMessage{
		Type:      "thought",
		Content:   "I wonder about things",
		Timestamp: time.Now(),
	})
	if len(m.thoughts) != 1 {
		t.Errorf("expected 1 thought after thought message, got %d", len(m.thoughts))
	}
}

func TestModel_AddThought_ThoughtWithTriggerIncludesTrigger(t *testing.T) {
	m := &Model{}

	m.addThought(server.ServerMessage{
		Type:      "thought",
		Content:   "That hurt",
		Trigger:   "pain_spike",
		Timestamp: time.Now(),
	})
	if len(m.thoughts) != 1 {
		t.Fatalf("expected 1 thought, got %d", len(m.thoughts))
	}
	if m.thoughts[0] == "" {
		t.Error("thought line must not be empty")
	}
}
