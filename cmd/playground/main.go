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

	img1, _ := cockpit_stream.OpenPng("test_data/xor_prev.png")
	img2, _ := cockpit_stream.OpenPng("test_data/xor_next.png")

	buffer := cockpit_stream.NewBufferWithSize(len(img1.Pix))

	for i := 0; i < 10; i++ {
		start := time.Now()
		cockpit_stream.CalculateBitmask(img1, img2, buffer)
		end := time.Now().Sub(start)
		fmt.Printf("single: %fs -- %dms\n", end.Seconds(), end.Milliseconds())

		start = time.Now()
		cockpit_stream.CalculateBitmaskProcPerRow(img1, img2, buffer)
		end = time.Now().Sub(start)
		fmt.Printf("multi: %fs -- %dms\n", end.Seconds(), end.Milliseconds())

		start = time.Now()
		cockpit_stream.CalculateBitmask4ProcPerRow(img1, img2, buffer)
		end = time.Now().Sub(start)
		fmt.Printf("multi2: %fs -- %dms\n", end.Seconds(), end.Milliseconds())

		fmt.Println("-----------")
	}

}
