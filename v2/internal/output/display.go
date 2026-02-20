package output

import (
	"fmt"
	"math"

	"github.com/marczahn/person/v2/internal/motivation"
)

type DriveChange struct {
	Name     motivation.Drive
	Previous float64
	Current  float64
}

func SignificantDriveChanges(previous, current motivation.MotivationState, threshold float64) []DriveChange {
	if threshold < 0 {
		threshold = 0
	}

	var out []DriveChange
	for _, drive := range []motivation.Drive{
		motivation.DriveEnergy,
		motivation.DriveSocialConnection,
		motivation.DriveStimulation,
		motivation.DriveSafety,
		motivation.DriveIdentityCoherence,
	} {
		prevValue := driveValue(previous, drive)
		currValue := driveValue(current, drive)
		if math.Abs(currValue-prevValue) >= threshold {
			out = append(out, DriveChange{
				Name:     drive,
				Previous: prevValue,
				Current:  currValue,
			})
		}
	}
	return out
}

func FormatDriveChangeLines(changes []DriveChange) []string {
	out := make([]string, 0, len(changes))
	for _, change := range changes {
		diff := change.Current - change.Previous
		out = append(out, FormatTaggedLine(
			SourceDRIVES,
			fmt.Sprintf("%s: %.2f -> %.2f (%+.2f)", change.Name, change.Previous, change.Current, diff),
		))
	}
	return out
}

func driveValue(state motivation.MotivationState, drive motivation.Drive) float64 {
	switch drive {
	case motivation.DriveEnergy:
		return state.EnergyUrgency
	case motivation.DriveSocialConnection:
		return state.SocialUrgency
	case motivation.DriveStimulation:
		return state.StimulationUrgency
	case motivation.DriveSafety:
		return state.SafetyUrgency
	case motivation.DriveIdentityCoherence:
		return state.IdentityUrgency
	default:
		return 0
	}
}
