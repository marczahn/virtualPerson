package simulation

import (
	"sync"
	"time"
)

// Clock tracks simulation time. Currently 1:1 with real time,
// but the abstraction allows future time-scaling.
type Clock struct {
	mu        sync.Mutex
	startReal time.Time // real time when simulation started
	startSim  time.Time // simulation time when started
	lastTick  time.Time // real time of the last tick
	paused    bool
	pausedAt  time.Time // real time when paused
	pauseAccum time.Duration // total time spent paused
}

// NewClock creates a clock starting at the given simulation time.
// simStart is the in-world time (e.g., 8:00 AM on day 1).
func NewClock(simStart time.Time) *Clock {
	now := time.Now()
	return &Clock{
		startReal: now,
		startSim:  simStart,
		lastTick:  now,
	}
}

// Tick advances the clock and returns the elapsed simulation seconds
// since the last tick. Returns 0 if paused.
func (c *Clock) Tick() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.paused {
		return 0
	}

	now := time.Now()
	dt := now.Sub(c.lastTick)
	c.lastTick = now

	return dt.Seconds()
}

// Now returns the current simulation time.
func (c *Clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.paused {
		elapsed := c.pausedAt.Sub(c.startReal) - c.pauseAccum
		return c.startSim.Add(elapsed)
	}

	elapsed := time.Since(c.startReal) - c.pauseAccum
	return c.startSim.Add(elapsed)
}

// Pause stops the clock. Subsequent Tick() calls return 0.
func (c *Clock) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.paused {
		return
	}
	c.paused = true
	c.pausedAt = time.Now()
}

// Resume restarts the clock after a pause.
func (c *Clock) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.paused {
		return
	}
	c.pauseAccum += time.Since(c.pausedAt)
	c.lastTick = time.Now()
	c.paused = false
}

// Paused returns whether the clock is currently paused.
func (c *Clock) Paused() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.paused
}

// ElapsedSim returns the total elapsed simulation time.
func (c *Clock) ElapsedSim() time.Duration {
	return c.Now().Sub(c.startSim)
}
