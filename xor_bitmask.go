package cockpit_stream

import (
	"image"
)

type XorBitmask struct {
	buffer *Buffer
}

func NewXorBitmask(w int, h int) *XorBitmask {
	return &XorBitmask{
		buffer: NewBuffer(w, h),
	}
}

func CalculateBitmask(prev *image.RGBA, next *image.RGBA, buffer *Buffer) {
	// bitmask length is same as buffer
	buffer.Size = len(buffer.Bytes)
	for i := 0; i < buffer.Size; i++ {
		buffer.Bytes[i] = next.Pix[i] ^ prev.Pix[i]
	}
}
