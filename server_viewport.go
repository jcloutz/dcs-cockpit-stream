package cockpit_stream

import (
	"fmt"
	"image"
	"math"
	"time"
)

type ServerViewport struct {
	Name string
	// True bounds of the viewpoint
	Bounds image.Rectangle

	// Adjusted bounds within the screen manager
	// capture rectangle this is set by the screen
	// manager
	AdjBounds image.Rectangle

	// Image of the viewport
	image *image.RGBA
}

func NewServerViewport(name string, x int, y int, width int, height int) *ServerViewport {
	bounds := image.Rect(x, y, x+width, y+width)

	return &ServerViewport{
		Name:      name,
		Bounds:    bounds,
		AdjBounds: image.Rectangle{},
		image:     image.NewRGBA(bounds),
	}
}

type ServerViewportManager struct {
	viewports      map[string]*ServerViewport
	capture        ScreenCapper
	computedBounds image.Rectangle
	done           chan bool
	timer          *time.Ticker
}

func NewServerViewportManager(screenCapper ScreenCapper) *ServerViewportManager {
	return &ServerViewportManager{
		viewports: make(map[string]*ServerViewport),
		capture:   screenCapper,
	}
}

func (vm *ServerViewportManager) AddNewViewport(name string, x int, y int, width int, height int) *ServerViewportManager {
	vm.viewports[name] = NewServerViewport(name, x, y, width, height)

	vm.recomputeBounds()

	return vm
}

func (vm *ServerViewportManager) AddViewport(viewport *ServerViewport) *ServerViewportManager {
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
		viewport.AdjBounds = viewport.Bounds.Sub(image.Point{
			X: minX,
			Y: minY,
		})
	}
}

func (vm *ServerViewportManager) updateViewports() {
	for name, vp := range vm.viewports {
		vm.capture.Slice(vp.image, vp.image.Bounds(), vp.AdjBounds.Min)
		save(vp.image, fmt.Sprintf("output/%s.png", name))
	}
}
func (vm *ServerViewportManager) benchmark() int64 {
	iterations := 60
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = vm.capture.Update(vm.computedBounds)
	}
	end := time.Now()

	f := end.Sub(start).Milliseconds()
	f += 1
	avg := end.Sub(start).Milliseconds() / int64(iterations)

	return avg

}

func (vm *ServerViewportManager) run() {
	for {
		select {
		case <-vm.timer.C:
			_ = vm.capture.Update(vm.computedBounds)
			go vm.updateViewports()
		case <-vm.done:
			close(vm.done)
			vm.timer.Stop()
			return
		}
	}

}

func (vm *ServerViewportManager) Run() {
	fmt.Println("Running benchmark")
	result := vm.benchmark()
	fmt.Printf("\tAvg Capture time: %dms\n", result)
	fmt.Printf("\tMax FPS: %dfps\n", 1000/(result+10))

	vm.done = make(chan bool)
	vm.timer = time.NewTicker(time.Duration(result+10) * time.Millisecond)
	go vm.run()
}
