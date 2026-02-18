package psychology

// CopingResources represents the person's available resources for coping.
type CopingResources struct {
	Cognitive  float64 // 0-1, available cognitive capacity (1 - cognitive_load)
	Energy     float64 // 0-1, available physical/mental energy
	Regulation float64 // 0-1, available emotional regulation capacity
}

// SelectCoping implements the coping strategy decision tree from the
// psychologist advisory. It returns a primary and optional secondary strategy.
//
// stress: 0-1, current stress level.
// controllability: 0-1, how controllable the stressor is.
func SelectCoping(stress float64, p Personality, res CopingResources, controllability float64) []CopingStrategy {
	// Under extreme stress, higher-order coping collapses.
	if stress > 0.85 {
		if res.Cognitive < 0.2 {
			return []CopingStrategy{Denial}
		}
		if p.Neuroticism > 0.7 {
			// High neuroticism under extreme stress: catastrophizing-driven
			// rumination or pure rumination.
			return []CopingStrategy{Rumination}
		}
	}

	// Is the stressor controllable?
	if controllability > 0.5 {
		if res.Cognitive > 0.4 && res.Energy > 0.3 {
			if p.Conscientiousness > 0.5 {
				return []CopingStrategy{ProblemSolving}
			}
			return []CopingStrategy{ProblemSolving, Distraction}
		}
		// Want to solve it but can't right now.
		return []CopingStrategy{Distraction}
	}

	// Stressor is uncontrollable.
	if res.Cognitive > 0.4 {
		if p.Openness > 0.5 && res.Regulation > 0.4 {
			return []CopingStrategy{Reappraisal}
		}
		if p.Agreeableness > 0.6 {
			return []CopingStrategy{Acceptance}
		}
		return []CopingStrategy{Distraction}
	}

	// Low resources, uncontrollable stressor.
	if p.Neuroticism > 0.6 {
		return []CopingStrategy{Rumination, Suppression}
	}
	return []CopingStrategy{Distraction, Acceptance}
}
