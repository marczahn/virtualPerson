package consciousness

import (
	"context"
	"fmt"
	"time"

	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/psychology"
)

// LLM abstracts the language model API for testability.
type LLM interface {
	// Complete sends a system prompt and user message, returning the response.
	Complete(ctx context.Context, systemPrompt, userMessage string) (string, error)
}

// Engine is the consciousness layer that generates first-person experience
// from psychological state using an LLM.
type Engine struct {
	llm             LLM
	promptBuilder   *PromptBuilder
	salience        *SalienceCalculator
	queue           *ThoughtQueue
	identity        *memory.IdentityCore
	contextSelector *memory.ContextSelector
	memories        []memory.EpisodicMemory

	// Rate limiting: minimum interval between LLM calls.
	minInterval     time.Duration
	lastCallTime    time.Time

	// Spontaneous thought interval.
	spontaneousInterval time.Duration
	lastSpontaneous     time.Time
}

// EngineConfig holds configuration for the consciousness engine.
type EngineConfig struct {
	LLM                 LLM
	Identity            *memory.IdentityCore
	MaxPromptTokens     int
	MaxContextMemories  int
	MinCallInterval     time.Duration // 0 means no rate limit
	SpontaneousInterval time.Duration // 0 means no rate limit
}

// NewEngine creates a consciousness engine.
func NewEngine(cfg EngineConfig) *Engine {
	if cfg.MaxPromptTokens == 0 {
		cfg.MaxPromptTokens = 2000
	}
	if cfg.MaxContextMemories == 0 {
		cfg.MaxContextMemories = 5
	}
	// MinCallInterval and SpontaneousInterval default to 0 (no rate limit).
	// The simulation loop should set appropriate values for production use.

	return &Engine{
		llm:                 cfg.LLM,
		promptBuilder:       NewPromptBuilder(cfg.MaxPromptTokens),
		salience:            NewSalienceCalculator(),
		queue:               NewThoughtQueue(),
		identity:            cfg.Identity,
		contextSelector:     memory.NewContextSelector(cfg.MaxContextMemories),
		minInterval:         cfg.MinCallInterval,
		spontaneousInterval: cfg.SpontaneousInterval,
	}
}

// React checks if the current psychological state change is salient enough
// to trigger a reactive conscious thought. Returns nil if no thought is triggered.
func (e *Engine) React(ctx context.Context, ps *psychology.State, dt float64) (*Thought, error) {
	result := e.salience.Compute(ps, dt)
	if !result.Exceeded {
		return nil, nil
	}

	if !e.canCall() {
		return nil, nil
	}

	trigger := fmt.Sprintf("%s changed significantly", result.Trigger)
	distCtx := DistortionContext(ps.ActiveDistortions)
	relevant := e.selectMemories(ps)

	systemPrompt := e.promptBuilder.SystemPrompt(e.identity)
	userMessage := e.promptBuilder.ReactivePrompt(ps, trigger, relevant, distCtx)

	response, err := e.llm.Complete(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("reactive thought: %w", err)
	}

	e.lastCallTime = time.Now()
	e.queue.ExitAbsorption()

	feedback := ParseFeedback(response)

	return &Thought{
		Type:      Reactive,
		Priority:  PriorityPredictionError,
		Content:   response,
		Trigger:   trigger,
		Timestamp: time.Now(),
		Feedback:  feedback,
	}, nil
}

// Spontaneous generates a spontaneous thought if enough time has passed
// since the last one. Returns nil if no thought is generated.
func (e *Engine) Spontaneous(ctx context.Context, ps *psychology.State) (*Thought, error) {
	if time.Since(e.lastSpontaneous) < e.spontaneousInterval {
		return nil, nil
	}

	if !e.canCall() {
		return nil, nil
	}

	e.queue.UpdateNeeds(ps)
	candidate := e.queue.SelectSpontaneous(ps)
	if candidate == nil {
		return nil, nil
	}

	distCtx := DistortionContext(ps.ActiveDistortions)
	relevant := e.selectMemories(ps)

	systemPrompt := e.promptBuilder.SystemPrompt(e.identity)
	userMessage := e.promptBuilder.SpontaneousPrompt(ps, candidate, relevant, distCtx)

	response, err := e.llm.Complete(ctx, systemPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("spontaneous thought: %w", err)
	}

	e.lastCallTime = time.Now()
	e.lastSpontaneous = time.Now()

	feedback := ParseFeedback(response)

	return &Thought{
		Type:      Spontaneous,
		Priority:  candidate.Priority,
		Content:   response,
		Trigger:   candidate.Category,
		Timestamp: time.Now(),
		Feedback:  feedback,
	}, nil
}

// UpdateMemories refreshes the engine's memory cache.
func (e *Engine) UpdateMemories(memories []memory.EpisodicMemory) {
	e.memories = memories
}

// UpdateIdentity refreshes the identity core.
func (e *Engine) UpdateIdentity(ic *memory.IdentityCore) {
	e.identity = ic
}

// Queue returns the thought queue for external manipulation.
func (e *Engine) Queue() *ThoughtQueue {
	return e.queue
}

// Salience returns the salience calculator for external inspection.
func (e *Engine) Salience() *SalienceCalculator {
	return e.salience
}

func (e *Engine) canCall() bool {
	return time.Since(e.lastCallTime) >= e.minInterval
}

func (e *Engine) selectMemories(ps *psychology.State) []memory.EpisodicMemory {
	if len(e.memories) == 0 {
		return nil
	}

	current := memory.BioSnapshot{
		Arousal: ps.Arousal,
		Valence: ps.Valence,
		Fatigue: 1.0 - ps.Energy,
	}

	return e.contextSelector.Select(e.memories, current)
}
