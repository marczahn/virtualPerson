package memory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/psychology"
)

// SQLiteStore implements Store using a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens or creates a SQLite database at the given path
// and initializes the schema.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS bio_state (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		body_temp REAL, heart_rate REAL, blood_pressure REAL, respiratory_rate REAL,
		hunger REAL, thirst REAL, fatigue REAL, pain REAL, muscle_tension REAL,
		blood_sugar REAL, cortisol REAL, adrenaline REAL, serotonin REAL,
		dopamine REAL, immune_response REAL, circadian_phase REAL,
		spo2 REAL, hydration REAL, glycogen REAL, endorphins REAL,
		cortisol_load REAL,
		last_update TEXT
	);

	CREATE TABLE IF NOT EXISTS personality (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		openness REAL, conscientiousness REAL, extraversion REAL,
		agreeableness REAL, neuroticism REAL
	);

	CREATE TABLE IF NOT EXISTS identity_core (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		self_narrative TEXT,
		disposition_traits TEXT,
		relational_markers TEXT,
		key_memories TEXT,
		emotional_patterns TEXT,
		values_commitments TEXT,
		last_updated TEXT
	);

	CREATE TABLE IF NOT EXISTS episodic_memories (
		id TEXT PRIMARY KEY,
		content TEXT,
		timestamp TEXT,
		emotional_valence REAL,
		importance REAL,
		bio_arousal REAL,
		bio_valence REAL,
		bio_body_temp REAL,
		bio_pain REAL,
		bio_fatigue REAL,
		bio_hunger REAL
	);

	CREATE TABLE IF NOT EXISTS emotional_memories (
		id TEXT PRIMARY KEY,
		stimulus TEXT,
		valence REAL,
		intensity REAL,
		created_at TEXT,
		traumatic INTEGER
	);
	`
	_, err := db.Exec(schema)
	return err
}

func (s *SQLiteStore) SaveBioState(st *biology.State) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO bio_state (
			id, body_temp, heart_rate, blood_pressure, respiratory_rate,
			hunger, thirst, fatigue, pain, muscle_tension,
			blood_sugar, cortisol, adrenaline, serotonin,
			dopamine, immune_response, circadian_phase,
			spo2, hydration, glycogen, endorphins,
			cortisol_load, last_update
		) VALUES (
			1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		st.BodyTemp, st.HeartRate, st.BloodPressure, st.RespiratoryRate,
		st.Hunger, st.Thirst, st.Fatigue, st.Pain, st.MuscleTension,
		st.BloodSugar, st.Cortisol, st.Adrenaline, st.Serotonin,
		st.Dopamine, st.ImmuneResponse, st.CircadianPhase,
		st.SpO2, st.Hydration, st.Glycogen, st.Endorphins,
		st.CortisolLoad, st.LastUpdate.Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteStore) LoadBioState() (*biology.State, error) {
	var st biology.State
	var lastUpdate string

	err := s.db.QueryRow(`
		SELECT body_temp, heart_rate, blood_pressure, respiratory_rate,
			hunger, thirst, fatigue, pain, muscle_tension,
			blood_sugar, cortisol, adrenaline, serotonin,
			dopamine, immune_response, circadian_phase,
			spo2, hydration, glycogen, endorphins,
			cortisol_load, last_update
		FROM bio_state WHERE id = 1
	`).Scan(
		&st.BodyTemp, &st.HeartRate, &st.BloodPressure, &st.RespiratoryRate,
		&st.Hunger, &st.Thirst, &st.Fatigue, &st.Pain, &st.MuscleTension,
		&st.BloodSugar, &st.Cortisol, &st.Adrenaline, &st.Serotonin,
		&st.Dopamine, &st.ImmuneResponse, &st.CircadianPhase,
		&st.SpO2, &st.Hydration, &st.Glycogen, &st.Endorphins,
		&st.CortisolLoad, &lastUpdate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	st.LastUpdate, err = time.Parse(time.RFC3339Nano, lastUpdate)
	if err != nil {
		return nil, fmt.Errorf("parse last_update: %w", err)
	}

	return &st, nil
}

func (s *SQLiteStore) SavePersonality(p *psychology.Personality) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO personality (id, openness, conscientiousness, extraversion, agreeableness, neuroticism)
		VALUES (1, ?, ?, ?, ?, ?)`,
		p.Openness, p.Conscientiousness, p.Extraversion, p.Agreeableness, p.Neuroticism,
	)
	return err
}

