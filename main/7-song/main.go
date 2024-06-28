package main

import (
	"bytes"
	"math"
	"time"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/pcm"
)

func main() {
	format := daw.DefaultFormat
	viz := daw.VisualWriter(format)

	data := make([]byte, 44100*2*10)

	//data := make([]byte, daw.BufferLength(format))
	pitches := []daw.Pitch{
		daw.D4,
		daw.E4,
		daw.F4,
		daw.F4,
		daw.E4,
		daw.E4,
		daw.E4,
		daw.E4,
		daw.E4,
		daw.E4,
	}
	vals := make([]int32, len(data)/8)
	i := 0
	for _, pitch := range pitches {
		pitch := pitch
		for j := 0; j < len(vals)/len(pitches); j++ {
			v := math.Sin(daw.ModPhase(pitch, i, format.SampleRate))
			vals[i] = daw.VolumeI32(v, .05)
			i++
		}
	}

	bytesPerI32 := int(format.Channels) * 4
	j := 0
	for i := 0; i+bytesPerI32 <= len(data); i += bytesPerI32 {
		i32 := vals[j]
		j++
		for c := 0; c < int(format.Channels); c++ {
			data[i+(4*c)] = byte(i32)
			data[i+(4*c)+1] = byte(i32 >> 8)
			data[i+(4*c)+2] = byte(i32 >> 16)
			data[i+(4*c)+3] = byte(i32 >> 24)
		}
	}

	go daw.Loop(viz, &pcm.IOReader{
		Format: format,
		Reader: bytes.NewReader(data),
	})
	time.Sleep(10 * time.Second)
}
