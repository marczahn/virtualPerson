package consciousness_test

import (
	"math"
	"strings"
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/consciousness"
	"github.com/marczahn/person/v2/internal/motivation"
)

func TestThoughtSelection_TickerDueSelectsHighestUrgencyDrive(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.2,
		SocialUrgency:      0.1,
		StimulationUrgency: 0.3,
		SafetyUrgency:      0.9,
		IdentityUrgency:    0.4,
		ActiveGoalDrive:    motivation.DriveSafety,
		ActiveGoalUrgency:  0.9,
	}

	thought, fired := consciousness.SelectSpontaneousThought(state, consciousness.TickSchedule{EveryTicks: 3}, 6)
	if !fired {
		t.Fatal("expected thought to fire when tick is due")
	}
	if thought.Category != consciousness.ThoughtCategoryDrive {
		t.Fatalf("expected drive category on high urgency, got %q", thought.Category)
	}
	if thought.Drive != motivation.DriveSafety {
		t.Fatalf("expected safety thought from highest urgency drive, got %q", thought.Drive)
	}
}

func TestThoughtSelection_BaselineStillFiresViaAssociativeDrift(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.0,
		SocialUrgency:      0.0,
		StimulationUrgency: 0.0,
		SafetyUrgency:      0.0,
		IdentityUrgency:    0.0,
		ActiveGoalDrive:    motivation.DriveSafety,
		ActiveGoalUrgency:  0.0,
	}

	thought, fired := consciousness.SelectSpontaneousThought(state, consciousness.TickSchedule{EveryTicks: 1}, 1)
	if !fired {
		t.Fatal("expected spontaneous thought to fire at baseline state")
	}
	if thought.Category != consciousness.ThoughtCategoryAssociativeDrift {
		t.Fatalf("expected associative drift at baseline, got %q", thought.Category)
	}
}

func TestThoughtSelection_NotDueDoesNotFire(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.9,
		SocialUrgency:      0.1,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.1,
		IdentityUrgency:    0.1,
		ActiveGoalDrive:    motivation.DriveEnergy,
		ActiveGoalUrgency:  0.9,
	}

	_, fired := consciousness.SelectSpontaneousThought(state, consciousness.TickSchedule{EveryTicks: 4}, 3)
	if fired {
		t.Fatal("expected no thought when ticker is not due")
	}
}

func TestBuildPromptContext_IncludesContinuityBuffer(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.2,
		SocialUrgency:      0.9,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.3,
		IdentityUrgency:    0.2,
		ActiveGoalDrive:    motivation.DriveSocialConnection,
		ActiveGoalUrgency:  0.9,
	}
	continuity := []consciousness.Thought{
		{Category: consciousness.ThoughtCategoryDrive, Drive: motivation.DriveSocialConnection, Text: "I keep reaching for contact."},
		{Category: consciousness.ThoughtCategoryAssociativeDrift, Text: "A stray memory of a corridor returns."},
	}

	ctx := consciousness.BuildPromptContextWithContinuity(state, continuity)
	if len(ctx.ContinuityBuffer) != 2 {
		t.Fatalf("expected continuity buffer entries in prompt context, got %d", len(ctx.ContinuityBuffer))
	}
	if ctx.ContinuityBuffer[0] != "I keep reaching for contact." {
		t.Fatalf("unexpected first continuity entry: %q", ctx.ContinuityBuffer[0])
	}
}

func TestContinuityBuffer_KeepsMostRecentBoundedThoughts(t *testing.T) {
	buffer := consciousness.NewContinuityBuffer(2)
	buffer.Add(consciousness.Thought{Text: "first"})
	buffer.Add(consciousness.Thought{Text: "second"})
	buffer.Add(consciousness.Thought{Text: "third"})

	items := buffer.Items()
	if len(items) != 2 {
		t.Fatalf("expected bounded continuity size, got %d", len(items))
	}
	if items[0].Text != "second" || items[1].Text != "third" {
		t.Fatalf("expected oldest item dropped, got %+v", items)
	}
}

