package psychology

import (
	"github.com/marczahn/person/internal/biology"
)

// Processor transforms biological state into psychological state.
// It applies personality modulation, emotional regulation, distortion
// activation, coping selection, and isolation effects.
type Processor struct {
	personality Personality
	regulation  RegulationState
	memories    *EmotionalMemoryStore
	isolation   IsolationState
}

// NewProcessor creates a psychology processor for a person with the given personality.
func NewProcessor(p Personality) *Processor {
	return &Processor{
		personality: p,
		regulation: RegulationState{
			Capacity: BaselineRegulation(p),
		},
		memories: NewEmotionalMemoryStore(),
	}
}

// Process transforms a biological state into a psychological state.
// dt is the elapsed time in seconds since the last processing cycle.
// stressorControllability is in [0,1] — how much the current stressor
// can be influenced by the person's actions (0 = helpless, 1 = fully controllable).
func (proc *Processor) Process(bio *biology.State, dt float64, stressorControllability float64) State {
	// Step 1: Compute raw affect dimensions from biology.
	rawArousal := computeArousal(bio)
	rawValence := computeValence(bio)
	energy := computeEnergy(bio)
	cogLoad := computeCognitiveLoad(bio)

	// Step 2: Apply personality modulation to negative signals.
	negMult := NegativeEmotionMultiplier(proc.personality)
	if rawValence < 0 {
		rawValence *= negMult
		rawValence = clamp(rawValence, -1, 1)
	}

	// Step 3: Compute stress level (used by regulation, coping, distortions).
	stress := computeStress(rawArousal, rawValence, cogLoad)

	// Step 4: Update emotional regulation (depletion/recovery).
	proc.regulation = UpdateRegulation(proc.regulation, proc.personality, stress, bio.Fatigue, dt)

	// Step 5: Apply regulation to affect dimensions.
	effReg := EffectiveCapacity(proc.regulation, bio.Fatigue)
	arousal := rawArousal * (1.0 - effReg*0.5)
	valence := rawValence * (1.0 - effReg*0.3)

	// Step 6: Activate cognitive distortions.
	distortions := ActivateDistortions(stress, effReg, proc.personality)

	// Step 7: Select coping strategies.
	resources := CopingResources{
		Cognitive:  1.0 - cogLoad,
		Energy:     energy,
		Regulation: proc.regulation.Capacity,
	}
	coping := SelectCoping(stress, proc.personality, resources, stressorControllability)

	// Step 8: Update isolation state.
	proc.isolation = UpdateIsolation(proc.isolation, proc.personality, dt)

	// Step 9: Query emotional memory activations.
	activations := proc.memories.QueryActivations(bio)

	return State{
		Arousal:            clamp(arousal, 0, 1),
		Valence:            clamp(valence, -1, 1),
		Energy:             energy,
		CognitiveLoad:      cogLoad,
		RegulationCapacity: proc.regulation.Capacity,
		ActiveDistortions:  distortions,
		ActiveCoping:       coping,
		ActivatedMemories:  activations,
		Isolation:          proc.isolation,
	}
}

// FeedbackChanges returns the biological state changes that the current
// psychological state feeds back into the body. These should be applied
// by the biology processor each cycle.
func (proc *Processor) FeedbackChanges(ps *State, dt float64) []biology.StateChange {
	var changes []biology.StateChange

	for _, c := range ps.ActiveCoping {
		switch c {
		case Rumination:
			changes = append(changes,
				biology.StateChange{Variable: biology.VarCortisol, Delta: 0.02 * dt, Source: "psych_rumination"},
				biology.StateChange{Variable: biology.VarSerotonin, Delta: -0.01 * dt, Source: "psych_rumination"},
			)
		case Acceptance, Reappraisal:
			changes = append(changes,
				biology.StateChange{Variable: biology.VarCortisol, Delta: -0.01 * dt, Source: "psych_acceptance"},
				biology.StateChange{Variable: biology.VarSerotonin, Delta: 0.005 * dt, Source: "psych_acceptance"},
			)
		}
	}

	for _, d := range ps.ActiveDistortions {
		if d == Catastrophizing {
			changes = append(changes,
				biology.StateChange{Variable: biology.VarAdrenaline, Delta: 0.03 * dt, Source: "psych_catastrophizing"},
				biology.StateChange{Variable: biology.VarCortisol, Delta: 0.02 * dt, Source: "psych_catastrophizing"},
			)
		}
	}

	if len(ps.ActiveDistortions) > 2 {
		changes = append(changes,
			biology.StateChange{Variable: biology.VarCortisol, Delta: 0.01 * dt, Source: "psych_distortion_load"},
		)
	}

	return changes
}

