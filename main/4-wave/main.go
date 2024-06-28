package main

import (
	"time"

	"github.com/200sc/daw"
)

func main() {
	data := make([]byte, daw.BufferLength(daw.DefaultFormat))
	v := 0
	dir := 1
	for i := range data {
		data[i] = byte(v)
		v += dir
		if v > 50 || v < -50 {
			dir *= -1
		}
	}
	viz := daw.VisualWriter(daw.DefaultFormat)
	viz.WritePCM(data)
	time.Sleep(5 * time.Second)
}
