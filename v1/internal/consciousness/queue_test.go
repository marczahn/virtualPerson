package consciousness

import (
	"testing"

	"github.com/marczahn/person/internal/psychology"
)

func TestThoughtQueue_EmptyQueueStillHasDrift(t *testing.T) {
	q := NewThoughtQueue()
	ps := &psychology.State{Energy: 0.5}

	candidate := q.SelectSpontaneousDeterministic(ps, 0.5)

	if candidate == nil {
		t.Fatal("expected at least associative drift candidate")
	}
	if candidate.Priority != PriorityAssociativeDrift {
		t.Errorf("expected drift priority, got %s", candidate.Priority)
	}
}

func TestThoughtQueue_PredictionErrorHighestWeight(t *testing.T) {
	q := NewThoughtQueue()
	q.AddPredictionError("unexpected noise")
	q.ActiveGoals = []string{"plan dinner"}

	ps := &psychology.State{Energy: 0.5}

	// With low random value, should select highest-weight item.
	candidate := q.SelectSpontaneousDeterministic(ps, 0.0)

	if candidate == nil {
		t.Fatal("expected a candidate")
	}
	if candidate.Priority != PriorityPredictionError {
		t.Errorf("expected prediction_error at low rand, got %s", candidate.Priority)
	}
}

func TestThoughtQueue_WeightedSelection_DriftAtHighRand(t *testing.T) {
	q := NewThoughtQueue()
	q.AddPredictionError("unexpected noise")

	ps := &psychology.State{Energy: 0.5}

	// With very high random value, should reach drift.
	candidate := q.SelectSpontaneousDeterministic(ps, 0.99)

	if candidate == nil {
		t.Fatal("expected a candidate")
	}
	if candidate.Priority != PriorityAssociativeDrift {
		t.Errorf("expected drift at high rand, got %s", candidate.Priority)
	}
}

func TestThoughtQueue_ResolvePredictionError(t *testing.T) {
	q := NewThoughtQueue()
	q.AddPredictionError("noise A")
	q.AddPredictionError("noise B")

	q.ResolvePredictionError("noise A")

	if len(q.UnresolvedErrors) != 1 || q.UnresolvedErrors[0] != "noise B" {
		t.Errorf("after resolving A, expected [noise B], got %v", q.UnresolvedErrors)
	}
}

func TestThoughtQueue_UpdateNeeds_Exhaustion(t *testing.T) {
	q := NewThoughtQueue()
	ps := &psychology.State{Energy: 0.1}

	q.UpdateNeeds(ps)

	if len(q.ActiveNeeds) == 0 {
		t.Error("expected exhaustion need at energy=0.1")
	}
}

func TestThoughtQueue_UpdateNeeds_Isolation(t *testing.T) {
	q := NewThoughtQueue()
	ps := &psychology.State{
		Energy: 0.5,
		Isolation: psychology.IsolationState{
			Phase: psychology.IsolationLoneliness,
		},
	}

	q.UpdateNeeds(ps)

	found := false
	for _, n := range q.ActiveNeeds {
		if n == "You feel deeply alone. You crave human contact." {
			found = true
		}
	}
	if !found {
		t.Errorf("expected isolation need, got %v", q.ActiveNeeds)
	}
}

func TestThoughtQueue_Absorption(t *testing.T) {
	q := NewThoughtQueue()

	q.EnterAbsorption(0.7)
	if !q.Absorbed || q.AbsorptionDepth != 0.7 {
		t.Errorf("expected absorbed at 0.7, got %v/%f", q.Absorbed, q.AbsorptionDepth)
	}

	q.ExitAbsorption()
	if q.Absorbed || q.AbsorptionDepth != 0 {
		t.Error("expected not absorbed after exit")
	}
}

func TestPriorityWeight_Ordering(t *testing.T) {
	weights := []struct {
		p Priority
		w float64
	}{
		{PriorityPredictionError, priorityWeight(PriorityPredictionError)},
		{PriorityBiologicalNeed, priorityWeight(PriorityBiologicalNeed)},
		{PriorityGoalRehearsal, priorityWeight(PriorityGoalRehearsal)},
		{PrioritySocialModeling, priorityWeight(PrioritySocialModeling)},
		{PriorityAssociativeDrift, priorityWeight(PriorityAssociativeDrift)},
	}

	for i := 0; i < len(weights)-1; i++ {
		if weights[i].w <= weights[i+1].w {
			t.Errorf("priority %s weight (%f) should be > %s weight (%f)",
				weights[i].p, weights[i].w, weights[i+1].p, weights[i+1].w)
		}
	}
}

func TestAssociativeDriftPrompt_LowEnergy(t *testing.T) {
	ps := &psychology.State{Energy: 0.1}
	got := associativeDriftPrompt(ps)

	if got == "" {
		t.Error("expected non-empty drift prompt for low energy")
	}
}

func TestAssociativeDriftPrompt_NegativeValence(t *testing.T) {
	ps := &psychology.State{Energy: 0.5, Valence: -0.5}
	got := associativeDriftPrompt(ps)

	if got == "" {
		t.Error("expected non-empty drift prompt for negative valence")
	}
}
