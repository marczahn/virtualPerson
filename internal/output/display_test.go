package output

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestDisplay_Show_PlainFormat(t *testing.T) {
	var buf bytes.Buffer
	d := NewDisplay(&buf, false)

	ts := time.Date(2024, 1, 1, 14, 30, 15, 0, time.UTC)
	d.Show(Entry{Source: Bio, Message: "heart rate elevated", Timestamp: ts})

	got := buf.String()
	if !strings.Contains(got, "14:30:15") {
		t.Errorf("expected timestamp in output, got: %s", got)
	}
	if !strings.Contains(got, "[BIO") {
		t.Errorf("expected BIO tag, got: %s", got)
	}
	if !strings.Contains(got, "heart rate elevated") {
		t.Errorf("expected message, got: %s", got)
	}
}

func TestDisplay_Show_ColorFormat(t *testing.T) {
	var buf bytes.Buffer
	d := NewDisplay(&buf, true)

	ts := time.Date(2024, 1, 1, 14, 30, 15, 0, time.UTC)
	d.Show(Entry{Source: Sense, Message: "cold wind detected", Timestamp: ts})

	got := buf.String()
	if !strings.Contains(got, colorCyan) {
		t.Errorf("expected cyan color code for SENSE, got: %s", got)
	}
	if !strings.Contains(got, "cold wind detected") {
		t.Errorf("expected message, got: %s", got)
	}
}

func TestDisplay_Show_AllSources(t *testing.T) {
	sources := []struct {
		source Source
		tag    string
	}{
		{Sense, "SENSE"},
		{Bio, "BIO"},
		{Psych, "PSYCH"},
		{Mind, "MIND"},
		{Review, "REVIEW"},
	}

	for _, s := range sources {
		var buf bytes.Buffer
		d := NewDisplay(&buf, false)
		d.Show(Entry{Source: s.source, Message: "test", Timestamp: time.Now()})
		if !strings.Contains(buf.String(), s.tag) {
			t.Errorf("expected tag %q for source %d, got: %s", s.tag, s.source, buf.String())
		}
	}
}

func TestDisplay_Show_EndsWithNewline(t *testing.T) {
	var buf bytes.Buffer
	d := NewDisplay(&buf, false)
	d.Show(Entry{Source: Mind, Message: "thinking...", Timestamp: time.Now()})

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("expected newline at end, got: %q", buf.String())
	}
}

func TestDisplay_ShowThought_PlainFormat(t *testing.T) {
	var buf bytes.Buffer
	d := NewDisplay(&buf, false)

	ts := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	d.ShowThought(Entry{Source: Mind, Message: "I feel uneasy...", Timestamp: ts})

	got := buf.String()
	if !strings.Contains(got, "10:00:00") {
		t.Errorf("expected timestamp, got: %s", got)
	}
	if !strings.Contains(got, "MIND") {
		t.Errorf("expected MIND tag, got: %s", got)
	}
	if !strings.Contains(got, "I feel uneasy") {
		t.Errorf("expected thought content, got: %s", got)
	}
}

func TestDisplay_ConcurrentWrites(t *testing.T) {
	var buf bytes.Buffer
	d := NewDisplay(&buf, false)

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			d.Show(Entry{Source: Bio, Message: "test", Timestamp: time.Now()})
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 10 {
		t.Errorf("expected 10 lines from concurrent writes, got %d", len(lines))
	}
}
