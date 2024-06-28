package main

import (
	"math"
	"time"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/synth"
)

func main() {
	format := daw.DefaultFormat
	viz := daw.VisualWriter(format)

	pitch := new(daw.Pitch)
	*pitch = daw.C5

	pr := &daw.PitchReader{
		Format: format,
		Pitch:  pitch,
		Volume: 0.05,
		WaveFunc: func(pr *daw.PitchReader) float64 {
			f := math.Sin(daw.ModPhase(*pr.Pitch, pr.Phase, pr.Format.SampleRate))
			return f * pr.Volume
		},
	}
	// detune down 20%
	pitch2 := *pitch
	halfDown := pitch2.Down(synth.HalfStep)
	rawDelta := float64(int16(pitch2) - int16(halfDown))
	delta := rawDelta * .2
	pitch2 = daw.Pitch(float64(pitch2) + delta)

	pr2 := &daw.PitchReader{
		Format: format,
		Pitch:  &pitch2,
		Volume: 0.05,
		WaveFunc: func(pr *daw.PitchReader) float64 {
			f := math.Sin(daw.ModPhase(*pr.Pitch, pr.Phase, pr.Format.SampleRate))
			return f * pr.Volume
		},
	}
	w2 := daw.NewWriter()
	go daw.Loop(w2, pr2)

	go daw.Loop(viz, pr)
	time.Sleep(20 * time.Second)
}
