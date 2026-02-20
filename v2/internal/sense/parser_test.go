package sense_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/sense"
)

func TestParser_ParseConventions(t *testing.T) {
	p := sense.NewParser()

	tests := []struct {
		name    string
		raw     string
		ok      bool
		kind    sense.InputKind
		content string
	}{
		{
			name:    "speech plain text",
			raw:     "hello there",
			ok:      true,
			kind:    sense.InputSpeech,
			content: "hello there",
		},
		{
			name:    "action wrapped in asterisks",
			raw:     "  *pushes gently*  ",
			ok:      true,
			kind:    sense.InputAction,
			content: "pushes gently",
		},
		{
			name:    "environment prefixed by tilde",
			raw:     "~cold wind",
			ok:      true,
			kind:    sense.InputEnvironment,
			content: "cold wind",
		},
		{
			name:    "bare tilde is environment with empty content",
			raw:     "~",
			ok:      true,
			kind:    sense.InputEnvironment,
			content: "",
		},
		{
			name:    "invalid empty action falls back to speech",
			raw:     "**",
			ok:      true,
			kind:    sense.InputSpeech,
			content: "**",
		},
		{
			name:    "whitespace only ignored",
			raw:     "   \t ",
			ok:      false,
			kind:    sense.InputSpeech,
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.Parse(tt.raw)
			if ok != tt.ok {
				t.Fatalf("ok mismatch: got=%v want=%v", ok, tt.ok)
			}
			if !ok {
				return
			}
			if got.Kind != tt.kind {
				t.Fatalf("kind mismatch: got=%q want=%q", got.Kind, tt.kind)
			}
			if got.Content != tt.content {
				t.Fatalf("content mismatch: got=%q want=%q", got.Content, tt.content)
			}
		})
	}
}
