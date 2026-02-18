package consciousness

import (
	"strings"
	"testing"

	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/psychology"
)

func TestSystemPrompt_ContainsIdentity(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ic := &memory.IdentityCore{
		SelfNarrative:     "I am a quiet person who values solitude.",
		DispositionTraits: []string{"introverted", "thoughtful"},
		RelationalMarkers: []string{"close to my sister"},
		KeyMemories:       []string{"childhood in the mountains"},
		EmotionalPatterns: []string{"tends to withdraw under stress"},
		ValuesCommitments: []string{"honesty above all"},
	}

	prompt := pb.SystemPrompt(ic, "")

	if !strings.Contains(prompt, "quiet person") {
		t.Error("system prompt should contain self-narrative")
	}
	if !strings.Contains(prompt, "introverted") {
		t.Error("system prompt should contain disposition traits")
	}
	if !strings.Contains(prompt, "close to my sister") {
		t.Error("system prompt should contain relational markers")
	}
	if !strings.Contains(prompt, "childhood in the mountains") {
		t.Error("system prompt should contain key memories")
	}
	if !strings.Contains(prompt, "honesty above all") {
		t.Error("system prompt should contain values")
	}
}

func TestSystemPrompt_NoIdentity(t *testing.T) {
	pb := NewPromptBuilder(2000)
	prompt := pb.SystemPrompt(nil, "")

	if !strings.Contains(prompt, "You are a person") {
		t.Error("system prompt should contain base consciousness framing")
	}
	if strings.Contains(prompt, "Who You Are") {
		t.Error("should not have identity section when nil")
	}
}

func TestSystemPrompt_NeverRevealsSimulation(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ic := &memory.IdentityCore{SelfNarrative: "I exist."}
	prompt := pb.SystemPrompt(ic, "")

	forbidden := []string{"simulation", "simulated", "LLM", "AI", "artificial", "program", "code"}
	for _, word := range forbidden {
		if strings.Contains(strings.ToLower(prompt), word) {
			t.Errorf("system prompt contains forbidden word %q", word)
		}
	}
}

func TestSystemPrompt_IncludesScenario(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ic := &memory.IdentityCore{SelfNarrative: "I exist."}

	prompt := pb.SystemPrompt(ic, "a small quiet apartment, soft morning light through the curtains")

	if !strings.Contains(prompt, "small quiet apartment") {
		t.Error("system prompt should contain scenario text")
	}
	if !strings.Contains(prompt, "Where you are") {
		t.Error("system prompt should contain scenario header")
	}
}

func TestSystemPrompt_EmptyScenarioOmitsBlock(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ic := &memory.IdentityCore{SelfNarrative: "I exist."}

	prompt := pb.SystemPrompt(ic, "")

	if strings.Contains(prompt, "Where you are") {
		t.Error("system prompt should not contain scenario header when scenario is empty")
	}
}

func TestReactivePrompt_ContainsTrigger(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Arousal: 0.8, Valence: -0.3, Energy: 0.5}

	prompt := pb.ReactivePrompt(ps, "arousal changed significantly", nil, "", nil)

	if !strings.Contains(prompt, "arousal changed significantly") {
		t.Error("reactive prompt should contain trigger")
	}
}

func TestReactivePrompt_ContainsDistortions(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Arousal: 0.5}
	distCtx := DistortionContext([]psychology.Distortion{psychology.Catastrophizing})

	prompt := pb.ReactivePrompt(ps, "test", nil, distCtx, nil)

	if !strings.Contains(prompt, "worst possible outcome") {
		t.Error("prompt should contain distortion description")
	}
}

func TestReactivePrompt_ContainsMemories(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{}
	memories := []memory.EpisodicMemory{
		{Content: "I felt the cold wind on my face"},
	}

	prompt := pb.ReactivePrompt(ps, "", memories, "", nil)

	if !strings.Contains(prompt, "cold wind") {
		t.Error("prompt should contain memory content")
	}
}

