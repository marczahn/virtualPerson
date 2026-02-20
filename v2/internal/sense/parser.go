package sense

import "strings"

type InputKind string

const (
	InputSpeech      InputKind = "speech"
	InputAction      InputKind = "action"
	InputEnvironment InputKind = "environment"
)

// ParsedInput is one classified operator input.
type ParsedInput struct {
	Kind    InputKind
	Content string
}

// Parser classifies one raw operator line into a typed external input.
type Parser interface {
	Parse(raw string) (ParsedInput, bool)
}

type ConventionParser struct{}

func NewParser() *ConventionParser {
	return &ConventionParser{}
}

// Parse applies v1 conventions:
// plain text -> speech, *text* -> action, ~text -> environment.
func (p *ConventionParser) Parse(raw string) (ParsedInput, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ParsedInput{}, false
	}

	if strings.HasPrefix(trimmed, "*") && strings.HasSuffix(trimmed, "*") && len(trimmed) > 2 {
		content := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
		if content != "" {
			return ParsedInput{Kind: InputAction, Content: content}, true
		}
	}

	if strings.HasPrefix(trimmed, "~") {
		return ParsedInput{
			Kind:    InputEnvironment,
			Content: strings.TrimSpace(trimmed[1:]),
		}, true
	}

	return ParsedInput{Kind: InputSpeech, Content: trimmed}, true
}
