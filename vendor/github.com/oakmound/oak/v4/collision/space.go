package collision

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/physics"
)

// ID Types constant
const (
	NONE = iota
	IDTypeCID
	IDTypePID
)

// A Space is a rectangle
// with a couple of ways of identifying
// an underlying object.
type Space struct {
	Location floatgeom.Rect3
	// A label can store type information.
	// Recommended to use with an enum.
	Label Label
	// A CID can be used to get the exact
	// entity which this rectangle belongs to.
	CID event.CallerID
	// Type represents which ID space the above ID
	// corresponds to.
	Type int
}

// Bounds satisfies the rtreego.Spatial interface.
func (s *Space) Bounds() floatgeom.Rect3 {
	return s.Location
}

// X returns a space's x position (leftmost)
func (s *Space) X() float64 {
	return s.Location.Min.X()
}

// Y returns a space's y position (upmost)
func (s *Space) Y() float64 {
	return s.Location.Min.Y()
}

// GetW returns a space's width (rightmost x - leftmost x)
// Deprecated: Use W instead
func (s *Space) GetW() float64 {
	return s.Location.W()
}

// GetH returns a space's height (upper y - lower y)
// Deprecated: Use H instead
func (s *Space) GetH() float64 {
	return s.Location.H()
}

// W returns a space's width (rightmost x - leftmost x)
func (s *Space) W() float64 {
	return s.Location.W()
}

// H returns a space's height (upper y - lower y)
func (s *Space) H() float64 {
	return s.Location.H()
}

// GetCenter returns the center point of the space
func (s *Space) GetCenter() (float64, float64) {
	return s.X() + s.GetW()/2, s.Y() + s.GetH()/2
}

// GetPos returns both y and x
func (s *Space) GetPos() (float64, float64) {
	return s.X(), s.Y()
}

// Above returns how much above this space another space is
// Important note: (10,10) is Above (10,20), because in oak's
// display, lower y values are higher than higher y values.
func (s *Space) Above(other *Space) float64 {
	return other.Y() - s.Y()
}

// Below returns how much below this space another space is,
// Equivalent to -1 * Above
func (s *Space) Below(other *Space) float64 {
	return s.Y() - other.Y()
}

// Contains returns whether this space contains another
func (s *Space) Contains(other *Space) bool {
	//You contain another space if it is fully inside your space
	//If you are the same size and location as the space you are checking then you both contain eachother
	if s.X() > other.X() || s.X()+s.GetW() < other.X()+other.GetW() ||
		s.Y() > other.Y() || s.Y()+s.GetH() < other.Y()+other.GetH() {
		return false
	}
	return true
}

// LeftOf returns how far to the left other is of this space
func (s *Space) LeftOf(other *Space) float64 {
	return other.X() - s.X()
}

// RightOf returns how far to the right other is of this space.
// Equivalent to -1 * LeftOf
func (s *Space) RightOf(other *Space) float64 {
	return s.X() - other.X()
}

// Overlap returns how much this space overlaps with another space
func (s *Space) Overlap(other *Space) (xOver, yOver float64) {
	if s.X() > other.X() {
		x2 := other.X() + other.GetW()
		if s.X() < x2 {
			xOver = s.X() - x2
		}
	} else {
		x2 := s.X() + s.GetW()
		if other.X() < x2 {
			xOver = x2 - other.X()
		}
	}
	if s.Y() > other.Y() {
		y2 := other.Y() + other.GetH()
		if s.Y() < y2 {
			yOver = s.Y() - y2
		}
	} else {
		y2 := s.Y() + s.GetH()
		if other.Y() < y2 {
			yOver = y2 - other.Y()
		}
	}
	return
}

// OverlapVector returns Overlap as a vector
func (s *Space) OverlapVector(other *Space) physics.Vector {
	xover, yover := s.Overlap(other)
	return physics.NewVector(xover, yover)
}

// SubtractRect removes a subrectangle from this rectangle and
// returns the rectangles remaining after the portion has been
// removed. The input x,y is relative to the original space:
// Example: removing 1,1 from 10,10 -> 12,12 is OK, but removing
// 11,11 from 10,10 -> 12,12 will not act as expected.
func (s *Space) SubtractRect(x2, y2, w2, h2 float64) []*Space {
	x1 := s.X()
	y1 := s.Y()
	w1 := s.GetW()
	h1 := s.GetH()

	// Left, Top, Right, Bottom
	// X, Y, W, H
	rects := [4][4]float64{}

	rects[0][0] = x1
	rects[0][1] = y1
	rects[0][2] = x2
	rects[0][3] = h1

	// Todo: these spaces overlap on the corners. We could remove that.
	rects[1][0] = x1
	rects[1][1] = y1
	rects[1][2] = w1
	rects[1][3] = y2

	rects[2][0] = x1 + x2 + w2
	rects[2][1] = y1
	rects[2][2] = w1 - (x2 + w2)
	rects[2][3] = h1

	rects[3][0] = x1
	rects[3][1] = y1 + y2 + h2
	rects[3][2] = w1
	rects[3][3] = h1 - (y2 + h2)

	var spaces []*Space

	for _, r := range rects {
		if r[2] > 0 && r[3] > 0 {
			spaces = append(spaces, NewFullSpace(r[0], r[1], r[2], r[3], s.Label, s.CID))
		}
	}

	return spaces
}

// NewUnassignedSpace returns a space that just has a rectangle
func NewUnassignedSpace(x, y, w, h float64) *Space {
	return NewLabeledSpace(x, y, w, h, NilLabel)
}

// NewSpace returns a space with an associated caller id
func NewSpace(x, y, w, h float64, cID event.CallerID) *Space {
	return NewFullSpace(x, y, w, h, NilLabel, cID)
}

// NewLabeledSpace returns a space with an associated integer label
func NewLabeledSpace(x, y, w, h float64, l Label) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		Location: rect,
		Label:    l,
		Type:     NONE,
	}
}

// NewFullSpace returns a space with both a label and a caller id
func NewFullSpace(x, y, w, h float64, l Label, cID event.CallerID) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		rect,
		l,
		cID,
		IDTypeCID,
	}
}

// NewRect2Space returns a space with an associated caller id from a rect2
func NewRect2Space(rect floatgeom.Rect2, cID event.CallerID) *Space {
	return NewSpace(rect.Min.X(), rect.Min.Y(), rect.W(), rect.H(), cID)
}

// NewRectSpace creates a colliison space with the specified 3D rectangle
func NewRectSpace(rect floatgeom.Rect3, l Label, cID event.CallerID) *Space {
	return &Space{
		rect,
		l,
		cID,
		IDTypeCID,
	}
}

// NewRect is a wrapper around rtreego.NewRect,
// casting the given x,y to an rtreego.Point.
// Used to not expose rtreego.Point to the user.
// Invalid widths and heights are converted to be valid.
// If zero width or height is given, it is replaced with 1.
// If a negative width or height is given, the rectangle is
// shifted to the left or up by that negative dimension and
// the dimension is made positive.
func NewRect(x, y, w, h float64) floatgeom.Rect3 {
	if w == 0 {
		w = 1
	}
	if h == 0 {
		h = 1
	}
	return floatgeom.NewRect3WH(x, y, 0, w, h, 1)
}

// SetZLayer sets a space's z layer.
func (s *Space) SetZLayer(z float64) {
	s.Location.Min[2] = z
	s.Location.Max[2] = z
}
