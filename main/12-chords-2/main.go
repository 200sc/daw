package main

import (
	"time"

	"github.com/200sc/daw"
)

func main() {
	format := daw.DefaultFormat

	pitches := daw.MinorMajorSeventh.WithRoot(daw.C5)
	for _, pitch := range pitches {
		pitch := pitch
		pr := &daw.PitchReader{
			Format:   format,
			Pitch:    &pitch,
			Volume:   0.05,
			WaveFunc: daw.SinFunc,
		}
		w := daw.NewWriter()
		go daw.Loop(w, pr)
	}
	time.Sleep(10 * time.Second)
}
