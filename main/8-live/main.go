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

	go digitalaudio.Loop(viz, &pitchReader{
		Format: format,
		pitch:  pitch,
	})
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "up":
			*pitch = (*pitch).Up(synth.HalfStep)
		case "down":
			*pitch = (*pitch).Down(synth.HalfStep)
		}
	}

}

type pitchReader struct {
	pitch *digitalaudio.Pitch
	phase int
	pcm.Format
}

func (pr *pitchReader) nextI32() int32 {
	pr.phase++
	v := math.Sin(digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.Format.SampleRate))
	return digitalaudio.VolumeI32(v, .05)
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