func TestBuildPromptContext_UsesTopTwoAsPrimary(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.25,
		SocialUrgency:      0.90,
		StimulationUrgency: 0.15,
		SafetyUrgency:      0.80,
		IdentityUrgency:    0.40,
		ActiveGoalDrive:    motivation.DriveSocialConnection,
		ActiveGoalUrgency:  0.90,
	}

	ctx := consciousness.BuildPromptContext(state)

	if len(ctx.Primary) != 2 {
		t.Fatalf("expected exactly 2 primary drives, got %d", len(ctx.Primary))
	}
	if ctx.Primary[0].Drive != motivation.DriveSocialConnection {
		t.Fatalf("expected top primary to be social, got %s", ctx.Primary[0].Drive)
	}
	if ctx.Primary[1].Drive != motivation.DriveSafety {
		t.Fatalf("expected second primary to be safety, got %s", ctx.Primary[1].Drive)
	}
	if len(ctx.Background) != 3 {
		t.Fatalf("expected 3 background drives, got %d", len(ctx.Background))
	}
}

func TestBuildPromptContext_FeltLanguageContainsNoRawNumbers(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.11,
		SocialUrgency:      0.22,
		StimulationUrgency: 0.33,
		SafetyUrgency:      0.44,
		IdentityUrgency:    0.55,
		ActiveGoalDrive:    motivation.DriveIdentityCoherence,
		ActiveGoalUrgency:  0.55,
	}

	ctx := consciousness.BuildPromptContext(state)
	lines := make([]string, 0, len(ctx.Primary)+len(ctx.Background)+1)
	for _, d := range ctx.Primary {
		lines = append(lines, d.Felt)
	}
	for _, d := range ctx.Background {
		lines = append(lines, d.Felt)
	}
	lines = append(lines, ctx.GoalPull)

	for _, line := range lines {
		for _, ch := range line {
			if ch >= '0' && ch <= '9' {
				t.Fatalf("prompt language must not expose raw numbers, got %q", line)
			}
		}
	}
}

func TestBuildPromptContext_ActiveGoalIsImplicitPull(t *testing.T) {
	state := motivation.MotivationState{
		EnergyUrgency:      0.20,
		SocialUrgency:      0.10,
		StimulationUrgency: 0.10,
		SafetyUrgency:      0.95,
		IdentityUrgency:    0.15,
		ActiveGoalDrive:    motivation.DriveSafety,
		ActiveGoalUrgency:  0.95,
	}

	ctx := consciousness.BuildPromptContext(state)
	if strings.Contains(strings.ToUpper(ctx.GoalPull), "GOAL:") {
		t.Fatalf("goal pull must be implicit, got %q", ctx.GoalPull)
	}
	if strings.Contains(strings.ToLower(ctx.GoalPull), "drive") {
		t.Fatalf("goal pull should not expose internal drive label, got %q", ctx.GoalPull)
	}
}

func TestParseResponse_ValidTagsParsedAndStripped(t *testing.T) {
	prior := consciousness.ParsedResponse{
		State:  consciousness.ParsedState{Arousal: 0.2, Valence: 0.1},
		Action: "journal",
		DriveOverrides: map[motivation.Drive]float64{
			motivation.DriveEnergy: 0.4,
		},
		Narrative: "old",
	}
	raw := "I need to steady myself.\n[STATE: arousal=0.8, valence=-0.4]\n[ACTION: breathe]\n[DRIVE: safety=0.9]"

	got := consciousness.ParseResponse(raw, prior)

	if got.State.Arousal != 0.8 || got.State.Valence != -0.4 {
		t.Fatalf("unexpected state parse: %+v", got.State)
	}
	if got.Action != "breathe" {
		t.Fatalf("unexpected action parse: %q", got.Action)
	}
	if got.DriveOverrides[motivation.DriveSafety] != 0.9 {
		t.Fatalf("expected safety override 0.9, got %v", got.DriveOverrides[motivation.DriveSafety])
	}
	if strings.Contains(got.Narrative, "[STATE:") || strings.Contains(got.Narrative, "[ACTION:") || strings.Contains(got.Narrative, "[DRIVE:") {
		t.Fatalf("tags must be stripped from narrative, got %q", got.Narrative)
	}
}

