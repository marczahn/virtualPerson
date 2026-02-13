package biology

import "math"

// CircadianModulation holds the target adjustments from circadian rhythm
// for a given phase (hour of the internal clock).
type CircadianModulation struct {
	CortisolBaseline   float64 // added to cortisol baseline
	BodyTempTarget     float64 // target body temperature
	BloodPressureShift float64 // added to BP baseline
	ImmuneMultiplier   float64 // multiplier on immune response
	SerotoninShift     float64 // added to serotonin baseline
	Alertness          float64 // 0-1, modifies perceived fatigue
}

// ComputeCircadian returns the circadian modulation for a given phase.
// Phase is in hours (0-24, where 0 = midnight).
func ComputeCircadian(phase float64) CircadianModulation {
	return CircadianModulation{
		CortisolBaseline:   circadianCortisol(phase),
		BodyTempTarget:     circadianBodyTemp(phase),
		BloodPressureShift: circadianBloodPressure(phase),
		ImmuneMultiplier:   circadianImmune(phase),
		SerotoninShift:     circadianSerotonin(phase),
		Alertness:          circadianAlertness(phase),
	}
}

// cortisol_circadian = 0.15 + 0.25 * max(0, cos((phase - 7) * π / 12))
// Peak: 06:00-08:00 (~0.35-0.40), Trough: 23:00-02:00 (~0.05-0.10)
func circadianCortisol(phase float64) float64 {
	return 0.15 + 0.25*math.Max(0, math.Cos((phase-7)*math.Pi/12))
}

// temp_circadian = 36.6 + 0.5 * cos((phase - 17) * π / 12)
// Peak: ~17:00 (37.1°C), Trough: ~05:00 (36.1°C)
func circadianBodyTemp(phase float64) float64 {
	return 36.6 + 0.5*math.Cos((phase-17)*math.Pi/12)
}

// Morning surge: +15-25 mmHg between 05:00-09:00, trough 02:00-04:00.
// Modeled as sine wave peaking at ~10:00.
func circadianBloodPressure(phase float64) float64 {
	return 10 * math.Cos((phase-10)*math.Pi/12)
}

// Immune: strongest 22:00-02:00, weakest 06:00-10:00.
// immune_circadian_multiplier = 1.0 + 0.2 * cos((phase - 0) * π / 12)
func circadianImmune(phase float64) float64 {
	return 1.0 + 0.2*math.Cos(phase*math.Pi/12)
}

// Serotonin: +0.05 during day (08:00-18:00), -0.05 at night.
// Smooth transition using cosine centered on 13:00.
func circadianSerotonin(phase float64) float64 {
	return 0.05 * math.Cos((phase-13)*math.Pi/12)
}

// Alertness combines a 24h fundamental (peak ~16:00, trough ~04:00) with a
// 12h harmonic that creates the afternoon dip (~14:30) and reinforces the
// nighttime low (~02:30). Peaks emerge at ~10:00 and ~20:00.
func circadianAlertness(phase float64) float64 {
	fundamental := 0.2 * math.Cos(2*math.Pi*(phase-16)/24)
	dip := 0.25 * math.Cos(2*math.Pi*(phase-14.5)/12)
	return Clamp(0.5+fundamental-dip, 0, 1)
}
