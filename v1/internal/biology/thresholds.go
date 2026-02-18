package biology

// CriticalState represents a qualitative shift in the person's condition
// when biological variables enter dangerous ranges.
type CriticalState int

const (
	Normal CriticalState = iota
	Impaired              // cognitive or physical impairment
	Critical              // life-threatening, requires intervention
	Unconscious           // loss of consciousness
	Lethal                // death without immediate intervention
)

var criticalStateNames = [...]string{
	"normal",
	"impaired",
	"critical",
	"unconscious",
	"lethal",
}

func (c CriticalState) String() string {
	if int(c) < len(criticalStateNames) {
		return criticalStateNames[c]
	}
	return "unknown"
}

// ThresholdResult describes a critical condition detected in the state.
type ThresholdResult struct {
	Condition   CriticalState
	System      string // which system is failing (e.g., "thermoregulation", "glycemic", "respiratory")
	Description string // human-readable explanation
}

// EvaluateThresholds checks the biological state for critical conditions.
// Returns all active threshold breaches, or nil if the person is in normal range.
func EvaluateThresholds(s *State) []ThresholdResult {
	var results []ThresholdResult

	results = append(results, evaluateHypothermia(s)...)
	results = append(results, evaluateHyperthermia(s)...)
	results = append(results, evaluateHypoglycemia(s)...)
	results = append(results, evaluateHyperglycemia(s)...)
	results = append(results, evaluateSpO2(s)...)
	results = append(results, evaluateBloodPressure(s)...)

	return results
}

// IsHypothermiaReversal returns true when the body has stopped shivering
// due to energy depletion, not because it recovered. This is the critical
// 33°C threshold where cessation of shivering means worsening, not improvement.
func IsHypothermiaReversal(s *State) bool {
	return s.BodyTemp < 33.0
}

func evaluateHypothermia(s *State) []ThresholdResult {
	var results []ThresholdResult

	switch {
	case s.BodyTemp < 28:
		results = append(results, ThresholdResult{
			Condition:   Lethal,
			System:      "thermoregulation",
			Description: "ventricular fibrillation risk, lethal without intervention",
		})
	case s.BodyTemp < 30:
		results = append(results, ThresholdResult{
			Condition:   Unconscious,
			System:      "thermoregulation",
			Description: "severe hypothermia, cardiac arrhythmia, loss of consciousness",
		})
	case s.BodyTemp < 33:
		results = append(results, ThresholdResult{
			Condition:   Critical,
			System:      "thermoregulation",
			Description: "moderate hypothermia, shivering stopped, confusion, HR dropping",
		})
	case s.BodyTemp < 35:
		results = append(results, ThresholdResult{
			Condition:   Impaired,
			System:      "thermoregulation",
			Description: "mild hypothermia, maximal shivering, cognitive impairment",
		})
	}

	return results
}

func evaluateHyperthermia(s *State) []ThresholdResult {
	var results []ThresholdResult

	switch {
	case s.BodyTemp > 42:
		results = append(results, ThresholdResult{
			Condition:   Lethal,
			System:      "thermoregulation",
			Description: "lethal hyperthermia, protein denaturation",
		})
	case s.BodyTemp > 41.5:
		results = append(results, ThresholdResult{
			Condition:   Critical,
			System:      "thermoregulation",
			Description: "organ damage, multi-organ failure cascade",
		})
	case s.BodyTemp > 40:
		results = append(results, ThresholdResult{
			Condition:   Impaired,
			System:      "thermoregulation",
			Description: "heat stroke onset, confusion, BP dropping",
		})
	}

	return results
}

func evaluateHypoglycemia(s *State) []ThresholdResult {
	var results []ThresholdResult

	switch {
	case s.BloodSugar < 25:
		results = append(results, ThresholdResult{
			Condition:   Lethal,
			System:      "glycemic",
			Description: "lethal hypoglycemia without intervention",
		})
	case s.BloodSugar < 35:
		results = append(results, ThresholdResult{
			Condition:   Unconscious,
			System:      "glycemic",
			Description: "loss of consciousness from hypoglycemia",
		})
	case s.BloodSugar < 45:
		results = append(results, ThresholdResult{
			Condition:   Critical,
			System:      "glycemic",
			Description: "seizure risk, behavioral changes, loss of motor control",
		})
	case s.BloodSugar < 55:
		results = append(results, ThresholdResult{
			Condition:   Impaired,
			System:      "glycemic",
			Description: "neuroglycopenia, confusion, slowed reaction time",
		})
	}

	return results
}

