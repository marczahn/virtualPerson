package i18n

import (
	"testing"
)

func TestLoad_English(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("failed to load English: %v", err)
	}

	tr := T()
	if tr.Consciousness.SystemPrompt == "" {
		t.Error("English consciousness system prompt is empty")
	}
	if tr.Sense.Fallback == "" {
		t.Error("English sense fallback is empty")
	}
	if len(tr.Sense.Keywords.Thermal) == 0 {
		t.Error("English thermal keywords are empty")
	}
	if tr.Defaults.SelfNarrative == "" {
		t.Error("English defaults self_narrative is empty")
	}
}

func TestLoad_German(t *testing.T) {
	if err := Load("de"); err != nil {
		t.Fatalf("failed to load German: %v", err)
	}

	tr := T()
	if tr.Consciousness.SystemPrompt == "" {
		t.Error("German consciousness system prompt is empty")
	}
	if tr.Sense.Fallback == "" {
		t.Error("German sense fallback is empty")
	}
	if len(tr.Sense.Keywords.Thermal) == 0 {
		t.Error("German thermal keywords are empty")
	}
	if tr.Defaults.SelfNarrative == "" {
		t.Error("German defaults self_narrative is empty")
	}
}

func TestLoad_InvalidLanguage(t *testing.T) {
	err := Load("xx")
	if err == nil {
		t.Error("expected error for invalid language, got nil")
	}
}

func TestAvailable_ContainsExpectedLanguages(t *testing.T) {
	langs := Available()
	found := map[string]bool{}
	for _, l := range langs {
		found[l] = true
	}
	if !found["en"] {
		t.Error("expected 'en' in available languages")
	}
	if !found["de"] {
		t.Error("expected 'de' in available languages")
	}
}

func TestLang_ReturnsLoadedLanguage(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	if got := Lang(); got != "en" {
		t.Errorf("Lang() = %q, want %q", got, "en")
	}
}

func TestEnglish_FeedbackKeywordsPresent(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	tr := T()

	if len(tr.Feedback.Distortions["catastrophizing"]) == 0 {
		t.Error("expected catastrophizing distortion keywords")
	}
	if len(tr.Feedback.Coping["rumination"]) == 0 {
		t.Error("expected rumination coping keywords")
	}
}

func TestEnglish_ConsciousnessDistortionsPresent(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	tr := T()

	expected := []string{"catastrophizing", "emotional_reasoning", "overgeneralization", "mind_reading", "personalization", "all_or_nothing"}
	for _, key := range expected {
		if tr.Consciousness.Distortions[key] == "" {
			t.Errorf("missing distortion description for %q", key)
		}
	}
}

func TestEnglish_AllSenseChannelsHaveKeywords(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	tr := T()

	if len(tr.Sense.Keywords.Pain) == 0 {
		t.Error("pain keywords empty")
	}
	if len(tr.Sense.Keywords.Auditory) == 0 {
		t.Error("auditory keywords empty")
	}
	if len(tr.Sense.Keywords.Visual) == 0 {
		t.Error("visual keywords empty")
	}
	if len(tr.Sense.Keywords.Tactile) == 0 {
		t.Error("tactile keywords empty")
	}
	if len(tr.Sense.Keywords.Olfactory) == 0 {
		t.Error("olfactory keywords empty")
	}
	if len(tr.Sense.Keywords.Gustatory) == 0 {
		t.Error("gustatory keywords empty")
	}
	if len(tr.Sense.Keywords.Interoceptive) == 0 {
		t.Error("interoceptive keywords empty")
	}
	if len(tr.Sense.Keywords.Vestibular) == 0 {
		t.Error("vestibular keywords empty")
	}
}

func TestTranslateBio_Found(t *testing.T) {
	m := map[string]string{"heart_rate": "Herzfrequenz"}
	got := TranslateBio(m, "heart_rate")
	if got != "Herzfrequenz" {
		t.Errorf("TranslateBio() = %q, want %q", got, "Herzfrequenz")
	}
}

func TestTranslateBio_Fallback(t *testing.T) {
	m := map[string]string{}
	got := TranslateBio(m, "unknown_key")
	if got != "unknown_key" {
		t.Errorf("TranslateBio() = %q, want %q", got, "unknown_key")
	}
}

func TestBiology_VariablesLoaded(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	tr := T()

	if len(tr.Biology.Variables) < 20 {
		t.Errorf("expected at least 20 biology variables, got %d", len(tr.Biology.Variables))
	}
	if len(tr.Biology.Sources) < 10 {
		t.Errorf("expected at least 10 biology sources, got %d", len(tr.Biology.Sources))
	}
	if len(tr.Biology.Conditions) < 5 {
		t.Errorf("expected at least 5 biology conditions, got %d", len(tr.Biology.Conditions))
	}
	if len(tr.Biology.Systems) < 4 {
		t.Errorf("expected at least 4 biology systems, got %d", len(tr.Biology.Systems))
	}
	if len(tr.Biology.Thresholds) < 18 {
		t.Errorf("expected at least 18 biology thresholds, got %d", len(tr.Biology.Thresholds))
	}
}

func TestGerman_BiologyLoaded(t *testing.T) {
	if err := Load("de"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	tr := T()

	if tr.Biology.Variables["heart_rate"] != "Herzfrequenz" {
		t.Errorf("German heart_rate = %q, want %q", tr.Biology.Variables["heart_rate"], "Herzfrequenz")
	}
	if tr.Biology.Conditions["impaired"] != "beeinträchtigt" {
		t.Errorf("German impaired = %q, want %q", tr.Biology.Conditions["impaired"], "beeinträchtigt")
	}
}

func TestGerman_HasMatchingStructure(t *testing.T) {
	if err := Load("de"); err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	tr := T()

	if len(tr.Feedback.Distortions["catastrophizing"]) == 0 {
		t.Error("German: expected catastrophizing distortion keywords")
	}
	if len(tr.Feedback.Coping["rumination"]) == 0 {
		t.Error("German: expected rumination coping keywords")
	}
	if len(tr.Consciousness.Distortions) != 6 {
		t.Errorf("German: expected 6 distortion descriptions, got %d", len(tr.Consciousness.Distortions))
	}
	if tr.Output.SourceLabels.Sense == "" {
		t.Error("German: missing sense source label")
	}
}
