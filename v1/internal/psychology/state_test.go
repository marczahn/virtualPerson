package psychology

import "testing"

func TestDistortion_String_AllValues(t *testing.T) {
	tests := []struct {
		d    Distortion
		want string
	}{
		{Catastrophizing, "catastrophizing"},
		{EmotionalReasoning, "emotional_reasoning"},
		{Overgeneralization, "overgeneralization"},
		{MindReading, "mind_reading"},
		{Personalization, "personalization"},
		{AllOrNothing, "all_or_nothing"},
	}

	for _, tt := range tests {
		got := tt.d.String()
		if got != tt.want {
			t.Errorf("Distortion(%d).String() = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestDistortion_String_OutOfRange(t *testing.T) {
	got := Distortion(99).String()
	if got != "unknown" {
		t.Errorf("Distortion(99).String() = %q, want %q", got, "unknown")
	}
}

func TestDistortion_EnumCount(t *testing.T) {
	if len(distortionNames) != int(AllOrNothing)+1 {
		t.Errorf("distortionNames has %d entries, expected %d", len(distortionNames), int(AllOrNothing)+1)
	}
}

func TestCopingStrategy_String_AllValues(t *testing.T) {
	tests := []struct {
		c    CopingStrategy
		want string
	}{
		{ProblemSolving, "problem_solving"},
		{Reappraisal, "reappraisal"},
		{Acceptance, "acceptance"},
		{Distraction, "distraction"},
		{Suppression, "suppression"},
		{Rumination, "rumination"},
		{Denial, "denial"},
	}

	for _, tt := range tests {
		got := tt.c.String()
		if got != tt.want {
			t.Errorf("CopingStrategy(%d).String() = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestCopingStrategy_String_OutOfRange(t *testing.T) {
	got := CopingStrategy(99).String()
	if got != "unknown" {
		t.Errorf("CopingStrategy(99).String() = %q, want %q", got, "unknown")
	}
}

func TestCopingStrategy_EnumCount(t *testing.T) {
	if len(copingNames) != int(Denial)+1 {
		t.Errorf("copingNames has %d entries, expected %d", len(copingNames), int(Denial)+1)
	}
}

func TestIsolationPhase_String_AllValues(t *testing.T) {
	tests := []struct {
		p    IsolationPhase
		want string
	}{
		{IsolationNone, "none"},
		{IsolationBoredom, "boredom"},
		{IsolationLoneliness, "loneliness"},
		{IsolationSignificant, "significant"},
		{IsolationDestabilizing, "destabilizing"},
		{IsolationSevere, "severe"},
	}

	for _, tt := range tests {
		got := tt.p.String()
		if got != tt.want {
			t.Errorf("IsolationPhase(%d).String() = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestIsolationPhase_String_OutOfRange(t *testing.T) {
	got := IsolationPhase(99).String()
	if got != "unknown" {
		t.Errorf("IsolationPhase(99).String() = %q, want %q", got, "unknown")
	}
}

func TestIsolationPhase_EnumCount(t *testing.T) {
	if len(isolationPhaseNames) != int(IsolationSevere)+1 {
		t.Errorf("isolationPhaseNames has %d entries, expected %d", len(isolationPhaseNames), int(IsolationSevere)+1)
	}
}
