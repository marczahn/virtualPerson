package biology

import "testing"

func TestCircadianBodyTemp_PeakEvening_TroughMorning(t *testing.T) {
	evening := circadianBodyTemp(17)   // 5 PM — peak
	earlyMorning := circadianBodyTemp(5) // 5 AM — trough

	if evening <= earlyMorning {
		t.Errorf("expected evening body temp (%f) > early morning (%f)", evening, earlyMorning)
	}

	// Peak should be around 37.1°C
	if evening < 37.0 || evening > 37.2 {
		t.Errorf("evening body temp = %f, expected ~37.1", evening)
	}

	// Trough should be around 36.1°C
	if earlyMorning < 36.0 || earlyMorning > 36.3 {
		t.Errorf("early morning body temp = %f, expected ~36.1", earlyMorning)
	}
}

func TestCircadianBloodPressure_MorningSurge(t *testing.T) {
	morning := circadianBloodPressure(10)   // 10 AM — peak
	nighttime := circadianBloodPressure(22) // 10 PM — trough

	if morning <= nighttime {
		t.Errorf("expected morning BP shift (%f) > nighttime (%f)", morning, nighttime)
	}

	// Morning surge should be positive
	if morning <= 0 {
		t.Errorf("expected positive morning BP shift, got %f", morning)
	}

	// Nighttime should be negative (dipping pattern)
	if nighttime >= 0 {
		t.Errorf("expected negative nighttime BP shift, got %f", nighttime)
	}
}

func TestCircadianImmune_StrongestAtNight(t *testing.T) {
	midnight := circadianImmune(0)  // midnight — peak
	morning := circadianImmune(8)    // 8 AM — approaching trough

	if midnight <= morning {
		t.Errorf("expected midnight immune (%f) > morning (%f)", midnight, morning)
	}

	// Midnight should be above 1.0 (enhanced)
	if midnight <= 1.0 {
		t.Errorf("expected midnight immune > 1.0, got %f", midnight)
	}
}

func TestCircadianSerotonin_HigherDuringDay(t *testing.T) {
	daytime := circadianSerotonin(13) // 1 PM — peak
	night := circadianSerotonin(1)     // 1 AM — trough

	if daytime <= night {
		t.Errorf("expected daytime serotonin shift (%f) > nighttime (%f)", daytime, night)
	}

	if daytime <= 0 {
		t.Errorf("expected positive daytime serotonin shift, got %f", daytime)
	}

	if night >= 0 {
		t.Errorf("expected negative nighttime serotonin shift, got %f", night)
	}
}

func TestCircadianCortisol_TroughAtNight(t *testing.T) {
	peak := circadianCortisol(7)     // 7 AM
	trough := circadianCortisol(23)   // 11 PM

	// Peak should be around 0.35-0.40
	if peak < 0.30 || peak > 0.45 {
		t.Errorf("morning cortisol = %f, expected 0.30-0.45", peak)
	}

	// Trough should be around 0.15 (no negative cos contribution, base only)
	if trough < 0.10 || trough > 0.20 {
		t.Errorf("night cortisol = %f, expected 0.10-0.20", trough)
	}
}

func TestCircadianAlertness_NighttimeLow(t *testing.T) {
	night := circadianAlertness(4) // 4 AM

	if night > 0.2 {
		t.Errorf("expected 4 AM alertness to be low (<0.2), got %f", night)
	}
}

func TestComputeCircadian_ReturnsAllFields(t *testing.T) {
	mod := ComputeCircadian(12)

	// Just verify all fields are non-zero (they should have some circadian signal).
	if mod.CortisolBaseline == 0 {
		t.Error("CortisolBaseline should not be zero at noon")
	}
	if mod.BodyTempTarget == 0 {
		t.Error("BodyTempTarget should not be zero")
	}
	// BloodPressureShift and SerotoninShift can be zero at specific phases,
	// but Alertness at noon should be non-zero.
	if mod.Alertness == 0 {
		t.Error("Alertness should not be zero at noon")
	}
}
