package infrastructure

import (
	"strings"
	"sync"
	"time"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/motivation"
	"github.com/marczahn/person/v2/internal/sense"
)

// InputAdapter stores raw operator lines and drains them once per tick.
type InputAdapter struct {
	mu     sync.Mutex
	queue  []string
	parser sense.Parser
	nowFn  func() int64
}

func NewInputAdapter(parser sense.Parser, nowFn func() int64) *InputAdapter {
	if parser == nil {
		panic("input adapter requires parser")
	}
	if nowFn == nil {
		nowFn = func() int64 { return time.Now().Unix() }
	}
	return &InputAdapter{
		parser: parser,
		nowFn:  nowFn,
	}
}

func (a *InputAdapter) Enqueue(raw string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.queue = append(a.queue, raw)
}

func (a *InputAdapter) Drain() TickInput {
	a.mu.Lock()
	rawItems := append([]string(nil), a.queue...)
	a.queue = nil
	a.mu.Unlock()

	out := TickInput{
		AllowedActions: defaultAllowedActions(),
		NowSeconds:     a.nowFn(),
	}

	external := make([]string, 0, len(rawItems))
	for _, raw := range rawItems {
		parsed, ok := a.parser.Parse(raw)
		if !ok {
			continue
		}

		external = append(external, strings.TrimSpace(raw))

		switch parsed.Kind {
		case sense.InputAction:
			applyActionInput(parsed.Content, &out)
		case sense.InputEnvironment:
			applyEnvironmentInput(parsed.Content, &out)
		}
	}

	out.ExternalText = strings.Join(external, "\n")
	return out
}

func defaultAllowedActions() map[string]bool {
	return map[string]bool{
		string(motivation.ActionRest):      true,
		string(motivation.ActionEat):       true,
		string(motivation.ActionHydrate):   true,
		string(motivation.ActionReachOut):  true,
		string(motivation.ActionJournal):   true,
		string(motivation.ActionBreathe):   true,
		string(motivation.ActionScanArea):  true,
		string(motivation.ActionSeekWarm):  true,
		string(motivation.ActionSeekCool):  true,
		string(motivation.ActionMicroTask): true,
	}
}

func applyActionInput(content string, out *TickInput) {
	lower := strings.ToLower(content)

	if containsAny(lower, "punch", "hit", "kick", "slap", "strike", "shove", "attack") {
		out.PreBioPulses = append(out.PreBioPulses,
			biology.BioPulse{Field: "stress", Amount: 0.20},
			biology.BioPulse{Field: "physical_tension", Amount: 0.15},
			biology.BioPulse{Field: "mood", Amount: -0.08},
		)
	}

	if containsAny(lower, "hug", "comfort", "reassure", "support", "care") {
		out.PreBioPulses = append(out.PreBioPulses,
			biology.BioPulse{Field: "stress", Amount: -0.12},
			biology.BioPulse{Field: "physical_tension", Amount: -0.08},
			biology.BioPulse{Field: "mood", Amount: 0.08},
		)
	}

	if containsAny(lower, "feed", "food", "meal", "snack") {
		out.PreBioPulses = append(out.PreBioPulses, biology.BioPulse{Field: "hunger", Amount: -0.20})
		out.AllowedActions[string(motivation.ActionEat)] = true
	}
}

func applyEnvironmentInput(content string, out *TickInput) {
	lower := strings.ToLower(content)

	if containsAny(lower, "cold", "freezing", "chilly", "frigid") {
		out.PreBioRates = append(out.PreBioRates, biology.BioRate{Field: "body_temp", PerSecond: -0.03})
	}
	if containsAny(lower, "hot", "heat", "scorching", "sweltering") {
		out.PreBioRates = append(out.PreBioRates, biology.BioRate{Field: "body_temp", PerSecond: 0.03})
	}
	if containsAny(lower, "loud", "crowd", "chaos", "sirens") {
		out.PreBioRates = append(out.PreBioRates, biology.BioRate{Field: "stress", PerSecond: 0.03})
	}
	if containsAny(lower, "quiet", "calm", "safe", "peaceful") {
		out.PreBioRates = append(out.PreBioRates, biology.BioRate{Field: "stress", PerSecond: -0.02})
	}

	if containsAny(lower, "no food", "without food", "food unavailable") {
		out.AllowedActions[string(motivation.ActionEat)] = false
	} else if containsAny(lower, "food available", "food is available", "has food", "meal nearby", "kitchen stocked") {
		out.AllowedActions[string(motivation.ActionEat)] = true
	}

	if containsAny(lower, "no water", "without water", "water unavailable") {
		out.AllowedActions[string(motivation.ActionHydrate)] = false
	} else if containsAny(lower, "water available", "drinkable water") {
		out.AllowedActions[string(motivation.ActionHydrate)] = true
	}

	if containsAny(lower, "no quiet space", "cannot rest", "rest impossible") {
		out.AllowedActions[string(motivation.ActionRest)] = false
	} else if containsAny(lower, "quiet space available", "can rest", "safe resting place") {
		out.AllowedActions[string(motivation.ActionRest)] = true
	}
}

func containsAny(s string, patterns ...string) bool {
	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
