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

type Viewport struct {
	name string

	// True bounds of the viewpoint
	bounds image.Rectangle

	// Adjusted bounds within the screen manager
	// capture rectangle this is set by the screen
	// manager
	adjOffset  image.Point
	realOffset image.Point

	mutex sync.RWMutex
}

func NewViewport(name string, x int, y int, width int, height int) *Viewport {
	return &Viewport{
		name:       name,
		bounds:     image.Rect(0, 0, width, height),
		realOffset: image.Point{X: x, Y: y},
	}
}

func (vp *Viewport) GetName() string {
	return vp.name
}

func (vp *Viewport) GetBounds() image.Rectangle {
	return vp.bounds
}

func (vp *Viewport) GetAdjOffset() image.Point {
	return vp.adjOffset
}

func (vp *Viewport) SetAdjOffset(pt image.Point) {
	vp.adjOffset = pt
}

func (vp *Viewport) GetRealOffset() image.Point {
	return vp.realOffset
}
