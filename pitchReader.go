package daw

import (
	"math"

	"github.com/oakmound/oak/v4/audio/pcm"
)

type PitchReader struct {
	Pitch    *Pitch
	Phase    int
	WaveFunc func(*PitchReader) float64
	Volume   float64
	pcm.Format
}

func (pr *PitchReader) nextI32() int32 {
	pr.Phase++
	return int32(pr.WaveFunc(pr) * math.MaxInt32)
}

func (pr *PitchReader) ReadPCM(data []byte) (n int, err error) {
	bytesPerI32 := int(pr.Format.Channels) * 4
	for i := 0; i+bytesPerI32 <= len(data); i += bytesPerI32 {
		i32 := pr.nextI32()
		for c := 0; c < int(pr.Format.Channels); c++ {
			data[i+(4*c)] = byte(i32)
			data[i+(4*c)+1] = byte(i32 >> 8)
			data[i+(4*c)+2] = byte(i32 >> 16)
			data[i+(4*c)+3] = byte(i32 >> 24)
		}
		n += 4 * int(pr.Format.Channels)
	}
	return n, nil
}
