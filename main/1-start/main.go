package main

import (
	"io"

	// DAW = Digital Audio Workstation
	"github.com/200sc/daw"
)

// We know: we can write binary data to a speaker

func main() {
	daw.Main(func(w io.Writer) {
		w.Write([]byte{5})
	})
}
