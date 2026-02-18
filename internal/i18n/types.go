package i18n

// Translations holds all translatable text for the simulation.
type Translations struct {
	Consciousness ConsciousnessTranslations `yaml:"consciousness"`
	Feedback      FeedbackTranslations      `yaml:"feedback"`
	Sense         SenseTranslations         `yaml:"sense"`
	Reviewer      ReviewerTranslations      `yaml:"reviewer"`
	Biology       BiologyTranslations       `yaml:"biology"`
	CLI           CLITranslations           `yaml:"cli"`
	Client        ClientTranslations        `yaml:"client"`
	Output        OutputTranslations        `yaml:"output"`
	Defaults      DefaultsTranslations      `yaml:"defaults"`
}

// ConsciousnessTranslations covers the LLM system prompt, state descriptions,
// distortion descriptions, and prompt fragments.
type ConsciousnessTranslations struct {
	SystemPrompt        string                `yaml:"system_prompt"`
	Identity            IdentityTranslations  `yaml:"identity"`
	ScenarioHeader      string                `yaml:"scenario_header"`
	State               StateTranslations     `yaml:"state"`
	Distortions         map[string]string     `yaml:"distortions"`
	Prompts             ConsciousnessPrompts  `yaml:"prompts"`
	EmotionalAnnotation string                `yaml:"emotional_annotation"`
}

// IdentityTranslations holds labels for the identity section of the system prompt.
type IdentityTranslations struct {
	Header       string `yaml:"header"`
	Tendencies   string `yaml:"tendencies"`
	Relationships string `yaml:"relationships"`
	Patterns     string `yaml:"patterns"`
	Values       string `yaml:"values"`
	Memories     string `yaml:"memories"`
}

// StateTranslations maps psychological affect dimensions to felt-experience descriptions.
type StateTranslations struct {
	ArousalHigh    string `yaml:"arousal_high"`
	ArousalMedium  string `yaml:"arousal_medium"`
	ArousalLow     string `yaml:"arousal_low"`
	ArousalVeryLow string `yaml:"arousal_very_low"`

	ValenceVeryPositive string `yaml:"valence_very_positive"`
	ValenceNeutral      string `yaml:"valence_neutral"`
	ValenceSlightNeg    string `yaml:"valence_slight_neg"`
	ValenceNegative     string `yaml:"valence_negative"`
	ValenceVeryNegative string `yaml:"valence_very_negative"`

	EnergyHigh    string `yaml:"energy_high"`
	EnergyMedium  string `yaml:"energy_medium"`
	EnergyLow     string `yaml:"energy_low"`
	EnergyVeryLow string `yaml:"energy_very_low"`

	CognitiveLoadHigh   string `yaml:"cognitive_load_high"`
	CognitiveLoadMedium string `yaml:"cognitive_load_medium"`

	RegulationLow     string `yaml:"regulation_low"`
	RegulationVeryLow string `yaml:"regulation_very_low"`

	IsolationBoredom       string `yaml:"isolation_boredom"`
	IsolationLoneliness    string `yaml:"isolation_loneliness"`
	IsolationSignificant   string `yaml:"isolation_significant"`
	IsolationDestabilizing string `yaml:"isolation_destabilizing"`
	IsolationSevere        string `yaml:"isolation_severe"`

	CurrentExperience string `yaml:"current_experience"`
	RecentThoughts    string `yaml:"recent_thoughts"`
	RecentExperiences string `yaml:"recent_experiences"`
}

// ConsciousnessPrompts holds the various prompt fragments used by the consciousness layer.
type ConsciousnessPrompts struct {
	ReactiveQuestion    string `yaml:"reactive_question"`
	SpontaneousQuestion string `yaml:"spontaneous_question"`
	ExternalQuestion    string `yaml:"external_question"`
	TriggerShifted      string `yaml:"trigger_shifted"`
	MindTurns           string `yaml:"mind_turns"`
	SpeechFraming       string `yaml:"speech_framing"`
	ActionFraming       string `yaml:"action_framing"`
	DistortionPrefix    string `yaml:"distortion_prefix"`
}

