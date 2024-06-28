package main

import (
	"bufio"
	"os"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/synth"
)

func main() {
	format := daw.DefaultFormat
	viz := daw.VisualWriter(format)

	pitch := new(daw.Pitch)
	*pitch = daw.C5

	go daw.Loop(viz, &daw.PitchReader{
		Format:   format,
		Pitch:    pitch,
		Volume:   0.05,
		WaveFunc: daw.SinFunc,
	})
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "up":
			*pitch = (*pitch).Up(synth.HalfStep)
		case "down":
			*pitch = (*pitch).Down(synth.HalfStep)
		}
	}
}
