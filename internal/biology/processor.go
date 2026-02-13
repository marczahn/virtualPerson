package biology

import (
	"math"
	"time"

	"github.com/marczahn/person/internal/sense"
)

// Processor manages biological state updates through interaction rules,
// circadian modulation, decay toward baseline, and threshold monitoring.
type Processor struct {
	rules []Rule
}

// NewProcessor creates a Processor with the full interaction rule set.
func NewProcessor() *Processor {
	return &Processor{
		rules: AllRules(),
	}
}

// TickResult contains everything that changed during a single simulation tick.
type TickResult struct {
	Changes    []StateChange
	Thresholds []ThresholdResult
}

// Tick advances the biological state by the elapsed time since the last update.
// It applies: natural decay/drift, circadian modulation, interaction rules,
// hypothermia overrides, cortisol load accumulation, and threshold evaluation.
func (p *Processor) Tick(s *State) TickResult {
	now := time.Now()
	dt := now.Sub(s.LastUpdate).Seconds()
	if dt <= 0 {
		return TickResult{}
	}

	// Cap dt to prevent huge jumps after pauses (e.g., 5 minutes max).
	if dt > 300 {
		dt = 300
	}

	var allChanges []StateChange

	// 1. Advance circadian clock.
	s.CircadianPhase = math.Mod(s.CircadianPhase+dt/3600, 24)

	// 2. Apply natural decay toward baselines.
	allChanges = append(allChanges, p.applyDecay(s, dt)...)

	// 3. Apply circadian modulation as gentle pull toward circadian targets.
	allChanges = append(allChanges, p.applyCircadian(s, dt)...)

	// 4. Apply interaction rules (max one pass to prevent runaway chains).
	allChanges = append(allChanges, p.applyInteractions(s, dt)...)

	// 5. Apply hypothermia overrides (these take precedence, modify state directly).
	overrides := ApplyHypothermiaOverrides(s, dt)
	allChanges = append(allChanges, overrides...)

	// 6. Accumulate cortisol load for immune suppression.
	if s.Cortisol > 0.3 {
		s.CortisolLoad += (s.Cortisol - 0.3) * dt
	}

	// 7. Apply cortisol load immune suppression.
	immuneFactor := CortisolLoadImmuneSuppressionFactor(s.CortisolLoad)
	if immuneFactor < 1.0 && s.ImmuneResponse > 0 {
		suppression := s.ImmuneResponse * (1 - immuneFactor) * 0.001 * dt
		s.ImmuneResponse = ClampVariable(VarImmuneResponse, s.ImmuneResponse-suppression)
		allChanges = append(allChanges, StateChange{
			Variable: VarImmuneResponse,
			Delta:    -suppression,
			Source:   "cortisol_load_suppression",
		})
	}

	// 8. Apply natural hydration depletion.
	hydrationLoss := 0.001 / 60 * dt // 0.001/min at rest
	s.Hydration = ClampVariable(VarHydration, s.Hydration-hydrationLoss)
	allChanges = append(allChanges, StateChange{
		Variable: VarHydration,
		Delta:    -hydrationLoss,
		Source:   "natural_water_loss",
	})

	// 9. Apply natural fatigue accumulation.
	fatigueDelta := 0.05 / 3600 * dt // 0.05/hr
	s.Fatigue = ClampVariable(VarFatigue, s.Fatigue+fatigueDelta)
	allChanges = append(allChanges, StateChange{
		Variable: VarFatigue,
		Delta:    fatigueDelta,
		Source:   "wakefulness",
	})

	// 10. Evaluate critical thresholds.
	thresholds := EvaluateThresholds(s)

	s.LastUpdate = now

	return TickResult{
		Changes:    allChanges,
		Thresholds: thresholds,
	}
}

