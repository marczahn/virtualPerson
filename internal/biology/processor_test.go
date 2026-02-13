package biology

import (
	"testing"
	"time"

	"github.com/marczahn/person/internal/sense"
)

func newTestState() *State {
	s := NewDefaultState()
	s.LastUpdate = time.Now().Add(-time.Second) // 1 second ago
	return &s
}

func TestColdChain_TemperatureDropTriggersShiveringAndCortisol(t *testing.T) {
	s := newTestState()
	s.BodyTemp = 34.5 // below 35.5 threshold
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	result := p.Tick(s)

	if s.MuscleTension <= 0 {
		t.Error("expected muscle tension increase from shivering, got 0")
	}
	if s.Cortisol <= 0.1 {
		t.Errorf("expected cortisol increase from cold stress, got %f", s.Cortisol)
	}

	_ = result // verify no panic
}

func TestHungerChain_LowBloodSugarDrivesHungerAndAdrenaline(t *testing.T) {
	s := newTestState()
	s.BloodSugar = 65 // below 70 threshold
	s.Hunger = 0.0
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Hunger <= 0 {
		t.Error("expected hunger increase from low blood sugar")
	}
	if s.Adrenaline <= 0 {
		t.Errorf("expected adrenaline increase from low blood sugar, got %f", s.Adrenaline)
	}
}

func TestAdrenalineDecay_HalfLifeApproximately150Seconds(t *testing.T) {
	s := newTestState()
	s.Adrenaline = 0.8
	s.LastUpdate = time.Now().Add(-150 * time.Second) // one half-life

	p := NewProcessor()
	p.Tick(s)

	// After one half-life, should be roughly halved.
	// Allow some tolerance due to interaction rules also affecting adrenaline.
	if s.Adrenaline > 0.55 {
		t.Errorf("expected adrenaline to decay significantly after 150s, got %f", s.Adrenaline)
	}
	if s.Adrenaline < 0.2 {
		t.Errorf("adrenaline decayed too fast, got %f", s.Adrenaline)
	}
}

func TestHypothermiaReversal_ShiveringStopsBelow33(t *testing.T) {
	s := newTestState()
	s.BodyTemp = 32.5  // below 33°C reversal point
	s.MuscleTension = 0.8
	s.HeartRate = 100
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Muscle tension should decrease (shivering stops)
	if s.MuscleTension >= 0.8 {
		t.Errorf("expected muscle tension to decrease at <33°C, got %f", s.MuscleTension)
	}

	// Heart rate should start dropping (bradycardia)
	if s.HeartRate >= 100 {
		t.Errorf("expected heart rate to decrease at <33°C, got %f", s.HeartRate)
	}
}

func TestCortisolLoadAccumulation(t *testing.T) {
	s := newTestState()
	s.Cortisol = 0.6 // above 0.3 threshold
	s.CortisolLoad = 0
	s.LastUpdate = time.Now().Add(-60 * time.Second) // 60 seconds

	p := NewProcessor()
	p.Tick(s)

	// Cortisol load should accumulate: (0.6 - 0.3) * 60 = 18 (approximately)
	// Actual value will differ because cortisol decays during tick.
	if s.CortisolLoad <= 0 {
		t.Error("expected cortisol load to accumulate when cortisol > 0.3")
	}
}

func TestInteractionRules_DoNotPanic(t *testing.T) {
	// Ensure all rules execute without panicking on default state.
	s := newTestState()
	s.LastUpdate = time.Now().Add(-time.Second)

	p := NewProcessor()
	result := p.Tick(s)
	_ = result
}

func TestInteractionRules_DoNotPanicOnExtremeState(t *testing.T) {
	// Push all variables to extremes and verify no panics or NaN.
	s := newTestState()
	s.BodyTemp = 34
	s.HeartRate = 180
	s.BloodPressure = 80
	s.RespiratoryRate = 38
	s.Hunger = 1
	s.Thirst = 1
	s.Fatigue = 0.95
	s.Pain = 0.9
	s.MuscleTension = 0.9
	s.BloodSugar = 55
	s.Cortisol = 0.9
	s.Adrenaline = 0.8
	s.Serotonin = 0.1
	s.Dopamine = 0.1
	s.ImmuneResponse = 0.7
	s.SpO2 = 80
	s.Hydration = 0.2
	s.Glycogen = 0.05
	s.Endorphins = 0.5
	s.LastUpdate = time.Now().Add(-5 * time.Second)

	p := NewProcessor()
	result := p.Tick(s)

	// Verify all variables are still in valid ranges.
	for v := VarBodyTemp; v <= VarEndorphins; v++ {
		val := s.Get(v)
		r := variableRanges[v]
		if val < r[0] || val > r[1] {
			t.Errorf("%s out of range: %f (expected %f-%f)", v, val, r[0], r[1])
		}
	}

	_ = result
}

