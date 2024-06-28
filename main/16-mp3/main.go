package main

import (
	"context"
	"os"

	"github.com/oakmound/oak/v4/audio"
	"github.com/oakmound/oak/v4/audio/format/mp3"
)

func main() {
	f, err := os.Open("ghosts.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r, err := mp3.Load(f)
	if err != nil {
		panic(err)
	}
	audio.InitDefault()
	err = audio.Play(context.Background(), r)
	if err != nil {
		panic(err)
	}
}

