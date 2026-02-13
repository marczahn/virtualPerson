package consciousness

import (
	"math/rand"

	"github.com/marczahn/person/internal/psychology"
)

// ThoughtQueue manages the priority queue for spontaneous thought generation.
// Higher-priority categories are more likely to fire, but lower-priority
// ones can still surface through weighted random selection.
type ThoughtQueue struct {
	// UnresolvedErrors tracks prediction errors that haven't been explained.
	UnresolvedErrors []string

	// ActiveNeeds tracks current biological needs above threshold.
	ActiveNeeds []string

	// ActiveGoals tracks ongoing plans and intentions.
	ActiveGoals []string

	// SocialConcerns tracks unresolved social modeling questions.
	SocialConcerns []string

	// Absorbed indicates the person is deep in a thought chain.
	// When absorbed, new triggers need higher salience to interrupt.
	Absorbed bool
	AbsorptionDepth float64 // 0-1, how deep in the current chain
}

// NewThoughtQueue creates an empty queue.
func NewThoughtQueue() *ThoughtQueue {
	return &ThoughtQueue{}
}

// ThoughtCandidate represents a possible spontaneous thought.
type ThoughtCandidate struct {
	Priority Priority
	Category string // human-readable category
	Prompt   string // the specific thought prompt to feed the LLM
}

// SelectSpontaneous picks a spontaneous thought from the queue using
// weighted random selection. Higher-priority items have more weight.
// Returns nil if the queue is empty.
func (q *ThoughtQueue) SelectSpontaneous(ps *psychology.State) *ThoughtCandidate {
	candidates := q.buildCandidates(ps)
	if len(candidates) == 0 {
		return nil
	}

	// Weight by priority: lower Priority value = higher weight.
	weights := make([]float64, len(candidates))
	var totalWeight float64
	for i, c := range candidates {
		w := priorityWeight(c.Priority)
		weights[i] = w
		totalWeight += w
	}

	// Weighted random selection.
	r := rand.Float64() * totalWeight
	var cumulative float64
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			return &candidates[i]
		}
	}
	return &candidates[len(candidates)-1]
}

// SelectSpontaneousDeterministic is like SelectSpontaneous but uses
// a provided random value in [0,1) for testability.
func (q *ThoughtQueue) SelectSpontaneousDeterministic(ps *psychology.State, randVal float64) *ThoughtCandidate {
	candidates := q.buildCandidates(ps)
	if len(candidates) == 0 {
		return nil
	}

	weights := make([]float64, len(candidates))
	var totalWeight float64
	for i, c := range candidates {
		w := priorityWeight(c.Priority)
		weights[i] = w
		totalWeight += w
	}

	r := randVal * totalWeight
	var cumulative float64
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			return &candidates[i]
		}
	}
	return &candidates[len(candidates)-1]
}

// buildCandidates gathers all possible spontaneous thoughts.
func (q *ThoughtQueue) buildCandidates(ps *psychology.State) []ThoughtCandidate {
	var candidates []ThoughtCandidate

	// 1. Unresolved prediction errors (highest priority).
	for _, e := range q.UnresolvedErrors {
		candidates = append(candidates, ThoughtCandidate{
			Priority: PriorityPredictionError,
			Category: "prediction_error",
			Prompt:   e,
		})
	}

	// 2. Active biological needs.
	for _, n := range q.ActiveNeeds {
		candidates = append(candidates, ThoughtCandidate{
			Priority: PriorityBiologicalNeed,
			Category: "biological_need",
			Prompt:   n,
		})
	}

	// 3. Goal rehearsal.
	for _, g := range q.ActiveGoals {
		candidates = append(candidates, ThoughtCandidate{
			Priority: PriorityGoalRehearsal,
			Category: "goal_rehearsal",
			Prompt:   g,
		})
	}

	// 4. Social modeling.
	for _, s := range q.SocialConcerns {
		candidates = append(candidates, ThoughtCandidate{
			Priority: PrioritySocialModeling,
			Category: "social_modeling",
			Prompt:   s,
		})
	}

	// 5. Associative drift (always available as fallback).
	// Generate based on current psychological state.
	driftPrompt := associativeDriftPrompt(ps)
	candidates = append(candidates, ThoughtCandidate{
		Priority: PriorityAssociativeDrift,
		Category: "associative_drift",
		Prompt:   driftPrompt,
	})

	return candidates
}

// priorityWeight returns the selection weight for a priority level.
// Higher priority = higher weight.
func priorityWeight(p Priority) float64 {
	switch p {
	case PriorityPredictionError:
		return 8.0
	case PriorityBiologicalNeed:
		return 5.0
	case PriorityGoalRehearsal:
		return 3.0
	case PrioritySocialModeling:
		return 2.0
	case PriorityAssociativeDrift:
		return 1.0
	default:
		return 1.0
	}
}

// associativeDriftPrompt generates a mind-wandering prompt based on
// the current psychological state.
func associativeDriftPrompt(ps *psychology.State) string {
	if ps.Energy < 0.3 {
		return "Your mind drifts, thoughts moving slowly. You feel drained."
	}
	if ps.Valence < -0.3 {
		return "Your thoughts wander into darker territory. Something nags at you."
	}
	if ps.Arousal < 0.2 && ps.Energy > 0.5 {
		return "Your mind is quiet and open. Thoughts arise and pass without urgency."
	}
	return "Your mind wanders freely, one thought leading to another."
}

// AddPredictionError records an unexpected event for later processing.
func (q *ThoughtQueue) AddPredictionError(description string) {
	q.UnresolvedErrors = append(q.UnresolvedErrors, description)
}

// ResolvePredictionError removes a resolved prediction error.
func (q *ThoughtQueue) ResolvePredictionError(description string) {
	for i, e := range q.UnresolvedErrors {
		if e == description {
			q.UnresolvedErrors = append(q.UnresolvedErrors[:i], q.UnresolvedErrors[i+1:]...)
			return
		}
	}
}

// UpdateNeeds refreshes the biological needs list based on current psych state.
func (q *ThoughtQueue) UpdateNeeds(ps *psychology.State) {
	q.ActiveNeeds = nil

	if ps.Energy < 0.2 {
		q.ActiveNeeds = append(q.ActiveNeeds, "You are exhausted. Your body demands rest.")
	}
	if ps.Valence < -0.5 && ps.Arousal > 0.6 {
		q.ActiveNeeds = append(q.ActiveNeeds, "Something is wrong. You feel distressed and agitated.")
	}

	// Isolation as a social need.
	if ps.Isolation.Phase >= psychology.IsolationLoneliness {
		q.ActiveNeeds = append(q.ActiveNeeds, "You feel deeply alone. You crave human contact.")
	}
}

// EnterAbsorption marks the person as absorbed in a thought chain.
func (q *ThoughtQueue) EnterAbsorption(depth float64) {
	q.Absorbed = true
	q.AbsorptionDepth = depth
}

// ExitAbsorption exits the absorbed state.
func (q *ThoughtQueue) ExitAbsorption() {
	q.Absorbed = false
	q.AbsorptionDepth = 0
}
