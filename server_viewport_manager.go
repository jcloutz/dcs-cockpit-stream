package cockpit_stream

import (
	"errors"
	"fmt"
	"image"
	"math"
	"sync"
	"time"
)

type ImageSlicerContainer interface {
	HasSlicer(name string) bool
	Slice(name string, dst *image.RGBA, bounds image.Rectangle, at image.Point) error
}

var _ ImageSlicerContainer = &ServerViewportManager{}

type ServerViewportManager struct {
	viewports      map[string]*ServerViewport
	viewportsMutex sync.RWMutex

	capture        ImageCapturer
	targetFps      int
	computedBounds image.Rectangle
	done           chan bool
	timer          *time.Ticker

	fps int
}

func NewServerViewportManager(screenCapper ImageCapturer, targetFps int) *ServerViewportManager {
	return &ServerViewportManager{
		viewports: make(map[string]*ServerViewport),
		capture:   screenCapper,
		targetFps: targetFps,
	}
}

func (vm *ServerViewportManager) AddNewViewport(name string, x int, y int, width int, height int) *ServerViewportManager {
	return vm.AddViewport(NewServerViewport(name, x, y, width, height))
}

func (vm *ServerViewportManager) AddViewport(viewport *ServerViewport) *ServerViewportManager {
	vm.viewportsMutex.Lock()
	defer vm.viewportsMutex.Unlock()
	vm.viewports[viewport.Name] = viewport

	vm.recomputeBounds()

	return vm
}

// recomputeBounds will adjust
func (vm *ServerViewportManager) recomputeBounds() {
	minX := math.MaxInt16
	maxX := math.MinInt16

	minY := math.MaxInt16
	maxY := math.MinInt16
	for _, viewport := range vm.viewports {
		if viewport.Bounds.Min.X < minX {
			minX = viewport.Bounds.Min.X
		}
		if viewport.Bounds.Min.Y < minY {
			minY = viewport.Bounds.Min.Y
		}

		if viewport.Bounds.Max.X > maxX {
			maxX = viewport.Bounds.Max.X
		}

		if viewport.Bounds.Max.Y > maxY {
			maxY = viewport.Bounds.Max.Y
		}
	}
	vm.computedBounds = image.Rect(minX, minY, maxX, maxY)

	for _, viewport := range vm.viewports {
		viewport.AdjPoint = viewport.Bounds.Sub(image.Point{
			X: minX,
			Y: minY,
		}).Min
	}
}

func (vm *ServerViewportManager) UpdateViewports(cap ImageSlicer) {
	vm.viewportsMutex.RLock()
	defer vm.viewportsMutex.RUnlock()

	for _, vp := range vm.viewports {
		vp.Update(cap)
	}
}

func (vm *ServerViewportManager) benchmark() (int, int) {
	iterations := 120
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = vm.capture.Capture()
	}
	end := time.Now()

	avgTime := end.Sub(start).Milliseconds() / int64(iterations)
	maxFps := 1000 / avgTime

	return int(avgTime), int(maxFps)
}

func (vm *ServerViewportManager) HasSlicer(name string) bool {
	vm.viewportsMutex.RLock()
	defer vm.viewportsMutex.RUnlock()
	_, ok := vm.viewports[name]
	return ok
}

func (vm *ServerViewportManager) Slice(name string, dst *image.RGBA, bounds image.Rectangle, at image.Point) error {
	if !vm.HasSlicer(name) {
		return errors.New(fmt.Sprintf("no viewport with name '%s' found", name))
	}
	vm.viewportsMutex.RLock()
	defer vm.viewportsMutex.RUnlock()

	vm.viewports[name].Slice(dst, bounds, at)

	return nil
}
