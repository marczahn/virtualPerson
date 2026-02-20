package consciousness

import (
	"strings"

	"github.com/marczahn/person/v2/internal/motivation"
)

func BuildPromptContext(state motivation.MotivationState) PromptContext {
	return BuildPromptContextWithContinuity(state, nil)
}

func BuildPromptContextWithContinuity(state motivation.MotivationState, continuity []Thought) PromptContext {
	ordered := rankedDrives(state)
	primaryCount := 2
	if len(ordered) < primaryCount {
		primaryCount = len(ordered)
	}

	primary := make([]PromptDrive, 0, primaryCount)
	for _, d := range ordered[:primaryCount] {
		primary = append(primary, PromptDrive{
			Drive: d.drive,
			Felt:  feltExperience(d.drive, d.urgency),
		})
	}

	background := make([]PromptDrive, 0, len(ordered)-primaryCount)
	for _, d := range ordered[primaryCount:] {
		background = append(background, PromptDrive{
			Drive: d.drive,
			Felt:  feltExperience(d.drive, d.urgency),
		})
	}

	return PromptContext{
		Primary:          primary,
		Background:       background,
		GoalPull:         implicitGoalPull(state.ActiveGoalDrive),
		ContinuityBuffer: continuityLines(continuity),
	}
}

func feltExperience(drive motivation.Drive, urgency float64) string {
	level := urgencyLevel(clamp01(urgency))
	switch drive {
	case motivation.DriveEnergy:
		return []string{
			"a faint pull toward rest and nourishment lingers in the background.",
			"a noticeable fatigue-and-hunger pull is starting to build.",
			"an insistent need for rest and nourishment is pressing into attention.",
			"an urgent depletion is dominating attention and demanding recovery now.",
		}[level]
	case motivation.DriveSocialConnection:
		return []string{
			"a light sense of distance from others is present.",
			"a growing wish for contact and response is becoming noticeable.",
			"an insistent loneliness is pressing for connection.",
			"an urgent need to reach someone is dominating focus.",
		}[level]
	case motivation.DriveStimulation:
		return []string{
			"a mild restlessness for novelty hums in the background.",
			"a noticeable urge for engagement and novelty is rising.",
			"an insistent need for stimulation is pushing for action.",
			"an urgent craving for meaningful engagement is taking over attention.",
		}[level]
	case motivation.DriveSafety:
		return []string{
			"a faint vigilance remains in the background.",
			"a noticeable need to check for safety is surfacing.",
			"an insistent threat-sensitivity is narrowing attention.",
			"an urgent need to secure safety is dominating attention.",
		}[level]
	case motivation.DriveIdentityCoherence:
		return []string{
			"a light pull to make sense of experience is present.",
			"a noticeable tension about self-coherence is forming.",
			"an insistent need to regain internal coherence is pressing.",
			"an urgent need to stabilize meaning and identity is overwhelming focus.",
		}[level]
	default:
		return "a general pressure is present."
	}
}

func implicitGoalPull(drive motivation.Drive) string {
	switch drive {
	case motivation.DriveEnergy:
		return "A pull toward food, water, and recovery keeps surfacing."
	case motivation.DriveSocialConnection:
		return "A pull toward contact and response keeps surfacing."
	case motivation.DriveStimulation:
		return "A pull toward something engaging and novel keeps surfacing."
	case motivation.DriveSafety:
		return "A pull toward checking safety and reducing threat keeps surfacing."
	case motivation.DriveIdentityCoherence:
		return "A pull toward making sense of experience keeps surfacing."
	default:
		return "A pull toward immediate regulation keeps surfacing."
	}
}

func driveThoughtText(drive motivation.Drive) string {
	switch drive {
	case motivation.DriveEnergy:
		return "Food and recovery keep intruding into thought."
	case motivation.DriveSocialConnection:
		return "The need for contact keeps returning to mind."
	case motivation.DriveStimulation:
		return "A search for novelty keeps tugging at attention."
	case motivation.DriveSafety:
		return "Threat-checking keeps cycling through awareness."
	case motivation.DriveIdentityCoherence:
		return "A need to make sense of self keeps pressing forward."
	default:
		return "A regulatory need keeps resurfacing."
	}
}

func continuityLines(continuity []Thought) []string {
	if len(continuity) == 0 {
		return nil
	}
	lines := make([]string, 0, len(continuity))
	for _, thought := range continuity {
		text := strings.TrimSpace(thought.Text)
		if text != "" {
			lines = append(lines, text)
		}
	}
	if len(lines) == 0 {
		return nil
	}
	return lines
}

func urgencyLevel(v float64) int {
	switch {
	case v >= 0.75:
		return 3
	case v >= 0.50:
		return 2
	case v >= 0.25:
		return 1
	default:
		return 0
	}
}
