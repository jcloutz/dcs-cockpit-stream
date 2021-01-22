package cockpit_stream

import (
	"image"
	"image/draw"
	"sync"

	"github.com/kbinani/screenshot"
)

type ImageSlicer interface {
	Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point)
}

type ImageCapturer interface {
	Capture() error
	SetBounds(bounds image.Rectangle)
	ImageSlicer
}

var _ ImageCapturer = &ScreenCapture{}

type ScreenCapture struct {
	screen *image.RGBA

	mutex sync.RWMutex

	bounds image.Rectangle
}

func NewScreenCapture() ImageCapturer {
	return &ScreenCapture{
		screen: image.NewRGBA(image.Rect(0, 0, 10, 10)),
		mutex:  sync.RWMutex{},
		bounds: image.Rectangle{},
	}
}

func (sc *ScreenCapture) Capture() error {
	img, err := screenshot.CaptureRect(sc.bounds)
	if err != nil {
		return err
	} else {
		sc.mutex.Lock()
		sc.screen = img
		sc.mutex.Unlock()
	}

	return nil
}

func (sc *ScreenCapture) SetBounds(bounds image.Rectangle) {
	sc.bounds = bounds
}

func (sc *ScreenCapture) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	draw.Draw(dst, bounds, sc.screen, at, draw.Src)
}