func TestStateBlock_HighArousal(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Arousal: 0.8, Energy: 0.5}

	block := pb.stateBlock(ps)

	if !strings.Contains(block, "heart pounding") {
		t.Error("high arousal should mention heart pounding")
	}
}

func TestStateBlock_LowEnergy(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Energy: 0.1}

	block := pb.stateBlock(ps)

	if !strings.Contains(block, "exhausted") {
		t.Error("very low energy should mention exhaustion")
	}
}

func TestStateBlock_NegativeValence(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Valence: -0.6}

	block := pb.stateBlock(ps)

	if !strings.Contains(block, "bad") {
		t.Error("very negative valence should describe things as bad")
	}
}

func TestStateBlock_IsolationPhases(t *testing.T) {
	pb := NewPromptBuilder(2000)

	tests := []struct {
		phase   psychology.IsolationPhase
		keyword string
	}{
		{psychology.IsolationBoredom, "restless"},
		{psychology.IsolationLoneliness, "lonely"},
		{psychology.IsolationSevere, "unbearable"},
	}

	for _, tt := range tests {
		ps := &psychology.State{
			Isolation: psychology.IsolationState{Phase: tt.phase},
		}
		block := pb.stateBlock(ps)
		if !strings.Contains(block, tt.keyword) {
			t.Errorf("phase %s: expected keyword %q in state block", tt.phase, tt.keyword)
		}
	}
}

func TestDistortionContext_Empty(t *testing.T) {
	got := DistortionContext(nil)
	if got != "" {
		t.Errorf("expected empty string for no distortions, got %q", got)
	}
}

func TestDistortionContext_Multiple(t *testing.T) {
	got := DistortionContext([]psychology.Distortion{
		psychology.Catastrophizing,
		psychology.EmotionalReasoning,
	})

	if !strings.Contains(got, "worst possible outcome") {
		t.Error("should contain catastrophizing description")
	}
	if !strings.Contains(got, "feelings as evidence") {
		t.Error("should contain emotional reasoning description")
	}
}

func TestSpontaneousPrompt_ContainsCandidate(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Energy: 0.5}
	candidate := &ThoughtCandidate{
		Priority: PriorityBiologicalNeed,
		Category: "biological_need",
		Prompt:   "You are hungry and your stomach is growling.",
	}

	prompt := pb.SpontaneousPrompt(ps, candidate, nil, "", nil)

	if !strings.Contains(prompt, "stomach is growling") {
		t.Error("spontaneous prompt should contain candidate prompt")
	}
}

func TestExternalInputPrompt_Speech(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Arousal: 0.3, Valence: 0.2, Energy: 0.5}
	input := ExternalInput{Type: InputSpeech, Content: "How are you feeling?"}

	prompt := pb.ExternalInputPrompt(ps, input, nil, "", nil)

	if !strings.Contains(prompt, "Someone says to you") {
		t.Error("speech prompt should contain speech framing")
	}
	if !strings.Contains(prompt, "How are you feeling?") {
		t.Error("speech prompt should contain the spoken words")
	}
	if !strings.Contains(prompt, "Your current experience") {
		t.Error("speech prompt should contain state block")
	}
}

func TestExternalInputPrompt_Action(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Arousal: 0.5, Valence: -0.2, Energy: 0.5}
	input := ExternalInput{Type: InputAction, Content: "gives you a warm blanket"}

	prompt := pb.ExternalInputPrompt(ps, input, nil, "", nil)

	if !strings.Contains(prompt, "Someone does this") {
		t.Error("action prompt should contain action framing")
	}
	if !strings.Contains(prompt, "gives you a warm blanket") {
		t.Error("action prompt should contain the action description")
	}
}

func TestExternalInputPrompt_IncludesMemories(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{}
	input := ExternalInput{Type: InputSpeech, Content: "hello"}
	memories := []memory.EpisodicMemory{
		{Content: "The stranger seemed kind"},
	}

	prompt := pb.ExternalInputPrompt(ps, input, memories, "", nil)

	if !strings.Contains(prompt, "stranger seemed kind") {
		t.Error("prompt should include memory content")
	}
}

