package biology

import "math"

// Rule defines a single interaction: when a condition is met,
// a target variable changes by a computed delta.
type Rule struct {
	Name      string
	Condition func(s *State) bool
	Target    Variable
	Delta     func(s *State, dt float64) float64 // dt in seconds
}

// AllRules returns the complete set of biological interaction rules.
// Magnitudes are calibrated per-second (dt=1). The processor scales by actual dt.
func AllRules() []Rule {
	return []Rule{
		// ── Body Temperature ──────────────────────────────────────
		{
			Name:      "cold_shivering",
			Condition: func(s *State) bool { return s.BodyTemp < 35.5 && s.BodyTemp >= 33.0 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return (35.5 - s.BodyTemp) * 0.15 * dt },
		},
		{
			Name:      "cold_tachycardia",
			Condition: func(s *State) bool { return s.BodyTemp < 35.5 && s.BodyTemp >= 33.0 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (35.5 - s.BodyTemp) * 8 * dt },
		},
		{
			Name:      "cold_cortisol",
			Condition: func(s *State) bool { return s.BodyTemp < 35.5 && s.BodyTemp >= 33.0 },
			Target:    VarCortisol,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},
		{
			Name:      "cold_adrenaline",
			Condition: func(s *State) bool { return s.BodyTemp < 35.0 && s.BodyTemp >= 33.0 },
			Target:    VarAdrenaline,
			Delta:     func(s *State, dt float64) float64 { return 0.005 * dt },
		},
		{
			Name:      "fever_tachycardia",
			Condition: func(s *State) bool { return s.BodyTemp > 38.5 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (s.BodyTemp - 38.5) * 10 * dt },
		},
		{
			Name:      "fever_immune_boost",
			Condition: func(s *State) bool { return s.BodyTemp > 38.5 },
			Target:    VarImmuneResponse,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},
		{
			Name:      "fever_fatigue",
			Condition: func(s *State) bool { return s.BodyTemp > 39.5 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.003 * dt },
		},

		// ── Heart Rate ───────────────────────────────────────────
		{
			Name:      "hr_bp_coupling",
			Condition: func(s *State) bool { return s.HeartRate > 100 },
			Target:    VarBloodPressure,
			Delta:     func(s *State, dt float64) float64 { return (s.HeartRate - 100) * 0.3 * dt },
		},
		{
			Name:      "hr_respiratory_demand",
			Condition: func(s *State) bool { return s.HeartRate > 100 },
			Target:    VarRespiratoryRate,
			Delta:     func(s *State, dt float64) float64 { return (s.HeartRate - 100) * 0.08 * dt },
		},
		{
			Name:      "hr_metabolic_demand",
			Condition: func(s *State) bool { return s.HeartRate > 100 },
			Target:    VarBloodSugar,
			Delta:     func(s *State, dt float64) float64 { return -0.02 * dt },
		},
		{
			Name:      "hr_sustained_fatigue",
			Condition: func(s *State) bool { return s.HeartRate > 150 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.005 * dt },
		},
		{
			Name:      "bradycardia_bp_drop",
			Condition: func(s *State) bool { return s.HeartRate < 50 },
			Target:    VarBloodPressure,
			Delta:     func(s *State, dt float64) float64 { return -(50 - s.HeartRate) * 0.5 * dt },
		},
		{
			Name:      "bradycardia_spo2_drop",
			Condition: func(s *State) bool { return s.HeartRate < 50 },
			Target:    VarSpO2,
			Delta:     func(s *State, dt float64) float64 { return -0.001 * dt },
		},

		// ── Blood Pressure ───────────────────────────────────────
		{
			Name:      "low_bp_compensatory_hr",
			Condition: func(s *State) bool { return s.BloodPressure < 90 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (90 - s.BloodPressure) * 0.5 * dt },
		},
		{
			Name:      "low_bp_fatigue",
			Condition: func(s *State) bool { return s.BloodPressure < 90 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},
		{
			Name:      "low_bp_spo2",
			Condition: func(s *State) bool { return s.BloodPressure < 80 },
			Target:    VarSpO2,
			Delta:     func(s *State, dt float64) float64 { return -0.003 * dt },
		},
		{
			Name:      "high_bp_headache",
			Condition: func(s *State) bool { return s.BloodPressure > 160 },
			Target:    VarPain,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},

		// ── Respiratory Rate ─────────────────────────────────────
		{
			Name:      "hyperventilation_fatigue",
			Condition: func(s *State) bool { return s.RespiratoryRate > 35 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},
		{
			Name:      "hyperventilation_tension",
			Condition: func(s *State) bool { return s.RespiratoryRate > 35 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},
		{
			Name:      "hypoventilation_spo2",
			Condition: func(s *State) bool { return s.RespiratoryRate < 10 },
			Target:    VarSpO2,
			Delta:     func(s *State, dt float64) float64 { return -0.005 * dt },
		},
		{
			Name:      "respiratory_failure_spo2",
			Condition: func(s *State) bool { return s.RespiratoryRate < 8 },
			Target:    VarSpO2,
			Delta:     func(s *State, dt float64) float64 { return -0.01 * dt },
		},

		// ── Blood Sugar ──────────────────────────────────────────
		{
			Name:      "low_bs_adrenaline",
			Condition: func(s *State) bool { return s.BloodSugar < 70 },
			Target:    VarAdrenaline,
			Delta:     func(s *State, dt float64) float64 { return 0.003 * dt },
		},
		{
			Name:      "low_bs_cortisol",
			Condition: func(s *State) bool { return s.BloodSugar < 70 },
			Target:    VarCortisol,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},
		{
			Name:      "low_bs_hunger",
			Condition: func(s *State) bool { return s.BloodSugar < 70 },
			Target:    VarHunger,
			Delta: func(s *State, dt float64) float64 {
				return math.Min(1.0-s.Hunger, 0.005*dt)
			},
		},
		{
			Name:      "very_low_bs_fatigue",
			Condition: func(s *State) bool { return s.BloodSugar < 60 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.004 * dt },
		},
		{
			Name:      "very_low_bs_weakness",
			Condition: func(s *State) bool { return s.BloodSugar < 55 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return -0.002 * dt },
		},
		{
			Name:      "high_bs_thirst",
			Condition: func(s *State) bool { return s.BloodSugar > 140 },
			Target:    VarThirst,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},
		{
			Name:      "high_bs_fatigue",
			Condition: func(s *State) bool { return s.BloodSugar > 160 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},

		// ── Cortisol ─────────────────────────────────────────────
		{
			Name:      "cortisol_gluconeogenesis",
			Condition: func(s *State) bool { return s.Cortisol > 0.5 },
			Target:    VarBloodSugar,
			Delta:     func(s *State, dt float64) float64 { return 0.05 * dt },
		},
		{
			Name:      "cortisol_serotonin_depletion",
			Condition: func(s *State) bool { return s.Cortisol > 0.5 },
			Target:    VarSerotonin,
			Delta:     func(s *State, dt float64) float64 { return -0.0005 * dt },
		},
		{
			Name:      "cortisol_tachycardia",
			Condition: func(s *State) bool { return s.Cortisol > 0.3 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (s.Cortisol - 0.3) * 15 * dt },
		},
		{
			Name:      "cortisol_hypertension",
			Condition: func(s *State) bool { return s.Cortisol > 0.3 },
			Target:    VarBloodPressure,
			Delta:     func(s *State, dt float64) float64 { return (s.Cortisol - 0.3) * 20 * dt },
		},

		// ── Adrenaline ──────────────────────────────────────────
		{
			Name:      "adrenaline_tachycardia",
			Condition: func(s *State) bool { return s.Adrenaline > 0.2 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return s.Adrenaline * 60 * dt },
		},
		{
			Name:      "adrenaline_hypertension",
			Condition: func(s *State) bool { return s.Adrenaline > 0.2 },
			Target:    VarBloodPressure,
			Delta:     func(s *State, dt float64) float64 { return s.Adrenaline * 40 * dt },
		},
		{
			Name:      "adrenaline_respiratory",
			Condition: func(s *State) bool { return s.Adrenaline > 0.2 },
			Target:    VarRespiratoryRate,
			Delta:     func(s *State, dt float64) float64 { return s.Adrenaline * 10 * dt },
		},
		{
			Name:      "adrenaline_glycogenolysis",
			Condition: func(s *State) bool { return s.Adrenaline > 0.1 },
			Target:    VarBloodSugar,
			Delta:     func(s *State, dt float64) float64 { return s.Adrenaline * 0.5 * dt },
		},
		{
			Name:      "adrenaline_analgesia",
			Condition: func(s *State) bool { return s.Adrenaline > 0.3 && s.Pain > 0 },
			Target:    VarPain,
			Delta: func(s *State, dt float64) float64 {
				// Reduce pain proportionally: pain *= (1 - adrenaline*0.5), applied as delta
				return -s.Pain * s.Adrenaline * 0.5 * dt
			},
		},
		{
			Name:      "adrenaline_muscle_readiness",
			Condition: func(s *State) bool { return s.Adrenaline > 0.3 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return s.Adrenaline * 0.3 * dt },
		},
		{
			Name:      "adrenaline_crash_fatigue",
			Condition: func(s *State) bool { return s.Adrenaline > 0.7 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.008 * dt },
		},

		// ── Serotonin ───────────────────────────────────────────
		{
			Name:      "low_serotonin_cortisol_shift",
			Condition: func(s *State) bool { return s.Serotonin < 0.2 },
			Target:    VarCortisol,
			Delta:     func(s *State, dt float64) float64 { return 0.0005 * dt },
		},

		// ── Dopamine ────────────────────────────────────────────
		{
			Name:      "dopamine_analgesia",
			Condition: func(s *State) bool { return s.Dopamine > 0.5 && s.Pain > 0 },
			Target:    VarPain,
			Delta:     func(s *State, dt float64) float64 { return -s.Pain * 0.2 * dt },
		},
		{
			Name:      "dopamine_tachycardia",
			Condition: func(s *State) bool { return s.Dopamine > 0.6 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (s.Dopamine - 0.6) * 20 * dt },
		},
		{
			Name:      "low_dopamine_fatigue",
			Condition: func(s *State) bool { return s.Dopamine < 0.15 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},

		// ── Pain ─────────────────────────────────────────────────
		{
			Name:      "pain_cortisol",
			Condition: func(s *State) bool { return s.Pain > 0.3 },
			Target:    VarCortisol,
			Delta:     func(s *State, dt float64) float64 { return s.Pain * 0.003 * dt },
		},
		{
			Name:      "pain_tachycardia",
			Condition: func(s *State) bool { return s.Pain > 0.3 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return s.Pain * 30 * dt },
		},
		{
			Name:      "pain_hypertension",
			Condition: func(s *State) bool { return s.Pain > 0.3 },
			Target:    VarBloodPressure,
			Delta:     func(s *State, dt float64) float64 { return s.Pain * 25 * dt },
		},
		{
			Name:      "pain_guarding",
			Condition: func(s *State) bool { return s.Pain > 0.3 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return s.Pain * 0.4 * dt },
		},
		{
			Name:      "pain_adrenaline",
			Condition: func(s *State) bool { return s.Pain > 0.5 },
			Target:    VarAdrenaline,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},
		{
			Name:      "pain_tachypnea",
			Condition: func(s *State) bool { return s.Pain > 0.7 },
			Target:    VarRespiratoryRate,
			Delta:     func(s *State, dt float64) float64 { return (s.Pain - 0.7) * 15 * dt },
		},
		{
			Name:      "pain_endorphin_release",
			Condition: func(s *State) bool { return s.Pain > 0.8 },
			Target:    VarEndorphins,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},

		// ── Fatigue ──────────────────────────────────────────────
		{
			Name:      "fatigue_hr_baseline_shift",
			Condition: func(s *State) bool { return s.Fatigue > 0.6 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return 5 * dt }, // +5 bpm shift
		},
		{
			Name:      "fatigue_cortisol",
			Condition: func(s *State) bool { return s.Fatigue > 0.7 },
			Target:    VarCortisol,
			Delta:     func(s *State, dt float64) float64 { return 0.0005 * dt },
		},
		{
			Name:      "fatigue_immune_suppression",
			Condition: func(s *State) bool { return s.Fatigue > 0.8 },
			Target:    VarImmuneResponse,
			Delta:     func(s *State, dt float64) float64 { return -0.0005 * dt },
		},
		{
			Name:      "fatigue_tension_collapse",
			Condition: func(s *State) bool { return s.Fatigue > 0.9 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return -0.001 * dt },
		},

		// ── Immune Response ──────────────────────────────────────
		{
			Name:      "immune_fever",
			Condition: func(s *State) bool { return s.ImmuneResponse > 0.4 },
			Target:    VarBodyTemp,
			Delta:     func(s *State, dt float64) float64 { return (s.ImmuneResponse - 0.4) * 0.005 * dt },
		},
		{
			Name:      "immune_fatigue",
			Condition: func(s *State) bool { return s.ImmuneResponse > 0.5 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.003 * dt },
		},
		{
			Name:      "immune_body_aches",
			Condition: func(s *State) bool { return s.ImmuneResponse > 0.5 },
			Target:    VarMuscleTension,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},
		{
			Name:      "immune_appetite_suppression",
			Condition: func(s *State) bool { return s.ImmuneResponse > 0.6 },
			Target:    VarHunger,
			Delta:     func(s *State, dt float64) float64 { return -0.002 * dt },
		},
		{
			Name:      "immune_serotonin_depletion",
			Condition: func(s *State) bool { return s.ImmuneResponse > 0.3 },
			Target:    VarSerotonin,
			Delta:     func(s *State, dt float64) float64 { return -0.0003 * dt },
		},

		// ── Muscle Tension ───────────────────────────────────────
		{
			Name:      "tension_pain",
			Condition: func(s *State) bool { return s.MuscleTension > 0.6 },
			Target:    VarPain,
			Delta:     func(s *State, dt float64) float64 { return 0.001 * dt },
		},
		{
			Name:      "tension_blood_sugar",
			Condition: func(s *State) bool { return s.MuscleTension > 0.5 },
			Target:    VarBloodSugar,
			Delta:     func(s *State, dt float64) float64 { return -0.01 * dt },
		},
		{
			Name:      "tension_fatigue",
			Condition: func(s *State) bool { return s.MuscleTension > 0.7 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.002 * dt },
		},

		// ── Hydration ────────────────────────────────────────────
		{
			Name:      "dehydration_tachycardia",
			Condition: func(s *State) bool { return s.Hydration < 0.6 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (0.6 - s.Hydration) * 30 * dt },
		},
		{
			Name:      "dehydration_hypotension",
			Condition: func(s *State) bool { return s.Hydration < 0.6 },
			Target:    VarBloodPressure,
			Delta:     func(s *State, dt float64) float64 { return -(0.6 - s.Hydration) * 40 * dt },
		},
		{
			Name:      "severe_dehydration_fatigue",
			Condition: func(s *State) bool { return s.Hydration < 0.4 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.005 * dt },
		},
		{
			Name:      "severe_dehydration_immune",
			Condition: func(s *State) bool { return s.Hydration < 0.3 },
			Target:    VarImmuneResponse,
			Delta:     func(s *State, dt float64) float64 { return -0.002 * dt },
		},
		{
			Name:      "dehydration_thirst",
			Condition: func(s *State) bool { return s.Hydration < 0.7 },
			Target:    VarThirst,
			Delta: func(s *State, dt float64) float64 {
				return math.Min(1.0-s.Thirst, (0.7-s.Hydration)*0.01*dt)
			},
		},
		{
			Name:      "hydrated_thirst_suppression",
			Condition: func(s *State) bool { return s.Hydration > 0.85 && s.Thirst > 0 },
			Target:    VarThirst,
			Delta:     func(s *State, dt float64) float64 { return -0.01 * dt },
		},

		// ── SpO2 ─────────────────────────────────────────────────
		{
			Name:      "hypoxia_tachycardia",
			Condition: func(s *State) bool { return s.SpO2 < 94 },
			Target:    VarHeartRate,
			Delta:     func(s *State, dt float64) float64 { return (94 - s.SpO2) * 3 * dt },
		},
		{
			Name:      "hypoxia_tachypnea",
			Condition: func(s *State) bool { return s.SpO2 < 94 },
			Target:    VarRespiratoryRate,
			Delta:     func(s *State, dt float64) float64 { return (94 - s.SpO2) * 2 * dt },
		},
		{
			Name:      "hypoxia_adrenaline",
			Condition: func(s *State) bool { return s.SpO2 < 90 },
			Target:    VarAdrenaline,
			Delta:     func(s *State, dt float64) float64 { return 0.005 * dt },
		},
		{
			Name:      "severe_hypoxia_fatigue",
			Condition: func(s *State) bool { return s.SpO2 < 88 },
			Target:    VarFatigue,
			Delta:     func(s *State, dt float64) float64 { return 0.01 * dt },
		},

		// ── Endorphins ───────────────────────────────────────────
		{
			Name:      "endorphin_analgesia",
			Condition: func(s *State) bool { return s.Endorphins > 0.3 && s.Pain > 0 },
			Target:    VarPain,
			Delta: func(s *State, dt float64) float64 {
				// Endorphins accelerate pain decay by factor 2
				return -s.Pain * s.Endorphins * 0.3 * dt
			},
		},

		// ── Glycogen → Blood Sugar buffering ─────────────────────
		{
			Name:      "glycogen_bs_buffer",
			Condition: func(s *State) bool { return s.BloodSugar < 80 && s.Glycogen > 0.05 },
			Target:    VarBloodSugar,
			Delta:     func(s *State, dt float64) float64 { return 0.1 * dt }, // slow glucose release
		},
		{
			Name:      "glycogen_depletion_from_low_bs",
			Condition: func(s *State) bool { return s.BloodSugar < 80 && s.Glycogen > 0.05 },
			Target:    VarGlycogen,
			Delta:     func(s *State, dt float64) float64 { return -0.0005 * dt },
		},

		// ── Hunger derived from blood sugar + glycogen ───────────
		{
			Name:      "hunger_suppression_normal_bs",
			Condition: func(s *State) bool { return s.BloodSugar > 110 && s.Hunger > 0 },
			Target:    VarHunger,
			Delta:     func(s *State, dt float64) float64 { return -0.003 * dt },
		},
	}
}

// Clamp restricts a value to [min, max].
func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// variableRanges defines the allowed range for each biological variable.
var variableRanges = map[Variable][2]float64{
	VarBodyTemp:        {25, 43},
	VarHeartRate:       {40, 200},
	VarBloodPressure:   {80, 200},
	VarRespiratoryRate: {8, 40},
	VarHunger:          {0, 1},
	VarThirst:          {0, 1},
	VarFatigue:         {0, 1},
	VarPain:            {0, 1},
	VarMuscleTension:   {0, 1},
	VarBloodSugar:      {50, 200},
	VarCortisol:        {0, 1},
	VarAdrenaline:      {0, 1},
	VarSerotonin:       {0, 1},
	VarDopamine:        {0, 1},
	VarImmuneResponse:  {0, 1},
	VarCircadianPhase:  {0, 24},
	VarSpO2:            {70, 100},
	VarHydration:       {0, 1},
	VarGlycogen:        {0, 1},
	VarEndorphins:      {0, 1},
}

// ClampVariable restricts a variable value to its valid range.
func ClampVariable(v Variable, val float64) float64 {
	r, ok := variableRanges[v]
	if !ok {
		return val
	}
	return Clamp(val, r[0], r[1])
}
