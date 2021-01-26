package cockpit_stream

import (
	"image"
	"image/draw"
	"sync"
)

type OffsetImage struct {
	*image.RGBA

	offset image.Point

	mutex sync.RWMutex
}

func NewOffsetImage(img *image.RGBA, offset image.Point) *OffsetImage {
	return &OffsetImage{
		RGBA:   img,
		offset: offset,
	}
}

func (i *OffsetImage) Offset() image.Point {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return i.offset
}

func (i *OffsetImage) CalcOffset(point image.Point) image.Point {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return point.Sub(i.offset)
}

func (i *OffsetImage) Slice(dst *image.RGBA, viewport *Viewport) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	offset := viewport.Position.Sub(i.offset)
	SavePng(i.RGBA, "output/src.png")

	draw.Draw(dst, viewport.Bounds, i, offset, draw.Src)
}

// Lock the Viewport for mutation
func (i *OffsetImage) Lock() {
	i.mutex.Lock()
}

// Unlock the Viewport after mutation
func (i *OffsetImage) Unlock() {
	i.mutex.Unlock()
}

// RLock the Viewport for reading
func (i *OffsetImage) RLock() {
	i.mutex.RLock()
}

// RUnlock the Viewport after reading
func (i *OffsetImage) RUnlock() {
	i.mutex.RUnlock()
}
