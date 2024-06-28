package digitalaudio

import "github.com/oakmound/oak/v4/audio/synth"

type Step = synth.Step

const (
	Perfect1  Step = 0
	HalfStep  Step = 1
	Minor2    Step = 1
	WholeStep Step = 2
	Major2    Step = 2
	Minor3    Step = 3
	Major3    Step = 4
	Perfect4  Step = 5
	Tritone   Step = 6
	Perfect5  Step = 7
	Minor6    Step = 8
	Major6    Step = 9
	Minor7    Step = 10
	Major7    Step = 11
	Octave    Step = 12
)
