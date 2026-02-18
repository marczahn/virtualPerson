package consciousness

import "time"

// ThoughtType categorizes the origin of a conscious thought.
type ThoughtType int

const (
	// Reactive thoughts are triggered by salience threshold breach.
	Reactive ThoughtType = iota

	// Spontaneous thoughts emerge from the priority queue when idle.
	Spontaneous

	// Conversational thoughts are triggered by external speech or actions.
	Conversational
)

var thoughtTypeNames = [...]string{
	"reactive",
	"spontaneous",
	"conversational",
}

// ExternalInputType distinguishes speech from physical actions.
type ExternalInputType int

const (
	InputSpeech ExternalInputType = iota
	InputAction
)

func (t ExternalInputType) String() string {
	switch t {
	case InputSpeech:
		return "speech"
	case InputAction:
		return "action"
	default:
		return "unknown"
	}
}

// ExternalInput represents communicative input directed at the person.
type ExternalInput struct {
	Type    ExternalInputType
	Content string
}

func (t ThoughtType) String() string {
	if int(t) < len(thoughtTypeNames) {
		return thoughtTypeNames[t]
	}
	return "unknown"
}

// Priority categories for spontaneous thought generation,
// ordered from highest to lowest priority.
type Priority int

const (
	PriorityPredictionError Priority = iota // "something unexpected happened"
	PriorityBiologicalNeed                  // hunger, pain, thermal discomfort
	PriorityGoalRehearsal                   // upcoming tasks, unfinished plans
	PrioritySocialModeling                  // "what did they think of me?"
	PriorityAssociativeDrift                // daydreaming, mind-wandering
)

var priorityNames = [...]string{
	"prediction_error",
	"biological_need",
	"goal_rehearsal",
	"social_modeling",
	"associative_drift",
}

func (p Priority) String() string {
	if int(p) < len(priorityNames) {
		return priorityNames[p]
	}
	return "unknown"
}

// Thought represents a single conscious experience — a thought, feeling,
// or realization that the person becomes aware of.
type Thought struct {
	Type      ThoughtType
	Priority  Priority
	Content   string    // the LLM-generated first-person experience
	Trigger   string    // what caused this thought (stimulus description or queue category)
	Timestamp time.Time

	// Feedback contains coping strategies or distortions detected in the thought,
	// which feed back into the biology layer.
	Feedback ThoughtFeedback
}

// EmotionalTag is the machine-readable emotional state the LLM emits at the
// end of each response. Arousal ∈ [0,1], Valence ∈ [-1,1].
type EmotionalTag struct {
	Arousal float64
	Valence float64
}

// ThoughtFeedback captures aspects of the thought that should
// modify the biological state (the consciousness→biology feedback loop).
type ThoughtFeedback struct {
	ActiveCoping      []string     // coping strategies present in the thought
	ActiveDistortions []string     // cognitive distortions present in the thought
	EmotionalState    EmotionalTag // parsed from [STATE: arousal=X, valence=Y] annotation
}
