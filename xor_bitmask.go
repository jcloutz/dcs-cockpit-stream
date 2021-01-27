package cockpit_stream

import (
	"image"
	"sync"
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
	for i := 0; i < buffer.Size; i++ {
		buffer.Bytes[i] = next.Pix[i] ^ prev.Pix[i]
	}

	return nil
}

func CalculateBitmaskProcPerRow(prev *image.RGBA, next *image.RGBA, buffer *Buffer) error {

	var wg sync.WaitGroup
	for rowOffset := 0; rowOffset < buffer.Size/prev.Stride; rowOffset += prev.Stride {
		wg.Add(1)
		go func(startIdx int, stride int) {
			for colOffset := startIdx; colOffset < startIdx+stride; colOffset++ {
				buffer.Bytes[colOffset] = next.Pix[colOffset] ^ prev.Pix[colOffset]
			}
			wg.Done()
		}(rowOffset, prev.Stride)
	}
	wg.Wait()
	return nil
}
