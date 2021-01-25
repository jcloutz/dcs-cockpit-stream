package cockpit_stream

import (
	"image"
)

type CaptureController interface {
	AddListener(listener CaptureResultHandler)
	RemoveListener(listener CaptureResultHandler)
	SetBounds(bounds image.Rectangle)
	Start()
	Stop()
}
