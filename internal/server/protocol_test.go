package server

import "testing"

func TestClientMessage_Validate_ValidTypes(t *testing.T) {
	for _, typ := range []string{"speech", "action", "environment"} {
		msg := ClientMessage{Type: typ, Content: "hello"}
		if err := msg.Validate(); err != nil {
			t.Errorf("Validate(%q) returned error: %v", typ, err)
		}
	}
}

func TestClientMessage_Validate_UnknownType(t *testing.T) {
	msg := ClientMessage{Type: "shout", Content: "hello"}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestClientMessage_Validate_EmptyContent(t *testing.T) {
	msg := ClientMessage{Type: "speech", Content: ""}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for empty content")
	}
}

func TestClientMessage_ToInputLine_Speech(t *testing.T) {
	msg := ClientMessage{Type: "speech", Content: "hello"}
	if got := msg.ToInputLine(); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestClientMessage_ToInputLine_Action(t *testing.T) {
	msg := ClientMessage{Type: "action", Content: "pushes you"}
	if got := msg.ToInputLine(); got != "*pushes you*" {
		t.Errorf("got %q, want %q", got, "*pushes you*")
	}
}

func TestClientMessage_ToInputLine_Environment(t *testing.T) {
	msg := ClientMessage{Type: "environment", Content: "cold wind"}
	if got := msg.ToInputLine(); got != "~cold wind" {
		t.Errorf("got %q, want %q", got, "~cold wind")
	}
}