func TestProcessStimulus_ColdExposure(t *testing.T) {
	s := newTestState()
	initialTemp := s.BodyTemp

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Thermal,
		Intensity: 0.1, // very cold (0.1 < 0.5)
		Timestamp: time.Now(),
	})

	if len(changes) == 0 {
		t.Fatal("expected state changes from cold stimulus")
	}

	if s.BodyTemp >= initialTemp {
		t.Errorf("expected body temp to drop from cold stimulus, got %f (was %f)", s.BodyTemp, initialTemp)
	}
}

func TestProcessStimulus_PainEvent(t *testing.T) {
	s := newTestState()

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Pain,
		Intensity: 0.6,
		Timestamp: time.Now(),
	})

	if len(changes) < 2 {
		t.Fatalf("expected at least 2 changes from pain stimulus, got %d", len(changes))
	}

	if s.Pain <= 0 {
		t.Error("expected pain to increase from pain stimulus")
	}
	if s.Adrenaline <= 0 {
		t.Error("expected adrenaline to increase from pain stimulus")
	}
}

func TestCircadianCortisol_PeakInMorning(t *testing.T) {
	morning := ComputeCircadian(7)  // 7 AM
	night := ComputeCircadian(23)    // 11 PM

	if morning.CortisolBaseline <= night.CortisolBaseline {
		t.Errorf("expected morning cortisol (%f) > night cortisol (%f)",
			morning.CortisolBaseline, night.CortisolBaseline)
	}
}

func TestCircadianAlertness_AfternoonDip(t *testing.T) {
	morning := circadianAlertness(10)    // 10 AM — morning peak
	afternoon := circadianAlertness(14.5) // 2:30 PM — afternoon dip
	evening := circadianAlertness(19)     // 7 PM — evening peak
	night := circadianAlertness(4)        // 4 AM — deepest low

	if afternoon >= morning {
		t.Errorf("expected afternoon dip: afternoon alertness (%f) should be < morning (%f)",
			afternoon, morning)
	}
	if afternoon >= evening {
		t.Errorf("expected afternoon dip: afternoon alertness (%f) should be < evening (%f)",
			afternoon, evening)
	}
	if night >= afternoon {
		t.Errorf("expected nighttime low (%f) < afternoon dip (%f)", night, afternoon)
	}
}

