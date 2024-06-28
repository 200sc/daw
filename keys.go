package digitalaudio

import (
	"github.com/oakmound/oak/v4/audio/synth"
)

type Key struct {
	Start   Pitch
	Pattern KeyPattern
}

type KeyPattern []Step

var MajorKey KeyPattern = []Step{
	synth.WholeStep,
	synth.WholeStep,
	synth.HalfStep,
	synth.WholeStep,
	synth.WholeStep,
	synth.WholeStep,
	synth.HalfStep,
}

var MinorKey KeyPattern = []Step{
	synth.WholeStep,
	synth.HalfStep,
	synth.WholeStep,
	synth.WholeStep,
	synth.HalfStep,
	synth.WholeStep,
	synth.WholeStep,
}

// Note also: Harmonic minor, melodic minor (and usage when descending / ascending)

var C5Major = Key{
	Start:   C5,
	Pattern: MajorKey,
}

func (k Key) Scale() []Pitch {
	ps := []Pitch{k.Start}
	next := k.Start
	for _, s := range k.Pattern {
		next = next.Up(s)
		ps = append(ps, next)
	}
	return ps
}
