package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jcloutz/cockpit_stream"
	"github.com/jcloutz/cockpit_stream/config"
)

func main() {
	cfg, err := config.New("config.json")
	if err != nil {
		log.Fatal(err)
	}

	screenCapture := cockpit_stream.NewScreenCapture()

	viewportManager := cockpit_stream.NewServerViewportManager(screenCapture, cfg.FramesPerSecond)
	for _, vp := range cfg.Viewports {
		viewportManager.AddNewViewport(
			vp.ID,
			vp.PosX,
			vp.PosY,
			vp.Width,
			vp.Height,
		)
	}
	//viewportManager.Run()

	listener := make(chan *cockpit_stream.ListenerResult)
	sm := cockpit_stream.NewHostScreenManager(func(smConfig *cockpit_stream.HostScreenManagerConfig) {
		smConfig.TargetCaptureFps = cfg.FramesPerSecond
		smConfig.ScreenCapper = screenCapture
		smConfig.ViewportManager = viewportManager
	})

	sm.OnCaptureUpdate(listener)
	sm.Start()

	go func() {
		start := time.Now()
		count := 0
		size := 500
		sizeRec := image.Rect(0, 0, size, size)
		leftImg := image.NewRGBA(sizeRec)
		centerImg := image.NewRGBA(sizeRec)
		rightImg := image.NewRGBA(sizeRec)
		for {
			select {
			case res := <-listener:
				count++
				//go func() {
				res.Slicer.Slice("left", leftImg, sizeRec, image.Point{X: 0, Y: 0})
				//}()
				//go func() {
				res.Slicer.Slice("right", centerImg, sizeRec, image.Point{X: 0, Y: 0})
				//}()
				//go func() {
				res.Slicer.Slice("center", rightImg, sizeRec, image.Point{X: 0, Y: 0})
				//}()

				if count%100 == 0 {
					elapsed := time.Now().Sub(res.T)
					elapsedTotal := time.Now().Sub(start)

					fmt.Println("-----------------")
					fmt.Printf("exec time: %dms, %fs\n", elapsed.Milliseconds(), elapsed.Seconds())
					fmt.Printf("avg cap time: %dms, %fs, FPS: %d\n", elapsedTotal.Milliseconds()/int64(count), elapsedTotal.Seconds()/float64(count), int64(count)/(elapsedTotal.Milliseconds()/1000))
					fmt.Printf("milliseconds: %d, Frames: %d\n", elapsedTotal.Milliseconds(), count)
				}
			}
		}
	}()

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Shutting down")
		close(done)
	}()

	<-done
}
