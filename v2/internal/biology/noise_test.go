package biology_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/marczahn/person/v2/internal/biology"
)

func TestApplyNoise_DtNonPositive_NoChange(t *testing.T) {
	s := biology.NewDefaultState()
	before := *s
	rng := rand.New(rand.NewSource(42))
	cfg := biology.NoiseConfig{Sigma: 0.5}

	biology.ApplyNoise(s, rng, cfg, 0)
	biology.ApplyNoise(s, rng, cfg, -1)

	if *s != before {
		t.Fatalf("state changed with non-positive dt: got %+v want %+v", *s, before)
	}
}

func TestApplyNoise_UsesSqrtDtAndBodyTempScale(t *testing.T) {
	seed := int64(7)
	dt := 4.0
	cfg := biology.NoiseConfig{Sigma: 0.002}
	s := biology.NewDefaultState()
	before := *s

	biology.ApplyNoise(s, rand.New(rand.NewSource(seed)), cfg, dt)

	sigma := cfg.Sigma * math.Sqrt(dt)
	expRng := rand.New(rand.NewSource(seed))
	wantEnergy := before.Energy + expRng.NormFloat64()*sigma
	wantStress := before.Stress + expRng.NormFloat64()*sigma
	wantCog := before.CognitiveCapacity + expRng.NormFloat64()*sigma
	wantMood := before.Mood + expRng.NormFloat64()*sigma
	wantTension := before.PhysicalTension + expRng.NormFloat64()*sigma
	wantHunger := before.Hunger + expRng.NormFloat64()*sigma
	wantSocial := before.SocialDeficit + expRng.NormFloat64()*sigma
	wantTemp := before.BodyTemp + expRng.NormFloat64()*sigma*0.1

	checks := []struct {
		name string
		got  float64
		want float64
	}{
		{"Energy", s.Energy, wantEnergy},
		{"Stress", s.Stress, wantStress},
		{"CognitiveCapacity", s.CognitiveCapacity, wantCog},
		{"Mood", s.Mood, wantMood},
		{"PhysicalTension", s.PhysicalTension, wantTension},
		{"Hunger", s.Hunger, wantHunger},
		{"SocialDeficit", s.SocialDeficit, wantSocial},
		{"BodyTemp", s.BodyTemp, wantTemp},
	}

	for _, c := range checks {
		if math.Abs(c.got-c.want) > 1e-12 {
			t.Fatalf("%s = %.15f want %.15f", c.name, c.got, c.want)
		}
	}
}
