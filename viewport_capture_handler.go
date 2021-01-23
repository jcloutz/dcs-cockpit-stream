package cockpit_stream

import (
	"image"
	"sync"
)

type HandlerViewport struct {
	// id of the viewport this belongs to
	ViewportId string

	// destination rectangle within the handlers image buffer
	DrawDestRect image.Rectangle

	// adjOffset at which to
	SrcOffset image.Point

	// position to render the screen on the client machine
	ClientRenderPosition image.Point
}
type ViewportCaptureHandler struct {
	id              string
	serverViewports *ViewportContainer

	handlerViewports map[string]*HandlerViewport

	prevImage *image.RGBA
	curImage  *image.RGBA
}

func NewViewportCaptureHandler(id string, container *ViewportContainer) *ViewportCaptureHandler {
	return &ViewportCaptureHandler{
		id:               id,
		serverViewports:  container,
		handlerViewports: make(map[string]*HandlerViewport),
	}
}

func (ch *ViewportCaptureHandler) Handle(result *ViewportCaptureResult) {
	var wg sync.WaitGroup
	wg.Add(len(ch.handlerViewports))
	for _, hVp := range ch.handlerViewports {
		go func() {
			result.Slice(ch.curImage, hVp.DrawDestRect, hVp.SrcOffset)
			wg.Done()
		}()
		wg.Wait()
	}

	// compute xor mask
	// compress mask
	// send

	// shuffle screens for reuse
	tmp := ch.prevImage
	ch.prevImage = ch.curImage
	ch.curImage = tmp
}

func (ch *ViewportCaptureHandler) RegisterViewport(name string, posX int, posY int) error {
	err := ch.serverViewports.One(name, func(viewport ViewportReader) {
		ch.handlerViewports[name] = &HandlerViewport{
			ViewportId: name,
			// get the top left corner of the viewport as represented in
			// the desktop capture from the ViewportCaptureResult in Handle()
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

func (ch *ViewportCaptureHandler) recalculateBounds() {
	height := 0
	width := 0

	for id, hViewport := range ch.handlerViewports {
		ch.serverViewports.One(id, func(viewport ViewportReader) {
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
