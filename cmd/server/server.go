package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jcloutz/cockpit_stream"
	"github.com/jcloutz/cockpit_stream/config"
)

func main() {
	cfg, err := config.New("config.json")
	if err != nil {
		log.Fatal(err)
	}

	screenCapture := cockpit_stream.NewScreenCapture()
	viewportManager := cockpit_stream.NewServerViewportManager(screenCapture)
	for _, vp := range cfg.Viewports {
		viewportManager.AddNewViewport(
			vp.ID,
			vp.PosX,
			vp.PosY,
			vp.Width,
			vp.Height,
		)
	}
	viewportManager.Run()

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