func TestThresholds_HypothermiaProgression(t *testing.T) {
	tests := []struct {
		temp      float64
		condition CriticalState
	}{
		{36.0, Normal},
		{34.5, Impaired},
		{32.0, Critical},
		{29.0, Unconscious},
		{27.0, Lethal},
	}

	for _, tt := range tests {
		s := NewDefaultState()
		s.BodyTemp = tt.temp

		results := EvaluateThresholds(&s)

		if tt.condition == Normal {
			// Should find no thermoregulation threshold
			for _, r := range results {
				if r.System == "thermoregulation" {
					t.Errorf("temp %f: expected no threshold, got %s", tt.temp, r.Condition)
				}
			}
			continue
		}

		found := false
		for _, r := range results {
			if r.System == "thermoregulation" && r.Condition == tt.condition {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("temp %f: expected %s threshold, not found in %v", tt.temp, tt.condition, results)
		}
	}
}

func TestSignificantChanges_FiltersSmallDeltas(t *testing.T) {
	changes := []StateChange{
		{VarHeartRate, 0.5, "test"},      // < 2 bpm, not significant
		{VarHeartRate, 5.0, "test"},      // >= 2 bpm, significant
		{VarCortisol, 0.005, "test"},     // < 0.01, not significant
		{VarCortisol, 0.05, "test"},      // >= 0.01, significant
	}

	significant := SignificantChanges(changes)
	if len(significant) != 2 {
		t.Errorf("expected 2 significant changes, got %d", len(significant))
	}
}

func TestGlycogenBuffersBloodSugar(t *testing.T) {
	s := newTestState()
	s.BloodSugar = 75 // below 80, glycogen should buffer
	s.Glycogen = 0.5
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Blood sugar should have risen from glycogen release.
	if s.BloodSugar <= 75 {
		t.Errorf("expected glycogen to buffer blood sugar above 75, got %f", s.BloodSugar)
	}
	// Glycogen should have decreased.
	if s.Glycogen >= 0.5 {
		t.Errorf("expected glycogen to decrease from buffering, got %f", s.Glycogen)
	}
}

func TestTick_ZeroDt_NoChanges(t *testing.T) {
	s := newTestState()
	// Set LastUpdate to a future time so dt will be <= 0.
	s.LastUpdate = time.Now().Add(time.Hour)

	p := NewProcessor()
	result := p.Tick(s)

	if len(result.Changes) != 0 {
		t.Errorf("expected no changes with zero/negative dt, got %d", len(result.Changes))
	}
}

func TestTick_LargeDt_CappedAt300Seconds(t *testing.T) {
	s := newTestState()
	initialFatigue := s.Fatigue
	s.LastUpdate = time.Now().Add(-1 * time.Hour) // 3600s, should be capped to 300

	p := NewProcessor()
	p.Tick(s)

	// Fatigue accumulates at 0.05/hr = 0.05/3600 per second.
	// With 300s cap: 0.05/3600 * 300 ≈ 0.00417
	// With full 3600s: 0.05/3600 * 3600 = 0.05
	// Fatigue should be much closer to the capped value.
	fatigueDelta := s.Fatigue - initialFatigue
	if fatigueDelta > 0.01 {
		t.Errorf("fatigue increased by %f, suggesting dt was not capped at 300s", fatigueDelta)
	}
}

func TestProcessStimulus_HeatExposure(t *testing.T) {
	s := newTestState()
	initialTemp := s.BodyTemp

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Thermal,
		Intensity: 0.9, // very hot (0.9 > 0.5)
		Timestamp: time.Now(),
	})

	if len(changes) == 0 {
		t.Fatal("expected state changes from heat stimulus")
	}

	if s.BodyTemp <= initialTemp {
		t.Errorf("expected body temp to rise from heat stimulus, got %f (was %f)", s.BodyTemp, initialTemp)
	}
}

func TestProcessStimulus_AuditoryStartle(t *testing.T) {
	s := newTestState()

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Auditory,
		Intensity: 0.9, // loud startling sound
		Timestamp: time.Now(),
	})

	if len(changes) == 0 {
		t.Fatal("expected state changes from loud auditory stimulus")
	}

	if s.Adrenaline <= 0 {
		t.Error("expected adrenaline increase from startle response")
	}
	if s.HeartRate <= 70 {
		t.Errorf("expected heart rate increase from startle, got %f", s.HeartRate)
	}
}

func TestProcessStimulus_AuditoryQuiet_NoStartle(t *testing.T) {
	s := newTestState()
	initialAdrenaline := s.Adrenaline
	initialHR := s.HeartRate

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Auditory,
		Intensity: 0.3, // quiet sound
		Timestamp: time.Now(),
	})

	if len(changes) != 0 {
		t.Errorf("expected no changes from quiet sound, got %d", len(changes))
	}
	if s.Adrenaline != initialAdrenaline {
		t.Errorf("quiet sound should not change adrenaline")
	}
	if s.HeartRate != initialHR {
		t.Errorf("quiet sound should not change heart rate")
	}
}

func TestProcessStimulus_VisualThreat(t *testing.T) {
	s := newTestState()

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Visual,
		Intensity: 0.95, // threatening visual
		Timestamp: time.Now(),
	})

	if len(changes) == 0 {
		t.Fatal("expected state changes from visual threat")
	}

	if s.Adrenaline <= 0 {
		t.Error("expected adrenaline increase from visual threat")
	}
}

