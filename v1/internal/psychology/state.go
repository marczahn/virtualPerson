package psychology

import "time"

// State holds the psychological affect dimensions and active processes.
// This is what the consciousness layer receives — NOT raw biology.
type State struct {
	// Core affect dimensions derived from biological state.
	Arousal       float64 // 0-1, from adrenaline/HR/cortisol
	Valence       float64 // -1 to 1, positive=pleasant, negative=unpleasant
	Energy        float64 // 0-1, from fatigue/blood sugar/circadian
	CognitiveLoad float64 // 0-1, from cortisol duration/fatigue/blood sugar

	// Emotional regulation is a depletable resource.
	RegulationCapacity float64 // 0-1, depletes under sustained stress

	// Active cognitive distortions (empty when unstressed).
	ActiveDistortions []Distortion

	// Currently active coping strategies.
	ActiveCoping []CopingStrategy

	// Emotional memories activated by current state similarity.
	ActivatedMemories []EmotionalMemoryActivation

	// Effects of social isolation over time.
	Isolation IsolationState
}

// Personality represents the Big Five traits, fixed for the lifetime of a person.
// Each trait is in [0,1] with population mean at 0.5.
type Personality struct {
	Openness          float64 // curiosity, imagination, willingness to try new things
	Conscientiousness float64 // organization, self-discipline, planning
	Extraversion      float64 // sociability, assertiveness, positive emotionality
	Agreeableness     float64 // cooperativeness, empathy, conflict avoidance
	Neuroticism       float64 // emotional instability, negative affect sensitivity
}

// Distortion identifies a cognitive distortion that is currently active.
type Distortion int

const (
	Catastrophizing    Distortion = iota // assuming worst possible outcome
	EmotionalReasoning                   // "I feel it, so it must be true"
	Overgeneralization                   // "this always happens"
	MindReading                          // assuming others' negative thoughts
	Personalization                      // blaming self for external events
	AllOrNothing                         // black-and-white thinking
)

var distortionNames = [...]string{
	"catastrophizing",
	"emotional_reasoning",
	"overgeneralization",
	"mind_reading",
	"personalization",
	"all_or_nothing",
}

func (d Distortion) String() string {
	if int(d) < len(distortionNames) {
		return distortionNames[d]
	}
	return "unknown"
}

// CopingStrategy identifies a coping mechanism the person is using.
type CopingStrategy int

const (
	ProblemSolving CopingStrategy = iota // active, constructive
	Reappraisal                          // reframing the situation
	Acceptance                           // acknowledging without fighting
	Distraction                          // shifting attention away
	Suppression                          // pushing emotions down (costly)
	Rumination                           // repetitive negative thinking
	Denial                               // refusing to acknowledge reality
)

var copingNames = [...]string{
	"problem_solving",
	"reappraisal",
	"acceptance",
	"distraction",
	"suppression",
	"rumination",
	"denial",
}

func (c CopingStrategy) String() string {
	if int(c) < len(copingNames) {
		return copingNames[c]
	}
	return "unknown"
}

// EmotionalMemoryActivation represents a past emotional memory
// triggered by similarity to the current biological/psychological state.
type EmotionalMemoryActivation struct {
	MemoryID    string  // reference to stored memory
	Similarity  float64 // 0-1, how closely current state matches the memory's state
	ValenceSign float64 // -1 or 1, whether the memory is negative or positive
	Intensity   float64 // 0-1, strength of the activated memory after decay
}

// IsolationState tracks the effects of social isolation over time.
type IsolationState struct {
	Duration        time.Duration // how long since last social contact
	LonelinessLevel float64       // 0-1, increases along timeline
	Phase           IsolationPhase
}

// IsolationPhase represents stages of isolation with distinct effects.
type IsolationPhase int

const (
	IsolationNone          IsolationPhase = iota // 0-2hr, no significant effect
	IsolationBoredom                             // 2-8hr, restlessness, mild discomfort
	IsolationLoneliness                          // 8-24hr, active loneliness
	IsolationSignificant                         // 1-3 days, cognitive effects, mood disruption
	IsolationDestabilizing                       // 3-7 days, identity disturbance, paranoia
	IsolationSevere                              // 7+ days, hallucinations, dissociation
)

var isolationPhaseNames = [...]string{
	"none",
	"boredom",
	"loneliness",
	"significant",
	"destabilizing",
	"severe",
}

func (p IsolationPhase) String() string {
	if int(p) < len(isolationPhaseNames) {
		return isolationPhaseNames[p]
	}
	return "unknown"
}
