package motivation

import (
	"github.com/marczahn/person/v2/internal/biology"
)

var drivePriority = []Drive{
	DriveSafety,
	DriveEnergy,
	DriveSocialConnection,
	DriveIdentityCoherence,
	DriveStimulation,
}

// Compute calculates drive urgencies and active goal deterministically.
func Compute(bio biology.State, personality Personality, chronic ChronicState) MotivationState {
	b := normalizedBio(bio)
	p := clampedPersonality(personality)
	c := clampedChronic(chronic)

	energyBase := clamp01(0.65*(1-b.Energy) + 0.35*b.Hunger + 0.15*c.FatiguePressure)
	socialBase := clamp01(b.SocialDeficit + 0.25*c.IsolationLoad)
	stimBase := clamp01(0.60*(1-b.CognitiveCapacity) + 0.40*(1-b.Mood))
	safetyBase := clamp01(0.50*b.Stress + 0.25*b.PhysicalTension + 0.25*tempDeviation(b.BodyTemp) + 0.20*c.ThreatLoad)
	identityBase := clamp01(0.55*(1-b.Mood) + 0.45*(1-b.CognitiveCapacity) + 0.20*c.IdentityStrain)

	energy := clamp01(energyBase * energyMultiplier(p))
	social := clamp01(socialBase * socialMultiplier(p))
	stim := clamp01(stimBase * stimulationMultiplier(p))
	safety := clamp01(safetyBase * safetyMultiplier(p))
	identity := clamp01(identityBase * identityMultiplier(p))

	state := MotivationState{
		EnergyUrgency:      energy,
		SocialUrgency:      social,
		StimulationUrgency: stim,
		SafetyUrgency:      safety,
		IdentityUrgency:    identity,
	}
	state.ActiveGoalDrive, state.ActiveGoalUrgency = selectActiveGoal(state)
	return state
}

func selectActiveGoal(m MotivationState) (Drive, float64) {
	driveValues := map[Drive]float64{
		DriveEnergy:            m.EnergyUrgency,
		DriveSocialConnection:  m.SocialUrgency,
		DriveStimulation:       m.StimulationUrgency,
		DriveSafety:            m.SafetyUrgency,
		DriveIdentityCoherence: m.IdentityUrgency,
	}

	best := drivePriority[0]
	bestValue := driveValues[best]
	for _, d := range drivePriority[1:] {
		if driveValues[d] > bestValue {
			best = d
			bestValue = driveValues[d]
		}
	}
	return best, bestValue
}

func normalizedBio(b biology.State) biology.State {
	b.Energy = clamp01(b.Energy)
	b.Stress = clamp01(b.Stress)
	b.CognitiveCapacity = clamp01(b.CognitiveCapacity)
	b.Mood = clamp01(b.Mood)
	b.PhysicalTension = clamp01(b.PhysicalTension)
	b.Hunger = clamp01(b.Hunger)
	b.SocialDeficit = clamp01(b.SocialDeficit)
	b.BodyTemp = clamp(b.BodyTemp, biology.Ranges.BodyTemp.Min, biology.Ranges.BodyTemp.Max)
	return b
}

func clampedPersonality(p Personality) Personality {
	p.StressSensitivity = clamp01(p.StressSensitivity)
	p.EnergyResilience = clamp01(p.EnergyResilience)
	p.Curiosity = clamp01(p.Curiosity)
	p.SelfObservation = clamp01(p.SelfObservation)
	p.FrustrationTolerance = clamp01(p.FrustrationTolerance)
	p.RiskAversion = clamp01(p.RiskAversion)
	p.SocialFactor = clamp01(p.SocialFactor)
	return p
}

func clampedChronic(c ChronicState) ChronicState {
	c.ThreatLoad = clamp01(c.ThreatLoad)
	c.IsolationLoad = clamp01(c.IsolationLoad)
	c.IdentityStrain = clamp01(c.IdentityStrain)
	c.FatiguePressure = clamp01(c.FatiguePressure)
	return c
}
