package cockpit_stream

type CaptureResultHandler interface {
	Handle(result *ScreenCaptureResult)
}
