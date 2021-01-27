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
	id              string
	viewports       *ViewportContainer
	serverViewports *ViewportContainer
	metricsService  *metrics.Service
	outputImage     bool
	outputPath      string
	prevImage       *image.RGBA
	curImage        *image.RGBA
	clientRenderPos map[string]image.Point

	xorMask           *Buffer
	compressionBuffer *Buffer

	mutex sync.RWMutex
}

func NewViewportStreamHandler(id string, container *ViewportContainer, metrics *metrics.Service) *ScreenCaptureHandler {
	return &ScreenCaptureHandler{
		id:              id,
		serverViewports: container,
		viewports:       NewViewportContainer(),
		metricsService:  metrics,
		prevImage:       &image.RGBA{},
		curImage:        &image.RGBA{},
		clientRenderPos: make(map[string]image.Point),
	}
}

func (sch *ScreenCaptureHandler) Handle(result *CaptureResult) {
	sch.mutex.Lock()
	defer sch.mutex.Unlock()

	defer sch.metricsService.MeasureTimeForClient(time.Now(), metrics.MetricSampleViewportHandler, sch.id)
	ctx, _ := result.GetCaptureContext()

	var wg sync.WaitGroup
	wg.Add(sch.viewports.Count())
	sch.viewports.Each(func(name string, viewport *Viewport) {
		serverViewport, err := sch.serverViewports.Get(name)
		if err != nil {
			// TODO: handle err
			return
		}

		go func() {
			result.Slice(sch.curImage, viewport.PositionRect(), serverViewport.Point())
			wg.Done()
		}()
	})
	wg.Wait()

	if sch.outputImage {
		SavePng(sch.curImage, path.Join(sch.outputPath, fmt.Sprintf("%s.png", sch.id)))
	}

	// encryption buffer
	// compression buffer

	// compute xor mask
	if err := CalculateBitmask(sch.prevImage, sch.curImage, sch.xorMask); err != nil {
		// TODO: handle error
	}
	// compress mask
	if _, err := CompressBuffer(sch.xorMask); err != nil {
		// TODO: handle error
	}

	// send

	// shuffle screens for reuse
	tmp := sch.prevImage
	sch.prevImage = sch.curImage
	sch.curImage = tmp

	sch.metricsService.MeasureTimeForClient(ctx.StartTime, metrics.MetricPipelineExecutionTime, sch.id)
}

func (sch *ScreenCaptureHandler) RegisterViewport(name string, posX int, posY int) error {
	sch.mutex.Lock()
	defer sch.mutex.Unlock()
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

	// create new xor mask
	sch.xorMask = NewBuffer(newCurImage.Rect.Dx(), newCurImage.Rect.Dy())

	sch.clientRenderPos[name] = image.Point{X: posX, Y: posY}

	sch.compressionBuffer = NewCompressionBuffer(len(sch.xorMask.Bytes))

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
