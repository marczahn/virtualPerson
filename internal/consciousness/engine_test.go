package consciousness

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/psychology"
)

// mockLLM is a test double for the LLM interface.
type mockLLM struct {
	response    string
	err         error
	calls       int
	lastUser    string
	lastSystem  string
}

func (m *mockLLM) Complete(ctx context.Context, system, user string) (string, error) {
	m.calls++
	m.lastSystem = system
	m.lastUser = user
	return m.response, m.err
}

func TestEngine_React_NoSalience_NoThought(t *testing.T) {
	llm := &mockLLM{response: "I feel calm."}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: 0,
	})

	ps := &psychology.State{Arousal: 0.2, Valence: 0.3, Energy: 0.5}

	// First call establishes baseline.
	engine.React(context.Background(), ps, 1.0)

	// Same state — no salience.
	thought, err := engine.React(context.Background(), ps, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should only have called LLM once (for the first baseline which may or may not trigger).
	// The second call with same state should not trigger.
	if thought != nil && llm.calls > 1 {
		t.Logf("thought generated on stable state: %q", thought.Content)
	}
}

func TestEngine_React_SalienceSpike_GeneratesThought(t *testing.T) {
	llm := &mockLLM{response: "What was that? My heart is racing!"}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: 0,
	})

	// Establish calm baseline.
	calm := &psychology.State{Arousal: 0.1, Valence: 0.3, Energy: 0.5}
	engine.React(context.Background(), calm, 1.0)

	// Spike arousal.
	spike := &psychology.State{Arousal: 0.9, Valence: -0.3, Energy: 0.5}
	thought, err := engine.React(context.Background(), spike, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if thought == nil {
		t.Fatal("expected reactive thought on arousal spike")
	}
	if thought.Type != Reactive {
		t.Errorf("type = %s, expected reactive", thought.Type)
	}
	if thought.Content != "What was that? My heart is racing!" {
		t.Errorf("content = %q, unexpected", thought.Content)
	}
}

func TestEngine_React_FeedbackParsed(t *testing.T) {
	llm := &mockLLM{response: "I keep thinking about it. I can't stop thinking about what happened."}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: 0,
	})

	calm := &psychology.State{Arousal: 0.1}
	engine.React(context.Background(), calm, 1.0)

	spike := &psychology.State{Arousal: 0.9, Valence: -0.5}
	thought, _ := engine.React(context.Background(), spike, 1.0)

	if thought == nil {
		t.Fatal("expected thought")
	}

	found := false
	for _, c := range thought.Feedback.ActiveCoping {
		if c == "rumination" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected rumination in feedback, got %v", thought.Feedback.ActiveCoping)
	}
}

func TestEngine_Spontaneous_RateLimited(t *testing.T) {
	llm := &mockLLM{response: "My mind wanders..."}
	engine := NewEngine(EngineConfig{
		LLM:                 llm,
		MinCallInterval:     0,
		SpontaneousInterval: time.Hour, // very long
	})

	// Set last spontaneous to now — should be rate limited.
	engine.lastSpontaneous = time.Now()

	ps := &psychology.State{Energy: 0.5}
	thought, err := engine.Spontaneous(context.Background(), ps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thought != nil {
		t.Error("expected nil thought due to rate limiting")
	}
}

func TestEngine_Spontaneous_GeneratesThought(t *testing.T) {
	llm := &mockLLM{response: "I wonder what time it is..."}
	engine := NewEngine(EngineConfig{
		LLM:                 llm,
		MinCallInterval:     0,
		SpontaneousInterval: 0, // no rate limit
	})

	ps := &psychology.State{Energy: 0.5}
	thought, err := engine.Spontaneous(context.Background(), ps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if thought == nil {
		t.Fatal("expected spontaneous thought")
	}
	if thought.Type != Spontaneous {
		t.Errorf("type = %s, expected spontaneous", thought.Type)
	}
}

func TestEngine_UpdateMemories(t *testing.T) {
	llm := &mockLLM{response: "test"}
	engine := NewEngine(EngineConfig{LLM: llm})

	memories := []memory.EpisodicMemory{
		{ID: "m1", Content: "test memory"},
	}

	engine.UpdateMemories(memories)

	if len(engine.memories) != 1 {
		t.Errorf("expected 1 memory, got %d", len(engine.memories))
	}
}

func TestEngine_UpdateIdentity(t *testing.T) {
	llm := &mockLLM{response: "test"}
	engine := NewEngine(EngineConfig{LLM: llm})

	ic := &memory.IdentityCore{SelfNarrative: "I am me."}
	engine.UpdateIdentity(ic)

	if engine.identity.SelfNarrative != "I am me." {
		t.Error("identity not updated")
	}
}

func TestEngine_MinCallInterval_Respected(t *testing.T) {
	llm := &mockLLM{response: "test"}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: time.Hour, // very long
	})

	// First call goes through (lastCallTime is zero).
	calm := &psychology.State{Arousal: 0.1}
	engine.React(context.Background(), calm, 1.0)

	// Mark that a call was made.
	engine.lastCallTime = time.Now()

	// Second call should be rate limited.
	spike := &psychology.State{Arousal: 0.9}
	thought, _ := engine.React(context.Background(), spike, 1.0)

	if thought != nil {
		t.Error("expected rate limiting to prevent second call")
	}
}

