package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
)

func main() {

	size := 500

	var img *image.RGBA
	for i := 0; i < 300; i++ {
		s := time.Now()
		img, _ = screenshot.Capture(0, 0, size*2, size*2)
		t := time.Now().Sub(s).Milliseconds()
		fmt.Printf("elapsed: %dms\n", t)
	}

	//fmt.Printf("avg: %dms\n", t/300)

	sizeRec := image.Rect(0, 0, size, size)

	rect1 := image.Rect(0, 0, size, size)
	rect2 := image.Rect(size, 0, size*2, size)
	rect3 := image.Rect(0, size, size, size*2)
	rect4 := image.Rect(size, size, size*2, size*2)

	img1 := image.NewRGBA(sizeRec)
	img2 := image.NewRGBA(sizeRec)
	img3 := image.NewRGBA(sizeRec)
	img4 := image.NewRGBA(sizeRec)

	start := time.Now()
	for i := 0; i < 1000; i++ {
		draw.Draw(img1, sizeRec, img, rect1.Min, draw.Src)
		draw.Draw(img2, sizeRec, img, rect2.Min, draw.Src)
		draw.Draw(img3, sizeRec, img, rect3.Min, draw.Src)
		draw.Draw(img4, sizeRec, img, rect4.Min, draw.Src)
	}
	elapsed := time.Now().Sub(start).Seconds()

	fmt.Printf("elapsed: %fs\n", elapsed)
	fmt.Printf("avg: %fs\n", elapsed/1000)

	save(img1, "output/image1.png")
	save(img2, "output/image2.png")
	save(img3, "output/image3.png")
	save(img4, "output/image4.png")

	//outputFile, err := os.Create("test.png")
	//if err != nil {
	//	// Handle error
	//}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(w, img)

	// Don't forget to close files
	//outputFile.Close()

}

func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
