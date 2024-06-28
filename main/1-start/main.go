package main

import (
	"io"

	"github.com/200sc/daw"
)

func main() {
	// DAW = Digital Audio Workstation
	daw.Main(func(w io.Writer) {
		w.Write([]byte{1})
	})
}
