package psychology

// NegativeEmotionMultiplier returns the multiplier for negative emotional
// intensity based on neuroticism. High neuroticism amplifies negative reactions.
// Formula: 1.0 + (neuroticism - 0.5) * 1.2
// Neuroticism 0.8 → 1.36 (36% stronger), 0.2 → 0.64 (36% weaker).
func NegativeEmotionMultiplier(p Personality) float64 {
	return 1.0 + (p.Neuroticism-0.5)*1.2
}

// IsolationDistressRate returns the rate multiplier for isolation distress
// accumulation. High extraversion increases distress rate.
// Formula: 1.0 + (extraversion - 0.5) * 1.5
func IsolationDistressRate(p Personality) float64 {
	return 1.0 + (p.Extraversion-0.5)*1.5
}

// StimulationSeeking returns the person's need for external stimulation.
// Formula: 0.3 + extraversion * 0.5
func StimulationSeeking(p Personality) float64 {
	return 0.3 + p.Extraversion*0.5
}

// SocialRewardSensitivity returns how much the person benefits from social interaction.
// Formula: 0.2 + extraversion * 0.6
func SocialRewardSensitivity(p Personality) float64 {
	return 0.2 + p.Extraversion*0.6
}

// DisorderTolerance returns how much disorder/chaos the person can tolerate.
// Low conscientiousness = high tolerance.
// Formula: 1.0 - conscientiousness * 0.8
func DisorderTolerance(p Personality) float64 {
	return 1.0 - p.Conscientiousness*0.8
}

// SelfRegulationBonus returns the additional regulation capacity from conscientiousness.
// Formula: conscientiousness * 0.3
func SelfRegulationBonus(p Personality) float64 {
	return p.Conscientiousness * 0.3
}

// PlanningUnderStress returns the ability to maintain structured coping under stress.
// Formula: conscientiousness * 0.4
func PlanningUnderStress(p Personality) float64 {
	return p.Conscientiousness * 0.4
}

// ReappraisalAbility returns the person's capacity for cognitive reappraisal.
// Formula: 0.2 + openness * 0.5
func ReappraisalAbility(p Personality) float64 {
	return 0.2 + p.Openness*0.5
}

// NoveltyAsThreat returns how much novel stimuli trigger anxiety vs curiosity.
// Low openness = novelty is threatening. High openness = novelty is interesting.
// Formula: 1.0 - openness * 0.7
func NoveltyAsThreat(p Personality) float64 {
	return 1.0 - p.Openness*0.7
}

// BaselineRegulation computes the trait-based regulation capacity.
// Range: roughly 0.3 to 0.85 depending on personality.
func BaselineRegulation(p Personality) float64 {
	r := 0.3 +
		p.Conscientiousness*0.2 +
		p.Openness*0.15 +
		(1-p.Neuroticism)*0.2

	return clamp(r, 0, 1)
}

// IsolationResilience computes how resistant the person is to isolation effects.
// Range: roughly 0.1 to 0.9.
func IsolationResilience(p Personality) float64 {
	return (1-p.Extraversion)*0.5 + (1-p.Neuroticism)*0.3 + p.Conscientiousness*0.2
}

// DestabilizationThresholdHours returns the number of hours of isolation
// before psychological destabilization begins.
func DestabilizationThresholdHours(p Personality) float64 {
	return 48 + IsolationResilience(p)*120
}