func TestProcessStimulus_VisualNonThreatening_NoEffect(t *testing.T) {
	s := newTestState()
	initialAdrenaline := s.Adrenaline

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Visual,
		Intensity: 0.5, // non-threatening
		Timestamp: time.Now(),
	})

	if len(changes) != 0 {
		t.Errorf("expected no changes from non-threatening visual, got %d", len(changes))
	}
	if s.Adrenaline != initialAdrenaline {
		t.Error("non-threatening visual should not change adrenaline")
	}
}

func TestProcessStimulus_UnhandledChannel_NoEffect(t *testing.T) {
	s := newTestState()

	p := NewProcessor()
	changes := p.ProcessStimulus(s, sense.Event{
		Channel:   sense.Olfactory, // not handled yet
		Intensity: 0.8,
		Timestamp: time.Now(),
	})

	if len(changes) != 0 {
		t.Errorf("expected no changes from unhandled olfactory channel, got %d", len(changes))
	}
}

func TestTick_NaturalHydrationDepletion(t *testing.T) {
	s := newTestState()
	initialHydration := s.Hydration
	s.LastUpdate = time.Now().Add(-60 * time.Second) // 1 minute

	p := NewProcessor()
	p.Tick(s)

	// Hydration depletes at 0.001/min.
	if s.Hydration >= initialHydration {
		t.Errorf("expected hydration to deplete over time, got %f (was %f)", s.Hydration, initialHydration)
	}
}

func TestTick_NaturalFatigueAccumulation(t *testing.T) {
	s := newTestState()
	initialFatigue := s.Fatigue
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Fatigue <= initialFatigue {
		t.Errorf("expected fatigue to accumulate over time, got %f (was %f)", s.Fatigue, initialFatigue)
	}
}

func TestTick_SpO2Recovery_AdequateRespiration(t *testing.T) {
	s := newTestState()
	s.SpO2 = 90 // below normal
	s.RespiratoryRate = 15 // adequate
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.SpO2 <= 90 {
		t.Errorf("expected SpO2 to recover with adequate respiration, got %f", s.SpO2)
	}
}

func TestTick_CortisolLoadImmuneIntegration(t *testing.T) {
	s := newTestState()
	s.Cortisol = 0.7
	s.CortisolLoad = 50 // significant accumulated load
	s.ImmuneResponse = 0.5
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	initialImmune := s.ImmuneResponse

	p := NewProcessor()
	p.Tick(s)

	// With cortisol load of 50, suppression factor = 1/(1+50*0.1) = 1/6 ≈ 0.17
	// Immune should decrease.
	if s.ImmuneResponse >= initialImmune {
		t.Errorf("expected immune suppression from cortisol load, got %f (was %f)",
			s.ImmuneResponse, initialImmune)
	}
}

func TestTick_CortisolLoad_NoAccumulationBelowThreshold(t *testing.T) {
	s := newTestState()
	s.Cortisol = 0.2 // below 0.3 threshold
	s.CortisolLoad = 0
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.CortisolLoad > 0 {
		t.Errorf("expected no cortisol load accumulation below 0.3, got %f", s.CortisolLoad)
	}
}

func TestTick_CircadianPhaseAdvances(t *testing.T) {
	s := newTestState()
	initialPhase := s.CircadianPhase // 8.0
	s.LastUpdate = time.Now().Add(-3600 * time.Second) // 1 hour, but capped to 300s

	p := NewProcessor()
	p.Tick(s)

	// Phase should advance by dt/3600 hours. With dt capped at 300s: 300/3600 ≈ 0.083 hours.
	if s.CircadianPhase <= initialPhase {
		t.Errorf("expected circadian phase to advance, got %f (was %f)", s.CircadianPhase, initialPhase)
	}
}

func TestTick_PainStressAmplificationLoop(t *testing.T) {
	// Pain → cortisol → sustained cortisol depletes serotonin →
	// low serotonin increases pain sensitivity.
	// This tests the pain-stress chain from the biologist advisory.
	s := newTestState()
	s.Pain = 0.6
	s.Cortisol = 0.1
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Pain > 0.3 should trigger cortisol increase.
	if s.Cortisol <= 0.1 {
		t.Errorf("expected cortisol increase from pain, got %f", s.Cortisol)
	}
}

