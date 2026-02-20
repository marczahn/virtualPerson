package biology

// ApplyThresholdCascades applies the cascade Deltas from all threshold events to s.
// After calling this, ClampAll must be called to keep values in range.
func ApplyThresholdCascades(s *State, events []ThresholdEvent) {
	for _, e := range events {
		for _, d := range e.Cascade {
			applyDelta(s, d)
		}
	}
}
