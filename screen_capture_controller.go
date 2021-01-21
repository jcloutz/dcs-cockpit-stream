package cockpit_stream

import (
	"fmt"
	"image"
	"image/draw"
	"sync"
	"time"

	"github.com/kbinani/screenshot"
)

type PartialCapturer interface {
	Capture(image image.RGBA)
}

type ScreenCaptureController struct {
	img             *image.RGBA
	viewportManager *ServerViewportManager
	imgMutex        sync.RWMutex
	bounds          image.Rectangle
	running         bool
	fps             int
	timeout         time.Duration
	screenIndex     int64
}

func New(bounds *image.Rectangle, fps int) *ScreenCaptureController {
	return &ScreenCaptureController{
		img:         image.NewRGBA(*bounds),
		imgMutex:    sync.RWMutex{},
		bounds:      *bounds,
		fps:         fps,
		timeout:     time.Duration(1000 / fps),
		screenIndex: 1,
	}
}

func (cc *ScreenCaptureController) Start() {
	cc.running = true
	cc.run()

}

func (cc *ScreenCaptureController) Stop() {
	cc.running = false
}

func (cc *ScreenCaptureController) run() {
	go func() {
		for cc.running {
			//startFrame := time.Now()
			img, err := screenshot.CaptureRect(cc.bounds)
			if err != nil {
				fmt.Println(err)
			} else {
				cc.imgMutex.Lock()
				cc.screenIndex++
				//fmt.Printf("--- CAPTURE FRAME %d --- \n", cc.screenIndex)
				cc.img = img
				cc.imgMutex.Unlock()
			}
			//elapsed := time.Duration(time.Now().Sub(startFrame).Milliseconds())
			//time.Sleep((cc.timeout - elapsed) * time.Millisecond)
		}
	}()
}

func (cc *ScreenCaptureController) GetScreen(dst *image.RGBA, bounds *image.Rectangle) int64 {
	cc.imgMutex.RLock()
	defer cc.imgMutex.RUnlock()
	draw.Draw(dst, image.Rect(0, 0, 50, 50), cc.img, image.Point{
		X: 0,
		Y: 0,
	}, draw.Src)
	return cc.screenIndex
}

func (cc *ScreenCaptureController) Index() int64 {
	return cc.screenIndex
}