func TestTick_DehydrationTachycardia(t *testing.T) {
	s := newTestState()
	s.Hydration = 0.4 // below 0.6 threshold
	initialHR := s.HeartRate
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Dehydration should cause compensatory tachycardia.
	if s.HeartRate <= initialHR {
		t.Errorf("expected HR increase from dehydration, got %f (was %f)", s.HeartRate, initialHR)
	}
}

func TestTick_HypoxiaResponse(t *testing.T) {
	s := newTestState()
	s.SpO2 = 88 // below 90 threshold
	s.RespiratoryRate = 8 // inadequate ventilation (keeps SpO2 dropping)
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Should trigger hypoxia adrenaline response.
	if s.Adrenaline <= 0 {
		t.Errorf("expected adrenaline increase from hypoxia, got %f", s.Adrenaline)
	}
}

func TestSignificantChanges_NegativeDeltas(t *testing.T) {
	changes := []StateChange{
		{VarHeartRate, -0.5, "test"},  // not significant
		{VarHeartRate, -5.0, "test"},  // significant (abs >= 2)
		{VarBodyTemp, -0.05, "test"},  // not significant
		{VarBodyTemp, -0.2, "test"},   // significant (abs >= 0.1)
	}

	significant := SignificantChanges(changes)
	if len(significant) != 2 {
		t.Errorf("expected 2 significant negative changes, got %d", len(significant))
	}
}

func TestSignificantChanges_AllVariableTypes(t *testing.T) {
	// Verify significance thresholds for all explicitly handled variable types.
	tests := []struct {
		variable  Variable
		delta     float64
		wantSignificant bool
	}{
		{VarHeartRate, 1.5, false},
		{VarHeartRate, 2.5, true},
		{VarBloodPressure, 2.0, false},
		{VarBloodPressure, 4.0, true},
		{VarBodyTemp, 0.05, false},
		{VarBodyTemp, 0.15, true},
		{VarRespiratoryRate, 0.5, false},
		{VarRespiratoryRate, 1.5, true},
		{VarBloodSugar, 1.0, false},
		{VarBloodSugar, 3.0, true},
		{VarSpO2, 0.3, false},
		{VarSpO2, 0.7, true},
		{VarCortisol, 0.005, false},  // ratio variable, threshold 0.01
		{VarCortisol, 0.02, true},
		{VarAdrenaline, 0.005, false},
		{VarAdrenaline, 0.02, true},
	}

	for _, tt := range tests {
		changes := []StateChange{{tt.variable, tt.delta, "test"}}
		significant := SignificantChanges(changes)
		got := len(significant) > 0
		if got != tt.wantSignificant {
			t.Errorf("SignificantChanges(%s, delta=%f): got significant=%v, want %v",
				tt.variable, tt.delta, got, tt.wantSignificant)
		}
	}
}

func TestSignificantChanges_EmptyInput(t *testing.T) {
	significant := SignificantChanges(nil)
	if significant != nil {
		t.Errorf("expected nil for nil input, got %v", significant)
	}

	significant = SignificantChanges([]StateChange{})
	if significant != nil {
		t.Errorf("expected nil for empty input, got %v", significant)
	}
}

func TestTick_EndorphinDecay(t *testing.T) {
	s := newTestState()
	s.Endorphins = 0.6 // well above baseline 0.1
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Endorphins >= 0.6 {
		t.Errorf("expected endorphins to decay toward baseline, got %f", s.Endorphins)
	}
	if s.Endorphins < 0.1 {
		t.Errorf("endorphins decayed below baseline, got %f", s.Endorphins)
	}
}

func TestTick_DopamineDecayTowardBaseline(t *testing.T) {
	s := newTestState()
	s.Dopamine = 0.8 // spike above baseline 0.3
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Dopamine >= 0.8 {
		t.Errorf("expected dopamine to decay toward 0.3, got %f", s.Dopamine)
	}
}

func TestTick_CortisolDecayTowardBaseline(t *testing.T) {
	s := newTestState()
	s.Cortisol = 0.8 // elevated
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Cortisol should decay toward 0.1 baseline.
	if s.Cortisol >= 0.8 {
		t.Errorf("expected cortisol to decay, got %f", s.Cortisol)
	}
}

func TestTick_BloodSugarInsulinHomeostasis(t *testing.T) {
	s := newTestState()
	s.BloodSugar = 150 // elevated above 95
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.BloodSugar >= 150 {
		t.Errorf("expected blood sugar to decrease from insulin homeostasis, got %f", s.BloodSugar)
	}
}