// ProcessStimulus applies the effects of a sensory event to the biological state.
// Returns the state changes produced.
func (p *Processor) ProcessStimulus(s *State, event sense.Event) []StateChange {
	var changes []StateChange

	switch event.Channel {
	case sense.Thermal:
		changes = append(changes, p.processThermal(s, event)...)
	case sense.Pain:
		changes = append(changes, p.processPain(s, event)...)
	case sense.Auditory:
		changes = append(changes, p.processAuditory(s, event)...)
	case sense.Visual:
		changes = append(changes, p.processVisual(s, event)...)
	}

	for _, c := range changes {
		s.Set(c.Variable, ClampVariable(c.Variable, s.Get(c.Variable)+c.Delta))
	}

	return changes
}

func (p *Processor) processThermal(s *State, event sense.Event) []StateChange {
	// Intensity maps to how extreme the temperature change is.
	// Positive intensity = heat, negative values would need separate handling.
	// For simplicity: intensity 0-0.5 = cold, 0.5-1 = hot.
	var changes []StateChange

	if event.Intensity < 0.5 {
		// Cold stimulus: drop body temp proportional to intensity
		coldDelta := -(0.5 - event.Intensity) * 4 // max -2°C immediate shift
		changes = append(changes, StateChange{
			Variable: VarBodyTemp,
			Delta:    coldDelta,
			Source:   "thermal_stimulus_cold",
		})
	} else {
		// Heat stimulus
		heatDelta := (event.Intensity - 0.5) * 4 // max +2°C immediate shift
		changes = append(changes, StateChange{
			Variable: VarBodyTemp,
			Delta:    heatDelta,
			Source:   "thermal_stimulus_heat",
		})
	}

	return changes
}

func (p *Processor) processPain(s *State, event sense.Event) []StateChange {
	return []StateChange{
		{
			Variable: VarPain,
			Delta:    event.Intensity,
			Source:   "pain_stimulus",
		},
		{
			Variable: VarAdrenaline,
			Delta:    event.Intensity * 0.3,
			Source:   "pain_adrenaline_response",
		},
	}
}

func (p *Processor) processAuditory(s *State, event sense.Event) []StateChange {
	// Loud/startling sounds trigger adrenaline.
	if event.Intensity > 0.7 {
		return []StateChange{
			{
				Variable: VarAdrenaline,
				Delta:    (event.Intensity - 0.7) * 0.5,
				Source:   "startle_response",
			},
			{
				Variable: VarHeartRate,
				Delta:    (event.Intensity - 0.7) * 20,
				Source:   "startle_response",
			},
		}
	}
	return nil
}

func (p *Processor) processVisual(s *State, event sense.Event) []StateChange {
	// Threatening visual stimuli trigger adrenaline.
	if event.Intensity > 0.8 {
		return []StateChange{
			{
				Variable: VarAdrenaline,
				Delta:    (event.Intensity - 0.8) * 0.4,
				Source:   "visual_threat",
			},
		}
	}
	return nil
}

