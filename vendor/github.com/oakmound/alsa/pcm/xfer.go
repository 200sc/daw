package pcm

import (
	"github.com/oakmound/alsa/alsatype"
)

type XferI struct {
	Result alsatype.Sframes
	Buf    uintptr
	Frames alsatype.Uframes
}

const XferISize = 24
