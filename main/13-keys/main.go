package main

import (
	"context"
	"math"
	"time"

	digitalaudio "github.com/200sc/digital-audio"
	"github.com/oakmound/oak/v4/audio/pcm"
)

func main() {
	format := digitalaudio.DefaultFormat

	root := digitalaudio.D5
	key := digitalaudio.Key{
		Start:   root,
		Pattern: digitalaudio.MajorKey,
	}

	pitches := key.Scale()
	for _, pitch := range pitches {
		pitch := pitch
		pr := &pitchReader{
			Format: format,
			pitch:  &pitch,
			volume: 0.05,
			waveFunc: func(pr *pitchReader) float64 {
				f := math.Sin(digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.Format.SampleRate))
				return f * pr.volume
			},
		}
		w := digitalaudio.NewWriter()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		go digitalaudio.LoopContext(ctx, w, pr)
		time.Sleep(230 * time.Millisecond)
		cancel()
	}
}

type pitchReader struct {
	pitch    *digitalaudio.Pitch
	phase    int
	volume   float64
	waveFunc func(*pitchReader) float64
	pcm.Format
}

func (pr *pitchReader) nextI32() int32 {
	pr.phase++
	return int32(pr.waveFunc(pr) * math.MaxInt32)
}

func (pr *pitchReader) ReadPCM(data []byte) (n int, err error) {
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