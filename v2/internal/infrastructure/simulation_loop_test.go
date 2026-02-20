package infrastructure_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/consciousness"
	"github.com/marczahn/person/v2/internal/infrastructure"
	"github.com/marczahn/person/v2/internal/motivation"
)

type fakeInputDrainer struct {
	order []string
	input infrastructure.TickInput
	calls int
}

func (f *fakeInputDrainer) Drain() infrastructure.TickInput {
	f.calls++
	f.order = append(f.order, "drain")
	return f.input
}

type fakeBioEngine struct {
	order []string
	calls int
}

func (f *fakeBioEngine) Tick(s *biology.State, dt float64) biology.TickResult {
	f.calls++
	f.order = append(f.order, "biology")
	return biology.TickResult{}
}

type fakeMotivationComputer struct {
	order   []string
	calls   int
	seenBio biology.State
	result  motivation.MotivationState
}

func (f *fakeMotivationComputer) Compute(
	bio biology.State,
	personality motivation.Personality,
	chronic motivation.ChronicState,
) motivation.MotivationState {
	f.calls++
	f.order = append(f.order, "motivation")
	f.seenBio = bio
	if f.result.ActiveGoalDrive == "" {
		f.result.ActiveGoalDrive = motivation.DriveEnergy
		f.result.ActiveGoalUrgency = 0.7
	}
	return f.result
}

type fakeMind struct {
	order      []string
	calls      int
	raw        string
	capturedIn infrastructure.MindRequest
}

func (f *fakeMind) Respond(in infrastructure.MindRequest) string {
	f.calls++
	f.order = append(f.order, "consciousness")
	f.capturedIn = in
	return f.raw
}

func TestSimulationLoop_TickRunsStagesSequentiallyOnce(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Hunger = 0.70

	drainer := &fakeInputDrainer{input: infrastructure.TickInput{
		AllowedActions: map[string]bool{"eat": true},
		NowSeconds:     42,
	}}
	bioEngine := &fakeBioEngine{}
	motivationComputer := &fakeMotivationComputer{result: motivation.MotivationState{
		EnergyUrgency:      0.8,
		SocialUrgency:      0.1,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.2,
		IdentityUrgency:    0.2,
		ActiveGoalDrive:    motivation.DriveEnergy,
		ActiveGoalUrgency:  0.8,
	}}
	mind := &fakeMind{raw: "Eat now. [STATE: arousal=0.0, valence=0.0] [ACTION: eat]"}

	loop := infrastructure.NewSimulationLoop(infrastructure.SimulationLoopDeps{
		Input:      drainer,
		Biology:    bioEngine,
		Motivation: motivationComputer,
		Mind:       mind,
	})

	state := infrastructure.SimulationState{Bio: *bio}
	result := loop.Tick(&state, 1.0)

	if drainer.calls != 1 || bioEngine.calls != 1 || motivationComputer.calls != 1 || mind.calls != 1 {
		t.Fatalf("expected exactly one call per stage, got drain=%d biology=%d motivation=%d consciousness=%d",
			drainer.calls, bioEngine.calls, motivationComputer.calls, mind.calls)
	}

	order := append(append(append(drainer.order, bioEngine.order...), motivationComputer.order...), mind.order...)
	expected := []string{"drain", "biology", "motivation", "consciousness"}
	if !reflect.DeepEqual(order, expected) {
		t.Fatalf("unexpected stage order: got=%v want=%v", order, expected)
	}

	if state.Bio.Hunger >= 0.70 {
		t.Fatalf("expected action feedback to apply by tick end, hunger did not drop: %f", state.Bio.Hunger)
	}
	if result.Parsed.Action != "eat" {
		t.Fatalf("expected parsed action to be retained in result, got %q", result.Parsed.Action)
	}
}

func TestSimulationLoop_DrainedInputPreBioPulsesAffectMotivationInput(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Stress = 0.10

	drainer := &fakeInputDrainer{input: infrastructure.TickInput{
		PreBioPulses: []biology.BioPulse{{Field: "stress", Amount: 0.25}},
	}}
	bioEngine := &fakeBioEngine{}
	motivationComputer := &fakeMotivationComputer{}
	mind := &fakeMind{raw: "steady [STATE: arousal=0.0, valence=0.0] [ACTION: breathe]"}

	loop := infrastructure.NewSimulationLoop(infrastructure.SimulationLoopDeps{
		Input:      drainer,
		Biology:    bioEngine,
		Motivation: motivationComputer,
		Mind:       mind,
	})

	state := infrastructure.SimulationState{Bio: *bio}
	loop.Tick(&state, 1.0)

	if math.Abs(motivationComputer.seenBio.Stress-0.35) > 0.000001 {
		t.Fatalf("expected pre-bio pulse to affect motivation input stress=0.35, got %f", motivationComputer.seenBio.Stress)
	}
}

func TestSimulationLoop_FeedbackAppliesAtTickEnd(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Hunger = 0.70

	drainer := &fakeInputDrainer{input: infrastructure.TickInput{
		AllowedActions: map[string]bool{"eat": true},
		NowSeconds:     100,
	}}
	bioEngine := &fakeBioEngine{}
	motivationComputer := &fakeMotivationComputer{result: motivation.MotivationState{
		EnergyUrgency:      0.7,
		SocialUrgency:      0.1,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.2,
		IdentityUrgency:    0.2,
		ActiveGoalDrive:    motivation.DriveEnergy,
		ActiveGoalUrgency:  0.7,
	}}
	mind := &fakeMind{raw: "need food [STATE: arousal=0.0, valence=0.0] [ACTION: eat]"}

	loop := infrastructure.NewSimulationLoop(infrastructure.SimulationLoopDeps{
		Input:      drainer,
		Biology:    bioEngine,
		Motivation: motivationComputer,
		Mind:       mind,
	})

	state := infrastructure.SimulationState{
		Bio: *bio,
		PriorParsed: consciousness.ParsedResponse{
			State: consciousness.ParsedState{Arousal: 0, Valence: 0},
		},
	}

	loop.Tick(&state, 1.0)

	if math.Abs(mind.capturedIn.Bio.Hunger-0.70) > 0.000001 {
		t.Fatalf("feedback must not mutate bio before consciousness stage, mind saw hunger=%f", mind.capturedIn.Bio.Hunger)
	}
	if math.Abs(state.Bio.Hunger-0.40) > 0.000001 {
		t.Fatalf("expected end-of-tick feedback to apply eat pulse (0.70 -> 0.40), got %f", state.Bio.Hunger)
	}
}