func TestParseResponse_MissingRequiredTagFallsBackToPrior(t *testing.T) {
	prior := consciousness.ParsedResponse{
		State:  consciousness.ParsedState{Arousal: 0.3, Valence: -0.2},
		Action: "rest",
		DriveOverrides: map[motivation.Drive]float64{
			motivation.DriveSocialConnection: 0.6,
		},
		Narrative: "old",
	}
	raw := "I should pause.\n[ACTION: breathe]"

	got := consciousness.ParseResponse(raw, prior)

	if got.State != prior.State {
		t.Fatalf("missing state tag must preserve prior state: got=%+v prior=%+v", got.State, prior.State)
	}
	if got.Action != prior.Action {
		t.Fatalf("missing state tag must preserve prior action: got=%q prior=%q", got.Action, prior.Action)
	}
	if got.DriveOverrides[motivation.DriveSocialConnection] != 0.6 {
		t.Fatalf("missing required tag must preserve prior overrides, got %v", got.DriveOverrides[motivation.DriveSocialConnection])
	}
}

func TestParseResponse_MissingActionTagFallsBackToPrior(t *testing.T) {
	prior := consciousness.ParsedResponse{
		State:  consciousness.ParsedState{Arousal: 0.3, Valence: -0.2},
		Action: "rest",
		DriveOverrides: map[motivation.Drive]float64{
			motivation.DriveSocialConnection: 0.6,
		},
		Narrative: "old",
	}
	raw := "I should pause.\n[STATE: arousal=0.5, valence=0.1]"

	got := consciousness.ParseResponse(raw, prior)

	if got.State != prior.State {
		t.Fatalf("missing action tag must preserve prior state: got=%+v prior=%+v", got.State, prior.State)
	}
	if got.Action != prior.Action {
		t.Fatalf("missing action tag must preserve prior action: got=%q prior=%q", got.Action, prior.Action)
	}
	if got.DriveOverrides[motivation.DriveSocialConnection] != 0.6 {
		t.Fatalf("missing required tag must preserve prior overrides, got %v", got.DriveOverrides[motivation.DriveSocialConnection])
	}
}

func TestParseResponse_MalformedTagFallsBackToPrior(t *testing.T) {
	prior := consciousness.ParsedResponse{
		State:     consciousness.ParsedState{Arousal: 0.3, Valence: -0.2},
		Action:    "rest",
		Narrative: "old",
	}
	raw := "Something is off.\n[STATE: arousal=abc, valence=-0.1]\n[ACTION: breathe]"

	got := consciousness.ParseResponse(raw, prior)
	if got.State != prior.State || got.Action != prior.Action {
		t.Fatalf("malformed state tag must preserve prior parse, got=%+v", got)
	}
}

func TestParseResponse_MalformedDriveTagFallsBackToPrior(t *testing.T) {
	prior := consciousness.ParsedResponse{
		State:  consciousness.ParsedState{Arousal: 0.3, Valence: -0.2},
		Action: "rest",
		DriveOverrides: map[motivation.Drive]float64{
			motivation.DriveSafety: 0.7,
		},
		Narrative: "old",
	}
	raw := "I can handle this.\n[STATE: arousal=0.4, valence=0.2]\n[ACTION: breathe]\n[DRIVE: safety=oops]"

	got := consciousness.ParseResponse(raw, prior)
	if got.State != prior.State || got.Action != prior.Action {
		t.Fatalf("malformed drive tag must trigger full fallback, got=%+v", got)
	}
	if got.DriveOverrides[motivation.DriveSafety] != 0.7 {
		t.Fatalf("malformed drive tag must preserve prior overrides, got=%+v", got.DriveOverrides)
	}
}

func TestParseResponse_DriveTagIsOptional(t *testing.T) {
	prior := consciousness.ParsedResponse{
		State:     consciousness.ParsedState{Arousal: 0.1, Valence: 0.1},
		Action:    "journal",
		Narrative: "old",
	}
	raw := "I can regulate this.\n[STATE: arousal=0.2, valence=0.3]\n[ACTION: breathe]"

	got := consciousness.ParseResponse(raw, prior)
	if got.State.Arousal != 0.2 || got.State.Valence != 0.3 {
		t.Fatalf("expected parsed state without drive tag, got=%+v", got.State)
	}
	if got.Action != "breathe" {
		t.Fatalf("expected parsed action without drive tag, got=%q", got.Action)
	}
	if len(got.DriveOverrides) != 0 {
		t.Fatalf("expected no drive overrides when tag is absent, got=%+v", got.DriveOverrides)
	}
}

