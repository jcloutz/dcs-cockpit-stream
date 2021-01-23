package cockpit_stream

import (
	"image"
	"image/draw"
	"sync"
)

type ViewportContainerReader interface {
	Get(name string) (ViewportReader, error)
}

type VPM struct {
	src       ViewportMutexMap
	srcBounds image.Rectangle

	dest ViewportMutexMap

	mutex sync.RWMutex
}

func (vpm *VPM) Slice(viewport string, dst *image.RGBA, bounds image.Rectangle, at image.Point) func(src *image.RGBA) {
	vpm.mutex.RLock()
	defer vpm.mutex.RUnlock()

	return func(src *image.RGBA) {
		draw.Draw(dst, bounds, src, at, draw.Src)
	}
}

func (vpm *VPM) AddSrc(name string, vp *Viewport) {
	vpm.mutex.Lock()
	defer vpm.mutex.Unlock()

	vpm.src.Set(name, vp)
}

func (vpm *VPM) AddDest(name string, vp *Viewport) {
	vpm.mutex.Lock()
	defer vpm.mutex.Unlock()

	vpm.dest.Set(name, vp)
}

//var _ ViewportContainerReader = &ServerViewportManager{}

//type ServerViewportManager struct {
//	viewports      *ViewportMutexMap
//	viewportsMutex sync.RWMutex
//
//	capture        ImageCapturer
//	targetFps      int
//	computedBounds image.Rectangle
//	done           chan bool
//	timer          *time.Ticker
//
//	fps int
//}

//func NewServerViewportManager(screenCapper ImageCapturer, targetFps int) *ServerViewportManager {
//	return &ServerViewportManager{
//		viewports: NewViewportMutextMap(),
//		capture:   screenCapper,
//		targetFps: targetFps,
//	}
//}

//func (vm *ServerViewportManager) AddNewViewport(name string, x int, y int, width int, height int) *ServerViewportManager {
//	return vm.AddViewport(NewServerViewport(name, x, y, width, height))
//}

//func (vm *ServerViewportManager) AddViewport(viewport *Viewport) *ServerViewportManager {
//	vm.viewportsMutex.Lock()
//	vm.viewports.Set(viewport.Name, viewport)
//	vm.viewportsMutex.Unlock()
//
//	vm.recomputeBounds()
//
//	return vm
//}

// recomputeBounds will adjust
//func (vm *ServerViewportManager) recomputeBounds() {
//
//	minX := math.MaxInt16
//	maxX := math.MinInt16
//
//	minY := math.MaxInt16
//	maxY := math.MinInt16
//
//	vm.viewportsMutex.Lock()
//	vm.viewports.Each(func(name string, viewport *Viewport) {
//		if viewport.Bounds.Min.X < minX {
//			minX = viewport.Bounds.Min.X
//		}
//		if viewport.Bounds.Min.Y < minY {
//			minY = viewport.Bounds.Min.Y
//		}
//
//		if viewport.Bounds.Max.X > maxX {
//			maxX = viewport.Bounds.Max.X
//		}
//
//		if viewport.Bounds.Max.Y > maxY {
//			maxY = viewport.Bounds.Max.Y
//		}
//	})
//	vm.viewportsMutex.Unlock()
//
//	vm.computedBounds = image.Rect(minX, minY, maxX, maxY)
//
//	vm.viewportsMutex.Lock()
//	vm.viewports.Each(func(name string, viewport *Viewport) {
//		viewport.AdjPoint = viewport.Bounds.Sub(image.Point{
//			X: minX,
//			Y: minY,
//		}).Min
//	})
//	vm.viewportsMutex.Unlock()
//}

//func (vm *ServerViewportManager) UpdateViewports(cap ImageSlicer) {
//	vm.viewportsMutex.Lock()
//	defer vm.viewportsMutex.Unlock()
//
//	vm.viewports.Each(func(name string, viewport *Viewport) {
//		viewport.Update(cap)
//	})
//}

//func (vm *ServerViewportManager) benchmark() (int, int) {
//	iterations := 120
//	start := time.Now()
//	for i := 0; i < iterations; i++ {
//		_ = vm.capture.Capture()
//	}
//	end := time.Now()
//
//	avgTime := end.Sub(start).Milliseconds() / int64(iterations)
//	maxFps := 1000 / avgTime
//
//	return int(avgTime), int(maxFps)
//}

//func (vm *ServerViewportManager) Get(name string) (ViewportReader, error) {
//	vm.viewportsMutex.RLock()
//	defer vm.viewportsMutex.RUnlock()
//
//	return vm.viewports.Get(name)
//}
