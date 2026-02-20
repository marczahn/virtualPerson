package infrastructure

import (
	"fmt"

	"github.com/marczahn/person/v2/internal/biology"
	"github.com/marczahn/person/v2/internal/motivation"
	"github.com/marczahn/person/v2/internal/output"
)

// BuildTaggedOutputLines converts one tick result to tagged display lines.
func BuildTaggedOutputLines(result TickResult, previous motivation.MotivationState, driveThreshold float64) []string {
	lines := []string{formatBIOLine(result.Bio)}
	lines = append(lines, output.FormatDriveChangeLines(
		output.SignificantDriveChanges(previous, result.Motivation, driveThreshold),
	)...)
	if result.Parsed.Narrative != "" {
		lines = append(lines, output.FormatTaggedLine(output.SourceMIND, result.Parsed.Narrative))
	}
	return lines
}

func formatBIOLine(bioTick biology.TickResult) string {
	if len(bioTick.Deltas) == 0 && len(bioTick.Thresholds) == 0 {
		return output.FormatTaggedLine(output.SourceBIO, "no significant biological deltas")
	}
	return output.FormatTaggedLine(
		output.SourceBIO,
		fmt.Sprintf("deltas=%d threshold_events=%d", len(bioTick.Deltas), len(bioTick.Thresholds)),
	)
}
