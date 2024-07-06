package main

import (
	"io"

	"github.com/200sc/daw"
)

func main() {
	daw.VisualMain(func(w io.Writer) {
		writeIncrementing(w, 10000)
	})
}

func writeIncrementing(w io.Writer, size int) {
	data := make([]byte, size)
	for i := 0; i < size/4; i++ {
		data[i] = 1
	}
	for i := size / 4; i < size/2; i++ {
		data[i] = 10
	}
	for i := size / 2; i < 3*size/4; i++ {
		data[i] = 20
	}
	for i := 3 * size / 4; i < size; i++ {
		data[i] = 30
	}
	w.Write(data)
}
