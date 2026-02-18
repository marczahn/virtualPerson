package biology

import "testing"

func TestCriticalState_String_AllValues(t *testing.T) {
	tests := []struct {
		cs   CriticalState
		want string
	}{
		{Normal, "normal"},
		{Impaired, "impaired"},
		{Critical, "critical"},
		{Unconscious, "unconscious"},
		{Lethal, "lethal"},
	}

	for _, tt := range tests {
		got := tt.cs.String()
		if got != tt.want {
			t.Errorf("CriticalState(%d).String() = %q, want %q", tt.cs, got, tt.want)
		}
	}
}

func TestCriticalState_String_OutOfRange(t *testing.T) {
	got := CriticalState(99).String()
	if got != "unknown" {
		t.Errorf("CriticalState(99).String() = %q, want %q", got, "unknown")
	}
}

func TestThresholds_HyperthermiaProgression(t *testing.T) {
	tests := []struct {
		temp      float64
		condition CriticalState
		system    string
	}{
		{37.0, Normal, ""},
		{40.5, Impaired, "thermoregulation"},
		{41.8, Critical, "thermoregulation"},
		{42.5, Lethal, "thermoregulation"},
	}

	for _, tt := range tests {
		s := NewDefaultState()
		s.BodyTemp = tt.temp

		results := EvaluateThresholds(&s)

		if tt.condition == Normal {
			for _, r := range results {
				if r.System == "thermoregulation" {
					t.Errorf("temp %f: expected no thermoregulation threshold, got %s", tt.temp, r.Condition)
				}
			}
			continue
		}

		found := false
		for _, r := range results {
			if r.System == tt.system && r.Condition == tt.condition {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("temp %f: expected %s/%s threshold, not found", tt.temp, tt.system, tt.condition)
		}
	}
}

func TestThresholds_HypoglycemiaProgression(t *testing.T) {
	tests := []struct {
		bs        float64
		condition CriticalState
	}{
		{90, Normal},
		{52, Impaired},
		{42, Critical},
		{32, Unconscious},
		{22, Lethal},
	}

	for _, tt := range tests {
		s := NewDefaultState()
		s.BloodSugar = tt.bs

		results := EvaluateThresholds(&s)

		if tt.condition == Normal {
			for _, r := range results {
				if r.System == "glycemic" {
					t.Errorf("BS %f: expected no glycemic threshold, got %s", tt.bs, r.Condition)
				}
			}
			continue
		}

		found := false
		for _, r := range results {
			if r.System == "glycemic" && r.Condition == tt.condition {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BS %f: expected glycemic/%s threshold, not found", tt.bs, tt.condition)
		}
	}
}

func TestThresholds_HyperglycemiaThresholds(t *testing.T) {
	tests := []struct {
		bs        float64
		condition CriticalState
	}{
		{150, Normal},
		{190, Impaired},
		{310, Critical},
	}

	for _, tt := range tests {
		s := NewDefaultState()
		s.BloodSugar = tt.bs

		results := EvaluateThresholds(&s)

		if tt.condition == Normal {
			for _, r := range results {
				if r.System == "glycemic" {
					t.Errorf("BS %f: expected no glycemic threshold, got %s", tt.bs, r.Condition)
				}
			}
			continue
		}

		found := false
		for _, r := range results {
			if r.System == "glycemic" && r.Condition == tt.condition {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BS %f: expected glycemic/%s threshold, not found", tt.bs, tt.condition)
		}
	}
}

func TestThresholds_SpO2Cascade(t *testing.T) {
	tests := []struct {
		spo2      float64
		condition CriticalState
	}{
		{97, Normal},
		{88, Impaired},
		{82, Critical},
		{72, Unconscious},
		{55, Lethal},
	}

	for _, tt := range tests {
		s := NewDefaultState()
		s.SpO2 = tt.spo2

		results := EvaluateThresholds(&s)

		if tt.condition == Normal {
			for _, r := range results {
				if r.System == "respiratory" {
					t.Errorf("SpO2 %f: expected no respiratory threshold, got %s", tt.spo2, r.Condition)
				}
			}
			continue
		}

		found := false
		for _, r := range results {
			if r.System == "respiratory" && r.Condition == tt.condition {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SpO2 %f: expected respiratory/%s threshold, not found", tt.spo2, tt.condition)
		}
	}
}

func TestThresholds_BloodPressureShock(t *testing.T) {
	s := NewDefaultState()
	s.BloodPressure = 65

	results := EvaluateThresholds(&s)

	found := false
	for _, r := range results {
		if r.System == "cardiovascular" && r.Condition == Critical {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected cardiovascular/critical threshold at BP 65")
	}
}

func TestThresholds_NormalState_NoThresholds(t *testing.T) {
	s := NewDefaultState()
	results := EvaluateThresholds(&s)

	if len(results) != 0 {
		t.Errorf("expected no thresholds on default state, got %d: %v", len(results), results)
	}
}

func TestThresholds_MultipleSimultaneous(t *testing.T) {
	s := NewDefaultState()
	s.BodyTemp = 34.0  // hypothermia impaired
	s.BloodSugar = 50  // hypoglycemia impaired
	s.SpO2 = 88        // respiratory impaired

	results := EvaluateThresholds(&s)

	systems := make(map[string]bool)
	for _, r := range results {
		systems[r.System] = true
	}

	if !systems["thermoregulation"] {
		t.Error("expected thermoregulation threshold")
	}
	if !systems["glycemic"] {
		t.Error("expected glycemic threshold")
	}
	if !systems["respiratory"] {
		t.Error("expected respiratory threshold")
	}
}

func TestCortisolLoadImmuneSuppressionFactor(t *testing.T) {
	tests := []struct {
		load     float64
		wantMin  float64
		wantMax  float64
	}{
		{0.0, 0.99, 1.01},    // no load → factor ~1.0
		{10.0, 0.4, 0.6},     // moderate load → around 0.5
		{100.0, 0.05, 0.15},  // heavy load → approaching 0
	}

	for _, tt := range tests {
		got := CortisolLoadImmuneSuppressionFactor(tt.load)
		if got < tt.wantMin || got > tt.wantMax {
			t.Errorf("CortisolLoadImmuneSuppressionFactor(%f) = %f, want [%f, %f]",
				tt.load, got, tt.wantMin, tt.wantMax)
		}
	}
}

func TestIsHypothermiaReversal(t *testing.T) {
	tests := []struct {
		temp float64
		want bool
	}{
		{36.6, false},
		{33.0, false}, // exactly at boundary
		{32.9, true},
		{28.0, true},
	}

	for _, tt := range tests {
		s := State{BodyTemp: tt.temp}
		got := IsHypothermiaReversal(&s)
		if got != tt.want {
			t.Errorf("IsHypothermiaReversal(temp=%f) = %v, want %v", tt.temp, got, tt.want)
		}
	}
}

func TestApplyHypothermiaOverrides_NotInReversal(t *testing.T) {
	s := NewDefaultState()
	s.BodyTemp = 34.0 // cold but not in reversal (<33)

	changes := ApplyHypothermiaOverrides(&s, 1.0)
	if len(changes) != 0 {
		t.Errorf("expected no overrides at 34°C, got %d changes", len(changes))
	}
}

func TestApplyHypothermiaOverrides_InReversal(t *testing.T) {
	s := NewDefaultState()
	s.BodyTemp = 31.0
	s.MuscleTension = 0.7
	s.HeartRate = 90
	s.Adrenaline = 0.3

	changes := ApplyHypothermiaOverrides(&s, 5.0)

	if len(changes) == 0 {
		t.Fatal("expected override changes at 31°C")
	}

	if s.MuscleTension >= 0.7 {
		t.Errorf("expected tension to drop from 0.7, got %f", s.MuscleTension)
	}
	if s.HeartRate >= 90 {
		t.Errorf("expected HR to drop from 90, got %f", s.HeartRate)
	}
	if s.Adrenaline >= 0.3 {
		t.Errorf("expected adrenaline to drop from 0.3, got %f", s.Adrenaline)
	}
}
