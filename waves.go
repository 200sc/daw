package daw

import "math"

var SinFunc = func(pr *PitchReader) float64 {
	return math.Sin(ModPhase(*pr.Pitch, pr.Phase, pr.Format.SampleRate)) * pr.Volume
}

var SawFunc = func(pr *PitchReader) float64 {
	return pr.Volume - (pr.Volume / math.Pi * ModPhase(*pr.Pitch, pr.Phase, pr.SampleRate))
}
