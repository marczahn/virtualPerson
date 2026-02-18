package output

import (
	"fmt"
	"io"
	"sync"
)

// ANSI color codes for each source.
const (
	colorReset   = "\033[0m"
	colorCyan    = "\033[36m" // SENSE
	colorYellow  = "\033[33m" // BIO
	colorBlue    = "\033[34m" // PSYCH
	colorGreen   = "\033[32m" // MIND
	colorMagenta = "\033[35m" // REVIEW
	colorGray    = "\033[90m" // timestamps
)

var sourceColors = [...]string{
	colorCyan,    // Sense
	colorYellow,  // Bio
	colorBlue,    // Psych
	colorGreen,   // Mind
	colorMagenta, // Review
}

// Display formats and writes simulation output to a writer.
type Display struct {
	mu       sync.Mutex
	writer   io.Writer
	useColor bool
	listener func(Entry)
}

// SetListener sets a callback that is invoked for every entry shown.
// The listener is called under the display's lock, so it must not call
// back into Display methods.
func (d *Display) SetListener(fn func(Entry)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listener = fn
}

// NewDisplay creates a display that writes to the given writer.
func NewDisplay(w io.Writer, useColor bool) *Display {
	return &Display{
		writer:   w,
		useColor: useColor,
	}
}

// Show formats and writes a single entry.
func (d *Display) Show(entry Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ts := entry.Timestamp.Format("15:04:05")
	tag := entry.Source.String()

	if d.useColor {
		tagColor := sourceColor(entry.Source)
		fmt.Fprintf(d.writer, "%s%s%s %s[%-6s]%s %s\n",
			colorGray, ts, colorReset,
			tagColor, tag, colorReset,
			entry.Message)
	} else {
		fmt.Fprintf(d.writer, "%s [%-6s] %s\n", ts, tag, entry.Message)
	}

	if d.listener != nil {
		d.listener(entry)
	}
}

// ShowThought formats a multi-line consciousness thought with a border.
func (d *Display) ShowThought(entry Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	ts := entry.Timestamp.Format("15:04:05")

	mindLabel := Mind.String()
	if d.useColor {
		fmt.Fprintf(d.writer, "%s%s%s %s[%-6s]%s %s\n",
			colorGray, ts, colorReset,
			colorGreen, mindLabel, colorReset,
			entry.Message)
	} else {
		fmt.Fprintf(d.writer, "%s [%-6s] %s\n", ts, mindLabel, entry.Message)
	}

	if d.listener != nil {
		d.listener(entry)
	}
}

func sourceColor(s Source) string {
	idx := int(s)
	if idx < len(sourceColors) {
		return sourceColors[idx]
	}
	return colorReset
}
