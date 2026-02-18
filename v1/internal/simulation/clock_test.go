package simulation

import (
	"testing"
	"time"
)

func TestNewClock_InitializesCorrectly(t *testing.T) {
	simStart := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	c := NewClock(simStart)

	if c.Paused() {
		t.Error("new clock should not be paused")
	}

	now := c.Now()
	// Should be very close to simStart (within a few ms of creation).
	diff := now.Sub(simStart)
	if diff < 0 || diff > 50*time.Millisecond {
		t.Errorf("Now() should be close to simStart, got diff=%v", diff)
	}
}

func TestClock_Tick_ReturnsPositiveDt(t *testing.T) {
	c := NewClock(time.Now())

	// Sleep briefly so there's measurable elapsed time.
	time.Sleep(10 * time.Millisecond)
	dt := c.Tick()

	if dt <= 0 {
		t.Errorf("Tick() should return positive dt, got %f", dt)
	}
	if dt > 1.0 {
		t.Errorf("Tick() returned unreasonably large dt=%f", dt)
	}
}

func TestClock_Tick_WhenPaused_ReturnsZero(t *testing.T) {
	c := NewClock(time.Now())
	c.Pause()

	time.Sleep(10 * time.Millisecond)
	dt := c.Tick()

	if dt != 0 {
		t.Errorf("Tick() while paused should return 0, got %f", dt)
	}
}

func TestClock_Pause_StopsTime(t *testing.T) {
	simStart := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	c := NewClock(simStart)

	time.Sleep(10 * time.Millisecond)
	c.Pause()
	pausedTime := c.Now()

	time.Sleep(20 * time.Millisecond)
	afterWait := c.Now()

	if !pausedTime.Equal(afterWait) {
		t.Errorf("time should not advance while paused: pausedTime=%v, afterWait=%v", pausedTime, afterWait)
	}
}

func TestClock_Resume_ContinuesFromPausePoint(t *testing.T) {
	simStart := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	c := NewClock(simStart)

	time.Sleep(10 * time.Millisecond)
	c.Pause()
	pausedTime := c.Now()

	time.Sleep(50 * time.Millisecond) // This time should not count.
	c.Resume()

	time.Sleep(10 * time.Millisecond)
	resumed := c.Now()

	// resumed should be slightly after pausedTime (the ~10ms after Resume),
	// NOT 50ms+ ahead (which would include the paused duration).
	diff := resumed.Sub(pausedTime)
	if diff < 0 {
		t.Errorf("time went backwards after resume: diff=%v", diff)
	}
	if diff > 30*time.Millisecond {
		t.Errorf("paused duration leaked into sim time: diff=%v (expected ~10ms)", diff)
	}
}

func TestClock_Resume_TickReturnsPositiveAfterResume(t *testing.T) {
	c := NewClock(time.Now())
	c.Pause()
	time.Sleep(10 * time.Millisecond)
	c.Resume()

	time.Sleep(10 * time.Millisecond)
	dt := c.Tick()

	if dt <= 0 {
		t.Errorf("Tick() after resume should return positive dt, got %f", dt)
	}
}

func TestClock_DoublePause_IsIdempotent(t *testing.T) {
	c := NewClock(time.Now())
	c.Pause()
	time.Sleep(10 * time.Millisecond)
	c.Pause() // second pause should be a no-op

	if !c.Paused() {
		t.Error("should still be paused after double Pause()")
	}
}

func TestClock_DoubleResume_IsIdempotent(t *testing.T) {
	c := NewClock(time.Now())
	// Resume when not paused should be a no-op.
	c.Resume()

	if c.Paused() {
		t.Error("should not be paused after Resume() on non-paused clock")
	}
}

func TestClock_ElapsedSim_MatchesNowMinusStart(t *testing.T) {
	simStart := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	c := NewClock(simStart)

	time.Sleep(10 * time.Millisecond)
	elapsed := c.ElapsedSim()

	if elapsed <= 0 {
		t.Errorf("ElapsedSim() should be positive, got %v", elapsed)
	}
	if elapsed > 100*time.Millisecond {
		t.Errorf("ElapsedSim() unreasonably large: %v", elapsed)
	}
}

func TestClock_ElapsedSim_ExcludesPausedTime(t *testing.T) {
	simStart := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	c := NewClock(simStart)

	time.Sleep(10 * time.Millisecond)
	c.Pause()
	time.Sleep(50 * time.Millisecond) // should not count
	c.Resume()
	time.Sleep(10 * time.Millisecond)

	elapsed := c.ElapsedSim()
	// Should be ~20ms (2x10ms active), not ~70ms (total wall time).
	if elapsed > 40*time.Millisecond {
		t.Errorf("ElapsedSim() should exclude paused time, got %v", elapsed)
	}
}
