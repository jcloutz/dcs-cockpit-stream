package cockpit_stream

import (
	"github.com/kbinani/screenshot"
	"image"
	"image/draw"
	"sync"
	"time"
)

type ScreenCaptureResult struct {
	T      time.Time
	screen *image.RGBA
	mutex  sync.RWMutex
}

func (scr *ScreenCaptureResult) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	scr.mutex.RLock()
	defer scr.mutex.RUnlock()
	draw.Draw(dst, bounds, scr.screen, at, draw.Src)
}

type ScreenCaptureController struct {
	captureFps int

	bounds      image.Rectangle
	boundsMutex sync.RWMutex

	ticker     *time.Ticker
	tickerDone chan bool

	listeners      []chan *ScreenCaptureResult
	listenersMutex sync.RWMutex
}

type ScreenCaptureControllerConfig struct {
	TargetCaptureFps int
	Bounds           image.Rectangle
}

func NewHostScreenManager(configure func(config *ScreenCaptureControllerConfig)) *ScreenCaptureController {
	cfg := ScreenCaptureControllerConfig{
		TargetCaptureFps: 30,
		Bounds:           image.Rect(0, 0, 100, 100),
	}

	configure(&cfg)

	return &ScreenCaptureController{
		captureFps: cfg.TargetCaptureFps,
		bounds:     cfg.Bounds,
	}
}

func (scc *ScreenCaptureController) run() {
	timeout := 1000 / scc.captureFps
	scc.ticker = time.NewTicker(time.Duration(timeout) * time.Millisecond)
	scc.tickerDone = make(chan bool)

	go func() {
		for {
			select {
			case <-scc.ticker.C:
				start := time.Now()
				scc.boundsMutex.RLock()
				img, err := screenshot.CaptureRect(scc.bounds)
				scc.boundsMutex.RUnlock()
				if err != nil {
					continue
				}

				// create result
				result := ScreenCaptureResult{
					screen: img,
					T:      start,
				}

				// notify listeners
				scc.listenersMutex.RLock()
				for _, listener := range scc.listeners {
					listener <- &result
				}
				scc.listenersMutex.RUnlock()
				//return
			case <-scc.tickerDone:
				close(scc.tickerDone)
				scc.ticker.Stop()

				return
			}
		}
	}()
}

func (scc *ScreenCaptureController) AddListener(listener chan *ScreenCaptureResult) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Unlock()

	scc.listeners = append(scc.listeners, listener)
}

func (scc *ScreenCaptureController) RemoveListener(listener chan *ScreenCaptureResult) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Lock()

	for i := range scc.listeners {
		if scc.listeners[i] == listener {
			scc.listeners = append(scc.listeners[:i], scc.listeners[i+1:]...)
		}
	}
}

func (scc *ScreenCaptureController) SetBounds(bounds image.Rectangle) {
	scc.boundsMutex.Lock()
	defer scc.boundsMutex.Unlock()

	scc.bounds = bounds
}

func (scc *ScreenCaptureController) Start() {
	scc.run()
}

func (scc *ScreenCaptureController) Stop() {
	scc.tickerDone <- true
}
