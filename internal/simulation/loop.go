package simulation

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/output"
	"github.com/marczahn/person/internal/psychology"
	"github.com/marczahn/person/internal/reviewer"
	"github.com/marczahn/person/internal/sense"
)

// InputType classifies user input into speech, action, or environment.
type InputType int

const (
	TypeSpeech      InputType = iota // plain text — someone talks to the person
	TypeAction                       // *text* — someone does something to/near the person
	TypeEnvironment                  // ~text — environmental change
)

// classifyInput determines the input type from text conventions:
// *text* → action, ~text → environment, plain text → speech.
func classifyInput(raw string) (InputType, string) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "*") && strings.HasSuffix(trimmed, "*") && len(trimmed) > 2 {
		return TypeAction, strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}
	if strings.HasPrefix(trimmed, "~") {
		return TypeEnvironment, strings.TrimSpace(trimmed[1:])
	}
	return TypeSpeech, trimmed
}

// Config holds all dependencies and settings for the simulation loop.
type Config struct {
	// Core components.
	BioProcessor   *biology.Processor
	PsychProcessor *psychology.Processor
	Consciousness  *consciousness.Engine
	SenseParser    sense.Parser
	Display        *output.Display
	Store          memory.Store

	// Optional meta-observer.
	Reviewer    *reviewer.Reviewer
	Personality *psychology.Personality

	// Initial state.
	BioState *biology.State
	Identity *memory.IdentityCore

	// Timing.
	TickInterval time.Duration // how often the main loop ticks (e.g., 100ms)
	SimStart     time.Time     // in-world start time

	// IO.
	Input io.Reader // stdin or test reader
}

// Loop orchestrates the simulation. It reads input, processes biology,
// computes psychology, triggers consciousness, and displays output.
type Loop struct {
	cfg   Config
	clock *Clock

	// Channel for input events from the reader goroutine.
	inputCh chan string
}

// NewLoop creates a simulation loop from the given configuration.
func NewLoop(cfg Config) *Loop {
	if cfg.TickInterval == 0 {
		cfg.TickInterval = 100 * time.Millisecond
	}
	return &Loop{
		cfg:     cfg,
		clock:   NewClock(cfg.SimStart),
		inputCh: make(chan string, 16),
	}
}

// Run starts the simulation loop. It blocks until the context is cancelled.
// On shutdown, it persists the current state.
func (l *Loop) Run(ctx context.Context) error {
	// Start input reader goroutine.
	go l.readInput(ctx)

	ticker := time.NewTicker(l.cfg.TickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return l.shutdown()
		case <-ticker.C:
			if l.clock.Paused() {
				continue
			}
			if err := l.tick(ctx); err != nil {
				l.cfg.Display.Show(output.Entry{
					Source:    output.Bio,
					Message:   fmt.Sprintf("tick error: %v", err),
					Timestamp: l.clock.Now(),
				})
			}
		}
	}
}

// Pause pauses the simulation clock.
func (l *Loop) Pause() {
	l.clock.Pause()
}

// Resume resumes the simulation clock.
func (l *Loop) Resume() {
	l.clock.Resume()
}

// tick performs one simulation cycle.
func (l *Loop) tick(ctx context.Context) error {
	dt := l.clock.Tick()
	if dt <= 0 {
		return nil
	}
	now := l.clock.Now()

	// 1. Drain pending input events.
	l.processInput(ctx, now)

	// 2. Biology tick (decay, circadian, interactions, thresholds).
	bioResult := l.cfg.BioProcessor.Tick(l.cfg.BioState)
	l.displayBioChanges(bioResult, now)

	// 3. Psychology: transform bio → psych state.
	psychState := l.cfg.PsychProcessor.Process(l.cfg.BioState, dt, 0.5)

	// 4. Consciousness: reactive (salience-gated).
	reactive, err := l.cfg.Consciousness.React(ctx, &psychState, dt)
	if err != nil {
		return fmt.Errorf("consciousness react: %w", err)
	}
	if reactive != nil {
		l.displayThought(reactive, now)
		l.applyFeedback(reactive, dt)
	}

	// 5. Consciousness: spontaneous thought.
	spontaneous, err := l.cfg.Consciousness.Spontaneous(ctx, &psychState)
	if err != nil {
		return fmt.Errorf("consciousness spontaneous: %w", err)
	}
	if spontaneous != nil {
		l.displayThought(spontaneous, now)
		l.applyFeedback(spontaneous, dt)
	}

	// 6. Psychologist reviewer (optional).
	l.runReviewer(ctx, reactive, spontaneous, &psychState, now)

	return nil
}

// processInput drains the input channel and routes each input based on type:
// speech/action → consciousness.Respond, environment → sensory parser only.
func (l *Loop) processInput(ctx context.Context, now time.Time) {
	for {
		select {
		case raw := <-l.inputCh:
			inputType, content := classifyInput(raw)
			switch inputType {
			case TypeSpeech:
				l.processSpeech(ctx, content, now)
			case TypeAction:
				l.processAction(ctx, content, now)
			case TypeEnvironment:
				l.processEnvironment(content, now)
			}
		default:
			return
		}
	}
}

// processSpeech handles spoken input — goes directly to consciousness.
func (l *Loop) processSpeech(ctx context.Context, content string, now time.Time) {
	l.cfg.Display.Show(output.Entry{
		Source:    output.Sense,
		Message:   fmt.Sprintf("speech: \"%s\"", content),
		Timestamp: now,
	})

	psychState := l.cfg.PsychProcessor.Process(l.cfg.BioState, 0, 0.5)
	input := consciousness.ExternalInput{
		Type:    consciousness.InputSpeech,
		Content: content,
	}

	thought, err := l.cfg.Consciousness.Respond(ctx, &psychState, input)
	if err != nil {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Mind,
			Message:   fmt.Sprintf("respond error: %v", err),
			Timestamp: now,
		})
		return
	}
	if thought != nil {
		l.displayThought(thought, now)
		l.applyFeedback(thought, 0)
	}
}

