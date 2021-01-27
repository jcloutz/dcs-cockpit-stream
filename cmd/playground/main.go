package main

import (
	"fmt"
	"time"

	"github.com/jcloutz/cockpit_stream"
)

func main() {
	//rect := image.Rect(0, 0, 10, 10)
	//
	//img1 := image.NewRGBA(rect)
	//img2 := image.NewRGBA(rect)
	//
	//buffer := cockpit_stream.NewBuffer(10, 10)
	//
	//cockpit_stream.CalculateBitmaskProcPerRow(img1, img2, buffer)

	img1, _ := cockpit_stream.OpenPng("output/client2.png")
	img2, _ := cockpit_stream.OpenPng("output/client2-prev.png")

	buffer := cockpit_stream.NewBufferWithSize(len(img1.Pix))

	for i := 0; i < 10; i++ {
		start := time.Now()
		cockpit_stream.CalculateBitmask(img1, img2, buffer)
		end := time.Now().Sub(start)
		fmt.Printf("single: %fs -- %dms -- %d micro s\n", end.Seconds(), end.Milliseconds(), end.Microseconds())

	}

}
