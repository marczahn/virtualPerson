package psychology

import "testing"

func TestSelectCoping_ExtremeStress_LowCognitive_Denial(t *testing.T) {
	p := Personality{Neuroticism: 0.5}
	res := CopingResources{Cognitive: 0.1, Energy: 0.2, Regulation: 0.1}

	coping := SelectCoping(0.9, p, res, 0.0)

	if len(coping) == 0 || coping[0] != Denial {
		t.Errorf("extreme stress + low cognitive → expected Denial, got %v", coping)
	}
}

func TestSelectCoping_ExtremeStress_HighNeuroticism_Rumination(t *testing.T) {
	p := Personality{Neuroticism: 0.8}
	res := CopingResources{Cognitive: 0.5, Energy: 0.5, Regulation: 0.3}

	coping := SelectCoping(0.9, p, res, 0.0)

	if len(coping) == 0 || coping[0] != Rumination {
		t.Errorf("extreme stress + high N → expected Rumination, got %v", coping)
	}
}

func TestSelectCoping_ControllableStressor_HighResources_ProblemSolving(t *testing.T) {
	p := Personality{Conscientiousness: 0.7}
	res := CopingResources{Cognitive: 0.6, Energy: 0.5, Regulation: 0.5}

	coping := SelectCoping(0.5, p, res, 0.8)

	if len(coping) == 0 || coping[0] != ProblemSolving {
		t.Errorf("controllable + resources + high C → expected ProblemSolving, got %v", coping)
	}
}

func TestSelectCoping_ControllableStressor_LowConscientious_IncludesDistraction(t *testing.T) {
	p := Personality{Conscientiousness: 0.3}
	res := CopingResources{Cognitive: 0.6, Energy: 0.5, Regulation: 0.5}

	coping := SelectCoping(0.5, p, res, 0.8)

	if len(coping) < 2 {
		t.Errorf("low C + controllable → expected 2 strategies, got %v", coping)
		return
	}
	if coping[0] != ProblemSolving || coping[1] != Distraction {
		t.Errorf("expected [ProblemSolving, Distraction], got %v", coping)
	}
}

func TestSelectCoping_UncontrollableStressor_HighOpenness_Reappraisal(t *testing.T) {
	p := Personality{Openness: 0.7}
	res := CopingResources{Cognitive: 0.6, Energy: 0.5, Regulation: 0.5}

	coping := SelectCoping(0.5, p, res, 0.2)

	if len(coping) == 0 || coping[0] != Reappraisal {
		t.Errorf("uncontrollable + open + regulated → expected Reappraisal, got %v", coping)
	}
}

func TestSelectCoping_UncontrollableStressor_HighAgreeableness_Acceptance(t *testing.T) {
	p := Personality{Openness: 0.3, Agreeableness: 0.8}
	res := CopingResources{Cognitive: 0.6, Energy: 0.5, Regulation: 0.3}

	coping := SelectCoping(0.5, p, res, 0.2)

	if len(coping) == 0 || coping[0] != Acceptance {
		t.Errorf("uncontrollable + agreeable → expected Acceptance, got %v", coping)
	}
}

func TestSelectCoping_LowResources_HighNeuroticism_Rumination(t *testing.T) {
	p := Personality{Neuroticism: 0.7}
	res := CopingResources{Cognitive: 0.2, Energy: 0.2, Regulation: 0.1}

	coping := SelectCoping(0.5, p, res, 0.2)

	if len(coping) < 2 || coping[0] != Rumination || coping[1] != Suppression {
		t.Errorf("low resources + high N → expected [Rumination, Suppression], got %v", coping)
	}
}

func TestSelectCoping_LowResources_LowNeuroticism_DistractionAcceptance(t *testing.T) {
	p := Personality{Neuroticism: 0.3}
	res := CopingResources{Cognitive: 0.2, Energy: 0.2, Regulation: 0.1}

	coping := SelectCoping(0.5, p, res, 0.2)

	if len(coping) < 2 || coping[0] != Distraction || coping[1] != Acceptance {
		t.Errorf("low resources + low N → expected [Distraction, Acceptance], got %v", coping)
	}
}
