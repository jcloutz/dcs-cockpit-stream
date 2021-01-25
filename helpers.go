package cockpit_stream

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func SavePng(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func OpenPng(filePath string) (*image.RGBA, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(srcImg.Bounds())

	draw.Draw(img, img.Bounds(), srcImg, image.Point{
		X: 0,
		Y: 0,
	}, draw.Src)

	return img, nil
}
