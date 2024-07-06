package main

import (
	"io"

	"github.com/200sc/daw"
)

func main() {
	daw.Main(func(w io.Writer) {
		size := 100000
		data := make([]byte, size)
		for i := 0; i < size; i++ {
			data[i] = 5
		}
		w.Write(data)
	})
}
