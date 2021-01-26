package cockpit_stream

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"path"
	"sync"
	"time"

	"github.com/jcloutz/cockpit_stream/metrics"
)

var _ CaptureResultHandler = &ScreenCaptureHandler{}

type ScreenCaptureHandler struct {
	id               string
	viewports        *ViewportContainer
	serverViewports  *ViewportContainer
	handlerViewports map[string]*Viewport
	MetricsService   *metrics.Service
	outputImage      bool
	outputPath       string
	prevImage        *image.RGBA
	curImage         *image.RGBA
	clientRenderPos  map[string]image.Point
}

func NewViewportStreamHandler(id string, container *ViewportContainer, metrics *metrics.Service) *ScreenCaptureHandler {
	return &ScreenCaptureHandler{
		id:              id,
		serverViewports: container,
		viewports:       NewViewportContainer(),
		MetricsService:  metrics,
		prevImage:       &image.RGBA{},
		curImage:        &image.RGBA{},
		clientRenderPos: make(map[string]image.Point),
	}
}

func (sch *ScreenCaptureHandler) Handle(result *CaptureResult) {
	defer sch.MetricsService.MeasureTimeForClient(time.Now(), metrics.MetricSampleViewportHandler, sch.id)
	ctx, _ := result.GetCaptureContext()

	var wg sync.WaitGroup
	wg.Add(sch.viewports.Count())

	sch.viewports.Each(func(name string, viewport *Viewport) {
		serverViewport, err := sch.serverViewports.Get(name)
		if err != nil {
			// TODO: handle err
			return
		}

		result.Slice(sch.curImage, viewport)
		wg.Done()
	})
	//for _, hVp := range sch.handlerViewports {
	//	go func(viewport *Viewport) {
	//		sVp, _ := sch.serverViewports.Get(hVp.Name())
	//		result.Slice(sch.curImage, sVp)
	//		wg.Done()
	//	}(hVp)
	//}
	wg.Wait()

	if sch.outputImage {
		SavePng(sch.curImage, path.Join(sch.outputPath, fmt.Sprintf("%s.png", sch.id)))
	}

	// imgPrev
	// imgCur
	// encryption buffer
	// compression buffer

	// compute xor mask
	// compress mask
	// send

	// shuffle screens for reuse
	tmp := sch.prevImage
	sch.prevImage = sch.curImage
	sch.curImage = tmp

	sch.MetricsService.MeasureTimeForClient(ctx.StartTime, metrics.MetricPipelineExecutionTime, sch.id)
}

func (sch *ScreenCaptureHandler) RegisterViewport(name string, posX int, posY int) error {
	// check to see if the viewport is already registered
	if exists := sch.viewports.Has(name); exists {
		return errors.New("viewport already registered with handler")
	}

	// fetch the viewport from the server container
	viewport, err := sch.serverViewports.Get(name)
	if err != nil {
		return errors.New("viewport not registered with server")
	}

	// create new image to expand the canvas
	// calculate new width and height
	maxX := sch.viewports.Bounds().Max.X
	sch.viewports.Add(name, maxX, 0, viewport.Width(), viewport.Height())

	// resize canvas and copy existing image
	newCurImage := image.NewRGBA(sch.viewports.Bounds())
	newPrevImage := image.NewRGBA(sch.viewports.Bounds())
	draw.Draw(newCurImage, sch.curImage.Rect, sch.curImage, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(newPrevImage, sch.prevImage.Rect, sch.prevImage, image.Point{X: 0, Y: 0}, draw.Src)

	// reassign image
	sch.curImage = newCurImage
	sch.prevImage = newPrevImage

	sch.viewports.Add(name, newCurImage.Rect.Dx(), 0, viewport.Width(), viewport.Height())

	sch.clientRenderPos[name] = image.Point{X: posX, Y: posY}

	return nil
}

//func (sch *ScreenCaptureHandler) recalculateBounds() error {
//	height := 0
//	width := 0
//
//	for id, hViewport := range sch.handlerViewports {
//		viewport, err := sch.serverViewports.Get(id)
//		if err != nil {
//			return err
//		}
//
//		viewport.RLock()
//		bounds := viewport.Bounds
//		viewport.RUnlock()
//
//		hViewport.Lock()
//		hViewport.Bounds = image.Rect(
//			width,
//			0,
//			width+bounds.Dx(),
//			bounds.Dy(),
//		)
//		hViewport.Unlock()
//
//		width += bounds.Dx()
//		if bounds.Dy() > height {
//			height = bounds.Dy()
//		}
//	}
//
//	newImageSize := image.Rect(0, 0, width, height)
//	sch.prevImage = image.NewRGBA(newImageSize)
//	sch.curImage = image.NewRGBA(newImageSize)
//
//	return nil
//}

func (sch *ScreenCaptureHandler) EnableOutput(path string) {
	sch.outputImage = true
	sch.outputPath = path
}

func (sch *ScreenCaptureHandler) DisableOutput() {
	sch.outputImage = false
	sch.outputPath = ""
}
