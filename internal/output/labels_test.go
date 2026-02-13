package output

import "testing"

func TestSource_String_AllValues(t *testing.T) {
	tests := []struct {
		s    Source
		want string
	}{
		{Sense, "SENSE"},
		{Bio, "BIO"},
		{Psych, "PSYCH"},
		{Mind, "MIND"},
		{Review, "REVIEW"},
	}

	for _, tt := range tests {
		got := tt.s.String()
		if got != tt.want {
			t.Errorf("Source(%d).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestSource_String_OutOfRange(t *testing.T) {
	got := Source(99).String()
	if got != "UNKNOWN" {
		t.Errorf("Source(99).String() = %q, want %q", got, "UNKNOWN")
	}
}

func TestSource_EnumCount(t *testing.T) {
	if len(sourceLabels) != int(Review)+1 {
		t.Errorf("sourceLabels has %d entries, expected %d", len(sourceLabels), int(Review)+1)
	}
}
