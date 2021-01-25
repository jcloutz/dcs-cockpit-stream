package cockpit_stream

import (
	"context"
	"errors"
	"github.com/jcloutz/cockpit_stream/metrics"
	"image"
	"image/draw"
	"sync"
	"time"

	"github.com/kbinani/screenshot"
)

const (
	CaptureContextKey = iota
)

type CaptureContext struct {
	StartTime time.Time
	Metric    *metrics.Service
	ClientId  string
}

func GetCaptureContext(ctx context.Context) (*CaptureContext, error) {
	capCtx, ok := ctx.Value(CaptureContextKey).(*CaptureContext)
	if !ok {
		return nil, errors.New("capture context not found")
	}

	return capCtx, nil
}

type ScreenCaptureResult struct {
	T      time.Time
	screen *image.RGBA
	mutex  sync.RWMutex
	Ctx    context.Context
}

func (vpr *ScreenCaptureResult) GetCaptureContext() (*CaptureContext, error) {
	return GetCaptureContext(vpr.Ctx)
}

func (scr *ScreenCaptureResult) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	scr.mutex.RLock()
	defer scr.mutex.RUnlock()

	draw.Draw(dst, bounds, scr.screen, at, draw.Src)
}

type DesktopCaptureController struct {
	captureFps int

	bounds      image.Rectangle
	boundsMutex sync.RWMutex

	ticker     *time.Ticker
	tickerDone chan bool

	listeners      []CaptureResultHandler
	listenersMutex sync.RWMutex
	metricsService *metrics.Service
}

type ViewCaptureControllerConfig struct {
	TargetCaptureFps int
	Bounds           image.Rectangle
	Metrics          *metrics.Service
}

func NewViewportCaptureController(configure func(config *ViewCaptureControllerConfig)) *DesktopCaptureController {
	cfg := ViewCaptureControllerConfig{
		TargetCaptureFps: 30,
		Bounds:           image.Rect(0, 0, 100, 100),
	}

	configure(&cfg)

	return &DesktopCaptureController{
		captureFps:     cfg.TargetCaptureFps,
		bounds:         cfg.Bounds,
		metricsService: cfg.Metrics,
	}
}

func (scc *DesktopCaptureController) run() {
	timeout := 1000 / scc.captureFps
	scc.ticker = time.NewTicker(time.Duration(timeout) * time.Millisecond)
	scc.tickerDone = make(chan bool)

	go func() {
		for {
			select {
			case <-scc.ticker.C:
				start := time.Now()

				ctx := context.Background()
				ctx = context.WithValue(ctx, CaptureContextKey, &CaptureContext{
					StartTime: time.Now(),
					Metric:    scc.metricsService,
					ClientId:  "",
				})

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
					Ctx:    ctx,
				}

				// notify listeners
				scc.listenersMutex.RLock()
				for _, listener := range scc.listeners {
					listener.Handle(&result)
				}
				scc.listenersMutex.RUnlock()
				scc.metricsService.MeasureTime(start, metrics.MetricSampleCaptureController)
				scc.metricsService.Count(metrics.MetricFrameCounter, 1)
			case <-scc.tickerDone:
				close(scc.tickerDone)
				scc.ticker.Stop()

				return
			}
		}
	}()
}

func (scc *DesktopCaptureController) AddListener(listener CaptureResultHandler) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Unlock()

	scc.listeners = append(scc.listeners, listener)
}

func (scc *DesktopCaptureController) RemoveListener(listener CaptureResultHandler) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Lock()

	for i := range scc.listeners {
		if scc.listeners[i] == listener {
			scc.listeners = append(scc.listeners[:i], scc.listeners[i+1:]...)
		}
	}
}

func (scc *DesktopCaptureController) SetBounds(bounds image.Rectangle) {
	scc.boundsMutex.Lock()
	defer scc.boundsMutex.Unlock()

	scc.bounds = bounds
}

func (scc *DesktopCaptureController) Start() {
	scc.run()
}

func (scc *DesktopCaptureController) Stop() {
	scc.tickerDone <- true
}
