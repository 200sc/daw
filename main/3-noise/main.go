package main

import (
	"math/rand"
	"time"

	digitalaudio "github.com/200sc/digital-audio"
)

func main() {
	data := make([]byte, digitalaudio.BufferLength(digitalaudio.DefaultFormat))
	for i := range data {
		data[i] = byte((rand.Float64() - .5) * 5)
	}
	viz := digitalaudio.VisualWriter(digitalaudio.DefaultFormat)
	viz.WritePCM(data)
	time.Sleep(5 * time.Second)
}

// SampleRate: 44100,
// Channels:   2,
// Bits:       32,
