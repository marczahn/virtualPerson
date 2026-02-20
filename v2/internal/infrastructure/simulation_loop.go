package infrastructure

import (
	"fmt"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/consciousness"
	"github.com/marczahn/person/v2/internal/motivation"
)

// InputDrainer drains pending external inputs once per tick.
type InputDrainer interface {
	Drain() TickInput
}

// BioEngine advances bio state by one tick.
type BioEngine interface {
	Tick(s *biology.State, dt float64) biology.TickResult
}

// MotivationComputer computes deterministic drive state from bio/personality/chronic inputs.
type MotivationComputer interface {
	Compute(bio biology.State, personality motivation.Personality, chronic motivation.ChronicState) motivation.MotivationState
}

// MindResponder returns one structured consciousness response for this tick.
type MindResponder interface {
	Respond(in MindRequest) string
}

// TickInput contains drained input effects and current action gating for one tick.
type TickInput struct {
	PreBioRates    []biology.BioRate
	PreBioPulses   []biology.BioPulse
	AllowedActions map[string]bool
	NowSeconds     int64
	ExternalText   string
}

// MindRequest is the consciousness-stage payload for one tick.
type MindRequest struct {
	Bio         biology.State
	Motivation  motivation.MotivationState
	Prompt      consciousness.PromptContext
	Input       TickInput
	PriorParsed consciousness.ParsedResponse
}

// SimulationState is the mutable simulation state carried across ticks.
type SimulationState struct {
	Bio           biology.State
	Personality   motivation.Personality
	Chronic       motivation.ChronicState
	PriorParsed   consciousness.ParsedResponse
	CooldownState consciousness.ActionCooldownState
	Continuity    *consciousness.ContinuityBuffer
}

// TickResult captures one fully-orchestrated INF-07 tick.
type TickResult struct {
	Input               TickInput
	Bio                 biology.TickResult
	Motivation          motivation.MotivationState
	PerceivedMotivation motivation.MotivationState
	Prompt              consciousness.PromptContext
	Raw                 string
	Parsed              consciousness.ParsedResponse
	ActionOutcome       consciousness.ActionOutcome
}

// SimulationLoopDeps wires infrastructure orchestration to layer contracts.
type SimulationLoopDeps struct {
	Input      InputDrainer
	Biology    BioEngine
	Motivation MotivationComputer
	Mind       MindResponder
	Cooldowns  consciousness.ActionCooldowns
}

// SimulationLoop orchestrates one sequential tick: input -> biology -> motivation -> consciousness -> feedback.
type SimulationLoop struct {
	input      InputDrainer
	biology    BioEngine
	motivation MotivationComputer
	mind       MindResponder
	cooldowns  consciousness.ActionCooldowns
}

func NewSimulationLoop(deps SimulationLoopDeps) *SimulationLoop {
	if deps.Input == nil {
		panic(fmt.Errorf("simulation loop requires InputDrainer"))
	}
	if deps.Biology == nil {
		panic(fmt.Errorf("simulation loop requires BioEngine"))
	}
	if deps.Motivation == nil {
		panic(fmt.Errorf("simulation loop requires MotivationComputer"))
	}
	if deps.Mind == nil {
		panic(fmt.Errorf("simulation loop requires MindResponder"))
	}

	return &SimulationLoop{
		input:      deps.Input,
		biology:    deps.Biology,
		motivation: deps.Motivation,
		mind:       deps.Mind,
		cooldowns:  deps.Cooldowns,
	}
}

func (l *SimulationLoop) Tick(state *SimulationState, dt float64) TickResult {
	if state == nil {
		panic(fmt.Errorf("simulation loop requires non-nil state"))
	}

	input := l.input.Drain()

	if len(input.PreBioRates) > 0 || len(input.PreBioPulses) > 0 {
		biology.ApplyFeedbackAtTickEnd(&state.Bio, dt, biology.FeedbackEnvelope{
			Rates:  input.PreBioRates,
			Pulses: input.PreBioPulses,
		})
	}

	bioResult := l.biology.Tick(&state.Bio, dt)
	motivationState := l.motivation.Compute(state.Bio, state.Personality, state.Chronic)

	var prompt consciousness.PromptContext
	if state.Continuity != nil {
		prompt = consciousness.BuildPromptContextWithContinuity(motivationState, state.Continuity.Items())
	} else {
		prompt = consciousness.BuildPromptContext(motivationState)
	}

	raw := l.mind.Respond(MindRequest{
		Bio:         state.Bio,
		Motivation:  motivationState,
		Prompt:      prompt,
		Input:       input,
		PriorParsed: state.PriorParsed,
	})

	parsed := consciousness.ParseResponse(raw, state.PriorParsed)
	perceived := consciousness.ApplyParsedDriveOverridesForNextTick(motivationState, parsed)

	allowed := false
	if input.AllowedActions != nil {
		allowed = input.AllowedActions[parsed.Action]
	}
	actionOutcome, nextCooldownState := consciousness.ResolveActionOutcomeWithCooldown(
		parsed.Action,
		allowed,
		input.NowSeconds,
		l.cooldowns,
		state.CooldownState,
	)

	var feedback biology.TickFeedbackBuffer
	feedback.AddPulses(consciousness.EmotionalPulseFromState(parsed.State))
	feedback.AddPulses(consciousness.ActionPulse(actionOutcome))
	feedback.ApplyAtTickEnd(&state.Bio, dt)

	state.PriorParsed = parsed
	state.CooldownState = nextCooldownState
	if state.Continuity != nil && parsed.Narrative != "" {
		state.Continuity.Add(consciousness.Thought{Text: parsed.Narrative})
	}

	return TickResult{
		Input:               input,
		Bio:                 bioResult,
		Motivation:          motivationState,
		PerceivedMotivation: perceived,
		Prompt:              prompt,
		Raw:                 raw,
		Parsed:              parsed,
		ActionOutcome:       actionOutcome,
	}
}
