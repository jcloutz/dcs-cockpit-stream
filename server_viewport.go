package cockpit_stream

import (
	"image"
	"image/draw"
	"sync"
)

type ServerViewport struct {
	Name string
	// True bounds of the viewpoint
	Bounds image.Rectangle

	// Adjusted bounds within the screen manager
	// capture rectangle this is set by the screen
	// manager
	AdjPoint image.Point

	// Image of the viewport
	image *image.RGBA

	mutex sync.RWMutex
}

func NewServerViewport(name string, x int, y int, width int, height int) *ServerViewport {
	sliceBounds := image.Rect(x, y, x+width, y+height)
	imageBounds := image.Rect(0, 0, width, height)

	return &ServerViewport{
		Name:     name,
		Bounds:   sliceBounds,
		AdjPoint: image.Point{},
		image:    image.NewRGBA(imageBounds),
	}
}
func (vp *ServerViewport) Update(cap ImageSlicer) {
	vp.mutex.Lock()
	defer vp.mutex.Unlock()

	cap.Slice(vp.image, vp.image.Bounds(), vp.AdjPoint)
}

func (vp *ServerViewport) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()
	draw.Draw(dst, bounds, vp.image, at, draw.Src)
}
