package infrastructure

import (
	"fmt"
	"strings"
	"sync"
)

// ScenarioInjector wraps an InputDrainer and injects active scenario effects per drain.
// Scenario effects are derived from environment descriptors and mapped through existing
// infrastructure environment rules.
type ScenarioInjector struct {
	base InputDrainer

	mu        sync.Mutex
	scenarios map[string][]string
	active    string
}

func NewScenarioInjector(base InputDrainer) *ScenarioInjector {
	if base == nil {
		panic(fmt.Errorf("scenario injector requires InputDrainer"))
	}
	return &ScenarioInjector{
		base:      base,
		scenarios: make(map[string][]string),
	}
}

func (s *ScenarioInjector) Register(name string, descriptors []string) error {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return fmt.Errorf("scenario name must not be empty")
	}

	clean := make([]string, 0, len(descriptors))
	for _, descriptor := range descriptors {
		descriptor = strings.TrimSpace(descriptor)
		if descriptor == "" {
			continue
		}
		clean = append(clean, descriptor)
	}
	if len(clean) == 0 {
		return fmt.Errorf("scenario %q requires at least one descriptor", trimmedName)
	}

	s.mu.Lock()
	s.scenarios[trimmedName] = clean
	s.mu.Unlock()
	return nil
}

func (s *ScenarioInjector) Activate(name string) bool {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scenarios[trimmedName]; !ok {
		return false
	}
	s.active = trimmedName
	return true
}

func (s *ScenarioInjector) Drain() TickInput {
	out := s.base.Drain()
	out.AllowedActions = cloneAllowedActions(out.AllowedActions)

	active, descriptors := s.snapshotActiveScenario()
	if active == "" {
		return out
	}

	for _, descriptor := range descriptors {
		applyEnvironmentInput(descriptor, &out)
	}

	scenarioText := "@scenario " + active + ": " + strings.Join(descriptors, "; ")
	if out.ExternalText == "" {
		out.ExternalText = scenarioText
	} else {
		out.ExternalText += "\n" + scenarioText
	}
	return out
}

func (s *ScenarioInjector) snapshotActiveScenario() (string, []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active == "" {
		return "", nil
	}
	descriptors, ok := s.scenarios[s.active]
	if !ok {
		return "", nil
	}
	return s.active, append([]string(nil), descriptors...)
}

func cloneAllowedActions(in map[string]bool) map[string]bool {
	if in == nil {
		return defaultAllowedActions()
	}
	out := make(map[string]bool, len(in))
	for action, allowed := range in {
		out[action] = allowed
	}
	return out
}
