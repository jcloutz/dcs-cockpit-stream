package cockpit_stream

import "image"

type HandlerViewport struct {
	// id of the viewport this belongs to
	ViewportId string

	// destination rectangle within the handlers image buffer
	DrawDestRect image.Rectangle

	// slicePosition at which to
	SrcOffset image.Point

	// position to render the screen on the client machine
	ClientRenderPosition image.Point
}
