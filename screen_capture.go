package cockpit_stream

import (
	"image"
	"image/draw"
	"sync"

	"github.com/kbinani/screenshot"
)

type ScreenCapper interface {
	Update(bounds image.Rectangle) error
	Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) int64
	HasNext(callerIdx int64) bool
}

type ScreenCapture struct {
	screen *image.RGBA

	mutex sync.RWMutex

	index int64
}

func NewScreenCapture() ScreenCapper {
	return &ScreenCapture{
		screen: image.NewRGBA(image.Rect(0, 0, 10, 10)),
		mutex:  sync.RWMutex{},
		index:  0,
	}
}

func (sc *ScreenCapture) Update(bounds image.Rectangle) error {
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return err
	} else {
		sc.mutex.Lock()
		sc.index++
		sc.screen = img
		sc.mutex.Unlock()
	}

	return nil
}

func (sc *ScreenCapture) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) int64 {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	draw.Draw(dst, bounds, sc.screen, at, draw.Src)
	return sc.index
}

func (sc *ScreenCapture) HasNext(callerIdx int64) bool {
	return sc.index > callerIdx
}
