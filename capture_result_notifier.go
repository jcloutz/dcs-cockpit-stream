package cockpit_stream

import (
	"sync"
)

type CaptureResultNotifier struct {
	listeners      []CaptureResultHandler
	listenersMutex sync.RWMutex
}

func NewCaptureResultNotifier() *CaptureResultNotifier {
	return &CaptureResultNotifier{}
}

func (scc *CaptureResultNotifier) AddListener(listener CaptureResultHandler) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Unlock()

	scc.listeners = append(scc.listeners, listener)
}

func (scc *CaptureResultNotifier) RemoveListener(listener CaptureResultHandler) {
	scc.listenersMutex.Lock()
	defer scc.listenersMutex.Unlock()

	for i := range scc.listeners {
		if scc.listeners[i] == listener {
			scc.listeners = append(scc.listeners[:i], scc.listeners[i+1:]...)
		}
	}
}

func (scc *CaptureResultNotifier) Notify(result *ScreenCaptureResult) {
	scc.listenersMutex.RLock()
	defer scc.listenersMutex.RUnlock()

	for _, listener := range scc.listeners {
		listener.Handle(result)
	}
}
