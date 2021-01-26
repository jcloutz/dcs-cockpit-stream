package cockpit_stream

import (
	"image"
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

	// x position of the viewport
	x int

	// y position of the viewport
	y int

	// width of the viewport
	width int

	// height of the viewport
	height int

	// one at a time please and thank you
	mutex sync.RWMutex
}

func NewViewport(name string, x int, y int, width int, height int) *Viewport {
	return &Viewport{
		name:   name,
		x:      x,
		y:      y,
		height: height,
		width:  width,
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

	return vp.width
}

func (vp *Viewport) Height() int {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return vp.height
}

// Point is the viewports location as it is rendered
// on screen. It is used to help define capture SizeRect
// for the screen capture controller
func (vp *Viewport) Point() image.Point {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return image.Point{
		X: vp.x,
		Y: vp.y,
	}
}

// SizeRect returns an image.Rectangle instance defining the
// dimensions of the viewport image
func (vp *Viewport) SizeRect() image.Rectangle {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return image.Rect(0, 0, vp.width, vp.height)
}

// PositionRect returns an image.Rectangle defining the viewport
// bounds with the original position offset applied to it
func (vp *Viewport) PositionRect() image.Rectangle {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	return image.Rect(vp.x, vp.y, vp.x+vp.width, vp.y+vp.height)
}
