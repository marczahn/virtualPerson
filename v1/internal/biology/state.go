package biology

import "time"

// State holds the 20 biological variables that model the person's body.
// All ratio fields are in [0,1] unless otherwise noted.
type State struct {
	BodyTemp        float64 // °C, baseline 36.6, range 34-42
	HeartRate       float64 // bpm, baseline 70, range 40-200
	BloodPressure   float64 // systolic mmHg, baseline 120, range 80-200
	RespiratoryRate float64 // breaths/min, baseline 15, range 8-40
	Hunger          float64 // 0-1, derived from blood sugar + glycogen
	Thirst          float64 // 0-1, derived from hydration level
	Fatigue         float64 // 0-1, accumulates ~0.05/hr waking
	Pain            float64 // 0-1, partially mental-state-influenced (+-30-40%)
	MuscleTension   float64 // 0-1, follows stress + cold + pain
	BloodSugar      float64 // mg/dL, baseline 90, range 50-200
	Cortisol        float64 // 0-1, peaks 15-30min after stressor, half-life 60-90min
	Adrenaline      float64 // 0-1, peaks in seconds, half-life 2-3min
	Serotonin       float64 // 0-1, slowest changes (hours-days)
	Dopamine        float64 // 0-1, phasic spikes (seconds) + tonic baseline (hours)
	ImmuneResponse  float64 // 0-1, slowest system (hours-days), suppressed by cortisol load
	CircadianPhase  float64 // hours 0-24, advances 1hr/hr
	SpO2            float64 // %, baseline 98, range 70-100
	Hydration       float64 // 0-1, depletes ~0.001/min resting
	Glycogen        float64 // 0-1, buffers blood sugar, depletes 12-24hr fasting
	Endorphins      float64 // 0-1, released after 20-30min sustained stress

	// Accumulated cortisol load for immune suppression calculation.
	// cortisol_load += max(0, cortisol - 0.3) * dt
	CortisolLoad float64

	// Timestamp of the last state update for delta-time calculations.
	LastUpdate time.Time
}

// NewDefaultState returns a State with healthy baseline values,
// as if the person just woke up at 8:00 AM, well-rested and fed.
func NewDefaultState() State {
	return State{
		BodyTemp:        36.6,
		HeartRate:       70,
		BloodPressure:   120,
		RespiratoryRate: 15,
		Hunger:          0.0,
		Thirst:          0.0,
		Fatigue:         0.0,
		Pain:            0.0,
		MuscleTension:   0.0,
		BloodSugar:      90,
		Cortisol:        0.1,
		Adrenaline:      0.0,
		Serotonin:       0.5,
		Dopamine:        0.3,
		ImmuneResponse:  0.1,
		CircadianPhase:  8.0,
		SpO2:            98,
		Hydration:       0.8,
		Glycogen:        0.7,
		Endorphins:      0.1,
		CortisolLoad:    0.0,
		LastUpdate:      time.Now(),
	}
}

// Variable identifies a single biological state variable by name.
type Variable int

const (
	VarBodyTemp Variable = iota
	VarHeartRate
	VarBloodPressure
	VarRespiratoryRate
	VarHunger
	VarThirst
	VarFatigue
	VarPain
	VarMuscleTension
	VarBloodSugar
	VarCortisol
	VarAdrenaline
	VarSerotonin
	VarDopamine
	VarImmuneResponse
	VarCircadianPhase
	VarSpO2
	VarHydration
	VarGlycogen
	VarEndorphins
)

var variableNames = [...]string{
	"body_temp",
	"heart_rate",
	"blood_pressure",
	"respiratory_rate",
	"hunger",
	"thirst",
	"fatigue",
	"pain",
	"muscle_tension",
	"blood_sugar",
	"cortisol",
	"adrenaline",
	"serotonin",
	"dopamine",
	"immune_response",
	"circadian_phase",
	"spo2",
	"hydration",
	"glycogen",
	"endorphins",
}

func (v Variable) String() string {
	if int(v) < len(variableNames) {
		return variableNames[v]
	}
	return "unknown"
}

// Get returns the current value of the given variable from the state.
func (s *State) Get(v Variable) float64 {
	switch v {
	case VarBodyTemp:
		return s.BodyTemp
	case VarHeartRate:
		return s.HeartRate
	case VarBloodPressure:
		return s.BloodPressure
	case VarRespiratoryRate:
		return s.RespiratoryRate
	case VarHunger:
		return s.Hunger
	case VarThirst:
		return s.Thirst
	case VarFatigue:
		return s.Fatigue
	case VarPain:
		return s.Pain
	case VarMuscleTension:
		return s.MuscleTension
	case VarBloodSugar:
		return s.BloodSugar
	case VarCortisol:
		return s.Cortisol
	case VarAdrenaline:
		return s.Adrenaline
	case VarSerotonin:
		return s.Serotonin
	case VarDopamine:
		return s.Dopamine
	case VarImmuneResponse:
		return s.ImmuneResponse
	case VarCircadianPhase:
		return s.CircadianPhase
	case VarSpO2:
		return s.SpO2
	case VarHydration:
		return s.Hydration
	case VarGlycogen:
		return s.Glycogen
	case VarEndorphins:
		return s.Endorphins
	default:
		return 0
	}
}

// Set updates the given variable in the state.
func (s *State) Set(v Variable, val float64) {
	switch v {
	case VarBodyTemp:
		s.BodyTemp = val
	case VarHeartRate:
		s.HeartRate = val
	case VarBloodPressure:
		s.BloodPressure = val
	case VarRespiratoryRate:
		s.RespiratoryRate = val
	case VarHunger:
		s.Hunger = val
	case VarThirst:
		s.Thirst = val
	case VarFatigue:
		s.Fatigue = val
	case VarPain:
		s.Pain = val
	case VarMuscleTension:
		s.MuscleTension = val
	case VarBloodSugar:
		s.BloodSugar = val
	case VarCortisol:
		s.Cortisol = val
	case VarAdrenaline:
		s.Adrenaline = val
	case VarSerotonin:
		s.Serotonin = val
	case VarDopamine:
		s.Dopamine = val
	case VarImmuneResponse:
		s.ImmuneResponse = val
	case VarCircadianPhase:
		s.CircadianPhase = val
	case VarSpO2:
		s.SpO2 = val
	case VarHydration:
		s.Hydration = val
	case VarGlycogen:
		s.Glycogen = val
	case VarEndorphins:
		s.Endorphins = val
	}
}

// StateChange represents a modification to a biological variable,
// produced by interaction rules or external stimuli.
type StateChange struct {
	Variable Variable
	Delta    float64 // additive change (can be negative)
	Source   string  // what caused this change (e.g., "circadian", "cold_exposure", "cortisol_interaction")
}
