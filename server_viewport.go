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
	targetFps      int
	computedBounds image.Rectangle
	done           chan bool
	timer          *time.Ticker

	fps int
}

func NewServerViewportManager(screenCapper ScreenCapper, targetFps int) *ServerViewportManager {
	return &ServerViewportManager{
		viewports: make(map[string]*ServerViewport),
		capture:   screenCapper,
		targetFps: targetFps,
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

func (vm *ServerViewportManager) benchmark() (int, int) {
	iterations := 120
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = vm.capture.Update(vm.computedBounds)
	}
	end := time.Now()

	avgTime := end.Sub(start).Milliseconds() / int64(iterations)
	maxFps := 1000 / avgTime

	return int(avgTime), int(maxFps)

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
	avgTime, maxFps := vm.benchmark()
	fmt.Printf("\tAvg Capture time: %dms\n", avgTime)
	fmt.Printf("\tMax FPS: %dfps\n", 1000/(avgTime))

	if maxFps > vm.targetFps {
		vm.fps = vm.targetFps
	} else {
		vm.fps = maxFps
	}
	fmt.Printf("\tTarget FPS: %dfps\n", vm.fps)

	vm.done = make(chan bool)
	vm.timer = time.NewTicker(time.Duration(1000/vm.fps) * time.Millisecond)
	go vm.run()
}