func TestResolveActionOutcome_AllowedActionExecutesAndSatisfies(t *testing.T) {
	outcome := consciousness.ResolveActionOutcome("eat", true)

	if !outcome.Executed {
		t.Fatalf("allowed action must execute: %+v", outcome)
	}
	if !outcome.Satisfied {
		t.Fatalf("successful action must satisfy drive: %+v", outcome)
	}
	if outcome.Action != "eat" {
		t.Fatalf("unexpected action in outcome: %+v", outcome)
	}
}

func TestResolveActionOutcome_BlockedActionDoesNotSatisfy(t *testing.T) {
	outcome := consciousness.ResolveActionOutcome("eat", false)

	if outcome.Executed {
		t.Fatalf("blocked action must not execute: %+v", outcome)
	}
	if outcome.Satisfied {
		t.Fatalf("blocked action must not satisfy drive: %+v", outcome)
	}
}

func TestResolveActionOutcome_IsDeterministic(t *testing.T) {
	a := consciousness.ResolveActionOutcome("breathe", true)
	b := consciousness.ResolveActionOutcome("breathe", true)
	if a != b {
		t.Fatalf("action outcome must be deterministic: a=%+v b=%+v", a, b)
	}
}

func TestEmotionalPulseFromState_UsesAbsolutePulseContract(t *testing.T) {
	state := consciousness.ParsedState{Arousal: 0.8, Valence: -0.5}

	deltas := consciousness.EmotionalPulseFromState(state)
	if len(deltas) == 0 {
		t.Fatal("expected emotional pulse deltas")
	}

	assertDeltaAmount(t, deltas, "stress", 0.136)
	assertDeltaAmount(t, deltas, "mood", -0.12)
	assertDeltaAmount(t, deltas, "physical_tension", 0.08)
	assertDeltaAmount(t, deltas, "cognitive_capacity", -0.048)
}

func TestActionPulse_SuccessfulActionProducesBioChanges(t *testing.T) {
	deltas := consciousness.ActionPulse(consciousness.ActionOutcome{Action: "eat", Executed: true, Satisfied: true})
	if len(deltas) == 0 {
		t.Fatal("expected bio deltas for successful action")
	}

	assertDeltaAmount(t, deltas, "hunger", -0.30)
	assertDeltaAmount(t, deltas, "energy", 0.08)
}

func TestActionPulse_BlockedActionProducesNoBioChanges(t *testing.T) {
	deltas := consciousness.ActionPulse(consciousness.ActionOutcome{Action: "eat", Executed: false, Satisfied: false})
	if len(deltas) != 0 {
		t.Fatalf("blocked action must not emit bio deltas, got %+v", deltas)
	}
}

func TestApplyDriveOverrides_RespectsHalfRawFloor(t *testing.T) {
	raw := motivation.MotivationState{
		EnergyUrgency:      0.8,
		SocialUrgency:      0.2,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.3,
		IdentityUrgency:    0.1,
		ActiveGoalDrive:    motivation.DriveEnergy,
		ActiveGoalUrgency:  0.8,
	}
	overrides := map[motivation.Drive]float64{
		motivation.DriveEnergy: 0.1,
	}

	got := consciousness.ApplyDriveOverrides(raw, overrides)
	if got.EnergyUrgency != 0.4 {
		t.Fatalf("expected floor at raw*0.5 (=0.4), got %f", got.EnergyUrgency)
	}
}

func TestApplyDriveOverrides_AllowsHigherReportedValue(t *testing.T) {
	raw := motivation.MotivationState{
		EnergyUrgency:      0.2,
		SocialUrgency:      0.3,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.1,
		IdentityUrgency:    0.1,
		ActiveGoalDrive:    motivation.DriveSocialConnection,
		ActiveGoalUrgency:  0.3,
	}
	overrides := map[motivation.Drive]float64{
		motivation.DriveEnergy: 0.9,
	}

	got := consciousness.ApplyDriveOverrides(raw, overrides)
	if got.EnergyUrgency != 0.9 {
		t.Fatalf("expected higher reported override to win, got %f", got.EnergyUrgency)
	}
	if got.ActiveGoalDrive != motivation.DriveEnergy {
		t.Fatalf("expected active goal to recompute to energy, got %s", got.ActiveGoalDrive)
	}
}

