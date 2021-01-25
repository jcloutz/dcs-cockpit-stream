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
	Name string

	// True Bounds of the viewport
	Bounds image.Rectangle

	// SlicePosition is the location that should be used
	// when slicing the viewport into or out of the current
	// screen capture. This value is calculated by the
	// ViewportManager and cached here in the viewport when
	// the capture Bounds are recalculated
	SlicePosition image.Point

	// Position is the viewports location as it is rendered
	// on screen. It is used to help define capture Bounds
	// for the screen capture controller
	Position image.Point

	mutex sync.RWMutex
}

func NewViewport(name string, x int, y int, width int, height int) *Viewport {
	return &Viewport{
		Name:     name,
		Bounds:   image.Rect(0, 0, width, height),
		Position: image.Point{X: x, Y: y},
	}
}

// Lock the Viewport for mutation
func (vp *Viewport) Lock() {
	vp.mutex.Lock()
}

// Unlock the Viewport after mutation
func (vp *Viewport) Unlock() {
	vp.mutex.Unlock()
}

// RLock the Viewport for reading
func (vp *Viewport) RLock() {
	vp.mutex.RLock()
}

// RUnlock the Viewport after reading
func (vp *Viewport) RUnlock() {
	vp.mutex.RUnlock()
}
