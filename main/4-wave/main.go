package main

import (
	"time"

	digitalaudio "github.com/200sc/digital-audio"
)

func main() {
	data := make([]byte, digitalaudio.BufferLength(digitalaudio.DefaultFormat))
	v := 0
	dir := 1
	for i := range data {
		data[i] = byte(v)
		v += dir
		if v > 50 || v < -50 {
			dir *= -1
		}
	}
	viz := digitalaudio.VisualWriter(digitalaudio.DefaultFormat)
	viz.WritePCM(data)
	time.Sleep(5 * time.Second)
}