func (s *SQLiteStore) LoadPersonality() (*psychology.Personality, error) {
	var p psychology.Personality
	err := s.db.QueryRow(`
		SELECT openness, conscientiousness, extraversion, agreeableness, neuroticism
		FROM personality WHERE id = 1
	`).Scan(&p.Openness, &p.Conscientiousness, &p.Extraversion, &p.Agreeableness, &p.Neuroticism)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *SQLiteStore) SaveIdentityCore(ic *IdentityCore) error {
	disposition, _ := json.Marshal(ic.DispositionTraits)
	relational, _ := json.Marshal(ic.RelationalMarkers)
	keyMem, _ := json.Marshal(ic.KeyMemories)
	emotional, _ := json.Marshal(ic.EmotionalPatterns)
	values, _ := json.Marshal(ic.ValuesCommitments)

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO identity_core (
			id, self_narrative, disposition_traits, relational_markers,
			key_memories, emotional_patterns, values_commitments, last_updated
		) VALUES (1, ?, ?, ?, ?, ?, ?, ?)`,
		ic.SelfNarrative, string(disposition), string(relational),
		string(keyMem), string(emotional), string(values),
		ic.LastUpdated.Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteStore) LoadIdentityCore() (*IdentityCore, error) {
	var ic IdentityCore
	var disposition, relational, keyMem, emotional, values, lastUpdated string

	err := s.db.QueryRow(`
		SELECT self_narrative, disposition_traits, relational_markers,
			key_memories, emotional_patterns, values_commitments, last_updated
		FROM identity_core WHERE id = 1
	`).Scan(
		&ic.SelfNarrative, &disposition, &relational,
		&keyMem, &emotional, &values, &lastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(disposition), &ic.DispositionTraits)
	json.Unmarshal([]byte(relational), &ic.RelationalMarkers)
	json.Unmarshal([]byte(keyMem), &ic.KeyMemories)
	json.Unmarshal([]byte(emotional), &ic.EmotionalPatterns)
	json.Unmarshal([]byte(values), &ic.ValuesCommitments)

	ic.LastUpdated, err = time.Parse(time.RFC3339Nano, lastUpdated)
	if err != nil {
		return nil, fmt.Errorf("parse last_updated: %w", err)
	}

	return &ic, nil
}

func (s *SQLiteStore) SaveMemory(m *EpisodicMemory) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO episodic_memories (
			id, content, timestamp, emotional_valence, importance,
			bio_arousal, bio_valence, bio_body_temp, bio_pain, bio_fatigue, bio_hunger
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID, m.Content, m.Timestamp.Format(time.RFC3339Nano),
		m.EmotionalValence, m.Importance,
		m.BioSnapshot.Arousal, m.BioSnapshot.Valence,
		m.BioSnapshot.BodyTemp, m.BioSnapshot.Pain,
		m.BioSnapshot.Fatigue, m.BioSnapshot.Hunger,
	)
	return err
}

func (s *SQLiteStore) LoadMemories() ([]EpisodicMemory, error) {
	rows, err := s.db.Query(`
		SELECT id, content, timestamp, emotional_valence, importance,
			bio_arousal, bio_valence, bio_body_temp, bio_pain, bio_fatigue, bio_hunger
		FROM episodic_memories ORDER BY timestamp ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []EpisodicMemory
	for rows.Next() {
		var m EpisodicMemory
		var ts string
		err := rows.Scan(
			&m.ID, &m.Content, &ts,
			&m.EmotionalValence, &m.Importance,
			&m.BioSnapshot.Arousal, &m.BioSnapshot.Valence,
			&m.BioSnapshot.BodyTemp, &m.BioSnapshot.Pain,
			&m.BioSnapshot.Fatigue, &m.BioSnapshot.Hunger,
		)
		if err != nil {
			return nil, err
		}
		m.Timestamp, _ = time.Parse(time.RFC3339Nano, ts)
		memories = append(memories, m)
	}
	return memories, rows.Err()
}

func (s *SQLiteStore) SaveEmotionalMemory(m *psychology.EmotionalMemory) error {
	traumatic := 0
	if m.Traumatic {
		traumatic = 1
	}
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO emotional_memories (id, stimulus, valence, intensity, created_at, traumatic)
		VALUES (?, ?, ?, ?, ?, ?)`,
		m.ID, m.Stimulus, m.Valence, m.Intensity,
		m.CreatedAt.Format(time.RFC3339Nano), traumatic,
	)
	return err
}

func (s *SQLiteStore) LoadEmotionalMemories() ([]psychology.EmotionalMemory, error) {
	rows, err := s.db.Query(`
		SELECT id, stimulus, valence, intensity, created_at, traumatic
		FROM emotional_memories ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []psychology.EmotionalMemory
	for rows.Next() {
		var m psychology.EmotionalMemory
		var createdAt string
		var traumatic int
		err := rows.Scan(&m.ID, &m.Stimulus, &m.Valence, &m.Intensity, &createdAt, &traumatic)
		if err != nil {
			return nil, err
		}
		m.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		m.Traumatic = traumatic != 0
		memories = append(memories, m)
	}
	return memories, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
