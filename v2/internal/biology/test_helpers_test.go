package biology_test

import "github.com/marczahn/person/v2/internal/biology"

// filterEvents filters threshold events by variable name.
func filterEvents(events []biology.ThresholdEvent, variable string) []biology.ThresholdEvent {
	var filtered []biology.ThresholdEvent
	for _, e := range events {
		if e.Variable == variable {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
