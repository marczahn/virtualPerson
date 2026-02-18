package sense

import "testing"

func TestChannel_String_AllChannels(t *testing.T) {
	tests := []struct {
		channel Channel
		want    string
	}{
		{Visual, "visual"},
		{Auditory, "auditory"},
		{Tactile, "tactile"},
		{Thermal, "thermal"},
		{Pain, "pain"},
		{Olfactory, "olfactory"},
		{Gustatory, "gustatory"},
		{Vestibular, "vestibular"},
		{Interoceptive, "interoceptive"},
	}

	for _, tt := range tests {
		got := tt.channel.String()
		if got != tt.want {
			t.Errorf("Channel(%d).String() = %q, want %q", tt.channel, got, tt.want)
		}
	}
}

func TestChannel_String_OutOfRange(t *testing.T) {
	got := Channel(99).String()
	if got != "unknown" {
		t.Errorf("Channel(99).String() = %q, want %q", got, "unknown")
	}
}

func TestChannel_EnumCount(t *testing.T) {
	// Ensure the channelNames array covers all defined channels.
	// Interoceptive is the last channel (iota=8), so there should be 9 names.
	if len(channelNames) != int(Interoceptive)+1 {
		t.Errorf("channelNames has %d entries, expected %d (one per channel through Interoceptive)",
			len(channelNames), int(Interoceptive)+1)
	}
}
