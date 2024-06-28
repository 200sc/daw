package main

import (
	"io"
	"math/rand"
	"time"

	"github.com/200sc/daw"
)

func main() {
	daw.VisualMain(func(w io.Writer) {
		data := make([]byte, daw.BufferLength(daw.DefaultFormat))
		const volume = 50
		for i := range data {
			data[i] = byte((rand.Float64() - .5) * volume)
		}
		w.Write(data)
		time.Sleep(10 * time.Second)
	})
}
