package sane

import "github.com/maitredede/puregolibs/sane/internal"

const (
	SANE_FIXED_SCALE_SHIFT = 16
)

func intFromSane(i internal.SANE_Word) int {
	return int(i)
}

func floatFromSane(f internal.SANE_Word) float64 {
	return float64(f) / (1 << SANE_FIXED_SCALE_SHIFT)
}

func floatToSane(f float64) internal.SANE_Word {
	return internal.SANE_Word(f * (1 << SANE_FIXED_SCALE_SHIFT))
}
