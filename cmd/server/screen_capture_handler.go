package main

import (
	"github.com/jcloutz/cockpit_stream"
	"github.com/jcloutz/cockpit_stream/metrics"
	"image"
	"sync"
	"time"
)

var _ cockpit_stream.CaptureResultHandler = &ViewportStream{}

type ViewportStream struct {
	id              string
	serverViewports *cockpit_stream.ViewportContainer

	handlerViewports map[string]*cockpit_stream.HandlerViewport

	MetricsService *metrics.Service

	prevImage *image.RGBA
	curImage  *image.RGBA
}

func NewViewportStreamHandler(id string, container *cockpit_stream.ViewportContainer, metrics *metrics.Service) *ViewportStream {
	return &ViewportStream{
		id:               id,
		serverViewports:  container,
		handlerViewports: make(map[string]*cockpit_stream.HandlerViewport),
		MetricsService:   metrics,
	}
}

func (ch *ViewportStream) Handle(result *cockpit_stream.ScreenCaptureResult) {
	defer ch.MetricsService.MeasureTimeForClient(time.Now(), metrics.MetricSampleViewportHandler, ch.id)
	ctx, _ := result.GetCaptureContext()

	var wg sync.WaitGroup
	wg.Add(len(ch.handlerViewports))
	for _, hVp := range ch.handlerViewports {
		go func() {
			result.Slice(ch.curImage, hVp.DrawDestRect, hVp.SrcOffset)
			wg.Done()
		}()
	}
	wg.Wait()

	// imgPrev
	// imgCur
	// encryption buffer
	// compression buffer

	// compute xor mask
	// compress mask
	// send

	// shuffle screens for reuse
	tmp := ch.prevImage
	ch.prevImage = ch.curImage
	ch.curImage = tmp

	ch.MetricsService.MeasureTimeForClient(ctx.StartTime, metrics.MetricPipelineExecutionTime, ch.id)
}

func (ch *ViewportStream) RegisterViewport(name string, posX int, posY int) error {
	err := ch.serverViewports.One(name, func(viewport cockpit_stream.ViewportReader) {
		ch.handlerViewports[name] = &cockpit_stream.HandlerViewport{
			ViewportId: name,
			// get the top left corner of the viewport as represented in
			// the desktop capture from the ScreenCaptureResult in Handle()
			SrcOffset:            viewport.GetAdjOffset(),
			ClientRenderPosition: image.Point{X: posX, Y: posY},
		}
	})
	if err != nil {
		return err
	}

	ch.recalculateBounds()

	return nil
}

func (ch *ViewportStream) recalculateBounds() {
	height := 0
	width := 0

	for id, hViewport := range ch.handlerViewports {
		ch.serverViewports.One(id, func(viewport cockpit_stream.ViewportReader) {
			bounds := viewport.GetBounds()

			hViewport.DrawDestRect = image.Rect(
				width,
				0,
				width+bounds.Dx(),
				bounds.Dy(),
			)

			width += bounds.Dx()
			if bounds.Dy() > height {
				height = bounds.Dy()
			}

		})
	}

	newImageSize := image.Rect(0, 0, width, height)
	ch.prevImage = image.NewRGBA(newImageSize)
	ch.curImage = image.NewRGBA(newImageSize)
}
