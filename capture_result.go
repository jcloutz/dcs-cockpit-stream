package cockpit_stream

import (
	"context"
	"errors"
	"image"
	"sync"
	"time"
)

type CaptureResult struct {
	T      time.Time
	screen *OffsetImage
	mutex  sync.RWMutex
	Ctx    context.Context
}

func (cr *CaptureResult) GetCaptureContext() (*CaptureContext, error) {
	capCtx, ok := cr.Ctx.Value(CaptureContextKey).(*CaptureContext)
	if !ok {
		return nil, errors.New("capture context not found")
	}

	return capCtx, nil
}

func (cr *CaptureResult) Slice(dst *image.RGBA, destRect image.Rectangle, srcPoint image.Point) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	cr.screen.Slice(dst, destRect, srcPoint)
}
