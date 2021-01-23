package main

import (
	"github.com/kbinani/screenshot"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/jcloutz/cockpit_stream"
)

func main() {
	srcImage, err := cockpit_stream.Open("test_data/capture_frame_2.png")
	if err != nil {
		log.Fatal(err)
	}
	// create white image
	destSize := image.Rect(0, 0, 150, 50)
	destImage := image.NewRGBA(destSize)
	{
		white := color.RGBA{255, 255, 255, 255}
		draw.Draw(destImage, destSize, &image.Uniform{white}, image.Point{X: 0, Y: 0}, draw.Src)
	}

	//srcRect := image.Rect(0, 0, 50, 50)
	// green
	destPt1 := image.Point{X: 150, Y: 50}
	destRect1 := image.Rect(100, 0, 150, 50)
	draw.Draw(destImage, destRect1, srcImage, destPt1, draw.Src)

	// red
	destPt2 := image.Point{X: 50, Y: 0}
	destRect2 := image.Rect(0, 0, 50, 50)
	draw.Draw(destImage, destRect2, srcImage, destPt2, draw.Src)

	cockpit_stream.Save(destImage, "output/test_out.png")

	// capture test
	desktop, _ := screenshot.Capture(0, 1440, 1680, 1050)

	centerRect := image.Rect(100, 100, 500, 500)
	centerPoint := image.Point{X: 1000, Y: 0}
	centerImg := image.NewRGBA(centerRect)

	draw.Draw(centerImg, centerRect, desktop, centerPoint, draw.Src)

	cockpit_stream.Save(centerImg, "output/center.png")
	cockpit_stream.Save(desktop, "output/desktop.png")
}