// RecordSocialContact resets isolation tracking when the person has social interaction.
func (proc *Processor) RecordSocialContact() {
	proc.isolation = IsolationState{}
}

// AddMemory stores a new emotional memory association.
func (proc *Processor) AddMemory(mem EmotionalMemory) {
	proc.memories.Add(mem)
}

// computeArousal derives arousal from adrenaline, heart rate, cortisol, fatigue.
// Formula from psychologist advisory section 1.
func computeArousal(bio *biology.State) float64 {
	// Normalize heart rate: 40-200 range, baseline 70.
	// We normalize relative to the baseline: 0 at 40, ~0.19 at 70, 1 at 200.
	normHR := (bio.HeartRate - 40) / 160

	a := 0.30*bio.Adrenaline +
		0.25*normHR +
		0.20*bio.Cortisol +
		0.15*0 + // norepinephrine not modeled; absorbed into adrenaline
		0.10*(1-bio.Fatigue)

	return clamp(a, 0, 1)
}

// computeValence derives valence bias from serotonin, dopamine, endorphins, cortisol.
func computeValence(bio *biology.State) float64 {
	v := 0.35*bio.Serotonin +
		0.30*bio.Dopamine +
		0.15*bio.Endorphins +
		-0.20*bio.Cortisol

	return clamp(v, -1, 1)
}

// computeEnergy derives energy from fatigue, blood sugar, circadian alertness, dopamine.
func computeEnergy(bio *biology.State) float64 {
	// Normalize blood sugar: 50-200 range.
	normBS := (bio.BloodSugar - 50) / 150

	// Get circadian alertness from current phase.
	circ := biology.ComputeCircadian(bio.CircadianPhase)

	e := 0.30*(1-bio.Fatigue) +
		0.25*normBS +
		0.25*circ.Alertness +
		0.20*bio.Dopamine

	return clamp(e, 0, 1)
}

// computeCognitiveLoad derives cognitive load from cortisol duration, fatigue, blood sugar.
func computeCognitiveLoad(bio *biology.State) float64 {
	// Normalize blood sugar inverted: low blood sugar = high cognitive load.
	normBSInv := 1.0 - (bio.BloodSugar-50)/150

	// Use cortisol load as proxy for stress duration.
	// CortisolLoad accumulates over time; normalize with a soft cap.
	normCortisolDuration := bio.CortisolLoad / (bio.CortisolLoad + 5.0) // sigmoid-ish, reaches 0.5 at load=5

	cl := 0.30*normCortisolDuration +
		0.25*bio.Fatigue +
		0.25*normBSInv +
		0.20*0 // sleep debt not modeled yet; placeholder

	return clamp(cl, 0, 1)
}

// computeStress derives an overall stress level used for coping/distortion decisions.
// High arousal + negative valence + high cognitive load = high stress.
func computeStress(arousal, valence, cogLoad float64) float64 {
	negValenceContrib := 0.0
	if valence < 0 {
		negValenceContrib = -valence // 0 to 1
	}
	s := 0.35*arousal + 0.35*negValenceContrib + 0.30*cogLoad
	return clamp(s, 0, 1)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
