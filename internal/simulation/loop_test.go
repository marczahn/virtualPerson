package simulation

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/output"
	"github.com/marczahn/person/internal/psychology"
	"github.com/marczahn/person/internal/reviewer"
	"github.com/marczahn/person/internal/sense"
)

type mockLLM struct {
	response string
}

func (m *mockLLM) Complete(_ context.Context, _, _ string) (string, error) {
	return m.response, nil
}

func ptrBioState(s biology.State) *biology.State { return &s }

func newTestLoop(inputReader io.Reader, displayBuf *bytes.Buffer) *Loop {
	bioState := ptrBioState(biology.NewDefaultState())
	personality := psychology.Personality{
		Openness:          0.5,
		Conscientiousness: 0.5,
		Extraversion:      0.5,
		Agreeableness:     0.5,
		Neuroticism:       0.5,
	}

	identity := &memory.IdentityCore{
		SelfNarrative: "I am a test person.",
	}

	consciousnessEngine := consciousness.NewEngine(consciousness.EngineConfig{
		LLM:                 &mockLLM{response: "I feel something changing."},
		Identity:            identity,
		MinCallInterval:     0,
		SpontaneousInterval: 10 * time.Second,
	})

	display := output.NewDisplay(displayBuf, false)

	cfg := Config{
		BioProcessor:   biology.NewProcessor(),
		PsychProcessor: psychology.NewProcessor(personality),
		Consciousness:  consciousnessEngine,
		SenseParser:    sense.NewKeywordParser(),
		Display:        display,
		BioState:       bioState,
		Identity:       identity,
		TickInterval:   50 * time.Millisecond,
		SimStart:       time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		Input:          inputReader,
	}

	return NewLoop(cfg)
}

func TestLoop_StartsAndStops(t *testing.T) {
	var buf bytes.Buffer
	reader := strings.NewReader("")
	loop := newTestLoop(reader, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error on clean shutdown, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("loop did not shut down within 2 seconds")
	}
}

func TestLoop_ProcessesInput(t *testing.T) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()
	loop := newTestLoop(pr, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	_, err := pw.Write([]byte("~a freezing wind blows\n"))
	if err != nil {
		t.Fatalf("write to pipe: %v", err)
	}

	time.Sleep(300 * time.Millisecond)
	cancel()
	pw.Close()
	<-done

	out := buf.String()
	if !strings.Contains(out, "SENSE") {
		t.Errorf("expected SENSE tag in output, got: %s", out)
	}
	if !strings.Contains(out, "cold") || !strings.Contains(out, "thermal") {
		t.Errorf("expected thermal/cold reference in output, got: %s", out)
	}
}

