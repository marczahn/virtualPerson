package reviewer

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/psychology"
)

type mockLLM struct {
	response   string
	lastSystem string
	lastUser   string
	callCount  int
}

func (m *mockLLM) Complete(_ context.Context, system, user string) (string, error) {
	m.lastSystem = system
	m.lastUser = user
	m.callCount++
	return m.response, nil
}

func TestAddThought_FillsBuffer(t *testing.T) {
	r := NewReviewer(ReviewerConfig{
		LLM:         &mockLLM{},
		MaxThoughts: 3,
	})

	for i := 0; i < 3; i++ {
		r.AddThought(consciousness.Thought{Content: "thought"})
	}

	if len(r.thoughts) != 3 {
		t.Errorf("expected 3 thoughts, got %d", len(r.thoughts))
	}
}

func TestAddThought_EvictsOldest(t *testing.T) {
	r := NewReviewer(ReviewerConfig{
		LLM:         &mockLLM{},
		MaxThoughts: 3,
	})

	r.AddThought(consciousness.Thought{Content: "first"})
	r.AddThought(consciousness.Thought{Content: "second"})
	r.AddThought(consciousness.Thought{Content: "third"})
	r.AddThought(consciousness.Thought{Content: "fourth"})

	if len(r.thoughts) != 3 {
		t.Fatalf("expected 3 thoughts after eviction, got %d", len(r.thoughts))
	}
	if r.thoughts[0].Content != "second" {
		t.Errorf("expected oldest remaining to be 'second', got %q", r.thoughts[0].Content)
	}
	if r.thoughts[2].Content != "fourth" {
		t.Errorf("expected newest to be 'fourth', got %q", r.thoughts[2].Content)
	}
}

func TestReview_GeneratesObservation(t *testing.T) {
	llm := &mockLLM{response: "Subject exhibits elevated arousal with negative valence."}
	r := NewReviewer(ReviewerConfig{LLM: llm})

	r.AddThought(consciousness.Thought{
		Type:    consciousness.Reactive,
		Trigger: "loud_noise",
		Content: "That startled me badly.",
	})

	ps := &psychology.State{Arousal: 0.8, Valence: -0.5}
	personality := &psychology.Personality{Neuroticism: 0.7}

	obs, err := r.Review(context.Background(), ps, personality)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs == nil {
		t.Fatal("expected observation, got nil")
	}
	if obs.Content != "Subject exhibits elevated arousal with negative valence." {
		t.Errorf("unexpected content: %q", obs.Content)
	}
	if obs.Timestamp.IsZero() {
		t.Error("observation timestamp should not be zero")
	}
}

func TestReview_ReturnsNilWhenBufferEmpty(t *testing.T) {
	r := NewReviewer(ReviewerConfig{LLM: &mockLLM{response: "something"}})

	obs, err := r.Review(context.Background(), &psychology.State{}, &psychology.Personality{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs != nil {
		t.Errorf("expected nil observation for empty buffer, got %+v", obs)
	}
}

func TestReview_ReturnsNilWhenRateLimited(t *testing.T) {
	llm := &mockLLM{response: "observation"}
	r := NewReviewer(ReviewerConfig{
		LLM:         llm,
		MinInterval: 1 * time.Hour,
	})

	r.AddThought(consciousness.Thought{Content: "something"})
	ps := &psychology.State{}
	personality := &psychology.Personality{}

	// First call succeeds.
	obs, err := r.Review(context.Background(), ps, personality)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs == nil {
		t.Fatal("first review should succeed")
	}

	// Second call should be rate-limited.
	r.AddThought(consciousness.Thought{Content: "another"})
	obs, err = r.Review(context.Background(), ps, personality)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs != nil {
		t.Error("expected nil observation when rate-limited")
	}
	if llm.callCount != 1 {
		t.Errorf("expected 1 LLM call, got %d", llm.callCount)
	}
}

func TestReview_NoRateLimitWhenIntervalZero(t *testing.T) {
	llm := &mockLLM{response: "observation"}
	r := NewReviewer(ReviewerConfig{
		LLM:         llm,
		MinInterval: 0,
	})

	ps := &psychology.State{}
	personality := &psychology.Personality{}

	r.AddThought(consciousness.Thought{Content: "first"})
	r.Review(context.Background(), ps, personality)

	r.AddThought(consciousness.Thought{Content: "second"})
	obs, err := r.Review(context.Background(), ps, personality)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs == nil {
		t.Error("expected observation with zero interval (no rate limit)")
	}
	if llm.callCount != 2 {
		t.Errorf("expected 2 LLM calls, got %d", llm.callCount)
	}
}

func TestReview_PromptContainsThoughtContent(t *testing.T) {
	llm := &mockLLM{response: "noted"}
	r := NewReviewer(ReviewerConfig{LLM: llm})

	r.AddThought(consciousness.Thought{
		Type:    consciousness.Spontaneous,
		Trigger: "rumination",
		Content: "I keep thinking about the argument.",
	})

	ps := &psychology.State{Arousal: 0.6, Valence: -0.4}
	personality := &psychology.Personality{Neuroticism: 0.8}

	r.Review(context.Background(), ps, personality)

	if !strings.Contains(llm.lastUser, "I keep thinking about the argument.") {
		t.Error("user prompt should contain thought content")
	}
	if !strings.Contains(llm.lastUser, "0.60") {
		t.Error("user prompt should contain arousal value")
	}
	if !strings.Contains(llm.lastUser, "0.80") {
		t.Error("user prompt should contain neuroticism value")
	}
	if !strings.Contains(llm.lastSystem, "clinical psychologist") {
		t.Error("system prompt should contain clinical role")
	}
}

func TestNewReviewer_DefaultMaxThoughts(t *testing.T) {
	r := NewReviewer(ReviewerConfig{LLM: &mockLLM{}})
	if r.maxThoughts != 20 {
		t.Errorf("expected default maxThoughts=20, got %d", r.maxThoughts)
	}
}
