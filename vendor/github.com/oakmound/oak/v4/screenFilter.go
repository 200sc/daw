package oak

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v4/render/mod"
)

// SetPalette tells oak to conform the screen to the input color palette before drawing.
func (w *Window) SetPalette(palette color.Palette) {
	w.SetDrawFilter(mod.ConformToPalette(palette))
}

// SetDrawFilter will filter the screen by the given modification function prior
// to publishing the screen's rgba to be displayed.
func (w *Window) SetDrawFilter(screenFilter mod.Filter) {
	w.prePublish = func(buf *image.RGBA) {
		screenFilter(buf)
	}
}

// ClearScreenFilter resets the draw function to no longer filter the screen before
// publishing it to the window.
func (w *Window) ClearScreenFilter() {
	w.prePublish = func(buf *image.RGBA) {}
}
