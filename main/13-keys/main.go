package main

import (
	"context"
	"time"

	"github.com/200sc/daw"
)

func main() {
	format := daw.DefaultFormat

	root := daw.D5
	key := daw.Key{
		Start:   root,
		Pattern: daw.MajorKey,
	}

	pitches := key.Scale()
	for _, pitch := range pitches {
		pr := &daw.PitchReader{
			Format:   format,
			Pitch:    &pitch,
			Volume:   0.05,
			WaveFunc: daw.SinFunc,
		}
		w := daw.NewWriter()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		go daw.LoopContext(ctx, w, pr)
		time.Sleep(230 * time.Millisecond)
		cancel()
	}
	for i := len(pitches) - 1; i >= 0; i-- {
		pitch := pitches[i]
		pr := &daw.PitchReader{
			Format:   format,
			Pitch:    &pitch,
			Volume:   0.05,
			WaveFunc: daw.SinFunc,
		}
		w := daw.NewWriter()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		go daw.LoopContext(ctx, w, pr)
		time.Sleep(230 * time.Millisecond)
		cancel()
	}
}
