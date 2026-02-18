package memory

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/psychology"
)

func tempDB(t *testing.T) *SQLiteStore {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	store, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestSQLiteStore_BioState_RoundTrip(t *testing.T) {
	store := tempDB(t)

	original := biology.NewDefaultState()
	original.Cortisol = 0.42
	original.CortisolLoad = 3.7
	original.BodyTemp = 35.2

	if err := store.SaveBioState(&original); err != nil {
		t.Fatalf("SaveBioState: %v", err)
	}

	loaded, err := store.LoadBioState()
	if err != nil {
		t.Fatalf("LoadBioState: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadBioState returned nil")
	}

	if loaded.Cortisol != 0.42 {
		t.Errorf("Cortisol = %f, want 0.42", loaded.Cortisol)
	}
	if loaded.CortisolLoad != 3.7 {
		t.Errorf("CortisolLoad = %f, want 3.7", loaded.CortisolLoad)
	}
	if loaded.BodyTemp != 35.2 {
		t.Errorf("BodyTemp = %f, want 35.2", loaded.BodyTemp)
	}
	if loaded.HeartRate != original.HeartRate {
		t.Errorf("HeartRate = %f, want %f", loaded.HeartRate, original.HeartRate)
	}

	// Timestamp should round-trip within a microsecond.
	diff := loaded.LastUpdate.Sub(original.LastUpdate)
	if diff < -time.Microsecond || diff > time.Microsecond {
		t.Errorf("LastUpdate diff = %v, expected < 1µs", diff)
	}
}

func TestSQLiteStore_BioState_LoadEmpty(t *testing.T) {
	store := tempDB(t)

	loaded, err := store.LoadBioState()
	if err != nil {
		t.Fatalf("LoadBioState: %v", err)
	}
	if loaded != nil {
		t.Errorf("expected nil from empty store, got %v", loaded)
	}
}

func TestSQLiteStore_BioState_Overwrite(t *testing.T) {
	store := tempDB(t)

	s1 := biology.NewDefaultState()
	s1.HeartRate = 80
	store.SaveBioState(&s1)

	s2 := biology.NewDefaultState()
	s2.HeartRate = 120
	store.SaveBioState(&s2)

	loaded, _ := store.LoadBioState()
	if loaded.HeartRate != 120 {
		t.Errorf("HeartRate = %f, want 120 (latest save)", loaded.HeartRate)
	}
}

func TestSQLiteStore_Personality_RoundTrip(t *testing.T) {
	store := tempDB(t)

	original := psychology.Personality{
		Openness:          0.7,
		Conscientiousness: 0.6,
		Extraversion:      0.4,
		Agreeableness:     0.8,
		Neuroticism:       0.3,
	}

	if err := store.SavePersonality(&original); err != nil {
		t.Fatalf("SavePersonality: %v", err)
	}

	loaded, err := store.LoadPersonality()
	if err != nil {
		t.Fatalf("LoadPersonality: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadPersonality returned nil")
	}

	if *loaded != original {
		t.Errorf("personality mismatch: got %+v, want %+v", *loaded, original)
	}
}

func TestSQLiteStore_Personality_LoadEmpty(t *testing.T) {
	store := tempDB(t)

	loaded, err := store.LoadPersonality()
	if err != nil {
		t.Fatalf("LoadPersonality: %v", err)
	}
	if loaded != nil {
		t.Errorf("expected nil from empty store")
	}
}

func TestSQLiteStore_IdentityCore_RoundTrip(t *testing.T) {
	store := tempDB(t)

	original := IdentityCore{
		SelfNarrative:     "I am a thoughtful person who values learning.",
		DispositionTraits: []string{"curious", "reserved", "empathetic"},
		RelationalMarkers: []string{"close to mother", "few but deep friendships"},
		KeyMemories:       []string{"first day at university", "meeting best friend"},
		EmotionalPatterns: []string{"tends to internalize anger", "cries when overwhelmed"},
		ValuesCommitments: []string{"honesty", "loyalty", "personal growth"},
		LastUpdated:       time.Now().Truncate(time.Microsecond),
	}

	if err := store.SaveIdentityCore(&original); err != nil {
		t.Fatalf("SaveIdentityCore: %v", err)
	}

	loaded, err := store.LoadIdentityCore()
	if err != nil {
		t.Fatalf("LoadIdentityCore: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadIdentityCore returned nil")
	}

	if loaded.SelfNarrative != original.SelfNarrative {
		t.Errorf("SelfNarrative = %q, want %q", loaded.SelfNarrative, original.SelfNarrative)
	}
	if len(loaded.DispositionTraits) != 3 || loaded.DispositionTraits[0] != "curious" {
		t.Errorf("DispositionTraits = %v, want %v", loaded.DispositionTraits, original.DispositionTraits)
	}
	if len(loaded.KeyMemories) != 2 {
		t.Errorf("KeyMemories = %v, want %v", loaded.KeyMemories, original.KeyMemories)
	}
}

func TestSQLiteStore_IdentityCore_LoadEmpty(t *testing.T) {
	store := tempDB(t)

	loaded, err := store.LoadIdentityCore()
	if err != nil {
		t.Fatalf("LoadIdentityCore: %v", err)
	}
	if loaded != nil {
		t.Errorf("expected nil from empty store")
	}
}

func TestSQLiteStore_EpisodicMemory_SaveAndLoad(t *testing.T) {
	store := tempDB(t)

	m1 := EpisodicMemory{
		ID:               "mem1",
		Content:          "Felt cold wind on my face",
		Timestamp:        time.Now().Add(-time.Hour).Truncate(time.Microsecond),
		EmotionalValence: -0.4,
		Importance:       0.6,
		BioSnapshot:      BioSnapshot{Arousal: 0.3, Valence: -0.2, BodyTemp: 34.5, Pain: 0, Fatigue: 0.1, Hunger: 0},
	}
	m2 := EpisodicMemory{
		ID:               "mem2",
		Content:          "Ate a warm meal",
		Timestamp:        time.Now().Truncate(time.Microsecond),
		EmotionalValence: 0.6,
		Importance:       0.4,
		BioSnapshot:      BioSnapshot{Arousal: 0.1, Valence: 0.5, BodyTemp: 36.6, Pain: 0, Fatigue: 0.2, Hunger: 0.8},
	}

	store.SaveMemory(&m1)
	store.SaveMemory(&m2)

	loaded, err := store.LoadMemories()
	if err != nil {
		t.Fatalf("LoadMemories: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 memories, got %d", len(loaded))
	}

	// Should be ordered by timestamp ASC.
	if loaded[0].ID != "mem1" {
		t.Errorf("first memory ID = %q, want mem1", loaded[0].ID)
	}
	if loaded[1].ID != "mem2" {
		t.Errorf("second memory ID = %q, want mem2", loaded[1].ID)
	}
	if loaded[0].BioSnapshot.BodyTemp != 34.5 {
		t.Errorf("bio snapshot body_temp = %f, want 34.5", loaded[0].BioSnapshot.BodyTemp)
	}
}

func TestSQLiteStore_EmotionalMemory_SaveAndLoad(t *testing.T) {
	store := tempDB(t)

	m := psychology.EmotionalMemory{
		ID:        "em1",
		Stimulus:  "cold",
		Valence:   -0.7,
		Intensity: 0.8,
		CreatedAt: time.Now().Add(-24 * time.Hour).Truncate(time.Microsecond),
		Traumatic: true,
	}

	if err := store.SaveEmotionalMemory(&m); err != nil {
		t.Fatalf("SaveEmotionalMemory: %v", err)
	}

	loaded, err := store.LoadEmotionalMemories()
	if err != nil {
		t.Fatalf("LoadEmotionalMemories: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 emotional memory, got %d", len(loaded))
	}

	got := loaded[0]
	if got.ID != "em1" {
		t.Errorf("ID = %q, want em1", got.ID)
	}
	if got.Stimulus != "cold" {
		t.Errorf("Stimulus = %q, want cold", got.Stimulus)
	}
	if got.Valence != -0.7 {
		t.Errorf("Valence = %f, want -0.7", got.Valence)
	}
	if !got.Traumatic {
		t.Error("Traumatic = false, want true")
	}
}

func TestSQLiteStore_BioState_AllFields_RoundTrip(t *testing.T) {
	store := tempDB(t)

	original := biology.State{
		BodyTemp:        35.5,
		HeartRate:       100,
		BloodPressure:   140,
		RespiratoryRate: 20,
		Hunger:          0.3,
		Thirst:          0.4,
		Fatigue:         0.5,
		Pain:            0.2,
		MuscleTension:   0.6,
		BloodSugar:      110,
		Cortisol:        0.4,
		Adrenaline:      0.3,
		Serotonin:       0.6,
		Dopamine:        0.4,
		ImmuneResponse:  0.2,
		CircadianPhase:  14.5,
		SpO2:            95,
		Hydration:       0.6,
		Glycogen:        0.5,
		Endorphins:      0.3,
		CortisolLoad:    5.5,
		LastUpdate:      time.Now().Truncate(time.Microsecond),
	}

	store.SaveBioState(&original)
	loaded, _ := store.LoadBioState()

	fields := []struct {
		name string
		got  float64
		want float64
	}{
		{"BodyTemp", loaded.BodyTemp, original.BodyTemp},
		{"HeartRate", loaded.HeartRate, original.HeartRate},
		{"BloodPressure", loaded.BloodPressure, original.BloodPressure},
		{"RespiratoryRate", loaded.RespiratoryRate, original.RespiratoryRate},
		{"Hunger", loaded.Hunger, original.Hunger},
		{"Thirst", loaded.Thirst, original.Thirst},
		{"Fatigue", loaded.Fatigue, original.Fatigue},
		{"Pain", loaded.Pain, original.Pain},
		{"MuscleTension", loaded.MuscleTension, original.MuscleTension},
		{"BloodSugar", loaded.BloodSugar, original.BloodSugar},
		{"Cortisol", loaded.Cortisol, original.Cortisol},
		{"Adrenaline", loaded.Adrenaline, original.Adrenaline},
		{"Serotonin", loaded.Serotonin, original.Serotonin},
		{"Dopamine", loaded.Dopamine, original.Dopamine},
		{"ImmuneResponse", loaded.ImmuneResponse, original.ImmuneResponse},
		{"CircadianPhase", loaded.CircadianPhase, original.CircadianPhase},
		{"SpO2", loaded.SpO2, original.SpO2},
		{"Hydration", loaded.Hydration, original.Hydration},
		{"Glycogen", loaded.Glycogen, original.Glycogen},
		{"Endorphins", loaded.Endorphins, original.Endorphins},
		{"CortisolLoad", loaded.CortisolLoad, original.CortisolLoad},
	}

	for _, f := range fields {
		if math.Abs(f.got-f.want) > 1e-10 {
			t.Errorf("%s = %f, want %f", f.name, f.got, f.want)
		}
	}
}

func TestSQLiteStore_PersistenceAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reopen.db")

	// Open, save, close.
	store1, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("open 1: %v", err)
	}

	bio := biology.NewDefaultState()
	bio.HeartRate = 99
	store1.SaveBioState(&bio)
	store1.Close()

	// Reopen and verify.
	store2, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("open 2: %v", err)
	}
	defer store2.Close()

	loaded, _ := store2.LoadBioState()
	if loaded == nil {
		t.Fatal("expected bio state after reopen")
	}
	if loaded.HeartRate != 99 {
		t.Errorf("HeartRate = %f after reopen, want 99", loaded.HeartRate)
	}
}

func TestSQLiteStore_InvalidPath_ReturnsError(t *testing.T) {
	_, err := NewSQLiteStore("/nonexistent/dir/test.db")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

// Verify the temp db cleanup works (no leftover files).
func TestSQLiteStore_TempDir_Cleanup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cleanup.db")

	store, _ := NewSQLiteStore(path)
	store.Close()

	// Dir should still exist (TempDir cleanup happens after test).
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("temp dir should still exist during test: %v", err)
	}
}
