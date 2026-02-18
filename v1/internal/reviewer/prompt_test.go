package reviewer

import (
	"strings"
	"testing"
	"time"

	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/psychology"
)

func TestSystemPrompt_ContainsClinicalRole(t *testing.T) {
	pb := NewPromptBuilder()
	sys := pb.SystemPrompt()

	if !strings.Contains(sys, "clinical psychologist") {
		t.Error("system prompt should establish clinical psychologist role")
	}
	if !strings.Contains(sys, "one-way mirror") {
		t.Error("system prompt should mention one-way mirror metaphor")
	}
	if !strings.Contains(sys, "Never address the person directly") {
		t.Error("system prompt should forbid direct address")
	}
}

func TestSystemPrompt_DoesNotRevealSimulation(t *testing.T) {
	pb := NewPromptBuilder()
	sys := strings.ToLower(pb.SystemPrompt())

	for _, word := range []string{"simulation", "simulated", "virtual", "artificial"} {
		if strings.Contains(sys, word) {
			t.Errorf("system prompt should not reveal simulation, found %q", word)
		}
	}
}

func TestUserPrompt_ContainsStateValues(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{
		Arousal:            0.75,
		Valence:            -0.30,
		Energy:             0.40,
		CognitiveLoad:      0.60,
		RegulationCapacity: 0.25,
	}
	personality := &psychology.Personality{
		Openness:          0.8,
		Conscientiousness: 0.3,
		Extraversion:      0.6,
		Agreeableness:     0.7,
		Neuroticism:       0.9,
	}

	prompt := pb.UserPrompt(ps, personality, nil)

	checks := []string{"0.75", "-0.30", "0.40", "0.60", "0.25"}
	for _, v := range checks {
		if !strings.Contains(prompt, v) {
			t.Errorf("user prompt should contain state value %s", v)
		}
	}
}

func TestUserPrompt_ContainsPersonalityValues(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{}
	personality := &psychology.Personality{
		Openness:          0.80,
		Conscientiousness: 0.30,
		Extraversion:      0.60,
		Agreeableness:     0.70,
		Neuroticism:       0.90,
	}

	prompt := pb.UserPrompt(ps, personality, nil)

	for _, label := range []string{"Openness", "Conscientiousness", "Extraversion", "Agreeableness", "Neuroticism"} {
		if !strings.Contains(prompt, label) {
			t.Errorf("user prompt should contain personality trait %s", label)
		}
	}
	if !strings.Contains(prompt, "0.90") {
		t.Error("user prompt should contain neuroticism value 0.90")
	}
}

func TestUserPrompt_ContainsThoughtContent(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{}
	personality := &psychology.Personality{Openness: 0.5, Conscientiousness: 0.5, Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5}

	thoughts := []consciousness.Thought{
		{Type: consciousness.Reactive, Trigger: "loud_noise", Content: "What was that?", Timestamp: time.Now()},
		{Type: consciousness.Spontaneous, Trigger: "goal_rehearsal", Content: "I need to finish the report.", Timestamp: time.Now()},
	}

	prompt := pb.UserPrompt(ps, personality, thoughts)

	if !strings.Contains(prompt, "What was that?") {
		t.Error("user prompt should contain first thought content")
	}
	if !strings.Contains(prompt, "I need to finish the report.") {
		t.Error("user prompt should contain second thought content")
	}
	if !strings.Contains(prompt, "reactive") {
		t.Error("user prompt should contain thought type")
	}
	if !strings.Contains(prompt, "loud_noise") {
		t.Error("user prompt should contain thought trigger")
	}
}

func TestUserPrompt_ContainsDistortions(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{
		ActiveDistortions: []psychology.Distortion{psychology.Catastrophizing, psychology.MindReading},
	}
	personality := &psychology.Personality{Openness: 0.5, Conscientiousness: 0.5, Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5}

	prompt := pb.UserPrompt(ps, personality, nil)

	if !strings.Contains(prompt, "catastrophizing") {
		t.Error("user prompt should contain catastrophizing distortion")
	}
	if !strings.Contains(prompt, "mind_reading") {
		t.Error("user prompt should contain mind_reading distortion")
	}
}

func TestUserPrompt_ContainsCopingStrategies(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{
		ActiveCoping: []psychology.CopingStrategy{psychology.Rumination, psychology.Suppression},
	}
	personality := &psychology.Personality{Openness: 0.5, Conscientiousness: 0.5, Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5}

	prompt := pb.UserPrompt(ps, personality, nil)

	if !strings.Contains(prompt, "rumination") {
		t.Error("user prompt should contain rumination coping strategy")
	}
	if !strings.Contains(prompt, "suppression") {
		t.Error("user prompt should contain suppression coping strategy")
	}
}

func TestUserPrompt_EmptyThoughtsShowsPlaceholder(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{}
	personality := &psychology.Personality{Openness: 0.5, Conscientiousness: 0.5, Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5}

	prompt := pb.UserPrompt(ps, personality, nil)

	if !strings.Contains(prompt, "no recent thoughts") {
		t.Error("user prompt should indicate no recent thoughts when buffer is empty")
	}
}

func TestUserPrompt_EndsWithQuestion(t *testing.T) {
	pb := NewPromptBuilder()
	ps := &psychology.State{}
	personality := &psychology.Personality{Openness: 0.5, Conscientiousness: 0.5, Extraversion: 0.5, Agreeableness: 0.5, Neuroticism: 0.5}

	prompt := pb.UserPrompt(ps, personality, nil)

	if !strings.Contains(prompt, "What patterns do you observe?") {
		t.Error("user prompt should end with the analysis question")
	}
}
