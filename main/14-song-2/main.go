package main

import (
	"context"
	"math"
	"time"

	digitalaudio "github.com/200sc/digital-audio"
	"github.com/oakmound/oak/v4/audio/pcm"
)

type Note struct {
	Pitch    digitalaudio.Pitch
	Duration time.Duration
}

const beatsPerMinute = 116

const (
	beatSixteenthNote = 1
	beatEighthNote    = 2
	beatQuarterNote   = 4
	beatHalfNote      = 8
	beatWholeNote     = 16
)

// assuming 4/4
// beatQuarterNote * beatsPerMinute = time.Minute
const sixteenthInterval = (time.Minute / beatQuarterNote) / beatsPerMinute

func beatToDuration(beat int) time.Duration {
	return sixteenthInterval * time.Duration(beat)
}

func chordNotes(root digitalaudio.Pitch, chord digitalaudio.Chord, beats int) []Note {
	pitches := chord.WithRoot(root)
	notes := make([]Note, len(pitches))
	for i, p := range pitches {
		notes[i] = Note{
			Pitch:    p,
			Duration: beatToDuration(beats),
		}
	}
	return notes
}

func rest(beats int) []Note {
	return []Note{
		{
			// No pitch == rest
			Duration: beatToDuration(beats),
		},
	}
}

func main() {
	format := digitalaudio.DefaultFormat

	notes := [][]Note{
		// Measure
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatQuarterNote),
		rest(beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote),
		// Measure
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.B5, digitalaudio.MinorTriad, beatEighthNote),
		chordNotes(digitalaudio.B5, digitalaudio.MinorTriad, beatQuarterNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatQuarterNote),
		// Measure
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatQuarterNote),
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote),
		chordNotes(digitalaudio.D5, digitalaudio.MajorTriad, beatEighthNote+beatWholeNote), // key I chord
		// Measure
		// (Measure in whole note above)
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatQuarterNote),
		rest(beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote),
		// Measure
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.D6, digitalaudio.MajorTriad, beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.Chord{digitalaudio.Perfect4, digitalaudio.Major6}, beatQuarterNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatQuarterNote),
		// Measure
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.G5, digitalaudio.MajorTriad, beatEighthNote+beatSixteenthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatEighthNote),
		chordNotes(digitalaudio.A5, digitalaudio.MajorTriad, beatQuarterNote),
		chordNotes(digitalaudio.C5s, digitalaudio.Chord{digitalaudio.Minor3, digitalaudio.Minor6}, beatQuarterNote),
		// Measure
		chordNotes(digitalaudio.D5s, digitalaudio.Chord{digitalaudio.Minor3, digitalaudio.Minor6}, beatWholeNote),
	}

	for _, ns := range notes {
		// assumption; all ns have same duration
		ctx, cancel := context.WithTimeout(context.Background(), ns[0].Duration-10*time.Millisecond)
		defer cancel()
		for _, n := range ns {
			pitch := n.Pitch
			if pitch == 0 {
				continue
			}
			pr := &pitchReader{
				Format: format,
				pitch:  &pitch,
				volume: 0.05,
				waveFunc: func(pr *pitchReader) float64 {
					return pr.volume - (pr.volume / math.Pi * digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.SampleRate))

					//f := math.Sin(digitalaudio.ModPhase(*pr.pitch, pr.phase, pr.Format.SampleRate))
					//return f * pr.volume
				},
			}
			w := digitalaudio.NewWriter()
			go digitalaudio.LoopContext(ctx, w, pr)
		}
		time.Sleep(ns[0].Duration)
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
