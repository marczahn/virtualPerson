package server

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"nhooyr.io/websocket"
)

// Client represents a single WebSocket connection to the hub.
type Client struct {
	conn   *websocket.Conn
	sendCh chan ServerMessage
}

// Hub manages WebSocket clients and broadcasts server messages to all of them.
// It also writes client input to the simulation's input pipe.
type Hub struct {
	mu          sync.RWMutex
	clients     map[*Client]bool
	inputWriter io.Writer
}

// NewHub creates a hub that forwards client input to the given writer.
func NewHub(inputWriter io.Writer) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		inputWriter: inputWriter,
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
}

// Unregister removes a client from the hub and closes its send channel.
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.sendCh)
	}
}

// Broadcast sends a message to all connected clients. Non-blocking per client:
// if a client's send buffer is full, the message is dropped for that client.
func (h *Hub) Broadcast(msg ServerMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c.sendCh <- msg:
		default:
			// Client too slow, drop message.
		}
	}
}

// HandleInput reads a client message and writes the translated input line
// to the simulation's input pipe.
func (h *Hub) HandleInput(msg ClientMessage) error {
	line := msg.ToInputLine() + "\n"
	_, err := io.WriteString(h.inputWriter, line)
	return err
}

// ServeClient runs the read and write loops for a single WebSocket client.
// It blocks until the connection is closed or the context is cancelled.
func (h *Hub) ServeClient(ctx context.Context, conn *websocket.Conn) {
	c := &Client{
		conn:   conn,
		sendCh: make(chan ServerMessage, 64),
	}
	h.Register(c)
	defer h.Unregister(c)

	// Write loop in a goroutine.
	done := make(chan struct{})
	go func() {
		defer close(done)
		h.writeLoop(ctx, c)
	}()

	// Read loop in the current goroutine.
	h.readLoop(ctx, c)

	// Wait for write loop to finish.
	<-done
}

func (h *Hub) readLoop(ctx context.Context, c *Client) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
		var msg ClientMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		if err := msg.Validate(); err != nil {
			continue
		}
		h.HandleInput(msg)
	}
}

func (h *Hub) writeLoop(ctx context.Context, c *Client) {
	for {
		select {
		case msg, ok := <-c.sendCh:
			if !ok {
				return
			}
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			if err := c.conn.Write(ctx, websocket.MessageText, data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
