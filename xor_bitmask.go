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
	//start := time.Now()
	for i := 0; i < buffer.Size; i++ {
		buffer.Bytes[i] = next.Pix[i] ^ prev.Pix[i]
	}

	//end := time.Now().Sub(start)
	//fmt.Printf("single: %fs -- %dms\n", end.Seconds(), end.Milliseconds())
	return nil
}

func CalculateBitmaskProcPerRow(prev *image.RGBA, next *image.RGBA, buffer *Buffer) error {
	//start := time.Now()
	var wg sync.WaitGroup

	for rowOffset := 0; rowOffset < buffer.Size/prev.Stride; rowOffset++ {
		wg.Add(1)
		//time.Sleep(time.Duration(counter) * time.Millisecond)
		go func(startIdx int, stride int) {
			colOffset := 0
			for colOffset = startIdx; colOffset < startIdx+stride; colOffset++ {
				buffer.Bytes[colOffset] = next.Pix[colOffset] ^ prev.Pix[colOffset]
			}
			//fmt.Printf("[range] %d-%d\n", startIdx, colOffset)
			wg.Done()
		}(rowOffset*prev.Stride, prev.Stride)
	}
	wg.Wait()
	//end := time.Now().Sub(start)
	//fmt.Printf("multi: %fs -- %dms\n", end.Seconds(), end.Milliseconds())
	return nil
}
func CalculateBitmask4ProcPerRow(prev *image.RGBA, next *image.RGBA, buffer *Buffer) error {
	var wg sync.WaitGroup

	procSpan := prev.Stride / 4
	for idx := 0; idx < buffer.Size; idx += procSpan {
		wg.Add(1)
		go func(startIdx int, endIdx int) {
			//fmt.Printf("args: %d-%d\n", startIdx, endIdx)
			c := 0
			for i := startIdx; i < endIdx; i++ {
				buffer.Bytes[i] = next.Pix[i] ^ prev.Pix[i]
				c++
			}
			wg.Done()
		}(idx, idx+procSpan-1)
	}

	wg.Wait()
	return nil
}
