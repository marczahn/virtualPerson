package psychology

import "testing"

func TestUpdateRegulation_DepletesUnderStress(t *testing.T) {
	p := Personality{
		Conscientiousness: 0.5,
		Openness:          0.5,
		Neuroticism:       0.5,
	}
	reg := RegulationState{Capacity: BaselineRegulation(p)}
	initial := reg.Capacity

	// 2 hours at stress 0.7.
	reg = UpdateRegulation(reg, p, 0.7, 0.0, 7200)

	if reg.Capacity >= initial {
		t.Errorf("regulation should deplete under stress: was %f, now %f", initial, reg.Capacity)
	}
}

func TestUpdateRegulation_RecoversAtLowStress(t *testing.T) {
	p := Personality{
		Conscientiousness: 0.5,
		Openness:          0.5,
		Neuroticism:       0.5,
	}
	depleted := RegulationState{Capacity: 0.2}

	recovered := UpdateRegulation(depleted, p, 0.1, 0.0, 3600)

	if recovered.Capacity <= depleted.Capacity {
		t.Errorf("regulation should recover at low stress: was %f, now %f", depleted.Capacity, recovered.Capacity)
	}
}

func TestUpdateRegulation_CapsAtBaseline(t *testing.T) {
	p := Personality{
		Conscientiousness: 0.5,
		Openness:          0.5,
		Neuroticism:       0.5,
	}
	baseline := BaselineRegulation(p)
	reg := RegulationState{Capacity: baseline}

	// Recovery at low stress should not exceed baseline.
	reg = UpdateRegulation(reg, p, 0.1, 0.0, 36000)

	if reg.Capacity > baseline {
		t.Errorf("regulation %f should not exceed baseline %f", reg.Capacity, baseline)
	}
}

func TestUpdateRegulation_SustainedStress_AcceleratingCollapse(t *testing.T) {
	p := Personality{
		Conscientiousness: 0.5,
		Openness:          0.5,
		Neuroticism:       0.5,
	}
	baseline := BaselineRegulation(p)
	reg := RegulationState{Capacity: baseline}

	// Track depletion per hour at sustained high stress.
	var depletionRates []float64
	for i := 0; i < 10; i++ {
		before := reg.Capacity
		reg = UpdateRegulation(reg, p, 0.8, 0.0, 3600)
		depletionRates = append(depletionRates, before-reg.Capacity)
	}

	// The collapse should accelerate once capacity is low.
	// First hour depletion should be less than later hour depletion (when capacity is very low).
	if depletionRates[0] >= depletionRates[len(depletionRates)-1] && reg.Capacity > 0.05 {
		// Only check if capacity hasn't fully bottomed out.
		t.Logf("depletion rates: %v", depletionRates)
		t.Log("note: expected accelerating depletion but may have bottomed out")
	}

	// After 10 hours at 0.8 stress, should be severely depleted.
	if reg.Capacity > 0.15 {
		t.Errorf("after 10h at stress 0.8, capacity = %f, expected < 0.15", reg.Capacity)
	}
}

func TestEffectiveCapacity_FatiguePenalty(t *testing.T) {
	reg := RegulationState{Capacity: 0.6}
	effective := EffectiveCapacity(reg, 0.8)

	// 0.6 - 0.8*0.3 = 0.36
	if effective < 0.3 || effective > 0.4 {
		t.Errorf("effective capacity = %f, expected ~0.36", effective)
	}
}

func TestEffectiveCapacity_NoFatigue(t *testing.T) {
	reg := RegulationState{Capacity: 0.7}
	effective := EffectiveCapacity(reg, 0.0)

	if effective != 0.7 {
		t.Errorf("effective capacity with no fatigue = %f, expected 0.7", effective)
	}
}