func TestLoop_PauseStopsTicks(t *testing.T) {
	var buf bytes.Buffer
	reader := strings.NewReader("")
	loop := newTestLoop(reader, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	loop.Pause()
	time.Sleep(50 * time.Millisecond) // Let any in-flight tick finish.
	outputAfterPause := buf.Len()

	time.Sleep(200 * time.Millisecond)
	outputAfterWait := buf.Len()

	cancel()
	<-done

	if outputAfterWait > outputAfterPause+200 {
		t.Errorf("output grew significantly while paused: before=%d, after=%d",
			outputAfterPause, outputAfterWait)
	}
}

func TestLoop_ResumeAfterPause(t *testing.T) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()
	loop := newTestLoop(pr, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	loop.Pause()
	time.Sleep(100 * time.Millisecond)
	outputWhilePaused := buf.Len()

	loop.Resume()
	// After resume, send input to guarantee output growth.
	pw.Write([]byte("~a loud explosion\n"))
	time.Sleep(200 * time.Millisecond)

	cancel()
	pw.Close()
	<-done

	if buf.Len() <= outputWhilePaused {
		t.Error("expected output to grow after resume with input")
	}
}

func TestLoop_ShutdownPersistsState(t *testing.T) {
	var buf bytes.Buffer
	reader := strings.NewReader("")

	bioState := ptrBioState(biology.NewDefaultState())
	bioState.Cortisol = 0.8

	personality := psychology.Personality{
		Openness: 0.5, Conscientiousness: 0.5,
		Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5,
	}
	identity := &memory.IdentityCore{SelfNarrative: "Test person."}

	store, err := memory.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	defer store.Close()

	cfg := Config{
		BioProcessor:   biology.NewProcessor(),
		PsychProcessor: psychology.NewProcessor(personality),
		Consciousness: consciousness.NewEngine(consciousness.EngineConfig{
			LLM:                 &mockLLM{response: "ok"},
			Identity:            identity,
			MinCallInterval:     0,
			SpontaneousInterval: 10 * time.Second,
		}),
		SenseParser:  sense.NewKeywordParser(),
		Display:      output.NewDisplay(&buf, false),
		Store:        store,
		BioState:     bioState,
		Identity:     identity,
		TickInterval: 50 * time.Millisecond,
		SimStart:     time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		Input:        reader,
	}

	loop := NewLoop(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	loaded, err := store.LoadBioState()
	if err != nil {
		t.Fatalf("load bio state: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected bio state to be persisted")
	}
	if loaded.Cortisol < 0.5 {
		t.Errorf("expected persisted cortisol to be elevated, got %f", loaded.Cortisol)
	}

	loadedIdentity, err := store.LoadIdentityCore()
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}
	if loadedIdentity == nil {
		t.Fatal("expected identity to be persisted")
	}
	if loadedIdentity.SelfNarrative != "Test person." {
		t.Errorf("expected identity narrative 'Test person.', got %q", loadedIdentity.SelfNarrative)
	}
}

func TestNewLoop_DefaultTickInterval(t *testing.T) {
	loop := NewLoop(Config{
		Input: strings.NewReader(""),
	})
	if loop.cfg.TickInterval != 100*time.Millisecond {
		t.Errorf("expected default tick interval 100ms, got %v", loop.cfg.TickInterval)
	}
}

func TestClassifyInput_Speech(t *testing.T) {
	inputType, content := classifyInput("Hello, how are you?")
	if inputType != TypeSpeech {
		t.Errorf("expected TypeSpeech, got %d", inputType)
	}
	if content != "Hello, how are you?" {
		t.Errorf("content = %q, expected unchanged text", content)
	}
}

func TestClassifyInput_Action(t *testing.T) {
	inputType, content := classifyInput("*pushes you gently*")
	if inputType != TypeAction {
		t.Errorf("expected TypeAction, got %d", inputType)
	}
	if content != "pushes you gently" {
		t.Errorf("content = %q, expected stripped asterisks", content)
	}
}

func TestClassifyInput_Environment(t *testing.T) {
	inputType, content := classifyInput("~a cold wind blows")
	if inputType != TypeEnvironment {
		t.Errorf("expected TypeEnvironment, got %d", inputType)
	}
	if content != "a cold wind blows" {
		t.Errorf("content = %q, expected stripped tilde", content)
	}
}

func TestClassifyInput_EmptyAsterisks(t *testing.T) {
	inputType, _ := classifyInput("**")
	if inputType != TypeSpeech {
		t.Errorf("expected TypeSpeech for empty asterisks, got %d", inputType)
	}
}

func TestClassifyInput_OnlyTilde(t *testing.T) {
	inputType, content := classifyInput("~")
	if inputType != TypeEnvironment {
		t.Errorf("expected TypeEnvironment for bare tilde, got %d", inputType)
	}
	if content != "" {
		t.Errorf("expected empty content for bare tilde, got %q", content)
	}
}

func TestClassifyInput_WhitespaceHandling(t *testing.T) {
	inputType, content := classifyInput("  *waves hello*  ")
	if inputType != TypeAction {
		t.Errorf("expected TypeAction with surrounding whitespace, got %d", inputType)
	}
	if content != "waves hello" {
		t.Errorf("content = %q, expected trimmed action text", content)
	}
}

func TestLoop_ReviewerObservationAppearsInOutput(t *testing.T) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()

	bioState := ptrBioState(biology.NewDefaultState())
	personality := psychology.Personality{
		Openness: 0.5, Conscientiousness: 0.5,
		Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5,
	}
	identity := &memory.IdentityCore{SelfNarrative: "I am a test person."}

	rev := reviewer.NewReviewer(reviewer.ReviewerConfig{
		LLM:         &mockLLM{response: "Subject shows signs of acute stress response."},
		MinInterval: 0,
		MaxThoughts: 10,
	})

	cfg := Config{
		BioProcessor:   biology.NewProcessor(),
		PsychProcessor: psychology.NewProcessor(personality),
		Consciousness: consciousness.NewEngine(consciousness.EngineConfig{
			LLM:                 &mockLLM{response: "I feel startled."},
			Identity:            identity,
			MinCallInterval:     0,
			SpontaneousInterval: 10 * time.Second,
		}),
		SenseParser: sense.NewKeywordParser(),
		Display:     output.NewDisplay(&buf, false),
		BioState:    bioState,
		Identity:    identity,
		Reviewer:    rev,
		Personality: &personality,
		TickInterval: 50 * time.Millisecond,
		SimStart:     time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		Input:        pr,
	}

	loop := NewLoop(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	// Send speech input to trigger a thought, which feeds the reviewer.
	pw.Write([]byte("Hello!\n"))
	time.Sleep(300 * time.Millisecond)
	cancel()
	pw.Close()
	<-done

	out := buf.String()
	if !strings.Contains(out, "REVIEW") {
		t.Errorf("expected REVIEW tag in output, got:\n%s", out)
	}
	if !strings.Contains(out, "acute stress response") {
		t.Errorf("expected reviewer observation in output, got:\n%s", out)
	}
}

func TestLoop_SpeechInput_ProducesMindOutput(t *testing.T) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()
	loop := newTestLoop(pr, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	pw.Write([]byte("Hello there!\n"))
	time.Sleep(300 * time.Millisecond)
	cancel()
	pw.Close()
	<-done

	out := buf.String()
	if !strings.Contains(out, "SENSE") {
		t.Errorf("expected SENSE tag for speech, got: %s", out)
	}
	if !strings.Contains(out, "MIND") {
		t.Errorf("expected MIND tag for speech response, got: %s", out)
	}
	if !strings.Contains(out, "speech") {
		t.Errorf("expected 'speech' in sense output, got: %s", out)
	}
}

func TestLoop_ActionInput_ProducesBioAndMindOutput(t *testing.T) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()
	loop := newTestLoop(pr, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	// Action with a keyword the sensory parser recognizes.
	pw.Write([]byte("*throws cold water on you*\n"))
	time.Sleep(300 * time.Millisecond)
	cancel()
	pw.Close()
	<-done

	out := buf.String()
	if !strings.Contains(out, "SENSE") {
		t.Errorf("expected SENSE tag for action, got: %s", out)
	}
	if !strings.Contains(out, "MIND") {
		t.Errorf("expected MIND tag for action consciousness response, got: %s", out)
	}
}

func TestLoop_EnvironmentInput_NoBioOnly(t *testing.T) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()
	loop := newTestLoop(pr, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- loop.Run(ctx)
	}()

	pw.Write([]byte("~a freezing wind blows\n"))
	time.Sleep(300 * time.Millisecond)
	cancel()
	pw.Close()
	<-done

	out := buf.String()
	if !strings.Contains(out, "SENSE") {
		t.Errorf("expected SENSE tag for environment, got: %s", out)
	}
	// Environment input should NOT directly trigger MIND output
	// (consciousness only reacts via salience in the tick, not here).
	// We can't easily guarantee no MIND output since ticks may trigger it,
	// but we verify the sensory path works.
	if !strings.Contains(out, "cold") || !strings.Contains(out, "thermal") {
		t.Errorf("expected thermal/cold reference in output, got: %s", out)
	}
}