func TestExternalInputPrompt_IncludesDistortions(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{}
	input := ExternalInput{Type: InputSpeech, Content: "hello"}
	distCtx := DistortionContext([]psychology.Distortion{psychology.MindReading})

	prompt := pb.ExternalInputPrompt(ps, input, nil, distCtx, nil)

	if !strings.Contains(prompt, "what others are thinking") {
		t.Error("prompt should contain distortion description")
	}
}

func TestTrimMemories_RespectsLimit(t *testing.T) {
	pb := NewPromptBuilder(300) // low budget → ~3 memories max

	memories := make([]memory.EpisodicMemory, 20)
	for i := range memories {
		memories[i] = memory.EpisodicMemory{Content: "memory"}
	}

	trimmed := pb.trimMemories(memories)
	if len(trimmed) > 10 {
		t.Errorf("trimmed to %d, expected <= 10", len(trimmed))
	}
}

func TestThoughtStreamBlock_Empty(t *testing.T) {
	pb := NewPromptBuilder(2000)
	got := pb.thoughtStreamBlock(nil)
	if got != "" {
		t.Errorf("expected empty string for nil thoughts, got %q", got)
	}
}

func TestThoughtStreamBlock_IncludesContent(t *testing.T) {
	pb := NewPromptBuilder(2000)
	thoughts := []Thought{
		{Content: "I feel uneasy about this"},
		{Content: "Maybe it's nothing..."},
	}

	got := pb.thoughtStreamBlock(thoughts)

	if !strings.Contains(got, "What you've been thinking") {
		t.Error("should contain section header")
	}
	if !strings.Contains(got, "I feel uneasy about this") {
		t.Error("should contain first thought")
	}
	if !strings.Contains(got, "Maybe it's nothing...") {
		t.Error("should contain second thought")
	}
}

func TestReactivePrompt_IncludesRecentThoughts(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Arousal: 0.5}
	thoughts := []Thought{
		{Content: "Something feels off today"},
	}

	prompt := pb.ReactivePrompt(ps, "test", nil, "", thoughts)

	if !strings.Contains(prompt, "Something feels off today") {
		t.Error("reactive prompt should include recent thoughts")
	}
}

func TestExternalInputPrompt_IncludesRecentThoughts(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{}
	input := ExternalInput{Type: InputSpeech, Content: "hello"}
	thoughts := []Thought{
		{Content: "I was just thinking about lunch"},
	}

	prompt := pb.ExternalInputPrompt(ps, input, nil, "", thoughts)

	if !strings.Contains(prompt, "thinking about lunch") {
		t.Error("external input prompt should include recent thoughts")
	}
}

func TestSpontaneousPrompt_IncludesRecentThoughts(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ps := &psychology.State{Energy: 0.5}
	thoughts := []Thought{
		{Content: "The silence is getting to me"},
	}

	prompt := pb.SpontaneousPrompt(ps, nil, nil, "", thoughts)

	if !strings.Contains(prompt, "silence is getting to me") {
		t.Error("spontaneous prompt should include recent thoughts")
	}
}

func TestSystemPrompt_ContainsAnnotationInstruction(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ic := &memory.IdentityCore{SelfNarrative: "I exist."}
	prompt := pb.SystemPrompt(ic, "")

	if !strings.Contains(prompt, "[STATE:") {
		t.Error("system prompt should contain annotation instruction with [STATE: tag")
	}
}

func TestSystemPrompt_NeverRevealsSimulation_StillPasses(t *testing.T) {
	pb := NewPromptBuilder(2000)
	ic := &memory.IdentityCore{SelfNarrative: "I exist."}
	prompt := pb.SystemPrompt(ic, "")

	forbidden := []string{"simulation", "simulated", "LLM", "AI", "artificial", "program", "code"}
	for _, word := range forbidden {
		if strings.Contains(strings.ToLower(prompt), word) {
			t.Errorf("system prompt contains forbidden word %q", word)
		}
	}
}
