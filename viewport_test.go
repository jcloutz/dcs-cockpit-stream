package cockpit_stream

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

const (
	// 200 x 200 px
	// arranged in 4 x 4 grid with 50px blocks
	CaptureFrame1 string = "test_data/capture_frame_1.png"
	CaptureFrame2 string = "test_data/capture_frame_2.png"
	CaptureFrame3 string = "test_data/capture_frame_3.png"
	CaptureFrame4 string = "test_data/capture_frame_4.png"

	// 50 x 50 px
	// solid color
	ViewportBlue  string = "test_data/viewport_blue.png"
	ViewportGreen string = "test_data/viewport_green.png"
	ViewportRed   string = "test_data/viewport_red.png"
)

func getNewRgba(width int, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func assertRgbaEqual(t *testing.T, img1, img2 *image.RGBA) {
	t.Helper()

	if !img1.Rect.Eq(img2.Rect) {
		t.Fatalf("image rectagles do not match: %v != %v", img1.Rect, img2.Rect)
	}

	for i := 0; i < len(img1.Pix); i++ {
		if img1.Pix[i] != img2.Pix[i] {
			t.Fatal("image pixels do not match")
		}
	}
}

func loadPng(t *testing.T, path string) *image.RGBA {
	t.Helper()
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal("unable to get current working directory", err)
	}
	p := filepath.Join(workingDir, path)
	img, err := OpenPng(p)
	if err != nil {
		t.Fatalf("failed open '%s': %s", p, err.Error())
	}

	return img
}

func Test_Viewport_Slice(t *testing.T) {
	for _, tc := range []struct {
		name     string
		src      string
		expected string
		slicePos image.Point
		width    int
		height   int
	}{
		{
			name:     "red viewport",
			src:      CaptureFrame2,
			expected: ViewportRed,
			slicePos: image.Point{X: 50, Y: 0},
			width:    50,
			height:   50,
		},
		{
			name:     "green viewport",
			src:      CaptureFrame2,
			expected: ViewportGreen,
			slicePos: image.Point{X: 150, Y: 50},
			width:    50,
			height:   50,
		},
		{
			name:     "blue viewport",
			src:      CaptureFrame2,
			expected: ViewportBlue,
			slicePos: image.Point{X: 0, Y: 100},
			width:    50,
			height:   50,
		},
	} {
		viewport := NewViewportWithSlice(tc.name, 0, 0, tc.width, tc.height, tc.slicePos)

		dest := getNewRgba(tc.width, tc.height)
		src := loadPng(t, tc.src)
		expected := loadPng(t, tc.expected)

		viewport.Slice(dest, src)

		assertRgbaEqual(t, dest, expected)
	}
}
