package main

import (
	"math/rand"
	"time"

	"github.com/200sc/daw"
)

func main() {
	data := make([]byte, daw.BufferLength(daw.DefaultFormat))
	for i := range data {
		data[i] = byte((rand.Float64() - .5) * 5)
	}
	viz := daw.VisualWriter(daw.DefaultFormat)
	viz.WritePCM(data)
	time.Sleep(5 * time.Second)
}

// SampleRate: 44100,
// Channels:   2,
// Bits:       32,