// processAction handles action input — biology effects + consciousness.
func (l *Loop) processAction(ctx context.Context, content string, now time.Time) {
	// Actions affect biology through the sensory parser.
	events := l.cfg.SenseParser.Parse(content)
	for _, event := range events {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Sense,
			Message:   fmt.Sprintf("%s: %s (intensity: %.1f)", event.Channel, event.Parsed, event.Intensity),
			Timestamp: now,
		})
		changes := l.cfg.BioProcessor.ProcessStimulus(l.cfg.BioState, event)
		for _, c := range biology.SignificantChanges(changes) {
			l.cfg.Display.Show(output.Entry{
				Source:    output.Bio,
				Message:   fmt.Sprintf("%s %+.2f (%s)", c.Variable, c.Delta, c.Source),
				Timestamp: now,
			})
		}
	}

	// Actions also trigger conscious response.
	psychState := l.cfg.PsychProcessor.Process(l.cfg.BioState, 0, 0.5)
	input := consciousness.ExternalInput{
		Type:    consciousness.InputAction,
		Content: content,
	}

	thought, err := l.cfg.Consciousness.Respond(ctx, &psychState, input)
	if err != nil {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Mind,
			Message:   fmt.Sprintf("respond error: %v", err),
			Timestamp: now,
		})
		return
	}
	if thought != nil {
		l.displayThought(thought, now)
		l.applyFeedback(thought, 0)
	}
}

// processEnvironment handles environmental input — sensory parser only,
// consciousness reacts via salience as before.
func (l *Loop) processEnvironment(content string, now time.Time) {
	events := l.cfg.SenseParser.Parse(content)
	for _, event := range events {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Sense,
			Message:   fmt.Sprintf("%s: %s (intensity: %.1f)", event.Channel, event.Parsed, event.Intensity),
			Timestamp: now,
		})
		changes := l.cfg.BioProcessor.ProcessStimulus(l.cfg.BioState, event)
		for _, c := range biology.SignificantChanges(changes) {
			l.cfg.Display.Show(output.Entry{
				Source:    output.Bio,
				Message:   fmt.Sprintf("%s %+.2f (%s)", c.Variable, c.Delta, c.Source),
				Timestamp: now,
			})
		}
	}
}

// runReviewer feeds thoughts to the reviewer and displays any observation.
func (l *Loop) runReviewer(
	ctx context.Context,
	reactive, spontaneous *consciousness.Thought,
	ps *psychology.State,
	now time.Time,
) {
	if l.cfg.Reviewer == nil {
		return
	}
	if reactive != nil {
		l.cfg.Reviewer.AddThought(*reactive)
	}
	if spontaneous != nil {
		l.cfg.Reviewer.AddThought(*spontaneous)
	}

	obs, err := l.cfg.Reviewer.Review(ctx, ps, l.cfg.Personality)
	if err != nil {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Review,
			Message:   fmt.Sprintf("review error: %v", err),
			Timestamp: now,
		})
		return
	}
	if obs != nil {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Review,
			Message:   obs.Content,
			Timestamp: now,
		})
	}
}

// displayBioChanges shows significant biological state changes.
func (l *Loop) displayBioChanges(result biology.TickResult, now time.Time) {
	for _, c := range biology.SignificantChanges(result.Changes) {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Bio,
			Message:   fmt.Sprintf("%s %+.2f (%s)", c.Variable, c.Delta, c.Source),
			Timestamp: now,
		})
	}
	for _, t := range result.Thresholds {
		l.cfg.Display.Show(output.Entry{
			Source:    output.Bio,
			Message:   fmt.Sprintf("THRESHOLD [%s] %s: %s", t.Condition, t.System, t.Description),
			Timestamp: now,
		})
	}
}

// displayThought shows a consciousness thought.
func (l *Loop) displayThought(thought *consciousness.Thought, now time.Time) {
	l.cfg.Display.ShowThought(output.Entry{
		Source:    output.Mind,
		Message:   fmt.Sprintf("[%s, trigger: %s] %s", thought.Type, thought.Trigger, thought.Content),
		Timestamp: now,
	})
}

// applyFeedback converts consciousness feedback into biological state changes.
func (l *Loop) applyFeedback(thought *consciousness.Thought, dt float64) {
	changes := consciousness.FeedbackToChanges(thought.Feedback)
	for _, c := range changes {
		delta := c.Delta * dt
		current := l.cfg.BioState.Get(c.Variable)
		l.cfg.BioState.Set(c.Variable, biology.ClampVariable(c.Variable, current+delta))
	}
}

// readInput reads lines from the configured input reader and sends them
// to the input channel. Exits when the context is cancelled or input ends.
func (l *Loop) readInput(ctx context.Context) {
	scanner := bufio.NewScanner(l.cfg.Input)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		select {
		case l.inputCh <- line:
		case <-ctx.Done():
			return
		}
	}
}

// shutdown persists the current state before exiting.
func (l *Loop) shutdown() error {
	if l.cfg.Store == nil {
		return nil
	}
	if err := l.cfg.Store.SaveBioState(l.cfg.BioState); err != nil {
		return fmt.Errorf("save bio state: %w", err)
	}
	if l.cfg.Identity != nil {
		if err := l.cfg.Store.SaveIdentityCore(l.cfg.Identity); err != nil {
			return fmt.Errorf("save identity: %w", err)
		}
	}
	return nil
}
