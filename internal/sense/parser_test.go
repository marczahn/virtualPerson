package sense

import (
	"testing"
)

func TestKeywordParser_EmptyInput_ReturnsNil(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("")
	if events != nil {
		t.Errorf("expected nil for empty input, got %d events", len(events))
	}
}

func TestKeywordParser_WhitespaceOnly_ReturnsNil(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("   \t  ")
	if events != nil {
		t.Errorf("expected nil for whitespace-only input, got %d events", len(events))
	}
}

func TestKeywordParser_UnrecognizedInput_DefaultsToAuditory(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("hello there, how are you?")

	if len(events) != 1 {
		t.Fatalf("expected 1 default event, got %d", len(events))
	}
	if events[0].Channel != Auditory {
		t.Errorf("expected Auditory, got %s", events[0].Channel)
	}
	if events[0].Intensity != 0.3 {
		t.Errorf("expected default intensity 0.3, got %f", events[0].Intensity)
	}
	if events[0].RawInput != "hello there, how are you?" {
		t.Errorf("RawInput not preserved: %q", events[0].RawInput)
	}
}

func TestKeywordParser_ColdThermal(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("a freezing wind blows")

	found := findChannel(events, Thermal)
	if found == nil {
		t.Fatal("expected Thermal event for 'freezing'")
	}
	if found.Intensity >= 0.5 {
		t.Errorf("cold stimulus should have intensity < 0.5, got %f", found.Intensity)
	}
}

func TestKeywordParser_HotThermal(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("the room is scorching hot")

	found := findChannel(events, Thermal)
	if found == nil {
		t.Fatal("expected Thermal event for 'scorching hot'")
	}
	// "scorching" = 0.9, "hot" = 0.75; max should be 0.9
	if found.Intensity < 0.5 {
		t.Errorf("hot stimulus should have intensity >= 0.5, got %f", found.Intensity)
	}
}

func TestKeywordParser_Pain(t *testing.T) {
	tests := []struct {
		input       string
		minIntensity float64
	}{
		{"a dull aching in the shoulder", 0.3},
		{"excruciating pain shoots through the leg", 0.9},
		{"a mild headache", 0.4},
		{"sharp pain in the chest", 0.8},
	}

	p := NewKeywordParser()
	for _, tt := range tests {
		events := p.Parse(tt.input)
		found := findChannel(events, Pain)
		if found == nil {
			t.Errorf("expected Pain event for %q", tt.input)
			continue
		}
		if found.Intensity < tt.minIntensity {
			t.Errorf("%q: expected intensity >= %f, got %f", tt.input, tt.minIntensity, found.Intensity)
		}
	}
}

func TestKeywordParser_LoudAuditory(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("a massive explosion nearby")

	found := findChannel(events, Auditory)
	if found == nil {
		t.Fatal("expected Auditory event for 'explosion'")
	}
	if found.Intensity < 0.9 {
		t.Errorf("explosion should be high intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_QuietAuditory(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("a whisper in the dark")

	found := findChannel(events, Auditory)
	if found == nil {
		t.Fatal("expected Auditory event for 'whisper'")
	}
	if found.Intensity > 0.3 {
		t.Errorf("whisper should be low intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_VisualThreat(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("you see someone with a knife")

	found := findChannel(events, Visual)
	if found == nil {
		t.Fatal("expected Visual event for 'knife'")
	}
	if found.Intensity < 0.8 {
		t.Errorf("knife should be high intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_Tactile(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("someone punches you hard")

	found := findChannel(events, Tactile)
	if found == nil {
		t.Fatal("expected Tactile event for 'punch'")
	}
	if found.Intensity < 0.7 {
		t.Errorf("punch should be high intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_MultipleChannels(t *testing.T) {
	p := NewKeywordParser()
	// This should trigger thermal (freezing) AND pain (aching).
	events := p.Parse("freezing wind causes aching joints")

	thermal := findChannel(events, Thermal)
	pain := findChannel(events, Pain)

	if thermal == nil {
		t.Error("expected Thermal event")
	}
	if pain == nil {
		t.Error("expected Pain event")
	}
}

func TestKeywordParser_Olfactory(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("a terrible stench fills the room")

	found := findChannel(events, Olfactory)
	if found == nil {
		t.Fatal("expected Olfactory event for 'stench'")
	}
	if found.Intensity < 0.6 {
		t.Errorf("stench should be moderate-high intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_Gustatory_Eating(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("eating some delicious food")

	found := findChannel(events, Gustatory)
	if found == nil {
		t.Fatal("expected Gustatory event for 'eating food'")
	}
}

func TestKeywordParser_Gustatory_Drinking(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("drinking cold water")

	found := findChannel(events, Gustatory)
	if found == nil {
		t.Fatal("expected Gustatory event for 'drinking water'")
	}
}

func TestKeywordParser_Interoceptive(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("can't breathe, suffocating")

	found := findChannel(events, Interoceptive)
	if found == nil {
		t.Fatal("expected Interoceptive event for 'suffocating'")
	}
	if found.Intensity < 0.8 {
		t.Errorf("suffocating should be high intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_Vestibular(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("the ground shakes violently, earthquake")

	found := findChannel(events, Vestibular)
	if found == nil {
		t.Fatal("expected Vestibular event for 'earthquake'")
	}
	if found.Intensity < 0.8 {
		t.Errorf("earthquake should be high intensity, got %f", found.Intensity)
	}
}

func TestKeywordParser_CaseInsensitive(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("A LOUD EXPLOSION ROCKS THE BUILDING")

	found := findChannel(events, Auditory)
	if found == nil {
		t.Fatal("expected Auditory event (case insensitive)")
	}
	if found.Intensity < 0.7 {
		t.Errorf("expected high intensity for LOUD EXPLOSION, got %f", found.Intensity)
	}
}

func TestKeywordParser_MaxIntensityWins(t *testing.T) {
	p := NewKeywordParser()
	// "excruciating" (0.95) and "pain" (0.5) both match pain channel.
	// Max should win.
	events := p.Parse("excruciating pain")

	found := findChannel(events, Pain)
	if found == nil {
		t.Fatal("expected Pain event")
	}
	if found.Intensity < 0.9 {
		t.Errorf("max keyword intensity should win, got %f", found.Intensity)
	}
}

func TestKeywordParser_TimestampSet(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("it's cold")

	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if events[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestKeywordParser_DescriptionSet(t *testing.T) {
	p := NewKeywordParser()
	events := p.Parse("freezing cold outside")

	found := findChannel(events, Thermal)
	if found == nil {
		t.Fatal("expected Thermal event")
	}
	if found.Parsed == "" {
		t.Error("expected non-empty Parsed description")
	}
}

func TestKeywordParser_ImplementsInterface(t *testing.T) {
	var _ Parser = NewKeywordParser()
}

func findChannel(events []Event, ch Channel) *Event {
	for i := range events {
		if events[i].Channel == ch {
			return &events[i]
		}
	}
	return nil
}
