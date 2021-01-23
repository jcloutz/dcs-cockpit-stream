package cockpit_stream

import (
	"errors"
	"image"
	"image/draw"
	"sync"
)

type ViewportReader interface {
	GetName() string
	GetBounds() image.Rectangle
	GetAdjPoint() image.Point
	Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point)
}

type Viewport struct {
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

func NewServerViewport(name string, x int, y int, width int, height int) *Viewport {
	sliceBounds := image.Rect(x, y, x+width, y+height)
	imageBounds := image.Rect(0, 0, width, height)

	return &Viewport{
		Name:     name,
		Bounds:   sliceBounds,
		AdjPoint: image.Point{},
		image:    image.NewRGBA(imageBounds),
	}
}
func (vp *Viewport) Update(cap ImageSlicer) {
	vp.mutex.Lock()
	defer vp.mutex.Unlock()

	cap.Slice(vp.image, vp.image.Bounds(), vp.AdjPoint)
}

func (vp *Viewport) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	vp.mutex.RLock()
	defer vp.mutex.RUnlock()
	draw.Draw(dst, bounds, vp.image, at, draw.Src)
}

type ViewportMutexMap struct {
	data  map[string]*Viewport
	mutex sync.RWMutex
}

func NewViewportMutextMap() *ViewportMutexMap {
	return &ViewportMutexMap{
		data: make(map[string]*Viewport),
	}
}
func (vm *ViewportMutexMap) Has(key string) bool {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	_, ok := vm.data[key]

	return ok
}

func (vm *ViewportMutexMap) Get(key string) (ViewportReader, error) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	viewport, ok := vm.data[key]
	if !ok {
		return nil, errors.New("viewport not found")
	}

	return viewport, nil
}

func (vm *ViewportMutexMap) Set(key string, viewport *Viewport) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	vm.data[key] = viewport
}

func (vm *ViewportMutexMap) Each(callback func(name string, viewport *Viewport)) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	for name, viewport := range vm.data {
		callback(name, viewport)
	}
}

func (vp *Viewport) GetName() string {
	return vp.Name
}

func (vp *Viewport) GetBounds() image.Rectangle {
	return vp.Bounds
}

func (vp *Viewport) GetAdjPoint() image.Point {
	return vp.AdjPoint
}
