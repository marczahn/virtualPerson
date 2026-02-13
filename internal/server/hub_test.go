package server

import (
	"bytes"
	"testing"
	"time"
)

func TestHub_Broadcast_DeliversToRegisteredClients(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})

	c1 := &Client{sendCh: make(chan ServerMessage, 8)}
	c2 := &Client{sendCh: make(chan ServerMessage, 8)}
	hub.Register(c1)
	hub.Register(c2)

	msg := ServerMessage{Type: "thought", Content: "hello", Timestamp: time.Now()}
	hub.Broadcast(msg)

	select {
	case got := <-c1.sendCh:
		if got.Content != "hello" {
			t.Errorf("c1 got content %q, want %q", got.Content, "hello")
		}
	default:
		t.Error("c1 did not receive message")
	}

	select {
	case got := <-c2.sendCh:
		if got.Content != "hello" {
			t.Errorf("c2 got content %q, want %q", got.Content, "hello")
		}
	default:
		t.Error("c2 did not receive message")
	}
}

func TestHub_Broadcast_SkipsUnregisteredClients(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})

	c := &Client{sendCh: make(chan ServerMessage, 8)}
	hub.Register(c)
	hub.Unregister(c)

	msg := ServerMessage{Type: "thought", Content: "hello", Timestamp: time.Now()}
	hub.Broadcast(msg) // Should not panic or send to unregistered client.
}

func TestHub_Broadcast_DropsMessageForSlowClient(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})

	// Buffer of 1 — fill it, then broadcast should drop.
	c := &Client{sendCh: make(chan ServerMessage, 1)}
	hub.Register(c)

	hub.Broadcast(ServerMessage{Type: "thought", Content: "first", Timestamp: time.Now()})
	hub.Broadcast(ServerMessage{Type: "thought", Content: "second", Timestamp: time.Now()})

	got := <-c.sendCh
	if got.Content != "first" {
		t.Errorf("expected first message, got %q", got.Content)
	}

	select {
	case msg := <-c.sendCh:
		t.Errorf("expected no second message, got %q", msg.Content)
	default:
		// Expected: second message was dropped.
	}
}

func TestHub_HandleInput_WritesToPipe(t *testing.T) {
	var buf bytes.Buffer
	hub := NewHub(&buf)

	err := hub.HandleInput(ClientMessage{Type: "speech", Content: "hello"})
	if err != nil {
		t.Fatalf("HandleInput error: %v", err)
	}
	if got := buf.String(); got != "hello\n" {
		t.Errorf("got %q, want %q", got, "hello\n")
	}
}

func TestHub_HandleInput_ActionFormat(t *testing.T) {
	var buf bytes.Buffer
	hub := NewHub(&buf)

	hub.HandleInput(ClientMessage{Type: "action", Content: "waves"})
	if got := buf.String(); got != "*waves*\n" {
		t.Errorf("got %q, want %q", got, "*waves*\n")
	}
}

func TestHub_HandleInput_EnvironmentFormat(t *testing.T) {
	var buf bytes.Buffer
	hub := NewHub(&buf)

	hub.HandleInput(ClientMessage{Type: "environment", Content: "rain starts"})
	if got := buf.String(); got != "~rain starts\n" {
		t.Errorf("got %q, want %q", got, "~rain starts\n")
	}
}

func TestHub_UnregisterClosesSendChannel(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})

	c := &Client{sendCh: make(chan ServerMessage, 8)}
	hub.Register(c)
	hub.Unregister(c)

	_, ok := <-c.sendCh
	if ok {
		t.Error("expected sendCh to be closed")
	}
}

func TestHub_DoubleUnregisterDoesNotPanic(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})

	c := &Client{sendCh: make(chan ServerMessage, 8)}
	hub.Register(c)
	hub.Unregister(c)
	hub.Unregister(c) // Should not panic.
}
