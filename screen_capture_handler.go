package cockpit_stream

import (
	"fmt"
	"image"
	"path"
	"sync"
	"time"

	"github.com/jcloutz/cockpit_stream/metrics"
)

var _ CaptureResultHandler = &ScreenCaptureHandler{}

type ScreenCaptureHandler struct {
	id               string
	serverViewports  *ViewportContainer
	handlerViewports map[string]*Viewport
	MetricsService   *metrics.Service
	outputImage      bool
	outputPath       string
	prevImage        *image.RGBA
	curImage         *image.RGBA

	container *ViewportContainer
}

func NewViewportStreamHandler(id string, container *ViewportContainer, metrics *metrics.Service) *ScreenCaptureHandler {
	return &ScreenCaptureHandler{
		id:               id,
		serverViewports:  container,
		handlerViewports: make(map[string]*Viewport),
		MetricsService:   metrics,
		container:        NewViewportContainer(),
	}
}

func (sch *ScreenCaptureHandler) Handle(result *CaptureResult) {
	defer sch.MetricsService.MeasureTimeForClient(time.Now(), metrics.MetricSampleViewportHandler, sch.id)
	ctx, _ := result.GetCaptureContext()

	var wg sync.WaitGroup
	wg.Add(len(sch.handlerViewports))
	for _, hVp := range sch.handlerViewports {
		go func(viewport *Viewport) {
			viewport.RLock()
			defer viewport.RUnlock()
			sVp, _ := sch.serverViewports.Get(hVp.Name)
			result.Slice(sch.curImage, sVp)
			wg.Done()
		}(hVp)
	}
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
	viewport, err := sch.serverViewports.Get(name)
	if err != nil {
		return err
	}

	viewport.RLock()
	defer viewport.RUnlock()
	sch.handlerViewports[name] = &Viewport{
		Name: name,
		// get the top left corner of the viewport as represented in
		// the desktop capture from the CaptureResult in Handle()
		Position: image.Point{X: posX, Y: posY},
	}

	if err = sch.recalculateBounds(); err != nil {
		return err
	}

	return nil
}

func (sch *ScreenCaptureHandler) recalculateBounds() error {
	height := 0
	width := 0

	for id, hViewport := range sch.handlerViewports {
		viewport, err := sch.serverViewports.Get(id)
		if err != nil {
			return err
		}

		viewport.RLock()
		bounds := viewport.Bounds
		viewport.RUnlock()

		hViewport.Lock()
		hViewport.Bounds = image.Rect(
			width,
			0,
			width+bounds.Dx(),
			bounds.Dy(),
		)
		hViewport.Unlock()

		width += bounds.Dx()
		if bounds.Dy() > height {
			height = bounds.Dy()
		}
	}

	newImageSize := image.Rect(0, 0, width, height)
	sch.prevImage = image.NewRGBA(newImageSize)
	sch.curImage = image.NewRGBA(newImageSize)

	return nil
}

func (sch *ScreenCaptureHandler) EnableOutput(path string) {
	sch.outputImage = true
	sch.outputPath = path
}

func (sch *ScreenCaptureHandler) DisableOutput() {
	sch.outputImage = false
	sch.outputPath = ""
}
