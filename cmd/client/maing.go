package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kbinani/screenshot"
)

var (
	leftImg     *gtk.Image
	rightImg    *gtk.Image
	leftBuffer  *gdk.Pixbuf
	rightBuffer *gdk.Pixbuf
)

func main() {
	gtk.Init(nil)

	tk := time.NewTicker((1000 / 30) * time.Millisecond)
	leftDone := make(chan bool)
	//rightDone := make(chan bool)

	win1, err := gtk.WindowNew(gtk.WINDOW_POPUP)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win1.SetResizable(false)
	win1.SetDecorated(false)
	win1.Move(1150, 0)

	win1.Connect("destroy", func() {
		//rightDone <- true
		gtk.MainQuit()
	})

	win2, err := gtk.WindowNew(gtk.WINDOW_POPUP)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win2.SetResizable(false)
	win2.SetDecorated(false)
	win2.Move(650, 0)

	win2.Connect("destroy", func() {
		leftDone <- true
		gtk.MainQuit()
	})

	win1.Add(rightWindowWidget())
	win2.Add(leftWindowWidget())

	go func() {
		for {
			select {
			case <-tk.C:
				go setLeftImage()
				go setRightImage()
			case <-leftDone:
				fmt.Println("Closing")
				tk.Stop()
				close(leftDone)
				return
			}
		}
	}()

	win1.ShowAll()
	win2.ShowAll()

	gtk.Main()
}

var elapsed int64 = 0
var frames int64 = 0

func setRightImage() {
	start := time.Now()
	screen, _ := screenshot.Capture(500, 1440, 500, 500)
	rightBuffer, _ = gdk.PixbufNewFromData(screen.Pix, gdk.COLORSPACE_RGB, true, 8, 500, 500, 2000)

	_, err := glib.IdleAdd(rightImg.SetFromPixbuf, rightBuffer)
	if err != nil {
		log.Fatal("IdleAdd() image failed:", err)
	}
	end := time.Now().Sub(start).Milliseconds()
	elapsed += end
	frames++

	if frames%100 == 0 {
		fmt.Printf("Right: avg frame time: %dms\n", elapsed/frames)
	}
}

func setLeftImage() {
	start := time.Now()
	screen, _ := screenshot.Capture(0, 1440, 500, 500)
	leftBuffer, _ = gdk.PixbufNewFromData(screen.Pix, gdk.COLORSPACE_RGB, true, 8, 500, 500, 2000)

	_, err := glib.IdleAdd(leftImg.SetFromPixbuf, leftBuffer)
	if err != nil {
		log.Fatal("IdleAdd() image failed:", err)
	}
	end := time.Now().Sub(start).Milliseconds()
	elapsed += end
	frames++

	if frames%100 == 0 {
		fmt.Printf("Left: avg frame time: %dms\n", elapsed/frames)
	}
}

func leftWindowWidget() *gtk.Widget {
	leftImg, _ = gtk.ImageNew()

	setLeftImage()

	leftImg.SetHExpand(true)
	leftImg.SetVExpand(true)

	return leftImg.ToWidget()
}

func rightWindowWidget() *gtk.Widget {
	rightImg, _ = gtk.ImageNew()

	setRightImage()

	rightImg.SetHExpand(true)
	rightImg.SetVExpand(true)

	return rightImg.ToWidget()
}
