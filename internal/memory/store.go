package memory

import (
	"time"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/psychology"
)

// Store defines the persistence contract for the person simulation.
// All state needed to stop and resume a simulation is saved through this interface.
type Store interface {
	// SaveBioState persists the current biological state.
	SaveBioState(s *biology.State) error

	// LoadBioState restores the most recently saved biological state.
	// Returns nil state and no error if no state has been saved yet.
	LoadBioState() (*biology.State, error)

	// SavePersonality persists the person's Big Five personality traits.
	SavePersonality(p *psychology.Personality) error

	// LoadPersonality restores the saved personality.
	LoadPersonality() (*psychology.Personality, error)

	// SaveIdentityCore persists the person's identity narrative.
	SaveIdentityCore(ic *IdentityCore) error

	// LoadIdentityCore restores the saved identity core.
	LoadIdentityCore() (*IdentityCore, error)

	// SaveMemory stores a single episodic memory.
	SaveMemory(m *EpisodicMemory) error

	// LoadMemories returns all stored episodic memories.
	LoadMemories() ([]EpisodicMemory, error)

	// SaveEmotionalMemory stores an emotional association.
	SaveEmotionalMemory(m *psychology.EmotionalMemory) error

	// LoadEmotionalMemories returns all stored emotional memories.
	LoadEmotionalMemories() ([]psychology.EmotionalMemory, error)

	// Close releases any resources held by the store.
	Close() error
}

// IdentityCore holds the persistent identity narrative that is fed to
// every consciousness prompt. It contains who the person believes they are.
type IdentityCore struct {
	SelfNarrative      string   // 2-3 sentences, may be biased/idealized
	DispositionTraits  []string // behavioral tendencies with context
	RelationalMarkers  []string // relationships that define the person
	KeyMemories        []string // 3-5 key autobiographical memory summaries
	EmotionalPatterns  []string // habitual emotional responses
	ValuesCommitments  []string // deeply held values and commitments
	LastUpdated        time.Time
}

// EpisodicMemory represents a single remembered experience — what happened,
// when, and the emotional residue it left.
type EpisodicMemory struct {
	ID              string
	Content         string    // summary of the experience
	Timestamp       time.Time // when the event occurred in simulation time
	EmotionalValence float64  // -1 to 1, how the experience felt
	Importance      float64   // 0 to 1, salience for future recall
	BioSnapshot     BioSnapshot // biological state at the time of the memory
}

// BioSnapshot captures the key biological variables at the time a memory
// was formed. Used for somatic similarity retrieval.
type BioSnapshot struct {
	Arousal    float64 // derived from adrenaline/HR at time of memory
	Valence    float64 // derived from serotonin/dopamine at time of memory
	BodyTemp   float64
	Pain       float64
	Fatigue    float64
	Hunger     float64
}
