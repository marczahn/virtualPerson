package motivation

import "math"

func tempDeviation(bodyTemp float64) float64 {
	return clamp01(math.Abs(bodyTemp-36.6) / 6.0)
}

func energyMultiplier(p Personality) float64 {
	// Lower resilience and lower tolerance increase perceived urgency.
	return clamp(1+0.6*(0.5-p.EnergyResilience)+0.3*(0.5-p.FrustrationTolerance), 0.5, 1.5)
}

func socialMultiplier(p Personality) float64 {
	return clamp(1+0.8*(p.SocialFactor-0.5), 0.6, 1.4)
}

func stimulationMultiplier(p Personality) float64 {
	return clamp(1+0.8*(p.Curiosity-0.5), 0.6, 1.4)
}

func safetyMultiplier(p Personality) float64 {
	return clamp(1+0.6*(p.StressSensitivity-0.5)+0.4*(p.RiskAversion-0.5), 0.5, 1.5)
}

func identityMultiplier(p Personality) float64 {
	return clamp(1+0.6*(p.SelfObservation-0.5)+0.3*(0.5-p.FrustrationTolerance), 0.5, 1.5)
}

func clamp01(v float64) float64 {
	return clamp(v, 0, 1)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
