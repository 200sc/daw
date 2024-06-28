package main

import (
	"context"
	"os"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/format/wav"
)

func main() {
	f, err := os.Open("midtown.wav")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r, err := wav.Load(f)
	if err != nil {
		panic(err)
	}

	err = daw.Play(context.Background(), r)
	if err != nil {
		panic(err)
	}
}
