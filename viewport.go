package cockpit_stream

import (
	"image"
	"image/draw"
	"sync"
)

// Viewport describes a portion of an overall screen capture.
// It is used slice out the specified portion of the screen
// for reuse.
//
// Any usage of Viewport should be handled with the proper
// mutex locks
type Viewport struct {
	// Name of the viewport
	name string

	// True Bounds of the viewport
	bounds image.Rectangle

	// Position is the viewports location as it is rendered
	// on screen. It is used to help define capture Bounds
	// for the screen capture controller
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

func (vp *Viewport) Name() string {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return vp.name
}

func (vp *Viewport) Width() int {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return vp.bounds.Dx()
}

func (vp *Viewport) Height() int {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return vp.bounds.Dy()
}

func (vp *Viewport) Position() image.Point {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return vp.position
}

func (vp *Viewport) Bounds() image.Rectangle {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return vp.bounds
}

func (vp *Viewport) Slice(dst *image.RGBA, src *image.RGBA, offset image.Point) {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	draw.Draw(dst, vp.bounds, src, offset, draw.Src)
}
