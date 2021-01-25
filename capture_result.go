package cockpit_stream

import (
	"context"
	"errors"
	"image"
	"image/draw"
	"sync"
	"time"
)

type CaptureResult struct {
	T      time.Time
	screen *image.RGBA
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

func (cr *CaptureResult) Slice(dst *image.RGBA, bounds image.Rectangle, at image.Point) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	draw.Draw(dst, bounds, cr.screen, at, draw.Src)
}
