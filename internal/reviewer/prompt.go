package reviewer

import (
	"fmt"
	"strings"

	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/i18n"
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
	return strings.TrimSpace(i18n.T().Reviewer.SystemPrompt)
}

// UserPrompt builds the user prompt from the current psychological state,
// personality profile, and recent thoughts.
func (pb *PromptBuilder) UserPrompt(
	ps *psychology.State,
	personality *psychology.Personality,
	thoughts []consciousness.Thought,
) string {
	tr := i18n.T()
	var b strings.Builder

	// Psychological state summary.
	b.WriteString(tr.Reviewer.Labels.CurrentState + "\n")
	fmt.Fprintf(&b, tr.Reviewer.Labels.Arousal+"\n", ps.Arousal)
	fmt.Fprintf(&b, tr.Reviewer.Labels.Valence+"\n", ps.Valence)
	fmt.Fprintf(&b, tr.Reviewer.Labels.Energy+"\n", ps.Energy)
	fmt.Fprintf(&b, tr.Reviewer.Labels.CognitiveLoad+"\n", ps.CognitiveLoad)
	fmt.Fprintf(&b, tr.Reviewer.Labels.Regulation+"\n", ps.RegulationCapacity)

	// Active distortions.
	if len(ps.ActiveDistortions) > 0 {
		b.WriteString(tr.Reviewer.Labels.Distortions + "\n")
		for _, d := range ps.ActiveDistortions {
			fmt.Fprintf(&b, "- %s\n", d)
		}
	}

	// Active coping strategies.
	if len(ps.ActiveCoping) > 0 {
		b.WriteString(tr.Reviewer.Labels.CopingStrategies + "\n")
		for _, c := range ps.ActiveCoping {
			fmt.Fprintf(&b, "- %s\n", c)
		}
	}

	// Personality profile.
	b.WriteString(tr.Reviewer.Labels.PersonalityProfile + "\n")
	fmt.Fprintf(&b, "- Openness: %.2f\n", personality.Openness)
	fmt.Fprintf(&b, "- Conscientiousness: %.2f\n", personality.Conscientiousness)
	fmt.Fprintf(&b, "- Extraversion: %.2f\n", personality.Extraversion)
	fmt.Fprintf(&b, "- Agreeableness: %.2f\n", personality.Agreeableness)
	fmt.Fprintf(&b, "- Neuroticism: %.2f\n", personality.Neuroticism)

	// Recent thoughts.
	b.WriteString(tr.Reviewer.Labels.RecentThoughts + "\n")
	if len(thoughts) == 0 {
		b.WriteString(tr.Reviewer.Labels.NoThoughts + "\n")
	} else {
		for i, t := range thoughts {
			fmt.Fprintf(&b, "%d. [%s, trigger: %s] %s\n", i+1, t.Type, t.Trigger, t.Content)
		}
	}

	b.WriteString(tr.Reviewer.Labels.AnalysisQuestion)
	return b.String()
}
