package cockpit_stream

type CaptureResultHandler interface {
	Handle(result *CaptureResult)
}

type CallbackCaptureHandler struct {
	callback func(result *CaptureResult)
}

func (ch *CallbackCaptureHandler) Handle(result *CaptureResult) {
	ch.callback(result)
}

func (ch *CallbackCaptureHandler) OnReceive(callback func(result *CaptureResult)) {
	ch.callback = callback
}
