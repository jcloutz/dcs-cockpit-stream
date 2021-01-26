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
	Name string

	// True Bounds of the viewport
	Bounds image.Rectangle

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

func (vp *Viewport) Slice(dst *image.RGBA, src *image.RGBA, offset image.Point) {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()

	draw.Draw(dst, vp.Bounds, src, offset, draw.Src)
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
