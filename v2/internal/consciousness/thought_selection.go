package consciousness

import "github.com/marczahn/person/v2/internal/motivation"

func SelectSpontaneousThought(state motivation.MotivationState, schedule TickSchedule, tick int) (Thought, bool) {
	if !schedule.due(tick) {
		return Thought{}, false
	}

	ordered := rankedDrives(state)
	if len(ordered) == 0 || ordered[0].urgency <= 0 {
		return Thought{
			Category: ThoughtCategoryAssociativeDrift,
			Text:     "A loose associative thread drifts into awareness.",
		}, true
	}

	top := ordered[0]
	return Thought{
		Category: ThoughtCategoryDrive,
		Drive:    top.drive,
		Text:     driveThoughtText(top.drive),
	}, true
}

func (s TickSchedule) due(tick int) bool {
	if s.EveryTicks <= 0 || tick <= 0 {
		return false
	}
	return tick%s.EveryTicks == 0
}