func TestTick_GlycogenNaturalDepletion(t *testing.T) {
	s := newTestState()
	initialGlycogen := s.Glycogen
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Glycogen >= initialGlycogen {
		t.Errorf("expected glycogen natural depletion, got %f (was %f)", s.Glycogen, initialGlycogen)
	}
}

func TestTick_HungerSuppressedByNormalBloodSugar(t *testing.T) {
	s := newTestState()
	s.Hunger = 0.5
	s.BloodSugar = 120 // above 110, hunger suppression rule fires
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Hunger >= 0.5 {
		t.Errorf("expected hunger to decrease with BS>110, got %f", s.Hunger)
	}
}

func TestTick_ThirstDerivedFromHydration(t *testing.T) {
	s := newTestState()
	s.Hydration = 0.5 // below 0.7 threshold
	s.Thirst = 0
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Thirst <= 0 {
		t.Error("expected thirst to increase when hydration < 0.7")
	}
}

func TestTick_ThirstSuppressedWhenWellHydrated(t *testing.T) {
	s := newTestState()
	s.Hydration = 0.9 // above 0.85 threshold
	s.Thirst = 0.3
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.Thirst >= 0.3 {
		t.Errorf("expected thirst to decrease when well hydrated, got %f", s.Thirst)
	}
}

