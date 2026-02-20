package motivation

type Drive string

const (
	DriveEnergy            Drive = "energy"
	DriveSocialConnection  Drive = "social_connection"
	DriveStimulation       Drive = "stimulation_novelty"
	DriveSafety            Drive = "safety"
	DriveIdentityCoherence Drive = "identity_coherence"
)

type Action string

const (
	ActionRest      Action = "rest"
	ActionEat       Action = "eat"
	ActionHydrate   Action = "hydrate"
	ActionReachOut  Action = "reach_out"
	ActionJournal   Action = "journal"
	ActionBreathe   Action = "breathe"
	ActionScanArea  Action = "scan_environment"
	ActionSeekWarm  Action = "seek_warmth"
	ActionSeekCool  Action = "seek_cooling"
	ActionMicroTask Action = "micro_task"
)

// Personality contains the 7 motivation multipliers.
// Values are expected in [0,1] and are clamped when consumed.
type Personality struct {
	StressSensitivity    float64
	EnergyResilience     float64
	Curiosity            float64
	SelfObservation      float64
	FrustrationTolerance float64
	RiskAversion         float64
	SocialFactor         float64
}

// ChronicState stores long-horizon pressure that can bias drives.
// These inputs are optional and default to zero value.
type ChronicState struct {
	ThreatLoad      float64
	IsolationLoad   float64
	IdentityStrain  float64
	FatiguePressure float64
}

type ActionConstraints struct {
	HasFood         bool
	HasPeopleNearby bool
	CanRest         bool
	CanExplore      bool
	HasQuietSpace   bool
}

type MotivationState struct {
	EnergyUrgency      float64
	SocialUrgency      float64
	StimulationUrgency float64
	SafetyUrgency      float64
	IdentityUrgency    float64

	ActiveGoalDrive   Drive
	ActiveGoalUrgency float64
}
