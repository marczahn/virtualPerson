package reviewer

import (
	"fmt"
	"strings"

	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/psychology"
)

// PromptBuilder constructs prompts for the psychologist reviewer LLM calls.
type PromptBuilder struct{}

// NewPromptBuilder creates a new reviewer prompt builder.
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// SystemPrompt returns the system prompt establishing the reviewer's persona.
func (pb *PromptBuilder) SystemPrompt() string {
	return "You are a clinical psychologist observing a person through a one-way mirror. " +
		"You see their inner thoughts and physiological state. " +
		"Provide brief, insightful observations about psychological patterns you notice. " +
		"Be concise — 2-3 sentences max. " +
		"Use clinical language but keep it accessible. " +
		"Never address the person directly."
}

// UserPrompt builds the user prompt from the current psychological state,
// personality profile, and recent thoughts.
func (pb *PromptBuilder) UserPrompt(
	ps *psychology.State,
	personality *psychology.Personality,
	thoughts []consciousness.Thought,
) string {
	var b strings.Builder

	// Psychological state summary.
	b.WriteString("## Current Psychological State\n")
	fmt.Fprintf(&b, "- Arousal: %.2f (0=calm, 1=highly activated)\n", ps.Arousal)
	fmt.Fprintf(&b, "- Valence: %.2f (-1=very negative, 1=very positive)\n", ps.Valence)
	fmt.Fprintf(&b, "- Energy: %.2f (0=exhausted, 1=energized)\n", ps.Energy)
	fmt.Fprintf(&b, "- Cognitive load: %.2f (0=clear, 1=overwhelmed)\n", ps.CognitiveLoad)
	fmt.Fprintf(&b, "- Regulation capacity: %.2f (0=depleted, 1=full)\n", ps.RegulationCapacity)

	// Active distortions.
	if len(ps.ActiveDistortions) > 0 {
		b.WriteString("\n## Active Cognitive Distortions\n")
		for _, d := range ps.ActiveDistortions {
			fmt.Fprintf(&b, "- %s\n", d)
		}
	}

	// Active coping strategies.
	if len(ps.ActiveCoping) > 0 {
		b.WriteString("\n## Active Coping Strategies\n")
		for _, c := range ps.ActiveCoping {
			fmt.Fprintf(&b, "- %s\n", c)
		}
	}

	// Personality profile.
	b.WriteString("\n## Personality Profile (Big Five)\n")
	fmt.Fprintf(&b, "- Openness: %.2f\n", personality.Openness)
	fmt.Fprintf(&b, "- Conscientiousness: %.2f\n", personality.Conscientiousness)
	fmt.Fprintf(&b, "- Extraversion: %.2f\n", personality.Extraversion)
	fmt.Fprintf(&b, "- Agreeableness: %.2f\n", personality.Agreeableness)
	fmt.Fprintf(&b, "- Neuroticism: %.2f\n", personality.Neuroticism)

	// Recent thoughts.
	b.WriteString("\n## Recent Thoughts\n")
	if len(thoughts) == 0 {
		b.WriteString("(no recent thoughts)\n")
	} else {
		for i, t := range thoughts {
			fmt.Fprintf(&b, "%d. [%s, trigger: %s] %s\n", i+1, t.Type, t.Trigger, t.Content)
		}
	}

	b.WriteString("\nWhat patterns do you observe? Any concerns?")
	return b.String()
}