func TestApplyDriveOverrides_ClampsInputsToRange(t *testing.T) {
	raw := motivation.MotivationState{
		EnergyUrgency:      1.5,
		SocialUrgency:      -0.2,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.1,
		IdentityUrgency:    0.1,
		ActiveGoalDrive:    motivation.DriveEnergy,
		ActiveGoalUrgency:  1.0,
	}
	overrides := map[motivation.Drive]float64{
		motivation.DriveEnergy:           -3,
		motivation.DriveSocialConnection: 3,
	}

	got := consciousness.ApplyDriveOverrides(raw, overrides)
	if got.EnergyUrgency != 0.5 {
		t.Fatalf("energy expected max(raw*0.5, reported) with clamp => 0.5, got %f", got.EnergyUrgency)
	}
	if got.SocialUrgency != 1.0 {
		t.Fatalf("social override should clamp to 1.0, got %f", got.SocialUrgency)
	}
}

func TestApplyParsedDriveOverridesForNextTick_UsesClampedPerceptionContract(t *testing.T) {
	raw := motivation.MotivationState{
		EnergyUrgency:      0.8,
		SocialUrgency:      0.2,
		StimulationUrgency: 0.1,
		SafetyUrgency:      0.3,
		IdentityUrgency:    0.1,
		ActiveGoalDrive:    motivation.DriveEnergy,
		ActiveGoalUrgency:  0.8,
	}
	parsed := consciousness.ParsedResponse{
		DriveOverrides: map[motivation.Drive]float64{
			motivation.DriveEnergy: 0.1,
		},
	}

	got := consciousness.ApplyParsedDriveOverridesForNextTick(raw, parsed)
	if got.EnergyUrgency != 0.4 {
		t.Fatalf("expected next-tick energy perception at raw floor 0.4, got %f", got.EnergyUrgency)
	}
}

func TestResolveActionOutcomeWithCooldown_RepeatedActionRejectedWithinWindow(t *testing.T) {
	cooldowns := consciousness.ActionCooldowns{
		string(motivation.ActionEat): 30,
	}

	first, state := consciousness.ResolveActionOutcomeWithCooldown("eat", true, 100, cooldowns, nil)
	if !first.Executed || !first.Satisfied {
		t.Fatalf("first action should execute and satisfy, got %+v", first)
	}

	second, state := consciousness.ResolveActionOutcomeWithCooldown("eat", true, 110, cooldowns, state)
	if second.Executed || second.Satisfied {
		t.Fatalf("repeated action in cooldown must be blocked and unsatisfied, got %+v", second)
	}

	third, _ := consciousness.ResolveActionOutcomeWithCooldown("eat", true, 130, cooldowns, state)
	if !third.Executed || !third.Satisfied {
		t.Fatalf("action should execute again after cooldown expiry, got %+v", third)
	}
}

func TestResolveActionOutcomeWithCooldown_EnvironmentBlockDoesNotStartCooldown(t *testing.T) {
	cooldowns := consciousness.ActionCooldowns{
		string(motivation.ActionEat): 30,
	}

	outcome, state := consciousness.ResolveActionOutcomeWithCooldown("eat", false, 200, cooldowns, nil)
	if outcome.Executed || outcome.Satisfied {
		t.Fatalf("blocked environment action must not execute/satisfy, got %+v", outcome)
	}
	if len(state) != 0 {
		t.Fatalf("environment-blocked action must not update cooldown state, got %+v", state)
	}
}

func assertDeltaAmount(t *testing.T, deltas []biology.BioPulse, field string, want float64) {
	t.Helper()
	for _, d := range deltas {
		if d.Field != field {
			continue
		}
		if math.Abs(d.Amount-want) > 0.000001 {
			t.Fatalf("unexpected delta for %s: got %f want %f", field, d.Amount, want)
		}
		return
	}
	t.Fatalf("expected delta for field %s not found in %+v", field, deltas)
}
