package main

import (
	"fmt"
	"math"
	"time"

	digitalaudio "github.com/200sc/digital-audio"
)

func main() {
	format := digitalaudio.DefaultFormat
	viz := digitalaudio.VisualWriter(format)

	data := make([]byte, digitalaudio.BufferLength(format))
	pitch := digitalaudio.C4
	samples := make([]int32, len(data)/4)
	for i := range samples {
		v := math.Sin(modPhase(pitch, i, format.SampleRate))
		samples[i] = int32((v * .05) * math.MaxInt32)
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

	viz.WritePCM(data)
	fmt.Println(data[:100])
	time.Sleep(5 * time.Second)
}

func phase(freq digitalaudio.Pitch, i int, sampleRate uint32) float64 {
	return float64(freq) * (float64(i) / float64(sampleRate)) * 2 * math.Pi
}

func modPhase(freq digitalaudio.Pitch, i int, sampleRate uint32) float64 {
	return math.Mod(phase(freq, i, sampleRate), 2*math.Pi)
}

// *= math.MaxInt32
// what happens if we execute?
// divide by 20
// phase -> modphase
// sin
