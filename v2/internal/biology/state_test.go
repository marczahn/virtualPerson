package biology_test

import (
	"math"
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestNewDefaultState_Baselines(t *testing.T) {
	s := biology.NewDefaultState()

	checks := []struct {
		name string
		got  float64
		want float64
	}{
		{"Energy", s.Energy, 0.80},
		{"Stress", s.Stress, 0.10},
		{"CognitiveCapacity", s.CognitiveCapacity, 1.00},
		{"Mood", s.Mood, 0.50},
		{"PhysicalTension", s.PhysicalTension, 0.05},
		{"Hunger", s.Hunger, 0.10},
		{"SocialDeficit", s.SocialDeficit, 0.00},
		{"BodyTemp", s.BodyTemp, 36.6},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("NewDefaultState().%s = %v, want %v", c.name, c.got, c.want)
		}
	}
	if s.UpdatedAt.IsZero() {
		t.Errorf("NewDefaultState().UpdatedAt should not be zero")
	}
}

func TestClampAll_EnforcesRanges(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(s *biology.State)
		field    string
		getField func(s *biology.State) float64
		want     float64
	}{
		// Energy bounds
		{
			name:     "Energy below min clamped to 0",
			setup:    func(s *biology.State) { s.Energy = -0.5 },
			field:    "Energy",
			getField: func(s *biology.State) float64 { return s.Energy },
			want:     0.0,
		},
		{
			name:     "Energy above max clamped to 1",
			setup:    func(s *biology.State) { s.Energy = 1.5 },
			field:    "Energy",
			getField: func(s *biology.State) float64 { return s.Energy },
			want:     1.0,
		},
		// Stress bounds
		{
			name:     "Stress below min clamped to 0",
			setup:    func(s *biology.State) { s.Stress = -0.1 },
			field:    "Stress",
			getField: func(s *biology.State) float64 { return s.Stress },
			want:     0.0,
		},
		{
			name:     "Stress above max clamped to 1",
			setup:    func(s *biology.State) { s.Stress = 1.5 },
			field:    "Stress",
			getField: func(s *biology.State) float64 { return s.Stress },
			want:     1.0,
		},
		// CognitiveCapacity bounds
		{
			name:     "CognitiveCapacity below min clamped to 0",
			setup:    func(s *biology.State) { s.CognitiveCapacity = -0.3 },
			field:    "CognitiveCapacity",
			getField: func(s *biology.State) float64 { return s.CognitiveCapacity },
			want:     0.0,
		},
		{
			name:     "CognitiveCapacity above max clamped to 1",
			setup:    func(s *biology.State) { s.CognitiveCapacity = 2.0 },
			field:    "CognitiveCapacity",
			getField: func(s *biology.State) float64 { return s.CognitiveCapacity },
			want:     1.0,
		},
		// Mood bounds
		{
			name:     "Mood below min clamped to 0",
			setup:    func(s *biology.State) { s.Mood = -0.2 },
			field:    "Mood",
			getField: func(s *biology.State) float64 { return s.Mood },
			want:     0.0,
		},
		{
			name:     "Mood above max clamped to 1",
			setup:    func(s *biology.State) { s.Mood = 1.1 },
			field:    "Mood",
			getField: func(s *biology.State) float64 { return s.Mood },
			want:     1.0,
		},
		// PhysicalTension bounds
		{
			name:     "PhysicalTension below min clamped to 0",
			setup:    func(s *biology.State) { s.PhysicalTension = -0.01 },
			field:    "PhysicalTension",
			getField: func(s *biology.State) float64 { return s.PhysicalTension },
			want:     0.0,
		},
		{
			name:     "PhysicalTension above max clamped to 1",
			setup:    func(s *biology.State) { s.PhysicalTension = 1.5 },
			field:    "PhysicalTension",
			getField: func(s *biology.State) float64 { return s.PhysicalTension },
			want:     1.0,
		},
		// Hunger bounds
		{
			name:     "Hunger below min clamped to 0",
			setup:    func(s *biology.State) { s.Hunger = -0.1 },
			field:    "Hunger",
			getField: func(s *biology.State) float64 { return s.Hunger },
			want:     0.0,
		},
		{
			name:     "Hunger above max clamped to 1",
			setup:    func(s *biology.State) { s.Hunger = 1.5 },
			field:    "Hunger",
			getField: func(s *biology.State) float64 { return s.Hunger },
			want:     1.0,
		},
		// SocialDeficit bounds
		{
			name:     "SocialDeficit below min clamped to 0",
			setup:    func(s *biology.State) { s.SocialDeficit = -0.5 },
			field:    "SocialDeficit",
			getField: func(s *biology.State) float64 { return s.SocialDeficit },
			want:     0.0,
		},
		{
			name:     "SocialDeficit above max clamped to 1",
			setup:    func(s *biology.State) { s.SocialDeficit = 1.2 },
			field:    "SocialDeficit",
			getField: func(s *biology.State) float64 { return s.SocialDeficit },
			want:     1.0,
		},
		// BodyTemp bounds (wider range: 25-43)
		{
			name:     "BodyTemp below min clamped to 25",
			setup:    func(s *biology.State) { s.BodyTemp = 20.0 },
			field:    "BodyTemp",
			getField: func(s *biology.State) float64 { return s.BodyTemp },
			want:     25.0,
		},
		{
			name:     "BodyTemp above max clamped to 43",
			setup:    func(s *biology.State) { s.BodyTemp = 50.0 },
			field:    "BodyTemp",
			getField: func(s *biology.State) float64 { return s.BodyTemp },
			want:     43.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := biology.NewDefaultState()
			tt.setup(s)
			biology.ClampAll(s)
			got := tt.getField(s)
			if got != tt.want {
				t.Errorf("%s = %v after clampAll, want %v", tt.field, got, tt.want)
			}
		})
	}
}

func TestClamp_AllVariables(t *testing.T) {
	tests := []struct {
		name   string
		v      float64
		lo, hi float64
		want   float64
	}{
		{"below range", -1.0, 0.0, 1.0, 0.0},
		{"above range", 2.0, 0.0, 1.0, 1.0},
		{"in range", 0.5, 0.0, 1.0, 0.5},
		{"at min", 0.0, 0.0, 1.0, 0.0},
		{"at max", 1.0, 0.0, 1.0, 1.0},
		{"body temp below", 20.0, 25.0, 43.0, 25.0},
		{"body temp above", 50.0, 25.0, 43.0, 43.0},
		{"body temp normal", 36.6, 25.0, 43.0, 36.6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := biology.Clamp(tt.v, tt.lo, tt.hi)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("clamp(%v, %v, %v) = %v, want %v", tt.v, tt.lo, tt.hi, got, tt.want)
			}
		})
	}
}
