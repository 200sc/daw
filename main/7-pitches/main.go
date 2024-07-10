package main

import (
	"math"
	"time"

	daw "github.com/200sc/daw"
)

func phase(freq daw.Pitch, sample int, sampleRate uint32) float64 {
	v := float64(freq) * (float64(sample) / float64(sampleRate)) * 2 * math.Pi
	return math.Mod(v, 2*math.Pi)
}

func main() {
	format := daw.DefaultFormat

	//
	pitch := daw.A4
	//

	data := make([]byte, daw.BufferLength(format))
	samples := make([]int32, len(data)/4)
	volume := .50 * math.MaxInt32
	for i := range samples {
		v := math.Sin(phase(pitch, i, format.SampleRate))
		samples[i] = int32((v * volume))
	}

	bytesPerI32 := int(format.Channels) * 4
	j := 0
	for i := 0; i+bytesPerI32 <= len(data); i += bytesPerI32 {
		i32 := samples[j]
		j++
		for c := 0; c < int(format.Channels); c++ {
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
		time.Sleep(5 * time.Second)
	}()
	daw.VisualWriter(format, ch)
}
