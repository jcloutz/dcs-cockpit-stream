package cockpit_stream

import (
	"image"
	"sync"
)

type ViewportReader interface {
	GetName() string
	GetBounds() image.Rectangle
	GetAdjOffset() image.Point
}

// Viewport describes a portion of an overall screen capture.
// It is used slice out the specified portion of the screen
// for reuse
type Viewport struct {
	// name of the viewport
	name string

	// True bounds of the viewport
	bounds image.Rectangle

	// slicePosition is the location that should be used
	// when slicing the viewport into or out of the current
	// screen capture. This value is calculated by the
	// ViewportManager and cached here in the viewport when
	// the capture bounds are recalculated
	slicePosition image.Point

	// position is the viewports location on the screen
	// as it is rendered on screen. It is used to
	// help define capture bounds for the screen capture
	// controller
	position image.Point

	mutex sync.RWMutex
}

func NewViewport(name string, x int, y int, width int, height int) *Viewport {
	return &Viewport{
		name:     name,
		bounds:   image.Rect(0, 0, width, height),
		position: image.Point{X: x, Y: y},
	}
}

func (vp *Viewport) GetName() string {
	return vp.name
}

func (vp *Viewport) GetBounds() image.Rectangle {
	return vp.bounds
}

func (vp *Viewport) GetAdjOffset() image.Point {
	return vp.slicePosition
}

func (vp *Viewport) SetAdjOffset(pt image.Point) {
	vp.slicePosition = pt
}

func (vp *Viewport) GetRealOffset() image.Point {
	return vp.position
}