// FeedbackTranslations holds keyword lists for detecting distortions and
// coping strategies in LLM output.
type FeedbackTranslations struct {
	Distortions map[string][]string `yaml:"distortions"`
	Coping      map[string][]string `yaml:"coping"`
}

// SenseTranslations holds sensory parser keyword rules and descriptions.
type SenseTranslations struct {
	Keywords          SenseKeywords          `yaml:"keywords"`
	DescriptionGroups SenseDescriptionGroups `yaml:"description_groups"`
	Descriptions      SenseDescriptions      `yaml:"descriptions"`
	Fallback          string                 `yaml:"fallback"`
}

// KeywordEntry pairs a phrase with an intensity value.
type KeywordEntry struct {
	Phrase    string  `yaml:"phrase"`
	Intensity float64 `yaml:"intensity"`
}

// SenseKeywords holds keyword entries per sensory channel.
type SenseKeywords struct {
	Thermal       []KeywordEntry `yaml:"thermal"`
	Pain          []KeywordEntry `yaml:"pain"`
	Auditory      []KeywordEntry `yaml:"auditory"`
	Visual        []KeywordEntry `yaml:"visual"`
	Tactile       []KeywordEntry `yaml:"tactile"`
	Olfactory     []KeywordEntry `yaml:"olfactory"`
	Gustatory     []KeywordEntry `yaml:"gustatory"`
	Interoceptive []KeywordEntry `yaml:"interoceptive"`
	Vestibular    []KeywordEntry `yaml:"vestibular"`
}

// DescriptionGroup maps a description key to the keywords that trigger it.
type DescriptionGroup struct {
	Phrases []string `yaml:"phrases"`
}

// SenseDescriptionGroups holds keyword groups that determine which description
// to use for each sensory channel's sub-categories.
type SenseDescriptionGroups struct {
	ThermalCold       []string `yaml:"thermal_cold"`
	AuditoryStartling []string `yaml:"auditory_startling"`
	AuditoryQuiet     []string `yaml:"auditory_quiet"`
	VisualDark        []string `yaml:"visual_dark"`
	VisualThreat      []string `yaml:"visual_threat"`
	TactileViolent    []string `yaml:"tactile_violent"`
	TactileGentle     []string `yaml:"tactile_gentle"`
	OlfactoryUnpleasant []string `yaml:"olfactory_unpleasant"`
	GustatoryEating   []string `yaml:"gustatory_eating"`
	GustatoryDrinking []string `yaml:"gustatory_drinking"`
	InteroBreathing   []string `yaml:"intero_breathing"`
	InteroDizzy       []string `yaml:"intero_dizzy"`
	InteroGastro      []string `yaml:"intero_gastro"`
	VestibularFalling []string `yaml:"vestibular_falling"`
	VestibularSpinning []string `yaml:"vestibular_spinning"`
}

// SenseDescriptions maps description keys to translated text.
// Used by the sensory parser to generate human-readable event descriptions.
type SenseDescriptions struct {
	FeelingCold           string `yaml:"feeling_cold"`
	FeelingHeat           string `yaml:"feeling_heat"`
	ExperiencingPain      string `yaml:"experiencing_pain"`
	HearingStartling      string `yaml:"hearing_startling"`
	HearingQuiet          string `yaml:"hearing_quiet"`
	HearingSound          string `yaml:"hearing_sound"`
	VisualDarkness        string `yaml:"visual_darkness"`
	SeeingThreatening     string `yaml:"seeing_threatening"`
	VisualPerception      string `yaml:"visual_perception"`
	PhysicallyStruck      string `yaml:"physically_struck"`
	GentleContact         string `yaml:"gentle_contact"`
	TactileSensation      string `yaml:"tactile_sensation"`
	SmellingUnpleasant    string `yaml:"smelling_unpleasant"`
	DetectingSmell        string `yaml:"detecting_smell"`
	TastingEating         string `yaml:"tasting_eating"`
	Drinking              string `yaml:"drinking"`
	GustatorySensation    string `yaml:"gustatory_sensation"`
	DifficultyBreathing   string `yaml:"difficulty_breathing"`
	FeelingDizzy          string `yaml:"feeling_dizzy"`
	GastrointestinalDistress string `yaml:"gastrointestinal_distress"`
	InternalBodySensation string `yaml:"internal_body_sensation"`
	SensationFalling      string `yaml:"sensation_falling"`
	RotationalDisorientation string `yaml:"rotational_disorientation"`
	BalanceDisruption     string `yaml:"balance_disruption"`
}

