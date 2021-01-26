package cockpit_stream

import (
	"image"
	"reflect"
	"testing"
)

func TestNewOffsetImage(t *testing.T) {
	type args struct {
		img    *image.RGBA
		offset image.Point
	}
	tests := []struct {
		name string
		args args
		want *OffsetImage
	}{
		{
			name: "can create offset image",
			args: args{
				img:    image.NewRGBA(image.Rect(0, 0, 10, 10)),
				offset: image.Point{},
			},
			want: &OffsetImage{
				RGBA:   image.NewRGBA(image.Rect(0, 0, 10, 10)),
				offset: image.Point{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOffsetImage(tt.args.img, tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOffsetImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetImage_CalcOffset(t *testing.T) {
	type fields struct {
		RGBA   *image.RGBA
		offset image.Point
	}
	type args struct {
		point image.Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   image.Point
	}{
		{
			name: "calc offset with 0,0 offset",
			fields: fields{
				RGBA:   &image.RGBA{},
				offset: image.Point{0, 0},
			},
			args: args{
				point: image.Point{X: 10, Y: 10},
			},
			want: image.Point{
				X: 10,
				Y: 10,
			},
		},
		{
			name: "calc offset with -10,-10 offset",
			fields: fields{
				RGBA:   &image.RGBA{},
				offset: image.Point{-10, -10},
			},
			args: args{
				point: image.Point{X: 10, Y: 10},
			},
			want: image.Point{
				X: 20,
				Y: 20,
			},
		},
		{
			name: "calc offset with 10,10 offset",
			fields: fields{
				RGBA:   &image.RGBA{},
				offset: image.Point{10, 10},
			},
			args: args{
				point: image.Point{X: 10, Y: 10},
			},
			want: image.Point{
				X: 0,
				Y: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &OffsetImage{
				RGBA:   tt.fields.RGBA,
				offset: tt.fields.offset,
			}
			if got := i.CalcOffset(tt.args.point); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetImage_Slice(t *testing.T) {
	type fields struct {
		RGBA   *image.RGBA
		offset image.Point
	}
	type args struct {
		dst      *image.RGBA
		viewport *Viewport
		expected *image.RGBA
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "slice image with 0, 0 offset",
			fields: fields{
				RGBA:   loadPng(t, CaptureFrame2),
				offset: image.Point{X: 0, Y: 0},
			},
			args: args{
				dst:      image.NewRGBA(image.Rect(0, 0, 50, 50)),
				viewport: NewViewport("test 1", 150, 50, 50, 50),
				expected: loadPng(t, ViewportGreen),
			},
		},
		{
			name: "slice image with -50, -50 offset",
			fields: fields{
				RGBA:   loadPng(t, CaptureFrame2),
				offset: image.Point{X: -50, Y: -50},
			},
			args: args{
				dst:      image.NewRGBA(image.Rect(0, 0, 50, 50)),
				viewport: NewViewport("test 2", 100, 0, 50, 50),
				expected: loadPng(t, ViewportGreen),
			},
		},
		{
			name: "slice image with 50, 50 offset",
			fields: fields{
				RGBA:   loadPng(t, CaptureFrame2),
				offset: image.Point{X: 50, Y: 50},
			},
			args: args{
				dst:      image.NewRGBA(image.Rect(0, 0, 50, 50)),
				viewport: NewViewport("test 2", 200, 100, 50, 50),
				expected: loadPng(t, ViewportGreen),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &OffsetImage{
				RGBA:   tt.fields.RGBA,
				offset: tt.fields.offset,
			}

			i.Slice(tt.args.dst, tt.args.viewport)

			assertRgbaEqual(t, tt.args.expected, tt.args.dst)
		})
	}
}
