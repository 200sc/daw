package floatgeom

import (
	"math/rand"

	"github.com/oakmound/oak/v4/alg/span"
)

// A Rect2 represents a span from one point in 2D space to another.
// If Min is less than max on any axis, it will return undefined results
// for methods.
type Rect2 struct {
	Min, Max Point2
}

// MaxDimensions reports that a Rect2 has only two dimensions of definition.
func (r Rect2) MaxDimensions() int {
	return 2
}

// A Rect3 represents a span from one point in 3D space to another.
// If Min is less than Max on any axis, it will return undefined results
// for methods.
type Rect3 struct {
	Min, Max Point3
}

// MaxDimensions reports that a Rect3 has only three dimensions of definition.
func (r Rect3) MaxDimensions() int {
	return 3
}

// NewRect2 returns an (X,Y):(X2,Y2) rectangle. This enforces that
// x < x2 and y < y2, and will swap the inputs if that is not true.
// If that enforcement is not desired, construct the struct manually.
func NewRect2(x, y, x2, y2 float64) Rect2 {
	if x > x2 {
		x, x2 = x2, x
	}
	if y > y2 {
		y, y2 = y2, y
	}
	return Rect2{
		Min: Point2{x, y},
		Max: Point2{x2, y2},
	}
}

// NewRect2WH returns an (X,Y):(X+W,Y+H) rectangle. This enforces that
// w and h are positive, and will decrease x and y respectively if that is not true.
func NewRect2WH(x, y, w, h float64) Rect2 {
	if w < 0 {
		x += w
		w *= -1
	}
	if h < 0 {
		y += h
		h *= -1
	}
	return Rect2{
		Min: Point2{x, y},
		Max: Point2{x + w, y + h},
	}
}

// NewBoundingRect2 will produce the minimal rectangle that contains all of
// the input points.
func NewBoundingRect2(pts ...Point2) Rect2 {
	return Rect2{
		Min: pts[0].LesserOf(pts...),
		Max: pts[0].GreaterOf(pts...),
	}
}

// NewRect3 returns an (X,Y,Z):(X2,Y2,Z2) rectangle. This enforces that
// x < x2, y < y2, and z < z2, and will swap the inputs if that is not true.
func NewRect3(x, y, z, x2, y2, z2 float64) Rect3 {
	if x > x2 {
		x, x2 = x2, x
	}
	if y > y2 {
		y, y2 = y2, y
	}
	if z > z2 {
		z, z2 = z2, z
	}
	return Rect3{
		Min: Point3{x, y, z},
		Max: Point3{x2, y2, z2},
	}
}

// NewRect3WH returns an (X,Y,Z):(X+W,Y+H,Z+D) rectangle. This enforces that
// w, h, and d and positive, and will decrease x, y, and z respectively if that
// is not true.
func NewRect3WH(x, y, z, w, h, d float64) Rect3 {
	if w < 0 {
		x += w
		w *= -1
	}
	if h < 0 {
		y += h
		h *= -1
	}
	if d < 0 {
		z += d
		d *= -1
	}
	return Rect3{
		Min: Point3{x, y, z},
		Max: Point3{x + w, y + h, z + d},
	}
}

// NewBoundingRect3 will produce the minimal rectangle that contains all of
// the input points.
func NewBoundingRect3(pts ...Point3) Rect3 {
	return Rect3{
		Min: pts[0].LesserOf(pts...),
		Max: pts[0].GreaterOf(pts...),
	}
}

// Shift moves the rectangle by a point returns a new instance
func (r Rect2) Shift(p Point2) Rect2 {
	r.Min = r.Min.Add(p)
	r.Max = r.Max.Add(p)
	return r
}

// Shift moves the rectangle by a point returns a new instance
func (r Rect3) Shift(p Point3) Rect3 {
	r.Min = r.Min.Add(p)
	r.Max = r.Max.Add(p)
	return r
}

// Area returns W * H.
func (r Rect2) Area() float64 {
	return r.W() * r.H()
}

// Span returns the span on this rectangle's ith axis.
func (r Rect2) Span(i int) float64 {
	return r.Max[i] - r.Min[i]
}

// W returns the width of this rectangle.
func (r Rect2) W() float64 {
	return r.Span(0)
}

// H returns the height of this rectangle.
func (r Rect2) H() float64 {
	return r.Span(1)
}

// Space returns W * H * D
func (r Rect3) Space() float64 {
	return r.W() * r.H() * r.D()
}

// Span returns the span on this rectangle's ith axis.
func (r Rect3) Span(i int) float64 {
	return r.Max[i] - r.Min[i]
}

// W returns the width of this rectangle.
func (r Rect3) W() float64 {
	return r.Span(0)
}

// H returns the height of this rectangle.
func (r Rect3) H() float64 {
	return r.Span(1)
}

