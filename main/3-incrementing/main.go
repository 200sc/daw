package main

import (
	"io"

	"github.com/200sc/daw"
)

func main() {
	daw.VisualMain(func(w io.Writer) {
		writeIncrementing(w, 100000)
	})
}

func writeIncrementing(w io.Writer, size int) {
	data := make([]byte, size)
	for i := 0; i < size/4; i++ {
		data[i] = 10
	}
	for i := size / 4; i < size/2; i++ {
		data[i] = 40
	}
	for i := size / 2; i < 3*size/4; i++ {
		data[i] = 70
	}
	for i := 3 * size / 4; i < size; i++ {
		data[i] = 100
	}
	w.Write(data)
}
