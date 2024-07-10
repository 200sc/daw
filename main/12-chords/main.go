package main

import (
	"time"

	"github.com/200sc/daw"
)

func main() {
	format := daw.DefaultFormat

	pitches := []daw.Pitch{
		daw.C5,
		daw.E5,
		daw.G5,
	}
	// pitches := daw.MinorMajorSeventh.WithRoot(daw.C5)
	for _, pitch := range pitches {
		pitch := pitch
		pr := &daw.PitchReader{
			Format:   format,
			Pitch:    &pitch,
			Volume:   0.50,
			WaveFunc: daw.SinFunc,
		}
		w := daw.NewWriter()
		go daw.Loop(w, pr)
	}
	time.Sleep(5 * time.Second)
}
