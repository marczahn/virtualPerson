package infrastructure_test

import (
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/infrastructure"
)

type staticDrainer struct {
	input infrastructure.TickInput
}

func (d *staticDrainer) Drain() infrastructure.TickInput {
	return d.input
}

func TestScenarioInjector_RegisterActivateAndDrainEffects(t *testing.T) {
	base := &staticDrainer{input: infrastructure.TickInput{
		AllowedActions: map[string]bool{"eat": true},
	}}
	injector := infrastructure.NewScenarioInjector(base)

	if err := injector.Register("cold_room", []string{"cold room", "no food available"}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if !injector.Activate("cold_room") {
		t.Fatalf("expected activate to succeed for registered scenario")
	}

	got := injector.Drain()

	if !containsRate(got.PreBioRates, "body_temp", -0.03) {
		t.Fatalf("expected cold scenario to inject body_temp decay rate")
	}
	if got.AllowedActions["eat"] {
		t.Fatalf("expected scenario environment descriptor to block eat action")
	}
}

func TestScenarioInjector_SwitchIsDeterministicLatestWins(t *testing.T) {
	base := &staticDrainer{input: infrastructure.TickInput{
		AllowedActions: map[string]bool{"eat": true},
	}}
	injector := infrastructure.NewScenarioInjector(base)
	if err := injector.Register("cold_room", []string{"cold room"}); err != nil {
		t.Fatalf("register cold scenario failed: %v", err)
	}
	if err := injector.Register("hot_room", []string{"hot room"}); err != nil {
		t.Fatalf("register hot scenario failed: %v", err)
	}

	injector.Activate("cold_room")
	injector.Activate("hot_room")
	injector.Activate("cold_room")

	got := injector.Drain()

	if !containsRate(got.PreBioRates, "body_temp", -0.03) {
		t.Fatalf("expected latest active scenario (cold_room) to control rate effects")
	}
	if containsRate(got.PreBioRates, "body_temp", 0.03) {
		t.Fatalf("did not expect hot scenario effects after switching back to cold")
	}
}

func TestScenarioInjector_EffectsAppearOnNextDrainAfterRuntimeSwitch(t *testing.T) {
	base := &staticDrainer{input: infrastructure.TickInput{
		AllowedActions: map[string]bool{"rest": true},
	}}
	injector := infrastructure.NewScenarioInjector(base)
	if err := injector.Register("calm_space", []string{"quiet and peaceful"}); err != nil {
		t.Fatalf("register scenario failed: %v", err)
	}

	first := injector.Drain()
	if containsRate(first.PreBioRates, "stress", -0.02) {
		t.Fatalf("did not expect scenario effects before activation")
	}

	if !injector.Activate("calm_space") {
		t.Fatalf("expected activate to succeed for registered scenario")
	}

	second := injector.Drain()
	if !containsRate(second.PreBioRates, "stress", -0.02) {
		t.Fatalf("expected scenario stress-calming rate in next drain cycle")
	}
}

func containsRate(rates []biology.BioRate, field string, perSecond float64) bool {
	for _, rate := range rates {
		if rate.Field == field && rate.PerSecond == perSecond {
			return true
		}
	}
	return false
}
