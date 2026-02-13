package psychology

import (
	"testing"

	"github.com/marczahn/person/internal/biology"
)

func TestComputeArousal_DefaultState_Low(t *testing.T) {
	bio := biology.NewDefaultState()
	a := computeArousal(&bio)

	// Default state: adrenaline=0, HR=70, cortisol=0.1, fatigue=0.
	// Expected: 0.3*0 + 0.25*(30/160) + 0.2*0.1 + 0 + 0.1*(1) ≈ 0.047 + 0.02 + 0.1 = 0.167
	if a < 0.1 || a > 0.25 {
		t.Errorf("arousal at default = %f, expected 0.1-0.25", a)
	}
}

func TestComputeArousal_HighAdrenaline(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Adrenaline = 0.9
	bio.HeartRate = 150

	a := computeArousal(&bio)

	// Should be high: 0.3*0.9 + 0.25*(110/160) + 0.2*0.1 + 0.1 ≈ 0.27 + 0.172 + 0.02 + 0.1 = 0.56
	if a < 0.4 {
		t.Errorf("arousal with high adrenaline = %f, expected > 0.4", a)
	}
}

func TestComputeValence_DefaultState_Positive(t *testing.T) {
	bio := biology.NewDefaultState()
	v := computeValence(&bio)

	// Default: serotonin=0.5, dopamine=0.3, endorphins=0.1, cortisol=0.1.
	// Expected: 0.35*0.5 + 0.30*0.3 + 0.15*0.1 - 0.20*0.1 = 0.175 + 0.09 + 0.015 - 0.02 = 0.26
	if v < 0.15 || v > 0.35 {
		t.Errorf("valence at default = %f, expected 0.15-0.35", v)
	}
}

func TestComputeValence_HighCortisol_LowSerotonin_Negative(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Cortisol = 0.8
	bio.Serotonin = 0.1
	bio.Dopamine = 0.1

	v := computeValence(&bio)

	// Expected: 0.35*0.1 + 0.30*0.1 + 0.15*0.1 - 0.20*0.8 = 0.035 + 0.03 + 0.015 - 0.16 = -0.08
	if v >= 0 {
		t.Errorf("valence with high cortisol/low serotonin = %f, expected negative", v)
	}
}

func TestComputeEnergy_DefaultState_Moderate(t *testing.T) {
	bio := biology.NewDefaultState()
	e := computeEnergy(&bio)

	// Default: fatigue=0, blood sugar=90, circadian phase=8 (morning), dopamine=0.3.
	// normBS = (90-50)/150 = 0.267
	// energy = 0.3*1 + 0.25*0.267 + 0.25*alertness(8) + 0.2*0.3
	if e < 0.4 || e > 0.8 {
		t.Errorf("energy at default = %f, expected 0.4-0.8", e)
	}
}

func TestComputeEnergy_HighFatigue_Low(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.Fatigue = 0.9
	bio.BloodSugar = 55
	bio.Dopamine = 0.1

	e := computeEnergy(&bio)

	if e > 0.3 {
		t.Errorf("energy with high fatigue/low BS = %f, expected < 0.3", e)
	}
}

func TestComputeCognitiveLoad_DefaultState_Low(t *testing.T) {
	bio := biology.NewDefaultState()
	cl := computeCognitiveLoad(&bio)

	// Default: cortisol load=0, fatigue=0, blood sugar=90.
	// normBSInv = 1 - 40/150 = 0.733
	// cl = 0.3*0 + 0.25*0 + 0.25*0.733 + 0 = 0.183
	if cl < 0.1 || cl > 0.3 {
		t.Errorf("cognitive load at default = %f, expected 0.1-0.3", cl)
	}
}

func TestComputeCognitiveLoad_SustainedStress_High(t *testing.T) {
	bio := biology.NewDefaultState()
	bio.CortisolLoad = 10 // substantial accumulated stress
	bio.Fatigue = 0.7
	bio.BloodSugar = 55

	cl := computeCognitiveLoad(&bio)

	if cl < 0.5 {
		t.Errorf("cognitive load under sustained stress = %f, expected > 0.5", cl)
	}
}

func TestComputeStress_CalmState(t *testing.T) {
	s := computeStress(0.1, 0.3, 0.1)
	if s > 0.15 {
		t.Errorf("stress in calm state = %f, expected < 0.15", s)
	}
}

func TestComputeStress_HighArousalNegativeValence(t *testing.T) {
	s := computeStress(0.8, -0.7, 0.6)
	if s < 0.5 {
		t.Errorf("stress with high arousal/neg valence = %f, expected > 0.5", s)
	}
}

