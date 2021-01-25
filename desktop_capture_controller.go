package cockpit_stream

import (
	"context"
	"image"
	"sync"
	"time"

	"github.com/jcloutz/cockpit_stream/metrics"

	"github.com/kbinani/screenshot"
)

type DesktopCaptureController struct {
	captureFps int

	bounds      image.Rectangle
	boundsMutex sync.RWMutex

	ticker     *time.Ticker
	tickerDone chan bool

	notifier       *CaptureResultNotifier
	metricsService *metrics.Service
}

type DesktopCaptureControllerConfig struct {
	TargetCaptureFps int
	Bounds           image.Rectangle
	Metrics          *metrics.Service
}

func NewViewportCaptureController(configure func(config *DesktopCaptureControllerConfig)) *DesktopCaptureController {
	cfg := DesktopCaptureControllerConfig{
		TargetCaptureFps: 30,
		Bounds:           image.Rect(0, 0, 100, 100),
	}

	configure(&cfg)

	return &DesktopCaptureController{
		captureFps:     cfg.TargetCaptureFps,
		bounds:         cfg.Bounds,
		metricsService: cfg.Metrics,
		notifier:       NewCaptureResultNotifier(),
	}
}

func (scc *DesktopCaptureController) execute() error {
	start := time.Now()
	// generate metrics
	defer scc.metricsService.MeasureTime(start, metrics.MetricSampleCaptureController)
	defer scc.metricsService.Count(metrics.MetricFrameCounter, 1)

	// create context
	ctx := context.Background()
	ctx = context.WithValue(ctx, CaptureContextKey, &CaptureContext{
		StartTime: time.Now(),
		Metric:    scc.metricsService,
		ClientId:  "",
	})

	// capture image
	scc.boundsMutex.RLock()
	img, err := screenshot.CaptureRect(scc.bounds)
	scc.boundsMutex.RUnlock()
	if err != nil {
		return err
	}

	// create result
	result := CaptureResult{
		screen: img,
		T:      start,
		Ctx:    ctx,
	}

	// notify notifier
	scc.notifier.Notify(&result)

	return nil
}

func (scc *DesktopCaptureController) run() {
	timeout := 1000 / scc.captureFps
	scc.ticker = time.NewTicker(time.Duration(timeout) * time.Millisecond)
	scc.tickerDone = make(chan bool)

	go func() {
		for {
			select {
			case <-scc.ticker.C:
				if err := scc.execute(); err != nil {
					continue
				}
			case <-scc.tickerDone:
				close(scc.tickerDone)
				scc.ticker.Stop()

				return
			}
		}
	}()
}

func (scc *DesktopCaptureController) RunOnce() error {
	return scc.execute()
}

func (scc *DesktopCaptureController) AddListener(listener CaptureResultHandler) {
	scc.notifier.AddListener(listener)
}

func (scc *DesktopCaptureController) RemoveListener(listener CaptureResultHandler) {
	scc.notifier.RemoveListener(listener)
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
