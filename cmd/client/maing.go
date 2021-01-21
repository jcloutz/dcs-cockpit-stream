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
	topLabel    *gtk.Label
	bottomLabel *gtk.Label
	img         *gtk.Image
	nSets       = 1
	buffer      *gdk.Pixbuf
)

func main() {
	gtk.Init(nil)

	win2, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win2.Move(500, 0)
	tk := time.NewTicker(30 * time.Millisecond)
	done := make(chan bool)

	win2.Connect("destroy", func() {
		done <- true
		gtk.MainQuit()
	})

	win2.Add(windowWidget2())

	go func() {
		for {
			select {
			case <-tk.C:
				go setImage()
			case <-done:
				fmt.Println("Closing")
				tk.Stop()
				close(done)
				return
			}
		}
	}()

	win2.ShowAll()

	gtk.Main()
}

var elapsed int64 = 0
var frames int64 = 0

func setImage() {
	start := time.Now()
	screen, _ := screenshot.Capture(0, 0, 500, 500)
	buffer, _ = gdk.PixbufNewFromData(screen.Pix, gdk.COLORSPACE_RGB, true, 8, 500, 500, 2000)

	_, err := glib.IdleAdd(img.SetFromPixbuf, buffer)
	if err != nil {
		log.Fatal("IdleAdd() image failed:", err)
	}
	end := time.Now().Sub(start).Milliseconds()
	elapsed += end
	frames++

	if frames%100 == 0 {
		fmt.Printf("avg frame time: %dms\n", elapsed/frames)
	}
}

func windowWidget2() *gtk.Widget {
	img, _ = gtk.ImageNew()

	setImage()

	img.SetHExpand(true)
	img.SetVExpand(true)

	return img.ToWidget()
}
