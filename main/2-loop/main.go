package main

import (
	"time"

	digitalaudio "github.com/200sc/digital-audio"
)

func main() {
	// 1, 10, 100
	size := 100000
	data := make([]byte, size)
	for i := 0; i < size/4; i++ {
		data[i] = 1
	}
	for i := size / 4; i < size/2; i++ {
		data[i] = 10
	}
	for i := size / 2; i < 3*size/4; i++ {
		data[i] = 50
	}
	for i := 3 * size / 4; i < size; i++ {
		data[i] = 75
	}
	viz := digitalaudio.VisualWriter(digitalaudio.DefaultFormat)
	viz.WritePCM(data)
	time.Sleep(5 * time.Second)
}
