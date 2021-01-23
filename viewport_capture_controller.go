package cockpit_stream

import (
	"github.com/kbinani/screenshot"
	"image"
	"image/draw"
	"sync"
	"time"
)

type ViewportCaptureResult struct {
	T      time.Time
	screen *image.RGBA
	mutex  sync.RWMutex
}

func (scr *ViewportCaptureResult) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	scr.mutex.RLock()
	defer scr.mutex.RUnlock()

	draw.Draw(dst, bounds, scr.screen, at, draw.Src)
}

type ViewportCaptureController struct {
	captureFps int

	bounds      image.Rectangle
	boundsMutex sync.RWMutex

	ticker     *time.Ticker
	tickerDone chan bool

	listeners      []chan *ViewportCaptureResult
	listenersMutex sync.RWMutex
}

type ViewCaptureControllerConfig struct {
	TargetCaptureFps int
	Bounds           image.Rectangle
}

func NewViewportCaptureController(configure func(config *ViewCaptureControllerConfig)) *ViewportCaptureController {
	cfg := ViewCaptureControllerConfig{
		TargetCaptureFps: 30,
		Bounds:           image.Rect(0, 0, 100, 100),
	}

	configure(&cfg)

	return &ViewportCaptureController{
		captureFps: cfg.TargetCaptureFps,
		bounds:     cfg.Bounds,
	}
}

func (scc *ViewportCaptureController) run() {
	timeout := 1000 / scc.captureFps
	scc.ticker = time.NewTicker(time.Duration(timeout) * time.Millisecond)
	scc.tickerDone = make(chan bool)

	go func() {
		for {
			select {
			case <-scc.ticker.C:
				//ctx := context.Background()
				//ctx :=
				start := time.Now()
				scc.boundsMutex.RLock()
				img, err := screenshot.CaptureRect(scc.bounds)
				scc.boundsMutex.RUnlock()
				if err != nil {
					continue
				}

				// create result
				result := ViewportCaptureResult{
					screen: img,
					T:      start,
				}

				// notify listeners
				scc.listenersMutex.RLock()
				for _, listener := range scc.listeners {
					listener <- &result
				}
				scc.listenersMutex.RUnlock()
			case <-scc.tickerDone:
				close(scc.tickerDone)
				scc.ticker.Stop()

				return
			}
		}
	}()
}

func (scc *ViewportCaptureController) AddListener(listener chan *ViewportCaptureResult) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Unlock()

	scc.listeners = append(scc.listeners, listener)
}

func (scc *ViewportCaptureController) RemoveListener(listener chan *ViewportCaptureResult) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Lock()

	for i := range scc.listeners {
		if scc.listeners[i] == listener {
			scc.listeners = append(scc.listeners[:i], scc.listeners[i+1:]...)
		}
	}
}

func (scc *ViewportCaptureController) SetBounds(bounds image.Rectangle) {
	scc.boundsMutex.Lock()
	defer scc.boundsMutex.Unlock()

	scc.bounds = bounds
}

func (scc *ViewportCaptureController) Start() {
	scc.run()
}

func (scc *ViewportCaptureController) Stop() {
	scc.tickerDone <- true
}
