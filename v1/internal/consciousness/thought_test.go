package consciousness

import "testing"

func TestThoughtType_String_AllValues(t *testing.T) {
	tests := []struct {
		tt   ThoughtType
		want string
	}{
		{Reactive, "reactive"},
		{Spontaneous, "spontaneous"},
		{Conversational, "conversational"},
	}

	for _, tc := range tests {
		got := tc.tt.String()
		if got != tc.want {
			t.Errorf("ThoughtType(%d).String() = %q, want %q", tc.tt, got, tc.want)
		}
	}
}

func TestThoughtType_String_OutOfRange(t *testing.T) {
	got := ThoughtType(99).String()
	if got != "unknown" {
		t.Errorf("ThoughtType(99).String() = %q, want %q", got, "unknown")
	}
}

func TestThoughtType_EnumCount(t *testing.T) {
	if len(thoughtTypeNames) != int(Conversational)+1 {
		t.Errorf("thoughtTypeNames has %d entries, expected %d", len(thoughtTypeNames), int(Conversational)+1)
	}
}

func TestPriority_String_AllValues(t *testing.T) {
	tests := []struct {
		p    Priority
		want string
	}{
		{PriorityPredictionError, "prediction_error"},
		{PriorityBiologicalNeed, "biological_need"},
		{PriorityGoalRehearsal, "goal_rehearsal"},
		{PrioritySocialModeling, "social_modeling"},
		{PriorityAssociativeDrift, "associative_drift"},
	}

	for _, tt := range tests {
		got := tt.p.String()
		if got != tt.want {
			t.Errorf("Priority(%d).String() = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestPriority_String_OutOfRange(t *testing.T) {
	got := Priority(99).String()
	if got != "unknown" {
		t.Errorf("Priority(99).String() = %q, want %q", got, "unknown")
	}
}

func TestPriority_EnumCount(t *testing.T) {
	if len(priorityNames) != int(PriorityAssociativeDrift)+1 {
		t.Errorf("priorityNames has %d entries, expected %d", len(priorityNames), int(PriorityAssociativeDrift)+1)
	}
}
