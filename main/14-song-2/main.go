package main

import (
	"context"
	"time"

	"github.com/200sc/daw"
)

const beatsPerMinute = 116

const (
	sixteenthNote = 1
	eighthNote    = 2
	quarterNote   = 4
	halfNote      = 8
	wholeNote     = 16
)

// assuming 4/4
// quarterNote * beatsPerMinute = time.Minute
const quarterInterval = time.Minute / beatsPerMinute
const sixteenthInterval = quarterInterval / 4

func beatToDuration(beat int) time.Duration {
	return sixteenthInterval * time.Duration(beat)
}

type Note struct {
	Pitch    daw.Pitch
	Duration time.Duration
}

func rest(beats int) []Note {
	return []Note{
		{
			// No pitch == rest
			Duration: beatToDuration(beats),
		},
	}
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
			pr := &daw.PitchReader{
				Format:   format,
				Pitch:    &pitch,
				Volume:   0.05,
				WaveFunc: daw.SawFunc,
			}
			w := daw.NewWriter()
			go daw.LoopContext(ctx, w, pr)
		}
		time.Sleep(ns[0].Duration)
	}
}
