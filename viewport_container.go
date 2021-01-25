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

func (vm *ViewportContainer) Get(key string) (ViewportReader, error) {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	viewport, ok := vm.data[key]
	if !ok {
		return nil, errors.New("viewport not found")
	}

	return viewport, nil
}

func (vm *ViewportContainer) Add(key string, viewport *Viewport) {
	vm.mutex.Lock()
	vm.data[key] = viewport
	vm.mutex.Unlock()

	vm.recomputeBounds()
}

func (vm *ViewportContainer) OneMutate(name string, callback func(viewport *Viewport)) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	viewport, ok := vm.data[name]
	if !ok {
		return errors.New("viewport does not exist")
	}

	callback(viewport)

	return nil
}

func (vm *ViewportContainer) One(name string, callback func(viewport ViewportReader)) error {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()

	viewport, ok := vm.data[name]
	if !ok {
		return errors.New("viewport does not exist")
	}

	callback(viewport)

	return nil
}

func (vm *ViewportContainer) EachMutate(callback func(name string, viewport *Viewport)) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	for name, viewport := range vm.data {
		callback(name, viewport)
	}
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

func (vm *ViewportContainer) GetBounds() image.Rectangle {
	return vm.bounds
}

// recomputeBounds will adjust
func (vm *ViewportContainer) recomputeBounds() {
	minX := math.MaxInt16
	maxX := math.MinInt16

	minY := math.MaxInt16
	maxY := math.MinInt16

	vm.Each(func(name string, viewport *Viewport) {
		bounds := viewport.GetBounds()
		offset := viewport.GetRealOffset()

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

	vm.bounds = image.Rect(minX, minY, maxX, maxY)

	vm.EachMutate(func(name string, viewport *Viewport) {
		offset := viewport.GetRealOffset()

		viewport.slicePosition = offset.Sub(image.Point{
			X: minX,
			Y: minY,
		})
	})
}
