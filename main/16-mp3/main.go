package main

import (
	"context"
	"os"

	"github.com/200sc/daw"

	"github.com/oakmound/oak/v4/audio"
	"github.com/oakmound/oak/v4/audio/format/mp3"
)

func main() {
	f, err := os.Open("../16a-mp3/ghosts.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r, err := mp3.Load(f)
	if err != nil {
		panic(err)
	}
	vizw := daw.VisualWriter(r.PCMFormat())
	err = audio.Play(context.Background(), r, func(po *audio.PlayOptions) {
		po.Destination = vizw
	})
	if err != nil {
		panic(err)
	}
}
