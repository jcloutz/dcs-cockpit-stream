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

func (i *OffsetImage) Slice(dst *image.RGBA, destRect image.Rectangle, at image.Point) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	offset := at.Sub(i.offset)

	draw.Draw(dst, destRect, i.RGBA, offset, draw.Src)
}

func (i *OffsetImage) SliceRaw(dst *image.RGBA, destRect image.Rectangle, at image.Point) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	draw.Draw(dst, destRect, i, at, draw.Src)
}
