package main

import (
	"io"
	"time"

	"github.com/200sc/daw"
)

func main() {
	daw.VisualMain(func(w io.Writer) {
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
		w.Write(data)
		time.Sleep(15 * time.Second)
	})
}
