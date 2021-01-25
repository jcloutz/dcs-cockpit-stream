package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jcloutz/cockpit_stream"
	"github.com/jcloutz/cockpit_stream/config"
	"github.com/jcloutz/cockpit_stream/metrics"
)

func main() {
	metricService := metrics.New()
	cfg, err := config.New("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// create viewports
	viewports := cockpit_stream.NewViewportContainer()
	for id, vp := range cfg.Viewports {
		viewports.Add(id, cockpit_stream.NewViewport(id, vp.PosX, vp.PosY, vp.Width, vp.Height))
	}

	// Create handlers
	var handlers []*cockpit_stream.ScreenCaptureHandler
	for id, client := range cfg.Clients {
		handler := cockpit_stream.NewViewportStreamHandler(id, viewports, metricService)
		handler.EnableOutput("output")

		for _, vpCfg := range client.Viewports {
			handler.RegisterViewport(vpCfg.ID, vpCfg.DisplayX, vpCfg.DisplayY)
		}

		handlers = append(handlers, handler)
	}

	// create capture controller
	viewportCaptureController := cockpit_stream.NewViewportCaptureController(func(smConfig *cockpit_stream.DesktopCaptureControllerConfig) {
		smConfig.TargetCaptureFps = cfg.FramesPerSecond
		smConfig.Metrics = metricService
	})
	viewportCaptureController.SetBounds(viewports.GetBounds())

	for _, handler := range handlers {
		viewportCaptureController.AddListener(handler)
	}

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	loggingListener := &cockpit_stream.CallbackCaptureHandler{}

	count := 0
	loggingListener.OnReceive(func(res *cockpit_stream.CaptureResult) {
		count++

		if count%100 == 0 {
			fmt.Println("-----------------")
			if capCtx, err := res.GetCaptureContext(); err == nil {
				data := capCtx.Metric.Data()
				fmt.Printf("captures frames: %.2f\n", data.GetCount(metrics.MetricFrameCounter).Sum)
				fmt.Printf("total screen cap time: %.2fs\n", data.GetSample(metrics.MetricSampleCaptureController).Sum/float64(time.Second))
				maxFramerate := data.GetCount(metrics.MetricFrameCounter).Sum / (data.GetSample(metrics.MetricSampleCaptureController).Sum / float64(time.Second))
				fmt.Printf("max possible framerate: %.2ffps\n", maxFramerate)

				fmt.Printf("screen cap rate: %.2f\n", data.GetSample(metrics.MetricSampleCaptureController).Rate/float64(time.Millisecond))
				fmt.Printf("avg screen cap time: %.2fms\n", data.GetSample(metrics.MetricSampleCaptureController).Mean()/float64(time.Millisecond))
				fmt.Printf("avg handle time [client1]: %.2fms\n", data.GetSampleForClient(metrics.MetricSampleViewportHandler, "client1").Mean()/float64(time.Millisecond))
				fmt.Printf("avg handle time [client2]: %.2fms\n", data.GetSampleForClient(metrics.MetricSampleViewportHandler, "client2").Mean()/float64(time.Millisecond))
				fmt.Printf("avg pipeline time [client1]: %.2fms\n", data.GetSampleForClient(metrics.MetricPipelineExecutionTime, "client1").Mean()/float64(time.Millisecond))
				fmt.Printf("avg pipeline time [client2]: %.2fms\n", data.GetSampleForClient(metrics.MetricPipelineExecutionTime, "client2").Mean()/float64(time.Millisecond))
			} else {
				log.Println("unable to get context")
			}
		}

	})
	viewportCaptureController.AddListener(loggingListener)

	err = viewportCaptureController.RunOnce()
	if err != nil {
		log.Fatal(err)
	}

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("Shutting down")
		close(done)
	}()

	<-done
}
