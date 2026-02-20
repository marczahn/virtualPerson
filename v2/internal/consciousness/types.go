package consciousness

import "github.com/marczahn/person/v2/internal/motivation"

type PromptDrive struct {
	Drive motivation.Drive
	Felt  string
}

type PromptContext struct {
	Primary          []PromptDrive
	Background       []PromptDrive
	GoalPull         string
	ContinuityBuffer []string
}

type ThoughtCategory string

const (
	ThoughtCategoryDrive            ThoughtCategory = "drive"
	ThoughtCategoryAssociativeDrift ThoughtCategory = "associative_drift"
)

type Thought struct {
	Category ThoughtCategory
	Drive    motivation.Drive
	Text     string
}

type TickSchedule struct {
	EveryTicks int
}

type ContinuityBuffer struct {
	capacity int
	thoughts []Thought
}

type ParsedState struct {
	Arousal float64
	Valence float64
}

type ParsedResponse struct {
	State          ParsedState
	Action         string
	DriveOverrides map[motivation.Drive]float64
	Narrative      string
}

type ActionOutcome struct {
	Action    string
	Executed  bool
	Satisfied bool
}

type ActionCooldowns map[string]int64

type ActionCooldownState map[string]int64
