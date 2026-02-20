package motivation

// ActionCandidatesFor emits deterministic candidates from active goal and constraints.
func ActionCandidatesFor(goal Drive, c ActionConstraints) []Action {
	switch goal {
	case DriveEnergy:
		actions := make([]Action, 0, 3)
		if c.HasFood {
			actions = append(actions, ActionEat)
		}
		if c.CanRest {
			actions = append(actions, ActionRest)
		}
		actions = append(actions, ActionHydrate)
		return actions
	case DriveSocialConnection:
		actions := make([]Action, 0, 2)
		if c.HasPeopleNearby {
			actions = append(actions, ActionReachOut)
		}
		actions = append(actions, ActionJournal)
		return actions
	case DriveStimulation:
		actions := make([]Action, 0, 2)
		if c.CanExplore {
			actions = append(actions, ActionMicroTask, ActionScanArea)
			return actions
		}
		return append(actions, ActionMicroTask)
	case DriveSafety:
		actions := make([]Action, 0, 3)
		actions = append(actions, ActionBreathe, ActionScanArea)
		if c.HasQuietSpace {
			actions = append(actions, ActionRest)
		}
		return actions
	case DriveIdentityCoherence:
		actions := make([]Action, 0, 2)
		actions = append(actions, ActionJournal)
		if c.HasQuietSpace {
			actions = append(actions, ActionBreathe)
		}
		return actions
	default:
		return []Action{ActionBreathe}
	}
}
