package consciousness

import (
	"fmt"
	"strings"

	"github.com/marczahn/person/internal/i18n"
	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/psychology"
)

// PromptBuilder constructs the LLM prompt from psychological state,
// memories, and identity core.
type PromptBuilder struct {
	maxTokenEstimate int // approximate token budget for the full prompt
}

// NewPromptBuilder creates a builder with the given token budget.
func NewPromptBuilder(maxTokens int) *PromptBuilder {
	return &PromptBuilder{maxTokenEstimate: maxTokens}
}

// SystemPrompt returns the system prompt that establishes the consciousness role.
// The simulated person must not know they are a simulation.
// scenario describes the physical environment the person is in; empty means no scenario block.
func (pb *PromptBuilder) SystemPrompt(identity *memory.IdentityCore, scenario string) string {
	tr := i18n.T()
	var b strings.Builder

	b.WriteString(tr.Consciousness.SystemPrompt)

	if identity != nil {
		b.WriteString(tr.Consciousness.Identity.Header)
		if identity.SelfNarrative != "" {
			b.WriteString(identity.SelfNarrative)
			b.WriteString("\n")
		}
		if len(identity.DispositionTraits) > 0 {
			b.WriteString(tr.Consciousness.Identity.Tendencies)
			b.WriteString(strings.Join(identity.DispositionTraits, ". "))
			b.WriteString(".\n")
		}
		if len(identity.RelationalMarkers) > 0 {
			b.WriteString(tr.Consciousness.Identity.Relationships)
			b.WriteString(strings.Join(identity.RelationalMarkers, ". "))
			b.WriteString(".\n")
		}
		if len(identity.EmotionalPatterns) > 0 {
			b.WriteString(tr.Consciousness.Identity.Patterns)
			b.WriteString(strings.Join(identity.EmotionalPatterns, ". "))
			b.WriteString(".\n")
		}
		if len(identity.ValuesCommitments) > 0 {
			b.WriteString(tr.Consciousness.Identity.Values)
			b.WriteString(strings.Join(identity.ValuesCommitments, ". "))
			b.WriteString(".\n")
		}
		if len(identity.KeyMemories) > 0 {
			b.WriteString(tr.Consciousness.Identity.Memories)
			for _, m := range identity.KeyMemories {
				b.WriteString("- ")
				b.WriteString(m)
				b.WriteString("\n")
			}
		}
	}

	if scenario != "" {
		b.WriteString(tr.Consciousness.ScenarioHeader)
		b.WriteString(scenario)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(tr.Consciousness.EmotionalAnnotation)

	return b.String()
}

// ReactivePrompt builds the user message for a reactive thought triggered
// by a salience breach.
func (pb *PromptBuilder) ReactivePrompt(
	ps *psychology.State,
	trigger string,
	recentMemories []memory.EpisodicMemory,
	distortionContext string,
	recentThoughts []Thought,
) string {
	tr := i18n.T()
	var b strings.Builder

	b.WriteString(pb.stateBlock(ps))
	b.WriteString("\n")

	b.WriteString(pb.thoughtStreamBlock(recentThoughts))

	if trigger != "" {
		b.WriteString(fmt.Sprintf(tr.Consciousness.Prompts.TriggerShifted, trigger))
	}

	if distortionContext != "" {
		b.WriteString(distortionContext)
		b.WriteString("\n")
	}

	if len(recentMemories) > 0 {
		b.WriteString(tr.Consciousness.State.RecentExperiences)
		for _, m := range pb.trimMemories(recentMemories) {
			b.WriteString(fmt.Sprintf("- %s\n", m.Content))
		}
		b.WriteString("\n")
	}

	b.WriteString(tr.Consciousness.Prompts.ReactiveQuestion)

	return b.String()
}

// SpontaneousPrompt builds the user message for a spontaneous thought.
func (pb *PromptBuilder) SpontaneousPrompt(
	ps *psychology.State,
	candidate *ThoughtCandidate,
	recentMemories []memory.EpisodicMemory,
	distortionContext string,
	recentThoughts []Thought,
) string {
	tr := i18n.T()
	var b strings.Builder

	b.WriteString(pb.stateBlock(ps))
	b.WriteString("\n")

	b.WriteString(pb.thoughtStreamBlock(recentThoughts))

	if candidate != nil {
		b.WriteString(fmt.Sprintf(tr.Consciousness.Prompts.MindTurns, candidate.Prompt))
	}

	if distortionContext != "" {
		b.WriteString(distortionContext)
		b.WriteString("\n")
	}

	if len(recentMemories) > 0 {
		b.WriteString(tr.Consciousness.State.RecentExperiences)
		for _, m := range pb.trimMemories(recentMemories) {
			b.WriteString(fmt.Sprintf("- %s\n", m.Content))
		}
		b.WriteString("\n")
	}

	b.WriteString(tr.Consciousness.Prompts.SpontaneousQuestion)

	return b.String()
}

// DistortionContext generates a prompt fragment describing active cognitive
// distortions, phrased as tendencies rather than labels.
func DistortionContext(distortions []psychology.Distortion) string {
	if len(distortions) == 0 {
		return ""
	}

	tr := i18n.T()
	descriptions := make([]string, 0, len(distortions))
	for _, d := range distortions {
		desc := distortionDescription(d)
		if desc != "" {
			descriptions = append(descriptions, desc)
		}
	}

	if len(descriptions) == 0 {
		return ""
	}

	return tr.Consciousness.Prompts.DistortionPrefix + strings.Join(descriptions, "; ") + ".\n"
}

func distortionDescription(d psychology.Distortion) string {
	tr := i18n.T()
	switch d {
	case psychology.Catastrophizing:
		return tr.Consciousness.Distortions["catastrophizing"]
	case psychology.EmotionalReasoning:
		return tr.Consciousness.Distortions["emotional_reasoning"]
	case psychology.Overgeneralization:
		return tr.Consciousness.Distortions["overgeneralization"]
	case psychology.MindReading:
		return tr.Consciousness.Distortions["mind_reading"]
	case psychology.Personalization:
		return tr.Consciousness.Distortions["personalization"]
	case psychology.AllOrNothing:
		return tr.Consciousness.Distortions["all_or_nothing"]
	default:
		return ""
	}
}

// ExternalInputPrompt builds the user message for a thought triggered by
// external speech or action directed at the person.
func (pb *PromptBuilder) ExternalInputPrompt(
	ps *psychology.State,
	input ExternalInput,
	recentMemories []memory.EpisodicMemory,
	distortionContext string,
	recentThoughts []Thought,
) string {
	tr := i18n.T()
	var b strings.Builder

	b.WriteString(pb.stateBlock(ps))
	b.WriteString("\n")

	b.WriteString(pb.thoughtStreamBlock(recentThoughts))

	switch input.Type {
	case InputSpeech:
		b.WriteString(fmt.Sprintf(tr.Consciousness.Prompts.SpeechFraming, input.Content))
	case InputAction:
		b.WriteString(fmt.Sprintf(tr.Consciousness.Prompts.ActionFraming, input.Content))
	}

	if distortionContext != "" {
		b.WriteString(distortionContext)
		b.WriteString("\n")
	}

	if len(recentMemories) > 0 {
		b.WriteString(tr.Consciousness.State.RecentExperiences)
		for _, m := range pb.trimMemories(recentMemories) {
			b.WriteString(fmt.Sprintf("- %s\n", m.Content))
		}
		b.WriteString("\n")
	}

	b.WriteString(tr.Consciousness.Prompts.ExternalQuestion)

	return b.String()
}

// stateBlock renders the psychological state as a felt-experience description.
// The LLM receives affect dimensions, NOT labeled emotions.
func (pb *PromptBuilder) stateBlock(ps *psychology.State) string {
	tr := i18n.T()
	var b strings.Builder

	b.WriteString(tr.Consciousness.State.CurrentExperience)

	// Arousal.
	switch {
	case ps.Arousal > 0.7:
		b.WriteString(tr.Consciousness.State.ArousalHigh)
	case ps.Arousal > 0.4:
		b.WriteString(tr.Consciousness.State.ArousalMedium)
	case ps.Arousal > 0.2:
		b.WriteString(tr.Consciousness.State.ArousalLow)
	default:
		b.WriteString(tr.Consciousness.State.ArousalVeryLow)
	}

	// Valence.
	switch {
	case ps.Valence > 0.4:
		b.WriteString(tr.Consciousness.State.ValenceVeryPositive)
	case ps.Valence > 0.1:
		b.WriteString(tr.Consciousness.State.ValenceNeutral)
	case ps.Valence > -0.2:
		b.WriteString(tr.Consciousness.State.ValenceSlightNeg)
	case ps.Valence > -0.5:
		b.WriteString(tr.Consciousness.State.ValenceNegative)
	default:
		b.WriteString(tr.Consciousness.State.ValenceVeryNegative)
	}

	// Energy.
	switch {
	case ps.Energy > 0.7:
		b.WriteString(tr.Consciousness.State.EnergyHigh)
	case ps.Energy > 0.4:
		b.WriteString(tr.Consciousness.State.EnergyMedium)
	case ps.Energy > 0.2:
		b.WriteString(tr.Consciousness.State.EnergyLow)
	default:
		b.WriteString(tr.Consciousness.State.EnergyVeryLow)
	}

	// Cognitive load.
	if ps.CognitiveLoad > 0.6 {
		b.WriteString(tr.Consciousness.State.CognitiveLoadHigh)
	} else if ps.CognitiveLoad > 0.3 {
		b.WriteString(tr.Consciousness.State.CognitiveLoadMedium)
	}

	// Regulation.
	if ps.RegulationCapacity < 0.2 {
		b.WriteString(tr.Consciousness.State.RegulationVeryLow)
	} else if ps.RegulationCapacity < 0.4 {
		b.WriteString(tr.Consciousness.State.RegulationLow)
	}

	// Isolation.
	switch ps.Isolation.Phase {
	case psychology.IsolationBoredom:
		b.WriteString(tr.Consciousness.State.IsolationBoredom)
	case psychology.IsolationLoneliness:
		b.WriteString(tr.Consciousness.State.IsolationLoneliness)
	case psychology.IsolationSignificant:
		b.WriteString(tr.Consciousness.State.IsolationSignificant)
	case psychology.IsolationDestabilizing:
		b.WriteString(tr.Consciousness.State.IsolationDestabilizing)
	case psychology.IsolationSevere:
		b.WriteString(tr.Consciousness.State.IsolationSevere)
	}

	return b.String()
}

// thoughtStreamBlock renders recent thoughts so the LLM has continuity
// with what it was just thinking. Returns empty string if no recent thoughts.
func (pb *PromptBuilder) thoughtStreamBlock(thoughts []Thought) string {
	if len(thoughts) == 0 {
		return ""
	}

	tr := i18n.T()
	var b strings.Builder
	b.WriteString(tr.Consciousness.State.RecentThoughts)
	for _, t := range thoughts {
		b.WriteString(fmt.Sprintf("- %s\n", t.Content))
	}
	b.WriteString("\n")
	return b.String()
}

// trimMemories limits the number of memories included in the prompt
// based on the token budget. Rough estimate: 1 memory ≈ 30 tokens.
func (pb *PromptBuilder) trimMemories(memories []memory.EpisodicMemory) []memory.EpisodicMemory {
	maxMemories := pb.maxTokenEstimate / 100 // rough allocation
	if maxMemories < 3 {
		maxMemories = 3
	}
	if maxMemories > 10 {
		maxMemories = 10
	}
	if len(memories) <= maxMemories {
		return memories
	}
	return memories[:maxMemories]
}
