package main

import (
	"context"
	"math"
	"time"

	"github.com/200sc/daw"
	"github.com/oakmound/oak/v4/audio/pcm"
)

type Note struct {
	Pitch    daw.Pitch
	Duration time.Duration
}

const beatsPerMinute = 116

const (
	sixteenthNote = 1
	eighthNote    = 2
	quarterNote   = 4
	halfNote      = 8
	wholeNote     = 16
)

// assuming 4/4
// beatQuarterNote * beatsPerMinute = time.Minute
const sixteenthInterval = (time.Minute / quarterNote) / beatsPerMinute

func beatToDuration(beat int) time.Duration {
	return sixteenthInterval * time.Duration(beat)
}

func chordNotes(root daw.Pitch, chord daw.Chord, beats int) []Note {
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
	format := daw.DefaultFormat

	notes := [][]Note{
		// Measure
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote),
		chordNotes(daw.A5, daw.MajorTriad, quarterNote),
		rest(eighthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote),
		// Measure
		chordNotes(daw.A5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.B5, daw.MinorTriad, eighthNote),
		chordNotes(daw.B5, daw.MinorTriad, quarterNote),
		chordNotes(daw.A5, daw.MajorTriad, quarterNote),
		// Measure
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote),
		chordNotes(daw.A5, daw.MajorTriad, quarterNote),
		chordNotes(daw.G5, daw.MajorTriad, eighthNote),
		chordNotes(daw.D5, daw.MajorTriad, eighthNote+wholeNote), // key I chord
		// Measure
		// (Measure in whole note above)
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote),
		chordNotes(daw.A5, daw.MajorTriad, quarterNote),
		rest(eighthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote),
		// Measure
		chordNotes(daw.A5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.D6, daw.MajorTriad, eighthNote),
		chordNotes(daw.A5, daw.Chord{daw.Perfect4, daw.Major6}, quarterNote),
		chordNotes(daw.A5, daw.MajorTriad, quarterNote),
		// Measure
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.G5, daw.MajorTriad, eighthNote+sixteenthNote),
		chordNotes(daw.A5, daw.MajorTriad, eighthNote),
		chordNotes(daw.A5, daw.MajorTriad, quarterNote),
		chordNotes(daw.C5s, daw.Chord{daw.Minor3, daw.Minor6}, quarterNote),
		// Measure
		chordNotes(daw.D5s, daw.Chord{daw.Minor3, daw.Minor6}, wholeNote),
	}

	for _, ns := range notes {
		// assumption; all notes within a chord have same duration
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
					return pr.volume - (pr.volume / math.Pi * daw.ModPhase(*pr.pitch, pr.phase, pr.SampleRate))

					//f := math.Sin(daw.ModPhase(*pr.pitch, pr.phase, pr.Format.SampleRate))
					//return f * pr.volume
				},
			}
			w := daw.NewWriter()
			go daw.LoopContext(ctx, w, pr)
		}
		time.Sleep(ns[0].Duration)
	}
}

type pitchReader struct {
	pitch    *daw.Pitch
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