func TestProcessor_Process_DefaultState(t *testing.T) {
	p := Personality{
		Openness:          0.5,
		Conscientiousness: 0.5,
		Extraversion:      0.5,
		Agreeableness:     0.5,
		Neuroticism:       0.5,
	}
	proc := NewProcessor(p)
	bio := biology.NewDefaultState()

	ps := proc.Process(&bio, 1.0, 0.5)

	if ps.Arousal < 0 || ps.Arousal > 0.3 {
		t.Errorf("arousal = %f, expected low for default state", ps.Arousal)
	}
	if ps.Valence < 0 {
		t.Errorf("valence = %f, expected positive for default state", ps.Valence)
	}
	if ps.Energy < 0.3 {
		t.Errorf("energy = %f, expected moderate for default state", ps.Energy)
	}
	if ps.RegulationCapacity < 0.3 {
		t.Errorf("regulation = %f, expected moderate-high", ps.RegulationCapacity)
	}
}

func TestProcessor_Process_StressedState(t *testing.T) {
	p := Personality{
		Openness:          0.3,
		Conscientiousness: 0.3,
		Extraversion:      0.5,
		Agreeableness:     0.5,
		Neuroticism:       0.8,
	}
	proc := NewProcessor(p)
	bio := biology.NewDefaultState()
	bio.Adrenaline = 0.7
	bio.Cortisol = 0.7
	bio.HeartRate = 130
	bio.Serotonin = 0.2
	bio.Dopamine = 0.1

	ps := proc.Process(&bio, 1.0, 0.0)

	if ps.Arousal < 0.2 {
		t.Errorf("arousal = %f, expected elevated under stress", ps.Arousal)
	}
	if ps.Valence > 0 {
		t.Errorf("valence = %f, expected negative under stress", ps.Valence)
	}
}

func TestProcessor_FeedbackChanges_Rumination(t *testing.T) {
	ps := State{
		ActiveCoping: []CopingStrategy{Rumination},
	}
	p := Personality{}
	proc := NewProcessor(p)

	changes := proc.FeedbackChanges(&ps, 1.0)

	var hasCortisol, hasSerotonin bool
	for _, c := range changes {
		if c.Variable == biology.VarCortisol && c.Delta > 0 {
			hasCortisol = true
		}
		if c.Variable == biology.VarSerotonin && c.Delta < 0 {
			hasSerotonin = true
		}
	}
	if !hasCortisol {
		t.Error("rumination should increase cortisol")
	}
	if !hasSerotonin {
		t.Error("rumination should decrease serotonin")
	}
}

func TestProcessor_FeedbackChanges_Acceptance(t *testing.T) {
	ps := State{
		ActiveCoping: []CopingStrategy{Acceptance},
	}
	p := Personality{}
	proc := NewProcessor(p)

	changes := proc.FeedbackChanges(&ps, 1.0)

	var hasCortisol, hasSerotonin bool
	for _, c := range changes {
		if c.Variable == biology.VarCortisol && c.Delta < 0 {
			hasCortisol = true
		}
		if c.Variable == biology.VarSerotonin && c.Delta > 0 {
			hasSerotonin = true
		}
	}
	if !hasCortisol {
		t.Error("acceptance should decrease cortisol")
	}
	if !hasSerotonin {
		t.Error("acceptance should increase serotonin")
	}
}

func TestProcessor_FeedbackChanges_Catastrophizing(t *testing.T) {
	ps := State{
		ActiveDistortions: []Distortion{Catastrophizing},
	}
	p := Personality{}
	proc := NewProcessor(p)

	changes := proc.FeedbackChanges(&ps, 1.0)

	var hasAdrenaline, hasCortisol bool
	for _, c := range changes {
		if c.Variable == biology.VarAdrenaline && c.Delta > 0 {
			hasAdrenaline = true
		}
		if c.Variable == biology.VarCortisol && c.Delta > 0 {
			hasCortisol = true
		}
	}
	if !hasAdrenaline {
		t.Error("catastrophizing should increase adrenaline")
	}
	if !hasCortisol {
		t.Error("catastrophizing should increase cortisol")
	}
}

func TestProcessor_FeedbackChanges_MultipleDistortions(t *testing.T) {
	ps := State{
		ActiveDistortions: []Distortion{Catastrophizing, EmotionalReasoning, Overgeneralization},
	}
	p := Personality{}
	proc := NewProcessor(p)

	changes := proc.FeedbackChanges(&ps, 1.0)

	distortionLoadFound := false
	for _, c := range changes {
		if c.Source == "psych_distortion_load" {
			distortionLoadFound = true
		}
	}
	if !distortionLoadFound {
		t.Error("3+ distortions should add distortion load cortisol")
	}
}
