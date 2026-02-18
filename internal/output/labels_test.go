package output

import (
	"os"
	"testing"

	"github.com/marczahn/person/internal/i18n"
)

func TestMain(m *testing.M) {
	if err := i18n.Load("en"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

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
