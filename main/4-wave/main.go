package main

import (
	"io"
	"time"

	"github.com/200sc/daw"
)

func main() {
	daw.VisualMain(func(w io.Writer) {
		data := make([]byte, daw.BufferLength(daw.DefaultFormat))
		v := 0.0
		delta := 0.5
		const volume = 25
		for i := range data {
			data[i] = byte(v)
			v += delta
			if v > volume || v < -volume {
				delta *= -1
			}
		}
		w.Write(data)
		time.Sleep(10 * time.Second)
	})
}