func TestTick_AdrenalineFightOrFlight(t *testing.T) {
	s := newTestState()
	s.Adrenaline = 0.6 // well above 0.2 threshold
	initialHR := s.HeartRate
	initialBP := s.BloodPressure
	initialRR := s.RespiratoryRate
	s.LastUpdate = time.Now().Add(-5 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Adrenaline > 0.2 should cause HR, BP, and RR to increase.
	if s.HeartRate <= initialHR {
		t.Errorf("expected HR increase from adrenaline, got %f (was %f)", s.HeartRate, initialHR)
	}
	if s.BloodPressure <= initialBP {
		t.Errorf("expected BP increase from adrenaline, got %f (was %f)", s.BloodPressure, initialBP)
	}
	if s.RespiratoryRate <= initialRR {
		t.Errorf("expected RR increase from adrenaline, got %f (was %f)", s.RespiratoryRate, initialRR)
	}
}

func TestTick_PainCascade_TachycardiaGuardingEndorphins(t *testing.T) {
	s := newTestState()
	s.Pain = 0.85 // above 0.8 for endorphin release, above 0.3 for tachycardia/guarding
	initialHR := s.HeartRate
	initialTension := s.MuscleTension
	initialEndorphins := s.Endorphins
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.HeartRate <= initialHR {
		t.Errorf("expected pain-driven tachycardia, got HR %f (was %f)", s.HeartRate, initialHR)
	}
	if s.MuscleTension <= initialTension {
		t.Errorf("expected pain guarding (muscle tension increase), got %f (was %f)", s.MuscleTension, initialTension)
	}
	if s.Endorphins <= initialEndorphins {
		t.Errorf("expected endorphin release at pain > 0.8, got %f (was %f)", s.Endorphins, initialEndorphins)
	}
}

func TestTick_EndorphinAnalgesia(t *testing.T) {
	// Test that endorphins reduce pain beyond normal decay.
	// Compare pain decay with and without endorphins.
	sWithEndorphins := newTestState()
	sWithEndorphins.Pain = 0.5
	sWithEndorphins.Endorphins = 0.5 // high endorphins
	sWithEndorphins.LastUpdate = time.Now().Add(-30 * time.Second)

	sWithout := newTestState()
	sWithout.Pain = 0.5
	sWithout.Endorphins = 0.1 // baseline endorphins
	sWithout.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(sWithEndorphins)
	p.Tick(sWithout)

	// Pain should be lower with endorphins (faster decay + analgesia rule).
	if sWithEndorphins.Pain >= sWithout.Pain {
		t.Errorf("expected endorphins to reduce pain more: with=%f, without=%f",
			sWithEndorphins.Pain, sWithout.Pain)
	}
}

func TestTick_ImmuneFeverCascade(t *testing.T) {
	s := newTestState()
	s.ImmuneResponse = 0.6 // above 0.4 for fever, above 0.5 for fatigue
	initialTemp := s.BodyTemp
	initialFatigue := s.Fatigue
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.BodyTemp <= initialTemp {
		t.Errorf("expected fever from immune response, got temp %f (was %f)", s.BodyTemp, initialTemp)
	}
	// Fatigue increase comes from immune > 0.5 rule + natural accumulation.
	if s.Fatigue <= initialFatigue {
		t.Errorf("expected fatigue from immune response, got %f (was %f)", s.Fatigue, initialFatigue)
	}
}

func TestTick_CortisolGluconeogenesis(t *testing.T) {
	s := newTestState()
	s.Cortisol = 0.7 // above 0.5 threshold
	initialBS := s.BloodSugar
	s.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// Cortisol > 0.5 should raise blood sugar via gluconeogenesis.
	if s.BloodSugar <= initialBS {
		t.Errorf("expected blood sugar increase from cortisol gluconeogenesis, got %f (was %f)",
			s.BloodSugar, initialBS)
	}
}

func TestTick_FatigueImmuneSuppression(t *testing.T) {
	s := newTestState()
	s.Fatigue = 0.85 // above 0.8 threshold
	s.ImmuneResponse = 0.3
	initialImmune := s.ImmuneResponse
	s.LastUpdate = time.Now().Add(-30 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.ImmuneResponse >= initialImmune {
		t.Errorf("expected immune suppression from high fatigue, got %f (was %f)",
			s.ImmuneResponse, initialImmune)
	}
}

func TestTick_PainDecayFasterWithEndorphins(t *testing.T) {
	// The decay function uses half-life 1800s normally, 900s with endorphins > 0.2.
	// This tests that the endorphin-accelerated decay path works.
	sHigh := newTestState()
	sHigh.Pain = 0.6
	sHigh.Endorphins = 0.4 // above 0.2 threshold for accelerated decay
	sHigh.LastUpdate = time.Now().Add(-60 * time.Second)

	sLow := newTestState()
	sLow.Pain = 0.6
	sLow.Endorphins = 0.05 // below threshold
	sLow.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(sHigh)
	p.Tick(sLow)

	// With endorphins, pain should decay more (lower value = more decay).
	if sHigh.Pain >= sLow.Pain {
		t.Errorf("expected faster pain decay with endorphins: high_endorphin_pain=%f, low_endorphin_pain=%f",
			sHigh.Pain, sLow.Pain)
	}
}

func TestTick_HeartRateDecayTowardBaseline(t *testing.T) {
	s := newTestState()
	s.HeartRate = 130 // elevated
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	// HR should decay toward baseline 70.
	if s.HeartRate >= 130 {
		t.Errorf("expected HR to decay toward 70, got %f", s.HeartRate)
	}
	// Should not overshoot below baseline.
	if s.HeartRate < 70 {
		t.Errorf("HR decayed below baseline, got %f", s.HeartRate)
	}
}

func TestTick_MuscleTensionDecay(t *testing.T) {
	s := newTestState()
	s.MuscleTension = 0.5
	s.LastUpdate = time.Now().Add(-60 * time.Second)

	p := NewProcessor()
	p.Tick(s)

	if s.MuscleTension >= 0.5 {
		t.Errorf("expected muscle tension to decay toward 0, got %f", s.MuscleTension)
	}
}

func TestTick_SpO2RecoveryFasterWithHyperventilation(t *testing.T) {
	sNormal := newTestState()
	sNormal.SpO2 = 92
	sNormal.RespiratoryRate = 15 // normal
	sNormal.LastUpdate = time.Now().Add(-10 * time.Second)

	sHyper := newTestState()
	sHyper.SpO2 = 92
	sHyper.RespiratoryRate = 30 // hyperventilation > 25
	sHyper.LastUpdate = time.Now().Add(-10 * time.Second)

	p := NewProcessor()
	p.Tick(sNormal)
	p.Tick(sHyper)

	// Hyperventilation should recover SpO2 faster.
	if sHyper.SpO2 <= sNormal.SpO2 {
		t.Errorf("expected faster SpO2 recovery with hyperventilation: hyper=%f, normal=%f",
			sHyper.SpO2, sNormal.SpO2)
	}
}
