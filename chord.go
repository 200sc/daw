package digitalaudio

// A chord is a series of intervals implicitly including the root interval
type Chord []Step

var (
	MajorTriad       Chord = []Step{Major3, Perfect5}
	MajorSixth       Chord = []Step{Major3, Perfect5, Major6}
	DominantSeventh  Chord = []Step{Major3, Perfect5, Minor7}
	AugmentedTriad   Chord = []Step{Major3, Minor6, Major7}
	AugmentedSeventh Chord = []Step{Major3, Minor6, Minor7}

	MinorTriad            Chord = []Step{Minor3, Perfect5}
	MinorSixth            Chord = []Step{Minor3, Perfect5, Major6}
	MinorSeventh          Chord = []Step{Minor3, Perfect5, Minor7}
	MinorMajorSeventh     Chord = []Step{Minor3, Perfect5, Major7}
	DiminishedTriad       Chord = []Step{Minor3, Tritone}
	DiminishedSeventh     Chord = []Step{Minor3, Tritone, Major6}
	HalfDiminishedSeventh Chord = []Step{Minor3, Tritone, Minor7}
)

func (c Chord) WithRoot(root Pitch) []Pitch {
	out := []Pitch{root}
	for _, s := range c {
		out = append(out, root.Up(s))
	}
	return out
}
