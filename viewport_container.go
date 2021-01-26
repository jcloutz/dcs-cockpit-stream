package cockpit_stream

import (
	"errors"
	"image"
	"math"
	"sync"
)

type ViewportContainer struct {
	// holds a map of all serverViewports registered with the application
	data map[string]*Viewport

	// handle concurrent access
	mutex sync.RWMutex

	// rectangle containing all serverViewports registered with this
	// instance
	bounds image.Rectangle

	boundsOffset image.Point
}

func NewViewportContainer() *ViewportContainer {
	return &ViewportContainer{
		data: make(map[string]*Viewport),
	}
}

func (vm *ViewportContainer) Has(key string) bool {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	_, ok := vm.data[key]

	return ok
}
func (vm *ViewportContainer) Add(id string, x, y, width, height int) {
	vm.mutex.Lock()
	vm.data[id] = NewViewport(id, x, y, width, height)
	vm.mutex.Unlock()

	vm.recomputeBounds()
}

func (vm *ViewportContainer) AddViewport(viewport *Viewport) {
	vm.mutex.Lock()
	vm.data[viewport.Name] = viewport
	vm.mutex.Unlock()

	vm.recomputeBounds()
}

func (vm *ViewportContainer) Get(key string) (*Viewport, error) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	viewport, ok := vm.data[key]
	if !ok {
		return nil, errors.New("viewport not found")
	}

	return viewport, nil
}

func (vm *ViewportContainer) Each(callback func(name string, viewport *Viewport)) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	for name, viewport := range vm.data {
		callback(name, viewport)
	}
}

func (vm *ViewportContainer) Count() int {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()
	return len(vm.data)
}

func (vm *ViewportContainer) Bounds() image.Rectangle {
	return vm.bounds
}

func (vm *ViewportContainer) Offset() image.Point {
	return vm.boundsOffset
}

func (vm *ViewportContainer) ViewportOffset(viewport *Viewport) (image.Point, error) {
	viewport.RLock()
	defer viewport.RUnlock()
	return viewport.Position.Sub(vm.boundsOffset), nil
}

// recomputeBounds will adjust
func (vm *ViewportContainer) recomputeBounds() {
	minX := math.MaxInt16
	maxX := math.MinInt16

	minY := math.MaxInt16
	maxY := math.MinInt16

	vm.Each(func(name string, viewport *Viewport) {
		viewport.RLock()
		defer viewport.RUnlock()

		bounds := viewport.Bounds
		offset := viewport.Position

		if offset.X < minX {
			minX = offset.X
		}
		if offset.Y < minY {
			minY = offset.Y
		}

		maxOffsetX := offset.X + bounds.Dx()
		if maxOffsetX > maxX {
			maxX = maxOffsetX
		}

		maxOffsetY := offset.Y + bounds.Dy()
		if maxOffsetY > maxY {
			maxY = maxOffsetY
		}
	})

	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	vm.bounds = image.Rect(minX, minY, maxX, maxY)
	vm.boundsOffset = image.Point{X: minX, Y: minY}
}
