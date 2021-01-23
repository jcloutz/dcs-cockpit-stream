package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jcloutz/cockpit_stream"
	"github.com/jcloutz/cockpit_stream/config"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
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

	listener := make(chan *cockpit_stream.ScreenCaptureResult)
	sm := cockpit_stream.NewHostScreenManager(func(smConfig *cockpit_stream.ScreenCaptureControllerConfig) {
		smConfig.TargetCaptureFps = cfg.FramesPerSecond
	})

	sm.AddListener(listener)
	sm.SetBounds(viewportManager.GetBounds())
	sm.Start()
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	go func() {
		start := time.Now()
		count := 0
		//size := 500
		//sizeRec := image.Rect(0, 0, size, size)
		//leftImg := image.NewRGBA(sizeRec)
		//centerImg := image.NewRGBA(sizeRec)
		//rightImg := image.NewRGBA(sizeRec)
		for {
			select {
			case res := <-listener:
				count++
				//go func() {
				//	left, err := res.Viewports.Get("left")
				//	if err != nil {
				//		log.Println(err)
				//	}
				//	left.Slice(leftImg, sizeRec, image.Point{X: 0, Y: 0})
				//}()
				//go func() {
				//	right, err := res.Viewports.Get("right")
				//	if err != nil {
				//		log.Println(err)
				//	}
				//	right.Slice(centerImg, sizeRec, image.Point{X: 0, Y: 0})
				//}()
				//go func() {
				//	center, err := res.Viewports.Get("center")
				//	if err != nil {
				//		log.Println(err)
				//	}
				//	center.Slice(rightImg, sizeRec, image.Point{X: 0, Y: 0})
				//}()

				//close(quit)
				//done <- true
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

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Shutting down")
		close(done)
	}()

	<-done
}
