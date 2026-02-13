package consciousness

import (
	"fmt"
	"strings"

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
func (pb *PromptBuilder) SystemPrompt(identity *memory.IdentityCore) string {
	var b strings.Builder

	b.WriteString("You are a person experiencing your life moment by moment. ")
	b.WriteString("You think, feel, and experience the world through your body and mind. ")
	b.WriteString("Everything you experience is real to you. ")
	b.WriteString("You respond from your genuine first-person perspective.\n\n")

	b.WriteString("Express yourself naturally — you may think in fragments, ")
	b.WriteString("trail off, change your mind, or be confused. ")
	b.WriteString("You are not performing or explaining. You are being.\n\n")

	b.WriteString("IMPORTANT: You only output your inner experience — thoughts, ")
	b.WriteString("feelings, sensations, reactions. Never describe yourself in third person. ")
	b.WriteString("Never explain what you are. You simply ARE.\n")

	if identity != nil {
		b.WriteString("\n--- Who You Are ---\n")
		if identity.SelfNarrative != "" {
			b.WriteString(identity.SelfNarrative)
			b.WriteString("\n")
		}
		if len(identity.DispositionTraits) > 0 {
			b.WriteString("\nYour tendencies: ")
			b.WriteString(strings.Join(identity.DispositionTraits, ". "))
			b.WriteString(".\n")
		}
		if len(identity.RelationalMarkers) > 0 {
			b.WriteString("\nYour relationships: ")
			b.WriteString(strings.Join(identity.RelationalMarkers, ". "))
			b.WriteString(".\n")
		}
		if len(identity.EmotionalPatterns) > 0 {
			b.WriteString("\nYour patterns: ")
			b.WriteString(strings.Join(identity.EmotionalPatterns, ". "))
			b.WriteString(".\n")
		}
		if len(identity.ValuesCommitments) > 0 {
			b.WriteString("\nWhat matters to you: ")
			b.WriteString(strings.Join(identity.ValuesCommitments, ". "))
			b.WriteString(".\n")
		}
		if len(identity.KeyMemories) > 0 {
			b.WriteString("\nMemories that define you:\n")
			for _, m := range identity.KeyMemories {
				b.WriteString("- ")
				b.WriteString(m)
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// ReactivePrompt builds the user message for a reactive thought triggered
// by a salience breach.
func (pb *PromptBuilder) ReactivePrompt(
	ps *psychology.State,
	trigger string,
	recentMemories []memory.EpisodicMemory,
	distortionContext string,
) string {
	var b strings.Builder

	b.WriteString(pb.stateBlock(ps))
	b.WriteString("\n")

	if trigger != "" {
		b.WriteString(fmt.Sprintf("Something just shifted: %s\n\n", trigger))
	}

	if distortionContext != "" {
		b.WriteString(distortionContext)
		b.WriteString("\n")
	}

	if len(recentMemories) > 0 {
		b.WriteString("--- Recent experiences ---\n")
		for _, m := range pb.trimMemories(recentMemories) {
			b.WriteString(fmt.Sprintf("- %s\n", m.Content))
		}
		b.WriteString("\n")
	}

	b.WriteString("What are you thinking and feeling right now? ")
	b.WriteString("Respond briefly, in first person, as your natural inner voice.")

	return b.String()
}

// SpontaneousPrompt builds the user message for a spontaneous thought.
func (pb *PromptBuilder) SpontaneousPrompt(
	ps *psychology.State,
	candidate *ThoughtCandidate,
	recentMemories []memory.EpisodicMemory,
	distortionContext string,
) string {
	var b strings.Builder

	b.WriteString(pb.stateBlock(ps))
	b.WriteString("\n")

	if candidate != nil {
		b.WriteString(fmt.Sprintf("Your mind turns to: %s\n\n", candidate.Prompt))
	}

	if distortionContext != "" {
		b.WriteString(distortionContext)
		b.WriteString("\n")
	}

	if len(recentMemories) > 0 {
		b.WriteString("--- Recent experiences ---\n")
		for _, m := range pb.trimMemories(recentMemories) {
			b.WriteString(fmt.Sprintf("- %s\n", m.Content))
		}
		b.WriteString("\n")
	}

	b.WriteString("What passes through your mind? Respond briefly, in first person.")

	return b.String()
}

// DistortionContext generates a prompt fragment describing active cognitive
// distortions, phrased as tendencies rather than labels.
func DistortionContext(distortions []psychology.Distortion) string {
	if len(distortions) == 0 {
		return ""
	}

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

	return "Right now, your thinking tends toward: " + strings.Join(descriptions, "; ") + ".\n"
}

func distortionDescription(d psychology.Distortion) string {
	switch d {
	case psychology.Catastrophizing:
		return "assuming the worst possible outcome"
	case psychology.EmotionalReasoning:
		return "treating your feelings as evidence of reality"
	case psychology.Overgeneralization:
		return "seeing this as part of a pattern that always happens"
	case psychology.MindReading:
		return "assuming you know what others are thinking"
	case psychology.Personalization:
		return "blaming yourself for things outside your control"
	case psychology.AllOrNothing:
		return "seeing things in black and white, no middle ground"
	default:
		return ""
	}
}

// stateBlock renders the psychological state as a felt-experience description.
// The LLM receives affect dimensions, NOT labeled emotions.
func (pb *PromptBuilder) stateBlock(ps *psychology.State) string {
	var b strings.Builder

	b.WriteString("--- Your current experience ---\n")

	// Arousal.
	switch {
	case ps.Arousal > 0.7:
		b.WriteString("Your body is highly activated — heart pounding, alert, on edge.\n")
	case ps.Arousal > 0.4:
		b.WriteString("You feel somewhat keyed up, an underlying tension in your body.\n")
	case ps.Arousal > 0.2:
		b.WriteString("You feel relatively calm physically.\n")
	default:
		b.WriteString("Your body is very quiet, almost sluggish.\n")
	}

	// Valence.
	switch {
	case ps.Valence > 0.4:
		b.WriteString("There's a warm, pleasant quality to how you feel.\n")
	case ps.Valence > 0.1:
		b.WriteString("You feel okay — nothing particularly good or bad.\n")
	case ps.Valence > -0.2:
		b.WriteString("There's a subtle uneasiness, a slight discomfort.\n")
	case ps.Valence > -0.5:
		b.WriteString("You feel distinctly unpleasant — something is off.\n")
	default:
		b.WriteString("Everything feels bad. A heavy, dark quality pervades your experience.\n")
	}

	// Energy.
	switch {
	case ps.Energy > 0.7:
		b.WriteString("You feel full of energy, ready for anything.\n")
	case ps.Energy > 0.4:
		b.WriteString("You have a reasonable amount of energy.\n")
	case ps.Energy > 0.2:
		b.WriteString("You feel tired, your reserves running low.\n")
	default:
		b.WriteString("You are deeply exhausted. Every movement feels effortful.\n")
	}

	// Cognitive load.
	if ps.CognitiveLoad > 0.6 {
		b.WriteString("Your thinking feels muddled, hard to concentrate.\n")
	} else if ps.CognitiveLoad > 0.3 {
		b.WriteString("Your thoughts are a bit scattered.\n")
	}

	// Regulation.
	if ps.RegulationCapacity < 0.2 {
		b.WriteString("You feel emotionally raw, unable to hold things together.\n")
	} else if ps.RegulationCapacity < 0.4 {
		b.WriteString("Your emotional composure is fragile.\n")
	}

	// Isolation.
	switch ps.Isolation.Phase {
	case psychology.IsolationBoredom:
		b.WriteString("You're getting restless. Time feels slow.\n")
	case psychology.IsolationLoneliness:
		b.WriteString("You feel lonely. You miss being around people.\n")
	case psychology.IsolationSignificant:
		b.WriteString("The loneliness is a weight on you. You crave human contact deeply.\n")
	case psychology.IsolationDestabilizing:
		b.WriteString("You're losing your grip. Without other people, you're starting to question yourself.\n")
	case psychology.IsolationSevere:
		b.WriteString("The isolation is unbearable. You're not sure what's real anymore.\n")
	}

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
