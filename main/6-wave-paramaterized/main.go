package main

import (
	"math"

	"github.com/200sc/daw"
)

// Volume = Wave Height = Amplitude
// Pitch = Pattern Size = Frequency

func main() {
	format := daw.DefaultFormat

	data := make([]byte, daw.BufferLength(format))
	v := int32(0)
	delta := int32(12800000) // 6400000, 7600000
	samples := make([]int32, len(data)/4)
	volume := int32(math.MaxInt32 / 4)
	for i := range samples {
		samples[i] = v
		v += delta
		if v > volume || v < -volume {
			delta *= -1
		}
	}

	const i32Size = 4
	bytesPerI32 := int(format.Channels) * i32Size
	j := 0
	for i := 0; i+bytesPerI32 <= len(data); i += bytesPerI32 {
		i32 := samples[j] // []int32
		j++
		for c := 0; c < int(format.Channels); c++ {
			// if c == 1 {
			// 	i32 = 0 // to silence channel 1
			// }
			data[i+(i32Size*c)] = byte(i32)
			data[i+(i32Size*c)+1] = byte(i32 >> 8)
			data[i+(i32Size*c)+2] = byte(i32 >> 16)
			data[i+(i32Size*c)+3] = byte(i32 >> 24)
		}
	}

	ch := make(chan daw.Writer)
	go func() {
		w := <-ch
		w.WritePCM(data)
	}()
	daw.VisualWriter(format, ch)
}