// applyDecay moves variables toward their baselines based on known half-lives.
func (p *Processor) applyDecay(s *State, dt float64) []StateChange {
	var changes []StateChange

	// Adrenaline: half-life 2-3 min (~150s)
	if s.Adrenaline > 0.01 {
		decay := s.Adrenaline * (1 - math.Exp(-0.693/150*dt))
		s.Adrenaline = ClampVariable(VarAdrenaline, s.Adrenaline-decay)
		changes = append(changes, StateChange{VarAdrenaline, -decay, "natural_decay"})
	}

	// Cortisol: half-life 75 min (~4500s)
	if s.Cortisol > 0.1 {
		decay := (s.Cortisol - 0.1) * (1 - math.Exp(-0.693/4500*dt))
		s.Cortisol = ClampVariable(VarCortisol, s.Cortisol-decay)
		changes = append(changes, StateChange{VarCortisol, -decay, "natural_decay"})
	}

	// Heart rate: half-life ~75s toward baseline 70
	if math.Abs(s.HeartRate-70) > 1 {
		pull := (s.HeartRate - 70) * (1 - math.Exp(-0.693/75*dt))
		s.HeartRate = ClampVariable(VarHeartRate, s.HeartRate-pull)
		changes = append(changes, StateChange{VarHeartRate, -pull, "natural_decay"})
	}

	// Blood pressure: half-life ~180s toward baseline 120
	if math.Abs(s.BloodPressure-120) > 1 {
		pull := (s.BloodPressure - 120) * (1 - math.Exp(-0.693/180*dt))
		s.BloodPressure = ClampVariable(VarBloodPressure, s.BloodPressure-pull)
		changes = append(changes, StateChange{VarBloodPressure, -pull, "natural_decay"})
	}

	// Respiratory rate: half-life ~45s toward baseline 15
	if math.Abs(s.RespiratoryRate-15) > 0.5 {
		pull := (s.RespiratoryRate - 15) * (1 - math.Exp(-0.693/45*dt))
		s.RespiratoryRate = ClampVariable(VarRespiratoryRate, s.RespiratoryRate-pull)
		changes = append(changes, StateChange{VarRespiratoryRate, -pull, "natural_decay"})
	}

	// Body temperature: recovery ~0.1°C per 15 min toward circadian target
	circadian := ComputeCircadian(s.CircadianPhase)
	tempTarget := circadian.BodyTempTarget
	if math.Abs(s.BodyTemp-tempTarget) > 0.05 {
		pull := (s.BodyTemp - tempTarget) * (1 - math.Exp(-0.693/900*dt)) // ~15min half-life
		s.BodyTemp = ClampVariable(VarBodyTemp, s.BodyTemp-pull)
		changes = append(changes, StateChange{VarBodyTemp, -pull, "thermoregulation"})
	}

	// Muscle tension: half-life ~7.5 min (~450s)
	if s.MuscleTension > 0.01 {
		decay := s.MuscleTension * (1 - math.Exp(-0.693/450*dt))
		s.MuscleTension = ClampVariable(VarMuscleTension, s.MuscleTension-decay)
		changes = append(changes, StateChange{VarMuscleTension, -decay, "natural_decay"})
	}

	// Pain (acute): half-life ~30 min (~1800s)
	// Endorphins accelerate decay by factor 2.
	if s.Pain > 0.01 {
		halfLife := 1800.0
		if s.Endorphins > 0.2 {
			halfLife /= 2
		}
		decay := s.Pain * (1 - math.Exp(-0.693/halfLife*dt))
		s.Pain = ClampVariable(VarPain, s.Pain-decay)
		changes = append(changes, StateChange{VarPain, -decay, "natural_decay"})
	}

	// Endorphins: half-life ~25 min (~1500s)
	if s.Endorphins > 0.1 {
		decay := (s.Endorphins - 0.1) * (1 - math.Exp(-0.693/1500*dt))
		s.Endorphins = ClampVariable(VarEndorphins, s.Endorphins-decay)
		changes = append(changes, StateChange{VarEndorphins, -decay, "natural_decay"})
	}

	// Dopamine: tonic decay toward 0.3, half-life ~15 min for spikes
	if math.Abs(s.Dopamine-0.3) > 0.01 {
		pull := (s.Dopamine - 0.3) * (1 - math.Exp(-0.693/900*dt))
		s.Dopamine = ClampVariable(VarDopamine, s.Dopamine-pull)
		changes = append(changes, StateChange{VarDopamine, -pull, "natural_decay"})
	}

	// SpO2: recovery toward 98 when respiratory rate is adequate.
	if s.SpO2 < 98 && s.RespiratoryRate >= 10 {
		recoveryRate := 0.693 / 30 // ~30s half-life for recovery
		if s.RespiratoryRate > 25 {
			recoveryRate *= 1.5 // hyperventilation compensates faster
		}
		recovery := (98 - s.SpO2) * (1 - math.Exp(-recoveryRate*dt))
		s.SpO2 = ClampVariable(VarSpO2, s.SpO2+recovery)
		changes = append(changes, StateChange{VarSpO2, recovery, "respiratory_recovery"})
	}

	// Blood sugar: slow drift toward 90 (insulin homeostasis)
	if s.BloodSugar > 95 {
		pull := (s.BloodSugar - 90) * (1 - math.Exp(-0.693/5400*dt)) // ~90min
		s.BloodSugar = ClampVariable(VarBloodSugar, s.BloodSugar-pull)
		changes = append(changes, StateChange{VarBloodSugar, -pull, "insulin_homeostasis"})
	}

	// Glycogen: slow natural depletion during fasting (~0.0005/min)
	glycogenLoss := 0.0005 / 60 * dt
	s.Glycogen = ClampVariable(VarGlycogen, s.Glycogen-glycogenLoss)
	changes = append(changes, StateChange{VarGlycogen, -glycogenLoss, "basal_metabolism"})

	return changes
}

