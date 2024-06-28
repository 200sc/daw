package daw

import (
	"math"
)

func Phase(freq Pitch, i int, sampleRate uint32) float64 {
	return float64(freq) * (float64(i) / float64(sampleRate)) * 2 * math.Pi
}

func ModPhase(freq Pitch, i int, sampleRate uint32) float64 {
	return math.Mod(Phase(freq, i, sampleRate), 2*math.Pi)
}

func VolumeI32(v float64, volume float64) int32 {
	return int32((v * volume) * math.MaxInt32)
}