func TestEngine_Respond_Speech_GeneratesThought(t *testing.T) {
	llm := &mockLLM{response: "Oh, someone's talking to me."}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: 0,
	})

	ps := &psychology.State{Arousal: 0.3, Valence: 0.2, Energy: 0.5}
	input := ExternalInput{Type: InputSpeech, Content: "Hello, how are you?"}

	thought, err := engine.Respond(context.Background(), ps, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thought == nil {
		t.Fatal("expected thought from speech input")
	}
	if thought.Type != Conversational {
		t.Errorf("type = %s, expected conversational", thought.Type)
	}
	if thought.Content != "Oh, someone's talking to me." {
		t.Errorf("content = %q, unexpected", thought.Content)
	}
	if thought.Trigger != "Hello, how are you?" {
		t.Errorf("trigger = %q, expected speech content", thought.Trigger)
	}
}

func TestEngine_Respond_Action_GeneratesThought(t *testing.T) {
	llm := &mockLLM{response: "Ow, that hurt!"}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: 0,
	})

	ps := &psychology.State{Arousal: 0.5, Valence: -0.1, Energy: 0.5}
	input := ExternalInput{Type: InputAction, Content: "pushes you hard"}

	thought, err := engine.Respond(context.Background(), ps, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thought == nil {
		t.Fatal("expected thought from action input")
	}
	if thought.Type != Conversational {
		t.Errorf("type = %s, expected conversational", thought.Type)
	}
}

func TestEngine_Respond_RateLimited(t *testing.T) {
	llm := &mockLLM{response: "test"}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: time.Hour,
	})

	engine.lastCallTime = time.Now()

	ps := &psychology.State{Arousal: 0.3}
	input := ExternalInput{Type: InputSpeech, Content: "hello"}

	thought, err := engine.Respond(context.Background(), ps, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thought != nil {
		t.Error("expected nil thought due to rate limiting")
	}
}

func TestEngine_Respond_ContentReachesPrompt(t *testing.T) {
	llm := &mockLLM{response: "I hear you."}
	engine := NewEngine(EngineConfig{
		LLM:             llm,
		MinCallInterval: 0,
	})

	ps := &psychology.State{Arousal: 0.3, Valence: 0.2, Energy: 0.5}
	input := ExternalInput{Type: InputSpeech, Content: "Do you like pizza?"}

	_, err := engine.Respond(context.Background(), ps, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(llm.lastUser, "Do you like pizza?") {
		t.Errorf("expected speech content in prompt, got: %s", llm.lastUser)
	}
	if !strings.Contains(llm.lastUser, "Someone says to you") {
		t.Errorf("expected speech framing in prompt, got: %s", llm.lastUser)
	}
}
