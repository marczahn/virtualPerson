package consciousness

func clamp01(v float64) float64 {
	return clamp(v, 0, 1)
}

func clampSigned(v float64) float64 {
	return clamp(v, -1, 1)
}

func max0(v float64) float64 {
	if v < 0 {
		return 0
	}
	return v
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