// D returns the depth of this rectangle.
func (r Rect3) D() float64 {
	return r.Span(2)
}

// Midpoint returns the midpoint of this rectangle's span over a given dimension.
func (r Rect2) Midpoint(i int) float64 {
	return (r.Min[i] + r.Max[i]) / 2
}

// Midpoint returns the midpoint of this rectangle's span over a given dimension.
func (r Rect3) Midpoint(i int) float64 {
	return (r.Min[i] + r.Max[i]) / 2
}

// Center returns the center of this rectangle
func (r Rect2) Center() Point2 {
	return Point2{
		r.Midpoint(0),
		r.Midpoint(1),
	}
}

// Center returns the center of this rectangle
func (r Rect3) Center() Point3 {
	return Point3{
		r.Midpoint(0),
		r.Midpoint(1),
		r.Midpoint(2),
	}
}

// Perimeter computes the sum of the edge lengths of a rectangle.
func (r Rect2) Perimeter() float64 {
	// The number of edges in an n-dimensional rectangle is n * 2^(n-1)
	// (http://en.wikipedia.org/wiki/Hypercube_graph).  Thus the number
	// of edges of length (ai - bi), where the rectangle is determined
	// by p = (a1, a2, ..., an) and q = (b1, b2, ..., bn), is 2^(n-1).
	//
	// The margin of the rectangle, then, is given by the formula
	// 2^(n-1) * [(b1 - a1) + (b2 - a2) + ... + (bn - an)].
	return 2 * (r.W() + r.H())
}

// Margin computes the sum of the edge lengths of a rectangle.
func (r Rect3) Margin() float64 {
	return 4 * (r.W() + r.H() + r.D())
}

// Contains tests whether p is located inside or on the boundary of r.
func (r Rect2) Contains(p Point2) bool {
	return (p.X() >= r.Min.X() &&
		p.X() <= r.Max.X() &&
		p.Y() >= r.Min.Y() &&
		p.Y() <= r.Max.Y())
}

// Contains tests whether p is located inside or on the boundary of r.
func (r Rect3) Contains(p Point3) bool {
	return (p.X() >= r.Min.X() &&
		p.X() <= r.Max.X() &&
		p.Y() >= r.Min.Y() &&
		p.Y() <= r.Max.Y() &&
		p.Z() >= r.Min.Z() &&
		p.Z() <= r.Max.Z())
}

// ContainsRect tests whether r2 is is located inside r1.
func (r Rect2) ContainsRect(r2 Rect2) bool {
	return (r2.Min.X() >= r.Min.X() &&
		r2.Max.X() <= r.Max.X() &&
		r2.Min.Y() >= r.Min.Y() &&
		r2.Max.Y() <= r.Max.Y())
}

// ContainsRect tests whether r2 is is located inside r1.
func (r Rect3) ContainsRect(r2 Rect3) bool {
	return (r2.Min.X() >= r.Min.X() &&
		r2.Max.X() <= r.Max.X() &&
		r2.Min.Y() >= r.Min.Y() &&
		r2.Max.Y() <= r.Max.Y() &&
		r2.Min.Z() >= r.Min.Z() &&
		r2.Max.Z() <= r.Max.Z())
}

// GreaterOf returns a rectangle formed of the lowest values on each
// dimension for Min, and the highest for Max.
func (r Rect2) GreaterOf(r2 Rect2) Rect2 {
	r.Min = r.Min.LesserOf(r2.Min)
	r.Max = r.Max.GreaterOf(r2.Max)
	return r
}

// GreaterOf returns a rectangle formed of the lowest values on each
// dimension for Min, and the highest for Max.
func (r Rect3) GreaterOf(r2 Rect3) Rect3 {
	r.Min = r.Min.LesserOf(r2.Min)
	r.Max = r.Max.GreaterOf(r2.Max)
	return r
}

// Intersects returns whether the two rectangles intersect.
func (r Rect3) Intersects(r2 Rect3) bool {
	// There are four cases of overlap:
	//
	//     1.  a1------------b1
	//              a2------------b2
	//              p--------q
	//
	//     2.       a1------------b1
	//         a2------------b2
	//              p--------q
	//
	//     3.  a1-----------------b1
	//              a2-------b2
	//              p--------q
	//
	//     4.       a1-------b1
	//         a2-----------------b2
	//              p--------q
	//
	// Thus there are only two cases of non-overlap:
	//
	//     1. a1------b1
	//                    a2------b2
	//
	//     2.             a1------b1
	//        a2------b2
	//
	// Enforced by constructor: a1 <= b1 and a2 <= b2.  So we can just
	// check the endpoints.

	return !((r2.Max.X() <= r.Min.X() || r.Max.X() <= r2.Min.X()) ||
		(r2.Max.Y() <= r.Min.Y() || r.Max.Y() <= r2.Min.Y()) ||
		(r2.Max.Z() <= r.Min.Z() || r.Max.Z() <= r2.Min.Z()))
}

