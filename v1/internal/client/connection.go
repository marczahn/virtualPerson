package client

import (
	"context"
	"encoding/json"

	"github.com/marczahn/person/internal/server"
	"nhooyr.io/websocket"
)

// Connection manages a WebSocket connection to the simulation server.
type Connection struct {
	conn      *websocket.Conn
	receiveCh chan server.ServerMessage
}

// Dial connects to the simulation server at the given WebSocket URL.
func Dial(ctx context.Context, url string) (*Connection, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return &Connection{
		conn:      conn,
		receiveCh: make(chan server.ServerMessage, 64),
	}, nil
}

// Run reads messages from the server and pushes them to the receive channel.
// It blocks until the connection is closed or the context is cancelled.
func (c *Connection) Run(ctx context.Context) {
	defer close(c.receiveCh)
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
		var msg server.ServerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		select {
		case c.receiveCh <- msg:
		case <-ctx.Done():
			return
		}
	}
}

// Send sends a client message to the server.
func (c *Connection) Send(ctx context.Context, msg server.ClientMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.conn.Write(ctx, websocket.MessageText, data)
}

// Messages returns the channel of server messages.
func (c *Connection) Messages() <-chan server.ServerMessage {
	return c.receiveCh
}

// Close closes the WebSocket connection.
func (c *Connection) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}
