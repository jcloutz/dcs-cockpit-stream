package cockpit_stream

import (
	"image"
	"testing"
)

func Test_Viewport_Slice(t *testing.T) {
	for _, tc := range []struct {
		name     string
		src      string
		expected string
		width    int
		height   int
		offset   image.Point
	}{
		{
			name:     "red viewport",
			src:      CaptureFrame2,
			expected: ViewportRed,
			width:    50,
			height:   50,
			offset:   image.Point{X: 50, Y: 0},
		},
		{
			name:     "green viewport",
			src:      CaptureFrame2,
			expected: ViewportGreen,
			width:    50,
			height:   50,
			offset:   image.Point{X: 150, Y: 50},
		},
		{
			name:     "blue viewport",
			src:      CaptureFrame2,
			expected: ViewportBlue,
			width:    50,
			height:   50,
			offset:   image.Point{X: 0, Y: 100},
		},
	} {
		viewport := NewViewport(tc.name, 0, 0, tc.width, tc.height)

		dest := getNewRgba(tc.width, tc.height)
		src := loadPng(t, tc.src)
		expected := loadPng(t, tc.expected)

		viewport.Slice(dest, src, tc.offset)

		assertRgbaEqual(t, dest, expected)
	}
}
