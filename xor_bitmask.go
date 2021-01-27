package cockpit_stream

import (
	"errors"
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

func CalculateBitmask(prev *image.RGBA, next *image.RGBA, buffer *Buffer) error {
	// bitmask length is same as mask
	prevLen := len(prev.Pix)
	nextLen := len(next.Pix)
	bufferLen := len(buffer.Bytes)

	isValid := prevLen == nextLen && prevLen == bufferLen
	if !isValid {
		return errors.New("prev, next, and mask images must be the same size to calculate the xor mask value")
	}

	for i := 0; i < bufferLen; i++ {
		buffer.Bytes[i] = next.Pix[i] ^ prev.Pix[i]
	}

	return nil
}
