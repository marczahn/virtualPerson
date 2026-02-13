package reviewer

import (
	"context"
	"time"

	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/psychology"
)

// Observation is a single clinical observation from the reviewer.
type Observation struct {
	Content   string
	Timestamp time.Time
}

// ReviewerConfig holds the settings for creating a Reviewer.
type ReviewerConfig struct {
	LLM         consciousness.LLM
	MinInterval time.Duration // minimum time between reviews; 0 = no rate limit
	MaxThoughts int           // rolling buffer capacity; default 20
}

// Reviewer is a periodic LLM-based clinical observer that analyzes
// recent thoughts and psychological state to produce observations.
type Reviewer struct {
	llm           consciousness.LLM
	promptBuilder *PromptBuilder
	minInterval   time.Duration
	lastReview    time.Time
	thoughts      []consciousness.Thought
	maxThoughts   int
}

// NewReviewer creates a Reviewer from the given configuration.
func NewReviewer(cfg ReviewerConfig) *Reviewer {
	max := cfg.MaxThoughts
	if max == 0 {
		max = 20
	}
	return &Reviewer{
		llm:           cfg.LLM,
		promptBuilder: NewPromptBuilder(),
		minInterval:   cfg.MinInterval,
		maxThoughts:   max,
		thoughts:      make([]consciousness.Thought, 0, max),
	}
}

// AddThought appends a thought to the rolling buffer.
// When the buffer is full, the oldest thought is dropped.
func (r *Reviewer) AddThought(t consciousness.Thought) {
	if len(r.thoughts) >= r.maxThoughts {
		// Drop oldest: shift left by one.
		copy(r.thoughts, r.thoughts[1:])
		r.thoughts[len(r.thoughts)-1] = t
	} else {
		r.thoughts = append(r.thoughts, t)
	}
}

// Review builds a prompt from buffered thoughts and the current state,
// calls the LLM, and returns a clinical observation.
//
// Returns nil (no error) when rate-limited or when the buffer is empty.
func (r *Reviewer) Review(
	ctx context.Context,
	ps *psychology.State,
	personality *psychology.Personality,
) (*Observation, error) {
	if len(r.thoughts) == 0 {
		return nil, nil
	}
	if !r.canCall() {
		return nil, nil
	}

	system := r.promptBuilder.SystemPrompt()
	user := r.promptBuilder.UserPrompt(ps, personality, r.thoughts)

	response, err := r.llm.Complete(ctx, system, user)
	if err != nil {
		return nil, err
	}

	r.lastReview = time.Now()

	return &Observation{
		Content:   response,
		Timestamp: time.Now(),
	}, nil
}

func (r *Reviewer) canCall() bool {
	if r.minInterval == 0 {
		return true
	}
	return time.Since(r.lastReview) >= r.minInterval
}
