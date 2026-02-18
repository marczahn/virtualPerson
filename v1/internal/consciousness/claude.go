package consciousness

import (
	"context"
	"fmt"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

var _ LLM = (*ClaudeAdapter)(nil)

// ClaudeAdapter implements the LLM interface using the Anthropic Claude API.
type ClaudeAdapter struct {
	client *anthropic.Client
	model  anthropic.Model
}

// NewClaudeAdapter creates an adapter that calls the Claude API.
// The API key is read from the ANTHROPIC_API_KEY environment variable
// by the SDK automatically.
func NewClaudeAdapter(model anthropic.Model) *ClaudeAdapter {
	client := anthropic.NewClient()
	return &ClaudeAdapter{
		client: &client,
		model:  model,
	}
}

// Complete sends a system prompt and user message to Claude, returning the text response.
func (ca *ClaudeAdapter) Complete(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	resp, err := ca.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     ca.model,
		MaxTokens: 1024,
		Temperature: param.NewOpt(0.9),
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude API call: %w", err)
	}

	return extractText(resp), nil
}

// extractText concatenates all text blocks from a Claude response.
func extractText(msg *anthropic.Message) string {
	var parts []string
	for _, block := range msg.Content {
		if block.Type == "text" {
			parts = append(parts, block.AsText().Text)
		}
	}
	return strings.Join(parts, "")
}