func evaluateHyperglycemia(s *State) []ThresholdResult {
	var results []ThresholdResult

	if s.BloodSugar > 300 {
		results = append(results, ThresholdResult{
			Condition:   Critical,
			System:      "glycemic",
			Description: "diabetic ketoacidosis territory, nausea, altered consciousness",
		})
	} else if s.BloodSugar > 180 {
		results = append(results, ThresholdResult{
			Condition:   Impaired,
			System:      "glycemic",
			Description: "hyperglycemia, increased thirst, dehydration acceleration",
		})
	}

	return results
}

func evaluateSpO2(s *State) []ThresholdResult {
	var results []ThresholdResult

	switch {
	case s.SpO2 < 60:
		results = append(results, ThresholdResult{
			Condition:   Lethal,
			System:      "respiratory",
			Description: "lethal hypoxia, organ damage",
		})
	case s.SpO2 < 75:
		results = append(results, ThresholdResult{
			Condition:   Unconscious,
			System:      "respiratory",
			Description: "severe hypoxia, cardiac arrhythmia risk",
		})
	case s.SpO2 < 85:
		results = append(results, ThresholdResult{
			Condition:   Critical,
			System:      "respiratory",
			Description: "confusion, agitation, loss of coordination",
		})
	case s.SpO2 < 90:
		results = append(results, ThresholdResult{
			Condition:   Impaired,
			System:      "respiratory",
			Description: "peripheral cyanosis, cognitive impairment, adrenaline spike",
		})
	}

	return results
}

func evaluateBloodPressure(s *State) []ThresholdResult {
	var results []ThresholdResult

	if s.BloodPressure < 70 {
		results = append(results, ThresholdResult{
			Condition:   Critical,
			System:      "cardiovascular",
			Description: "shock territory, cognitive impairment",
		})
	}

	return results
}

// ApplyHypothermiaOverrides directly modifies the state to enforce hypothermia
// reversal effects. At <33°C, shivering stops and HR drops — the person feels
// better while dying. This must be called AFTER all interaction rules, because
// it overrides their effects with absolute targets.
func ApplyHypothermiaOverrides(s *State, dt float64) []StateChange {
	if !IsHypothermiaReversal(s) {
		return nil
	}

	var changes []StateChange

	// Shivering stops — force muscle tension toward 0.3 (from whatever it is).
	if s.MuscleTension > 0.3 {
		drop := (s.MuscleTension - 0.3) * 0.1 * dt // 10%/sec convergence
		if drop > s.MuscleTension-0.3 {
			drop = s.MuscleTension - 0.3
		}
		s.MuscleTension -= drop
		changes = append(changes, StateChange{VarMuscleTension, -drop, "hypothermia_reversal"})
	}

	// HR drops toward 50 then lower. Force convergence regardless of other rules.
	target := 40.0 + (s.BodyTemp-28)*2 // 40 bpm at 28°C, 50 at 33°C
	if s.HeartRate > target {
		drop := (s.HeartRate - target) * 0.05 * dt // 5%/sec convergence
		if drop > s.HeartRate-target {
			drop = s.HeartRate - target
		}
		s.HeartRate -= drop
		changes = append(changes, StateChange{VarHeartRate, -drop, "hypothermia_reversal"})
	}

	// Adrenaline response weakens in severe hypothermia.
	if s.Adrenaline > 0.05 {
		drop := s.Adrenaline * 0.1 * dt
		s.Adrenaline -= drop
		changes = append(changes, StateChange{VarAdrenaline, -drop, "hypothermia_reversal"})
	}

	return changes
}

// CortisolLoadImmuneSuppressionFactor returns the immune efficiency multiplier
// based on accumulated cortisol load.
// Formula: 1.0 / (1.0 + cortisol_load * 0.1)
// Returns 1.0 when no load, approaches 0 with extreme load.
func CortisolLoadImmuneSuppressionFactor(cortisolLoad float64) float64 {
	return 1.0 / (1.0 + cortisolLoad*0.1)
}
