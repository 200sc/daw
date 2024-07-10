package main

import (
	"bufio"
	"math"
	"os"
	"time"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/synth"
)

func main() {
	format := daw.DefaultFormat

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

	ch := make(chan daw.Writer)
	go func() {
		w := <-ch
		go daw.Loop(w, pr)
		time.Sleep(10 * time.Second)
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			switch scanner.Text() {
			case "up":
				*pitch = (*pitch).Up(synth.HalfStep)
			case "down":
				*pitch = (*pitch).Down(synth.HalfStep)
			case "triangle":
				pr.WaveFunc = func(pr *daw.PitchReader) float64 {
					p := daw.ModPhase(*pr.Pitch, pr.Phase, pr.Format.SampleRate)
					m := p * (2 * pr.Volume / math.Pi)
					if math.Sin(p) > 0 {
						return -pr.Volume + m
					}
					return 3*pr.Volume - m
				}
			case "square":
				// pulse with ratio of 2
				pr.WaveFunc = func(pr *daw.PitchReader) float64 {
					if math.Sin(daw.Phase(*pr.Pitch, pr.Phase, pr.SampleRate)) > 0 {
						return pr.Volume
					}
					return -pr.Volume
				}
			case "saw":
				pr.WaveFunc = func(pr *daw.PitchReader) float64 {
					p := daw.ModPhase(*pr.Pitch, pr.Phase, pr.SampleRate)
					return pr.Volume - (pr.Volume / math.Pi * p)
				}
			case "sin":
				pr.WaveFunc = func(pr *daw.PitchReader) float64 {
					f := math.Sin(daw.ModPhase(*pr.Pitch, pr.Phase, pr.Format.SampleRate))
					return f * pr.Volume
				}
			}
		}
	}()
	daw.VisualWriter(format, ch)
}
