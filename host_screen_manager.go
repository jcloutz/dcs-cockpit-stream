package cockpit_stream

import (
	"time"
)

type ViewportListenerResult struct {
	Viewports ViewportContainerReader
	T         time.Time
}
type HostScreenManager struct {
	captureFps int

	screenCap ImageCapturer

	viewportManger *ServerViewportManager

	ticker     *time.Ticker
	tickerDone chan bool

	listeners []chan *ViewportListenerResult
}

type HostScreenManagerConfig struct {
	TargetCaptureFps int
	ScreenCapper     ImageCapturer
	ViewportManager  *ServerViewportManager
}

func NewHostScreenManager(configure func(config *HostScreenManagerConfig)) *HostScreenManager {
	cfg := HostScreenManagerConfig{
		TargetCaptureFps: 30,
		ScreenCapper:     nil,
		ViewportManager:  nil,
	}

	configure(&cfg)

	return &HostScreenManager{
		captureFps:     cfg.TargetCaptureFps,
		screenCap:      cfg.ScreenCapper,
		viewportManger: cfg.ViewportManager,
	}
}

func (hsm *HostScreenManager) run() {
	hsm.screenCap.SetBounds(hsm.viewportManger.computedBounds)
	timeout := 1000 / hsm.captureFps
	hsm.ticker = time.NewTicker(time.Duration(timeout) * time.Millisecond)
	hsm.tickerDone = make(chan bool)

	go func() {
		for {
			select {
			case <-hsm.ticker.C:
				start := time.Now()
				if err := hsm.screenCap.Capture(); err != nil {
					continue
				}

				// go update viewports -> *screencap
				hsm.viewportManger.UpdateViewports(hsm.screenCap)
				// notify listeners
				for _, listener := range hsm.listeners {
					listener <- &ViewportListenerResult{
						Viewports: hsm.viewportManger,
						T:         start,
					}
				}
				return
			case <-hsm.tickerDone:
				close(hsm.tickerDone)
				hsm.ticker.Stop()

				return
			}
		}
	}()
}

func (hsm *HostScreenManager) OnCaptureUpdate(listener chan *ViewportListenerResult) {
	hsm.listeners = append(hsm.listeners, listener)
}

func (hsm *HostScreenManager) Start() {
	hsm.run()
}

func (hsm *HostScreenManager) Stop() {
	hsm.tickerDone <- true
}
