package server

import (
	"fmt"
	"time"
)

// ClientMessage is sent from the TUI client to the server.
type ClientMessage struct {
	Type    string `json:"type"`    // "speech", "action", "environment"
	Content string `json:"content"`
}

// Validate checks that the message has a known type and non-empty content.
func (m ClientMessage) Validate() error {
	switch m.Type {
	case "speech", "action", "environment":
	default:
		return fmt.Errorf("unknown message type: %q", m.Type)
	}
	if m.Content == "" {
		return fmt.Errorf("content must not be empty")
	}
	return nil
}

// ToInputLine converts a client message to the simulation's input format:
// speech → plain text, action → *text*, environment → ~text.
func (m ClientMessage) ToInputLine() string {
	switch m.Type {
	case "action":
		return "*" + m.Content + "*"
	case "environment":
		return "~" + m.Content
	default:
		return m.Content
	}
}

// ServerMessage is sent from the server to connected clients.
type ServerMessage struct {
	Type        string    `json:"type"`                   // "thought"
	Content     string    `json:"content"`
	ThoughtType string    `json:"thought_type,omitempty"` // "reactive", "spontaneous", "conversational"
	Trigger     string    `json:"trigger,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}
