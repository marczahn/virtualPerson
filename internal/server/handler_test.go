package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func TestHandler_WebSocketUpgrade(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})
	srv := httptest.NewServer(NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	conn.Close(websocket.StatusNormalClosure, "")
}

func TestHandler_ClientSendsMessage(t *testing.T) {
	var inputBuf bytes.Buffer
	hub := NewHub(&inputBuf)
	srv := httptest.NewServer(NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	msg := ClientMessage{Type: "speech", Content: "hello world"}
	data, _ := json.Marshal(msg)
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("write error: %v", err)
	}

	// Give the server a moment to process.
	time.Sleep(100 * time.Millisecond)

	got := inputBuf.String()
	if got != "hello world\n" {
		t.Errorf("input pipe got %q, want %q", got, "hello world\n")
	}
}

func TestHandler_ServerBroadcastsToClient(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})
	srv := httptest.NewServer(NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Wait for client registration.
	time.Sleep(50 * time.Millisecond)

	broadcast := ServerMessage{
		Type:        "thought",
		Content:     "I wonder...",
		ThoughtType: "spontaneous",
		Timestamp:   time.Now(),
	}
	hub.Broadcast(broadcast)

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	var got ServerMessage
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got.Content != "I wonder..." {
		t.Errorf("got content %q, want %q", got.Content, "I wonder...")
	}
	if got.ThoughtType != "spontaneous" {
		t.Errorf("got thought_type %q, want %q", got.ThoughtType, "spontaneous")
	}
}

func TestHandler_InvalidMessageIgnored(t *testing.T) {
	hub := NewHub(&bytes.Buffer{})
	srv := httptest.NewServer(NewHandler(hub))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Send invalid JSON — should not crash.
	conn.Write(ctx, websocket.MessageText, []byte("not json"))

	// Send message with unknown type — should be ignored.
	msg := ClientMessage{Type: "unknown", Content: "bad"}
	data, _ := json.Marshal(msg)
	conn.Write(ctx, websocket.MessageText, data)

	// Connection should still be alive — verify by sending a valid message.
	time.Sleep(50 * time.Millisecond)
	valid := ClientMessage{Type: "speech", Content: "still here"}
	data, _ = json.Marshal(valid)
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Errorf("expected connection to survive invalid messages, got error: %v", err)
	}
}
