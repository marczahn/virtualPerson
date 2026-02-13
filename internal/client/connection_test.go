package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/marczahn/person/internal/server"
)

func TestConnection_ReceivesServerMessages(t *testing.T) {
	hub := server.NewHub(&bytes.Buffer{})
	srv := httptest.NewServer(server.NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, err := Dial(ctx, wsURL)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	go conn.Run(ctx)

	// Wait for registration.
	time.Sleep(50 * time.Millisecond)

	hub.Broadcast(server.ServerMessage{
		Type:        "thought",
		Content:     "test thought",
		ThoughtType: "reactive",
		Timestamp:   time.Now(),
	})

	select {
	case msg := <-conn.Messages():
		if msg.Content != "test thought" {
			t.Errorf("got content %q, want %q", msg.Content, "test thought")
		}
		if msg.ThoughtType != "reactive" {
			t.Errorf("got thought_type %q, want %q", msg.ThoughtType, "reactive")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for message")
	}
}

func TestConnection_SendsClientMessages(t *testing.T) {
	var inputBuf bytes.Buffer
	hub := server.NewHub(&inputBuf)
	srv := httptest.NewServer(server.NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, err := Dial(ctx, wsURL)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	go conn.Run(ctx)

	err = conn.Send(ctx, server.ClientMessage{Type: "speech", Content: "hello"})
	if err != nil {
		t.Fatalf("send error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	got := inputBuf.String()
	if got != "hello\n" {
		t.Errorf("server input got %q, want %q", got, "hello\n")
	}
}

func TestConnection_ChannelClosesOnDisconnect(t *testing.T) {
	hub := server.NewHub(&bytes.Buffer{})
	srv := httptest.NewServer(server.NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, err := Dial(ctx, wsURL)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}

	go conn.Run(ctx)
	time.Sleep(50 * time.Millisecond)

	conn.Close()

	// Messages channel should eventually close.
	select {
	case _, ok := <-conn.Messages():
		if ok {
			// Got a message before close, try again.
			_, ok = <-conn.Messages()
		}
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}

func TestDial_InvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := Dial(ctx, "ws://localhost:0/ws")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

// verifyJSON checks that ServerMessage round-trips through JSON correctly.
func TestServerMessage_JSONRoundTrip(t *testing.T) {
	msg := server.ServerMessage{
		Type:        "thought",
		Content:     "I feel something",
		ThoughtType: "spontaneous",
		Trigger:     "boredom",
		Timestamp:   time.Date(2024, 6, 15, 8, 15, 0, 0, time.UTC),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got server.ServerMessage
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.Content != msg.Content {
		t.Errorf("content: got %q, want %q", got.Content, msg.Content)
	}
	if got.ThoughtType != msg.ThoughtType {
		t.Errorf("thought_type: got %q, want %q", got.ThoughtType, msg.ThoughtType)
	}
}
