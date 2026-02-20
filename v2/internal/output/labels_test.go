package output

import "testing"

func TestFormatTaggedLine_UsesExpectedTagPrefix(t *testing.T) {
	cases := []struct {
		name string
		tag  SourceTag
		msg  string
		want string
	}{
		{name: "bio", tag: SourceBIO, msg: "stress threshold crossed", want: "[BIO] stress threshold crossed"},
		{name: "drives", tag: SourceDRIVES, msg: "energy +0.20", want: "[DRIVES] energy +0.20"},
		{name: "mind", tag: SourceMIND, msg: "I should rest for a moment.", want: "[MIND] I should rest for a moment."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatTaggedLine(tc.tag, tc.msg); got != tc.want {
				t.Fatalf("unexpected tagged line: got=%q want=%q", got, tc.want)
			}
		})
	}
}
