package consciousness

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/marczahn/person/v2/internal/motivation"
)

var (
	stateTagRe      = regexp.MustCompile(`(?i)\[STATE:\s*arousal=([-\d.]+),\s*valence=([-\d.]+)\]`)
	actionTagRe     = regexp.MustCompile(`(?i)\[ACTION:\s*([a-z_]+)\s*\]`)
	driveTagRe      = regexp.MustCompile(`(?i)\[DRIVE:\s*([a-z_]+)\s*=\s*([-\d.]+)\s*\]`)
	driveTagLooseRe = regexp.MustCompile(`(?i)\[DRIVE:[^\]]*\]`)
	stripTagsRe     = regexp.MustCompile(`(?i)\[(STATE|ACTION|DRIVE):[^\]]*\]`)
)

func ParseResponse(raw string, prior ParsedResponse) ParsedResponse {
	result := ParsedResponse{
		State:          prior.State,
		Action:         prior.Action,
		DriveOverrides: cloneOverrides(prior.DriveOverrides),
		Narrative:      stripTags(raw),
	}

	state, okState := parseState(raw)
	action, okAction := parseAction(raw)
	if !okState || !okAction {
		return result
	}

	result.State = state
	result.Action = action

	overrides, okOverrides := parseDriveOverrides(raw)
	if okOverrides {
		result.DriveOverrides = overrides
	} else if len(driveTagLooseRe.FindAllString(raw, -1)) == 0 {
		result.DriveOverrides = map[motivation.Drive]float64{}
	} else {
		// Malformed optional DRIVE tag still invalidates this parse turn.
		return ParsedResponse{
			State:          prior.State,
			Action:         prior.Action,
			DriveOverrides: cloneOverrides(prior.DriveOverrides),
			Narrative:      stripTags(raw),
		}
	}

	return result
}

func parseState(raw string) (ParsedState, bool) {
	match := stateTagRe.FindStringSubmatch(raw)
	if len(match) != 3 {
		return ParsedState{}, false
	}

	arousal, errA := strconv.ParseFloat(match[1], 64)
	valence, errV := strconv.ParseFloat(match[2], 64)
	if errA != nil || errV != nil {
		return ParsedState{}, false
	}

	return ParsedState{Arousal: arousal, Valence: valence}, true
}

func parseAction(raw string) (string, bool) {
	match := actionTagRe.FindStringSubmatch(raw)
	if len(match) != 2 {
		return "", false
	}
	return strings.ToLower(match[1]), true
}

func parseDriveOverrides(raw string) (map[motivation.Drive]float64, bool) {
	loose := driveTagLooseRe.FindAllString(raw, -1)
	if len(loose) == 0 {
		return nil, false
	}

	strict := driveTagRe.FindAllStringSubmatch(raw, -1)
	if len(strict) != len(loose) {
		return nil, false
	}

	overrides := make(map[motivation.Drive]float64, len(strict))
	for _, m := range strict {
		drive, ok := parseDriveName(strings.ToLower(m[1]))
		if !ok {
			return nil, false
		}
		value, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return nil, false
		}
		overrides[drive] = clamp01(value)
	}
	return overrides, true
}

func parseDriveName(name string) (motivation.Drive, bool) {
	switch name {
	case string(motivation.DriveEnergy):
		return motivation.DriveEnergy, true
	case string(motivation.DriveSocialConnection):
		return motivation.DriveSocialConnection, true
	case string(motivation.DriveStimulation):
		return motivation.DriveStimulation, true
	case string(motivation.DriveSafety):
		return motivation.DriveSafety, true
	case string(motivation.DriveIdentityCoherence):
		return motivation.DriveIdentityCoherence, true
	default:
		return "", false
	}
}

func stripTags(raw string) string {
	clean := stripTagsRe.ReplaceAllString(raw, "")
	lines := strings.Split(clean, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			kept = append(kept, trimmed)
		}
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

func cloneOverrides(in map[motivation.Drive]float64) map[motivation.Drive]float64 {
	if len(in) == 0 {
		return nil
	}
	out := make(map[motivation.Drive]float64, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
