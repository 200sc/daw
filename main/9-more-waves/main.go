package main

import (
	"bufio"
	"math"
	"os"

	digitalaudio "github.com/200sc/digital-audio"
	"github.com/oakmound/oak/v4/audio/pcm"
	"github.com/oakmound/oak/v4/audio/synth"
)

func main() {
	format := digitalaudio.DefaultFormat
	viz := digitalaudio.VisualWriter(format)

	pitch := new(digitalaudio.Pitch)
	*pitch = digitalaudio.C5

	pr := &pitchReader{
		Format: format,
		pitch:  pitch,
		volume: 0.05,
		waveFunc: func(pr *pitchReader) float64 {
			f := math.Sin(digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.Format.SampleRate))
			return f * pr.volume
		},
	}

	go digitalaudio.Loop(viz, pr)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "up":
			*pitch = (*pitch).Up(synth.HalfStep)
		case "down":
			*pitch = (*pitch).Down(synth.HalfStep)
		case "triangle":
			pr.waveFunc = func(pr *pitchReader) float64 {
				p := digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.Format.SampleRate)
				m := p * (2 * pr.volume / math.Pi)
				if math.Sin(p) > 0 {
					return -pr.volume + m
				}
				return 3*pr.volume - m
			}
		case "square":
			// pulse with ratio of 2 
			pr.waveFunc = func(pr *pitchReader) float64 {
				if math.Sin(digitalaudio.Phase(*pr.pitch, pr.phase, pr.SampleRate)) > 0 {
					return pr.volume
				}
				return -pr.volume
			}
		case "saw":
			pr.waveFunc = func(pr *pitchReader) float64 {
				return pr.volume - (pr.volume / math.Pi * digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.SampleRate))
			}
		}
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