// ReviewerTranslations holds prompts for the psychologist reviewer.
type ReviewerTranslations struct {
	SystemPrompt string         `yaml:"system_prompt"`
	Labels       ReviewerLabels `yaml:"labels"`
}

// ReviewerLabels holds section headers for the reviewer prompt.
type ReviewerLabels struct {
	CurrentState    string `yaml:"current_state"`
	Arousal         string `yaml:"arousal"`
	Valence         string `yaml:"valence"`
	Energy          string `yaml:"energy"`
	CognitiveLoad   string `yaml:"cognitive_load"`
	Regulation      string `yaml:"regulation"`
	Distortions     string `yaml:"distortions"`
	CopingStrategies string `yaml:"coping_strategies"`
	PersonalityProfile string `yaml:"personality_profile"`
	RecentThoughts  string `yaml:"recent_thoughts"`
	NoThoughts      string `yaml:"no_thoughts"`
	AnalysisQuestion string `yaml:"analysis_question"`
}

// CLITranslations holds messages shown in the CLI/simulation loop.
type CLITranslations struct {
	FreshState         string `yaml:"fresh_state"`
	ResumedState       string `yaml:"resumed_state"`
	CreatedPersonality string `yaml:"created_personality"`
	CreatedIdentity    string `yaml:"created_identity"`
	SimulationStarted  string `yaml:"simulation_started"`
	UsingModel         string `yaml:"using_model"`
	ShuttingDown       string `yaml:"shutting_down"`
	WebSocketListening string `yaml:"websocket_listening"`
	Speech             string `yaml:"speech"`
	ScenarioLoaded     string `yaml:"scenario_loaded"`
	ScenarioUpdated    string `yaml:"scenario_updated"`
}

// ClientTranslations holds TUI text.
type ClientTranslations struct {
	PlaceholderSpeech      string `yaml:"placeholder_speech"`
	PlaceholderAction      string `yaml:"placeholder_action"`
	PlaceholderEnvironment string `yaml:"placeholder_environment"`
	Connecting             string `yaml:"connecting"`
	ConnectionError        string `yaml:"connection_error"`
	ConnectionClosed       string `yaml:"connection_closed"`
}

// OutputTranslations holds display labels.
type OutputTranslations struct {
	SourceLabels SourceLabels `yaml:"source_labels"`
	Unknown      string       `yaml:"unknown"`
}

// SourceLabels maps source types to their display labels.
type SourceLabels struct {
	Sense  string `yaml:"sense"`
	Bio    string `yaml:"bio"`
	Psych  string `yaml:"psych"`
	Mind   string `yaml:"mind"`
	Review string `yaml:"review"`
}

// BiologyTranslations holds display names for biology layer output.
type BiologyTranslations struct {
	Variables  map[string]string `yaml:"variables"`
	Units      map[string]string `yaml:"units"`
	Sources    map[string]string `yaml:"sources"`
	Conditions map[string]string `yaml:"conditions"`
	Systems    map[string]string `yaml:"systems"`
	Thresholds map[string]string `yaml:"thresholds"`
}

// DefaultsTranslations holds default identity values.
type DefaultsTranslations struct {
	SelfNarrative     string   `yaml:"self_narrative"`
	DispositionTraits []string `yaml:"disposition_traits"`
	EmotionalPatterns []string `yaml:"emotional_patterns"`
	ValuesCommitments []string `yaml:"values_commitments"`
}