// Intersects returns whether the two rectangles intersect.
func (r Rect2) Intersects(r2 Rect2) bool {
	// There are four cases of overlap:
	//
	//     1.  a1------------b1
	//              a2------------b2
	//              p--------q
	//
	//     2.       a1------------b1
	//         a2------------b2
	//              p--------q
	//
	//     3.  a1-----------------b1
	//              a2-------b2
	//              p--------q
	//
	//     4.       a1-------b1
	//         a2-----------------b2
	//              p--------q
	//
	// Thus there are only two cases of non-overlap:
	//
	//     1. a1------b1
	//                    a2------b2
	//
	//     2.             a1------b1
	//        a2------b2
	//
	// Enforced by constructor: a1 <= b1 and a2 <= b2.  So we can just
	// check the endpoints.

	return !((r2.Max.X() <= r.Min.X() || r.Max.X() <= r2.Min.X()) ||
		(r2.Max.Y() <= r.Min.Y() || r.Max.Y() <= r2.Min.Y()))
}

// ProjectZ projects the Rect3 onto the z axis, removing it's
// z component and returning a Rect2
func (r Rect3) ProjectZ() Rect2 {
	return Rect2{
		Min: r.Min.ProjectZ(),
		Max: r.Max.ProjectZ(),
	}
}

// MulConst multiplies the boundary points of this rectangle by i.
func (r Rect2) MulConst(i float64) Rect2 {
	return Rect2{
		r.Min.MulConst(i),
		r.Max.MulConst(i),
	}
}

// Poll returns a pseudorandom point from within this rectangle
func (r Rect2) Poll() Point2 {
	return Point2{
		r.Min.X() + rand.Float64()*float64(r.W()),
		r.Min.Y() + rand.Float64()*float64(r.H()),
	}
}

// Clamp returns a version of the provided point such that it is contained within r. If it was already contained in
// r, it will not be changed.
func (r Rect2) Clamp(pt Point2) Point2 {
	for i := 0; i < r.MaxDimensions(); i++ {
		if pt[i] < r.Min[i] {
			pt[i] = r.Min[i]
		} else if pt[i] > r.Max[i] {
			pt[i] = r.Max[i]
		}
	}
	return pt
}

// Percentile returns a point within this rectangle along the vector from the top left to the bottom right of the
// rectangle, where for example, 0.0 will be r.Min, 1.0 will be r.Max, and 2.0 will be project the vector beyond r
// and return r.Min + {r.W()*2, r.H()*2}
func (r Rect2) Percentile(f float64) Point2 {
	return Point2{
		r.Min.X() + f*float64(r.W()),
		r.Min.Y() + f*float64(r.H()),
	}
}

// MulSpan returns this rectangle as a Point2 Span after multiplying the boundary points of the rectangle by f.
func (r Rect2) MulSpan(f float64) span.Span[Point2] {
	return r.MulConst(f)
}

// MulConst multiplies the boundary points of this rectangle by i.
func (r Rect3) MulConst(i float64) Rect3 {
	return Rect3{
		r.Min.MulConst(i),
		r.Max.MulConst(i),
	}
}

// Poll returns a pseudorandom point from within this rectangle
func (r Rect3) Poll() Point3 {
	return Point3{
		r.Min.X() + (rand.Float64() * float64(r.W())),
		r.Min.Y() + (rand.Float64() * float64(r.H())),
		r.Min.Z() + (rand.Float64() * float64(r.D())),
	}
}

// Clamp returns a version of the provided point such that it is contained within r. If it was already contained in
// r, it will not be changed.
func (r Rect3) Clamp(pt Point3) Point3 {
	for i := 0; i < r.MaxDimensions(); i++ {
		if pt[i] < r.Min[i] {
			pt[i] = r.Min[i]
		} else if pt[i] > r.Max[i] {
			pt[i] = r.Max[i]
		}
	}
	return pt
}

// Percentile returns a point within this rectangle along the vector from the top left to the bottom right of the
// rectangle, where for example, 0.0 will be r.Min, 1.0 will be r.Max, and 2.0 will be project the vector beyond r
// and return r.Min + {r.W()*2, r.H()*2, r.D()*2}
func (r Rect3) Percentile(f float64) Point3 {
	return Point3{
		r.Min.X() + (f * float64(r.W())),
		r.Min.Y() + (f * float64(r.H())),
		r.Min.Z() + (f * float64(r.D())),
	}
}

// MulConst multiplies the boundary points of this rectangle by i.
func (r Rect3) MulSpan(f float64) span.Span[Point3] {
	return r.MulConst(f)
}
