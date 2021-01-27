package cockpit_stream

import (
	"image"
	"sync"
)

func CalculateBitmask(prev *image.RGBA, next *image.RGBA, buffer *Buffer) (bool, error) {
	var wg sync.WaitGroup

	changed := false
	procSpan := prev.Stride / 2
	for idx := 0; idx < len(prev.Pix); idx += procSpan {
		wg.Add(1)
		go func(startIdx int, endIdx int) {
			//fmt.Printf("args: %d-%d\n", startIdx, endIdx)
			var pixel byte
			for i := startIdx; i < endIdx; i++ {
				pixel = next.Pix[i] ^ prev.Pix[i]
				buffer.Bytes[i] = pixel
				if pixel != 0 {
					changed = true
				}
			}
			wg.Done()
		}(idx, idx+procSpan-1)
	}

	wg.Wait()
	return changed, nil
}
