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
	// try lowering this
	//mod := int32(12800000)
	//mod := int32(6400000)
	//mod := int32(3200000)
	mod := int32(3900000)
	samples := make([]int32, len(data)/4)
	volume := int32(math.MaxInt32 / 4)
	for i := range samples {
		samples[i] = v
		v += mod
		if v > volume || v < -volume {
			mod *= -1
		}
	}

	bytesPerI32 := int(format.Channels) * 4
	j := 0
	for i := 0; i+bytesPerI32 <= len(data); i += bytesPerI32 {
		i32 := samples[j]
		j++
		for c := 0; c < int(format.Channels); c++ {
			// silence a channel
			// if c == 1 {
			// 	i32 = 0
			// }
			data[i+(4*c)] = byte(i32)
			data[i+(4*c)+1] = byte(i32 >> 8)
			data[i+(4*c)+2] = byte(i32 >> 16)
			data[i+(4*c)+3] = byte(i32 >> 24)
		}
	}

	ch := make(chan daw.Writer)
	go func() {
		w := <-ch
		w.WritePCM(data)
	}()
	daw.VisualWriter(format, ch)
}
