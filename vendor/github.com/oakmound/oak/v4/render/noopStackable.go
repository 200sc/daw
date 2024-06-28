package render

import (
	"image/draw"

	"github.com/oakmound/oak/v4/alg/intgeom"
)

// NoopStackable is a Stackable element where all methods are no-ops.
// Use for tests to disable rendering.
type NoopStackable struct{}

// PreDraw on a NoopStackable does nothing.
func (ns NoopStackable) PreDraw() {}

// Add on a NoopStackable does nothing. The input Renderable is still returned.
func (ns NoopStackable) Add(r Renderable, _ ...int) Renderable {
	return r
}

// Replace on a NoopStackable does nothing.
func (ns NoopStackable) Replace(Renderable, Renderable, int) {}

// Copy on a NoopStackable returns itself.
func (ns NoopStackable) Copy() Stackable {
	return ns
}

// DrawToScreen on a NoopStackable does nothing.
func (ns NoopStackable) DrawToScreen(draw.Image, *intgeom.Point2, int, int) {}

// Clear on a NoopStackable does nothing.
func (ns NoopStackable) Clear() {}
