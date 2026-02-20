package consciousness

import (
	"sort"

	"github.com/marczahn/person/v2/internal/motivation"
)

// ApplyParsedDriveOverridesForNextTick applies parsed DRIVE feedback as the
// perceived drive state for the next tick.
func ApplyParsedDriveOverridesForNextTick(raw motivation.MotivationState, parsed ParsedResponse) motivation.MotivationState {
	return ApplyDriveOverrides(raw, parsed.DriveOverrides)
}

func ApplyDriveOverrides(raw motivation.MotivationState, overrides map[motivation.Drive]float64) motivation.MotivationState {
	result := raw
	result.EnergyUrgency = effectiveDrive(raw.EnergyUrgency, overrides, motivation.DriveEnergy)
	result.SocialUrgency = effectiveDrive(raw.SocialUrgency, overrides, motivation.DriveSocialConnection)
	result.StimulationUrgency = effectiveDrive(raw.StimulationUrgency, overrides, motivation.DriveStimulation)
	result.SafetyUrgency = effectiveDrive(raw.SafetyUrgency, overrides, motivation.DriveSafety)
	result.IdentityUrgency = effectiveDrive(raw.IdentityUrgency, overrides, motivation.DriveIdentityCoherence)

	ordered := rankedDrives(result)
	if len(ordered) > 0 {
		result.ActiveGoalDrive = ordered[0].drive
		result.ActiveGoalUrgency = ordered[0].urgency
	}
	return result
}

type rankedDrive struct {
	drive    motivation.Drive
	urgency  float64
	priority int
}

var drivePriority = map[motivation.Drive]int{
	motivation.DriveSafety:            0,
	motivation.DriveEnergy:            1,
	motivation.DriveSocialConnection:  2,
	motivation.DriveIdentityCoherence: 3,
	motivation.DriveStimulation:       4,
}

func rankedDrives(state motivation.MotivationState) []rankedDrive {
	drives := []rankedDrive{
		{drive: motivation.DriveEnergy, urgency: clamp01(state.EnergyUrgency), priority: drivePriority[motivation.DriveEnergy]},
		{drive: motivation.DriveSocialConnection, urgency: clamp01(state.SocialUrgency), priority: drivePriority[motivation.DriveSocialConnection]},
		{drive: motivation.DriveStimulation, urgency: clamp01(state.StimulationUrgency), priority: drivePriority[motivation.DriveStimulation]},
		{drive: motivation.DriveSafety, urgency: clamp01(state.SafetyUrgency), priority: drivePriority[motivation.DriveSafety]},
		{drive: motivation.DriveIdentityCoherence, urgency: clamp01(state.IdentityUrgency), priority: drivePriority[motivation.DriveIdentityCoherence]},
	}

	sort.Slice(drives, func(i, j int) bool {
		if drives[i].urgency == drives[j].urgency {
			return drives[i].priority < drives[j].priority
		}
		return drives[i].urgency > drives[j].urgency
	})
	return drives
}

func effectiveDrive(raw float64, overrides map[motivation.Drive]float64, drive motivation.Drive) float64 {
	rawClamped := clamp01(raw)
	override, ok := overrides[drive]
	if !ok {
		return rawClamped
	}
	overrideClamped := clamp01(override)
	floor := rawClamped * 0.5
	if overrideClamped > floor {
		return overrideClamped
	}
	return floor
}
