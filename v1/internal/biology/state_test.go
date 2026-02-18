package biology

import (
	"math"
	"testing"
)

func TestNewDefaultState_BaselineValues(t *testing.T) {
	s := NewDefaultState()

	checks := []struct {
		name string
		got  float64
		want float64
	}{
		{"BodyTemp", s.BodyTemp, 36.6},
		{"HeartRate", s.HeartRate, 70},
		{"BloodPressure", s.BloodPressure, 120},
		{"RespiratoryRate", s.RespiratoryRate, 15},
		{"Hunger", s.Hunger, 0.0},
		{"Thirst", s.Thirst, 0.0},
		{"Fatigue", s.Fatigue, 0.0},
		{"Pain", s.Pain, 0.0},
		{"MuscleTension", s.MuscleTension, 0.0},
		{"BloodSugar", s.BloodSugar, 90},
		{"Cortisol", s.Cortisol, 0.1},
		{"Adrenaline", s.Adrenaline, 0.0},
		{"Serotonin", s.Serotonin, 0.5},
		{"Dopamine", s.Dopamine, 0.3},
		{"ImmuneResponse", s.ImmuneResponse, 0.1},
		{"CircadianPhase", s.CircadianPhase, 8.0},
		{"SpO2", s.SpO2, 98},
		{"Hydration", s.Hydration, 0.8},
		{"Glycogen", s.Glycogen, 0.7},
		{"Endorphins", s.Endorphins, 0.1},
		{"CortisolLoad", s.CortisolLoad, 0.0},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("NewDefaultState().%s = %f, want %f", c.name, c.got, c.want)
		}
	}

	if s.LastUpdate.IsZero() {
		t.Error("NewDefaultState().LastUpdate should not be zero")
	}
}

func TestState_GetSet_RoundTrip_AllVariables(t *testing.T) {
	s := NewDefaultState()
	testVal := 42.42

	for v := VarBodyTemp; v <= VarEndorphins; v++ {
		s.Set(v, testVal)
		got := s.Get(v)
		if got != testVal {
			t.Errorf("Get(%s) = %f after Set(%s, %f)", v, got, v, testVal)
		}
	}
}

func TestState_Get_UnknownVariable(t *testing.T) {
	s := NewDefaultState()
	got := s.Get(Variable(999))
	if got != 0 {
		t.Errorf("Get(999) = %f, want 0", got)
	}
}

func TestState_Set_UnknownVariable_NoPanic(t *testing.T) {
	s := NewDefaultState()
	// Should not panic.
	s.Set(Variable(999), 42)
}

func TestVariable_String_AllVariables(t *testing.T) {
	tests := []struct {
		v    Variable
		want string
	}{
		{VarBodyTemp, "body_temp"},
		{VarHeartRate, "heart_rate"},
		{VarBloodPressure, "blood_pressure"},
		{VarRespiratoryRate, "respiratory_rate"},
		{VarHunger, "hunger"},
		{VarThirst, "thirst"},
		{VarFatigue, "fatigue"},
		{VarPain, "pain"},
		{VarMuscleTension, "muscle_tension"},
		{VarBloodSugar, "blood_sugar"},
		{VarCortisol, "cortisol"},
		{VarAdrenaline, "adrenaline"},
		{VarSerotonin, "serotonin"},
		{VarDopamine, "dopamine"},
		{VarImmuneResponse, "immune_response"},
		{VarCircadianPhase, "circadian_phase"},
		{VarSpO2, "spo2"},
		{VarHydration, "hydration"},
		{VarGlycogen, "glycogen"},
		{VarEndorphins, "endorphins"},
	}

	for _, tt := range tests {
		got := tt.v.String()
		if got != tt.want {
			t.Errorf("Variable(%d).String() = %q, want %q", tt.v, got, tt.want)
		}
	}
}

func TestVariable_String_OutOfRange(t *testing.T) {
	got := Variable(99).String()
	if got != "unknown" {
		t.Errorf("Variable(99).String() = %q, want %q", got, "unknown")
	}
}

