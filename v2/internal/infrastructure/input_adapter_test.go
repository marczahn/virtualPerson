package infrastructure_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/infrastructure"
	"github.com/marczahn/person/v2/internal/sense"
)

func TestInputAdapter_DrainNoInputNoOp(t *testing.T) {
	adapter := infrastructure.NewInputAdapter(sense.NewParser(), func() int64 { return 77 })

	got := adapter.Drain()

	if got.ExternalText != "" {
		t.Fatalf("expected empty external text, got %q", got.ExternalText)
	}
	if len(got.PreBioRates) != 0 {
		t.Fatalf("expected no rates, got %d", len(got.PreBioRates))
	}
	if len(got.PreBioPulses) != 0 {
		t.Fatalf("expected no pulses, got %d", len(got.PreBioPulses))
	}
	if !got.AllowedActions["eat"] {
		t.Fatalf("expected default eat action to be allowed")
	}
	if got.NowSeconds != 77 {
		t.Fatalf("unexpected now seconds: got=%d want=77", got.NowSeconds)
	}
}

func TestInputAdapter_DrainParsesSpeechActionEnvironment(t *testing.T) {
	adapter := infrastructure.NewInputAdapter(sense.NewParser(), func() int64 { return 123 })
	adapter.Enqueue("hello operator")
	adapter.Enqueue("*someone punches you*")
	adapter.Enqueue("~cold room, no food available")

	got := adapter.Drain()

	wantText := "hello operator\n*someone punches you*\n~cold room, no food available"
	if got.ExternalText != wantText {
		t.Fatalf("unexpected external text: got=%q want=%q", got.ExternalText, wantText)
	}
	if len(got.PreBioPulses) == 0 {
		t.Fatalf("expected action pulse effects from punch input")
	}
	if len(got.PreBioRates) == 0 {
		t.Fatalf("expected environment rate effects from cold input")
	}
	if got.AllowedActions["eat"] {
		t.Fatalf("expected eat action blocked by 'no food' environment input")
	}
}

func TestInputAdapter_DrainDeterministicOrderingWithinCycle(t *testing.T) {
	adapter := infrastructure.NewInputAdapter(sense.NewParser(), func() int64 { return 9 })
	adapter.Enqueue("~no food available")
	adapter.Enqueue("~food is available")

	got := adapter.Drain()
	if !got.AllowedActions["eat"] {
		t.Fatalf("expected later environment input to deterministically override eat gate to allowed")
	}

	adapter.Enqueue("~food is available")
	adapter.Enqueue("~no food available")
	got = adapter.Drain()
	if got.AllowedActions["eat"] {
		t.Fatalf("expected later environment input to deterministically override eat gate to blocked")
	}
}
