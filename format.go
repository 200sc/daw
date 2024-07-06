package daw

import (
	"github.com/oakmound/oak/v4/audio"
	"github.com/oakmound/oak/v4/audio/pcm"
)

var DefaultFormat = pcm.Format{
	SampleRate: 44100,
	Channels:   2,
	// Bits does not reflect that some writers expect floats, and some expect ints.
	Bits: 32,
}

func NewWriter() Writer {
	return audio.MustNewWriter(DefaultFormat)
}
