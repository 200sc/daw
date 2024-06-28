package render

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v4/render/mod"
)

// A Sprite is a basic wrapper around image data and a point. The most basic Renderable.
type Sprite struct {
	LayeredPoint
	r *image.RGBA
}

// NewEmptySprite returns a sprite of the given dimensions with a blank RGBA
func NewEmptySprite(x, y float64, w, h int) *Sprite {
	r := image.NewRGBA(image.Rect(0, 0, w, h))
	return NewSprite(x, y, r)
}

// NewSprite creates a new sprite
func NewSprite(x, y float64, r *image.RGBA) *Sprite {
	return &Sprite{
		LayeredPoint: NewLayeredPoint(x, y, 0),
		r:            r,
	}
}

// GetRGBA returns the rgba behind this sprite
func (s *Sprite) GetRGBA() *image.RGBA {
	return s.r
}

// GetDims returns the dimensions of this sprite, or if this sprite has no
// defined RGBA returns default values.
func (s *Sprite) GetDims() (int, int) {
	if s.r == nil {
		return 1, 1
	}
	bds := s.r.Bounds()
	return bds.Max.X, bds.Max.Y
}

// SetRGBA will replace the rgba behind this sprite
func (s *Sprite) SetRGBA(r *image.RGBA) {
	s.r = r
}

// Bounds is an alternative to GetDims that alows a sprite
// to satisfy draw.Image.
func (s *Sprite) Bounds() image.Rectangle {
	return s.r.Bounds()
}

// ColorModel allows sprites to satisfy draw.Image. Returns
// color.RGBAModel.
func (s *Sprite) ColorModel() color.Model {
	return s.r.ColorModel()
}

// At returns the color of a given pixel location
func (s *Sprite) At(x, y int) color.Color {
	return s.r.At(x, y)
}

// Set sets a color of a given pixel location
func (s *Sprite) Set(x, y int, c color.Color) {
	s.r.Set(x, y, c)
}

// Draw draws this sprite at +xOff, +yOff
func (s *Sprite) Draw(buff draw.Image, xOff, yOff float64) {
	DrawImage(buff, s.r, int(s.X()+xOff), int(s.Y()+yOff))
}

// Copy returns a copy of this Sprite
func (s *Sprite) Copy() Modifiable {
	newS := new(Sprite)
	if s.r != nil {
		newS.r = rgbaCopy(s.r)
	}
	newS.LayeredPoint = s.LayeredPoint.Copy()
	return newS
}

func rgbaCopy(r *image.RGBA) *image.RGBA {
	newRgba := new(image.RGBA)
	newRgba.Rect = r.Rect
	newRgba.Stride = r.Stride
	newRgba.Pix = make([]uint8, len(r.Pix))
	copy(newRgba.Pix, r.Pix)
	return newRgba
}

// Modify takes in modifications (modify.go) and alters this sprite accordingly
func (s *Sprite) Modify(ms ...mod.Mod) Modifiable {
	for _, m := range ms {
		s.r = m(s.GetRGBA())
	}
	return s
}

// Filter filters this sprite's rgba on all the input filters
func (s *Sprite) Filter(fs ...mod.Filter) {
	for _, f := range fs {
		f(s.r)
	}
}

// OverlaySprites combines sprites together through masking to form a single sprite
func OverlaySprites(sps []*Sprite) *Sprite {
	tmpSprite := sps[len(sps)-1].Copy().(*Sprite)
	for i := len(sps) - 1; i > 0; i-- {
		mod.FillMask(*sps[i-1].GetRGBA())(tmpSprite.GetRGBA())
	}
	return tmpSprite
}
