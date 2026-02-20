package consciousness

import (
	"strings"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/motivation"
)

// EmotionalPulseFromState converts parsed [STATE] arousal/valence into absolute bio deltas.
// Contract: these are one-shot pulses and must not be dt-scaled by callers.
func EmotionalPulseFromState(state ParsedState) []biology.BioPulse {
	arousal := clampSigned(state.Arousal)
	valence := clampSigned(state.Valence)

	return []biology.BioPulse{
		{Field: "stress", Amount: 0.12*arousal - 0.08*valence},
		{Field: "mood", Amount: 0.16*valence - 0.05*arousal},
		{Field: "physical_tension", Amount: 0.10 * max0(arousal)},
		{Field: "cognitive_capacity", Amount: -0.06*max0(arousal) + 0.04*max0(valence)},
	}
}

// ActionPulse maps successful [ACTION] outcomes to absolute bio deltas.
// Contract: blocked/failed actions emit no bio effects.
func ActionPulse(outcome ActionOutcome) []biology.BioPulse {
	if !outcome.Executed || !outcome.Satisfied {
		return nil
	}

	switch strings.ToLower(strings.TrimSpace(outcome.Action)) {
	case string(motivation.ActionEat):
		return []biology.BioPulse{
			{Field: "hunger", Amount: -0.30},
			{Field: "energy", Amount: 0.08},
			{Field: "mood", Amount: 0.04},
		}
	case string(motivation.ActionHydrate):
		return []biology.BioPulse{
			{Field: "stress", Amount: -0.03},
			{Field: "mood", Amount: 0.01},
		}
	case string(motivation.ActionRest):
		return []biology.BioPulse{
			{Field: "energy", Amount: 0.18},
			{Field: "stress", Amount: -0.06},
			{Field: "physical_tension", Amount: -0.08},
		}
	case string(motivation.ActionReachOut):
		return []biology.BioPulse{
			{Field: "social_deficit", Amount: -0.20},
			{Field: "mood", Amount: 0.06},
		}
	case string(motivation.ActionJournal):
		return []biology.BioPulse{
			{Field: "stress", Amount: -0.02},
			{Field: "mood", Amount: 0.03},
		}
	case string(motivation.ActionBreathe):
		return []biology.BioPulse{
			{Field: "stress", Amount: -0.08},
			{Field: "physical_tension", Amount: -0.10},
		}
	case string(motivation.ActionScanArea):
		return []biology.BioPulse{
			{Field: "stress", Amount: -0.04},
		}
	case string(motivation.ActionSeekWarm):
		return []biology.BioPulse{
			{Field: "body_temp", Amount: 0.60},
			{Field: "stress", Amount: -0.02},
		}
	case string(motivation.ActionSeekCool):
		return []biology.BioPulse{
			{Field: "body_temp", Amount: -0.60},
			{Field: "stress", Amount: -0.02},
		}
	case string(motivation.ActionMicroTask):
		return []biology.BioPulse{
			{Field: "cognitive_capacity", Amount: 0.04},
			{Field: "mood", Amount: 0.02},
			{Field: "energy", Amount: -0.02},
		}
	default:
		return nil
	}
}

func ResolveActionOutcome(action string, allowed bool) ActionOutcome {
	normalized := strings.ToLower(strings.TrimSpace(action))
	return ActionOutcome{
		Action:    normalized,
		Executed:  allowed,
		Satisfied: allowed,
	}
}

// ResolveActionOutcomeWithCooldown applies environment gating plus per-action cooldowns.
// Contract: blocked/cooldown-rejected actions are unsatisfied and do not update cooldown state.
func ResolveActionOutcomeWithCooldown(
	action string,
	allowedByEnvironment bool,
	nowSeconds int64,
	cooldowns ActionCooldowns,
	state ActionCooldownState,
) (ActionOutcome, ActionCooldownState) {
	normalized := strings.ToLower(strings.TrimSpace(action))
	nextState := cloneCooldownState(state)
	if !allowedByEnvironment {
		return ActionOutcome{Action: normalized, Executed: false, Satisfied: false}, nextState
	}

	if until, ok := nextState[normalized]; ok && nowSeconds < until {
		return ActionOutcome{Action: normalized, Executed: false, Satisfied: false}, nextState
	}

	duration := cooldowns[normalized]
	if duration > 0 {
		nextState[normalized] = nowSeconds + duration
	}

	return ActionOutcome{Action: normalized, Executed: true, Satisfied: true}, nextState
}

func cloneCooldownState(in ActionCooldownState) ActionCooldownState {
	if len(in) == 0 {
		return ActionCooldownState{}
	}
	out := make(ActionCooldownState, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