// applyCircadian gently pulls relevant variables toward their circadian targets.
func (p *Processor) applyCircadian(s *State, dt float64) []StateChange {
	var changes []StateChange
	circ := ComputeCircadian(s.CircadianPhase)

	// Cortisol baseline pull (slow, ~30 min to shift fully).
	cortisolPull := (circ.CortisolBaseline - s.Cortisol) * 0.0005 * dt
	if math.Abs(cortisolPull) > 0.0001 {
		s.Cortisol = ClampVariable(VarCortisol, s.Cortisol+cortisolPull)
		changes = append(changes, StateChange{VarCortisol, cortisolPull, "circadian"})
	}

	// Serotonin circadian shift.
	serotoninTarget := 0.5 + circ.SerotoninShift
	serotoninPull := (serotoninTarget - s.Serotonin) * 0.0002 * dt
	if math.Abs(serotoninPull) > 0.00001 {
		s.Serotonin = ClampVariable(VarSerotonin, s.Serotonin+serotoninPull)
		changes = append(changes, StateChange{VarSerotonin, serotoninPull, "circadian"})
	}

	return changes
}

// applyInteractions evaluates all rules against current state and applies deltas.
// Single pass — no cascading within a single tick to prevent runaway feedback.
func (p *Processor) applyInteractions(s *State, dt float64) []StateChange {
	var changes []StateChange

	hypothermiaReversal := IsHypothermiaReversal(s)

	for _, rule := range p.rules {
		if !rule.Condition(s) {
			continue
		}

		delta := rule.Delta(s, dt)
		if math.Abs(delta) < 1e-10 {
			continue
		}

		// In hypothermia reversal, suppress rules that increase HR or muscle tension.
		// The body's compensatory mechanisms are failing.
		if hypothermiaReversal {
			if (rule.Target == VarHeartRate || rule.Target == VarMuscleTension) && delta > 0 {
				continue
			}
		}

		newVal := ClampVariable(rule.Target, s.Get(rule.Target)+delta)
		actualDelta := newVal - s.Get(rule.Target)
		if math.Abs(actualDelta) < 1e-10 {
			continue
		}

		s.Set(rule.Target, newVal)
		changes = append(changes, StateChange{
			Variable: rule.Target,
			Delta:    actualDelta,
			Source:   rule.Name,
		})
	}

	return changes
}

// SignificantChanges filters a list of state changes to only those that
// cross a significance threshold (to avoid noisy output).
func SignificantChanges(changes []StateChange) []StateChange {
	var significant []StateChange
	for _, c := range changes {
		if isSignificant(c) {
			significant = append(significant, c)
		}
	}
	return significant
}

func isSignificant(c StateChange) bool {
	switch c.Variable {
	case VarHeartRate:
		return math.Abs(c.Delta) >= 2
	case VarBloodPressure:
		return math.Abs(c.Delta) >= 3
	case VarBodyTemp:
		return math.Abs(c.Delta) >= 0.1
	case VarRespiratoryRate:
		return math.Abs(c.Delta) >= 1
	case VarBloodSugar:
		return math.Abs(c.Delta) >= 2
	case VarSpO2:
		return math.Abs(c.Delta) >= 0.5
	default:
		// Ratio variables (0-1): significance at 0.01
		return math.Abs(c.Delta) >= 0.01
	}
}
