package main

import (
	"fmt"
	"math"

	"github.com/kbinani/screenshot"

	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"sync"
	"time"

	metrics2 "github.com/armon/go-metrics"
	"github.com/jcloutz/cockpit_stream"
)

type foo struct {
	value *int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	memSink := metrics2.NewInmemSink(5*time.Second, 10*time.Second)

	start := time.Now()
	//_, err := metrics.NewGlobal(metrics.DefaultConfig("service-name"), inm)
	//if err != nil {
	//	log.Fatal(err)
	//}

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

	//print := func() {
	//	fmt.Println("--------------")
	//	fmt.Println("sample count", timer.Count())
	//	fmt.Println("mean", int64(timer.Mean()/float64(time.Millisecond)))
	//	fmt.Println("max", timer.Max()/int64(time.Millisecond))
	//	fmt.Println("min", timer.Min()/int64(time.Millisecond))
	//}

	for i := 0; i < 500; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			funcStart := time.Now()
			timeout := rand.Intn(29) + 1
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			memSink.AddSample([]string{fmt.Sprintf("exec")}, float32(time.Since(funcStart).Milliseconds()))
			wg.Done()
		}()

		wg.Wait()
	}

	elapsed := time.Since(start)
	fmt.Println("elapsed", elapsed.Milliseconds())
	dt := memSink.Data()
	fmt.Println("len", len(dt))
	for _, el := range dt {
		for name, sample := range el.Samples {
			fmt.Println(fmt.Sprintf("---- %s ----", name))
			fmt.Println("elapsed", elapsed.Milliseconds())

			fmt.Println("sample count", sample.Count)
			fmt.Println("mean", math.Round(sample.Sum/float64(sample.Count)))
			fmt.Println("max", sample.Max /*float64(time.Millisecond)*/)
			fmt.Println("min", sample.Min /*float64(time.Millisecond)*/)
		}
		fmt.Println("------------------------------------------------")
		fmt.Println("----------------- NEW DATASET ------------------")
		fmt.Println("------------------------------------------------")
	}
}