func TestVariable_EnumCount(t *testing.T) {
	if len(variableNames) != int(VarEndorphins)+1 {
		t.Errorf("variableNames has %d entries, expected %d", len(variableNames), int(VarEndorphins)+1)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		val, min, max, want float64
	}{
		{5, 0, 10, 5},     // within range
		{-1, 0, 10, 0},    // below min
		{15, 0, 10, 10},   // above max
		{0, 0, 10, 0},     // at min
		{10, 0, 10, 10},   // at max
		{0.5, 0, 1, 0.5},  // ratio within
	}

	for _, tt := range tests {
		got := Clamp(tt.val, tt.min, tt.max)
		if got != tt.want {
			t.Errorf("Clamp(%f, %f, %f) = %f, want %f", tt.val, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestClampVariable_AllVariables(t *testing.T) {
	// Verify every variable has a defined range and clamps correctly.
	for v := VarBodyTemp; v <= VarEndorphins; v++ {
		r, ok := variableRanges[v]
		if !ok {
			t.Errorf("variable %s has no defined range in variableRanges", v)
			continue
		}

		// Below min
		got := ClampVariable(v, r[0]-100)
		if got != r[0] {
			t.Errorf("ClampVariable(%s, below_min) = %f, want %f", v, got, r[0])
		}

		// Above max
		got = ClampVariable(v, r[1]+100)
		if got != r[1] {
			t.Errorf("ClampVariable(%s, above_max) = %f, want %f", v, got, r[1])
		}

		// Within range (midpoint)
		mid := (r[0] + r[1]) / 2
		got = ClampVariable(v, mid)
		if got != mid {
			t.Errorf("ClampVariable(%s, midpoint) = %f, want %f", v, got, mid)
		}
	}
}

func TestClampVariable_UnknownVariable(t *testing.T) {
	// Unknown variable should return value unchanged.
	got := ClampVariable(Variable(999), 42.0)
	if got != 42.0 {
		t.Errorf("ClampVariable(unknown, 42) = %f, want 42", got)
	}
}

func TestVariableRanges_Completeness(t *testing.T) {
	// Every variable from VarBodyTemp to VarEndorphins must have a range.
	for v := VarBodyTemp; v <= VarEndorphins; v++ {
		if _, ok := variableRanges[v]; !ok {
			t.Errorf("variableRanges missing entry for %s", v)
		}
	}
}

func TestVariableRanges_MinLessThanMax(t *testing.T) {
	for v, r := range variableRanges {
		if r[0] >= r[1] {
			t.Errorf("variableRanges[%s]: min (%f) >= max (%f)", Variable(v), r[0], r[1])
		}
	}
}

func TestNewDefaultState_AllValuesWithinRanges(t *testing.T) {
	s := NewDefaultState()
	for v := VarBodyTemp; v <= VarEndorphins; v++ {
		val := s.Get(v)
		r := variableRanges[v]
		if val < r[0] || val > r[1] {
			t.Errorf("default %s = %f, outside range [%f, %f]", v, val, r[0], r[1])
		}
	}
}

func TestState_Get_MatchesStructFields(t *testing.T) {
	// Verify that Get returns the same value as direct field access for each variable.
	s := State{
		BodyTemp:        36.6,
		HeartRate:       72,
		BloodPressure:   118,
		RespiratoryRate: 16,
		Hunger:          0.1,
		Thirst:          0.2,
		Fatigue:         0.3,
		Pain:            0.4,
		MuscleTension:   0.5,
		BloodSugar:      95,
		Cortisol:        0.15,
		Adrenaline:      0.05,
		Serotonin:       0.55,
		Dopamine:        0.35,
		ImmuneResponse:  0.12,
		CircadianPhase:  10.5,
		SpO2:            97,
		Hydration:       0.75,
		Glycogen:        0.65,
		Endorphins:      0.08,
	}

	fieldValues := map[Variable]float64{
		VarBodyTemp:        36.6,
		VarHeartRate:       72,
		VarBloodPressure:   118,
		VarRespiratoryRate: 16,
		VarHunger:          0.1,
		VarThirst:          0.2,
		VarFatigue:         0.3,
		VarPain:            0.4,
		VarMuscleTension:   0.5,
		VarBloodSugar:      95,
		VarCortisol:        0.15,
		VarAdrenaline:      0.05,
		VarSerotonin:       0.55,
		VarDopamine:        0.35,
		VarImmuneResponse:  0.12,
		VarCircadianPhase:  10.5,
		VarSpO2:            97,
		VarHydration:       0.75,
		VarGlycogen:        0.65,
		VarEndorphins:      0.08,
	}

	for v, want := range fieldValues {
		got := s.Get(v)
		if math.Abs(got-want) > 1e-10 {
			t.Errorf("Get(%s) = %f, want %f", v, got, want)
		}
	}
}
