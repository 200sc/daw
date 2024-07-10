package main

import (
	"bufio"
	"os"
	"time"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/synth"
)

func main() {
	format := daw.DefaultFormat

	pitch := new(daw.Pitch)
	*pitch = daw.C5

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			switch scanner.Text() {
			case "up":
				*pitch = (*pitch).Up(synth.HalfStep)
			case "down":
				*pitch = (*pitch).Down(synth.HalfStep)
			}
		}
	}()

	ch := make(chan daw.Writer)
	go func() {
		w := <-ch
		go daw.Loop(w, &daw.PitchReader{
			Format:   format,
			Pitch:    pitch,
			Volume:   0.50,
			WaveFunc: daw.SinFunc,
		})
		time.Sleep(10 * time.Second)
	}()
	daw.VisualWriter(format, ch)
}
