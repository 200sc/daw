package main

import (
	"io"
	"time"

	"github.com/200sc/daw"
)

func main() {
	daw.Main(func(w io.Writer) {
		size := 100000
		data := make([]byte, size)
		for i := 0; i < size; i++ {
			// 10, 50
			data[i] = 50
		}
		w.Write(data)
		time.Sleep(5 * time.Second)
	})
}
