package cockpit_stream

import (
	"image"
	"testing"
)

func applyMaskToImage(t *testing.T, dest *image.RGBA, mask *Buffer) {
	t.Helper()
	for i := 0; i < len(mask.Bytes); i++ {
		dest.Pix[i] = dest.Pix[i] ^ mask.Bytes[i]
	}
}

func TestCalculateBitmask(t *testing.T) {
	type args struct {
		prev string
		next string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected string
	}{
		{
			name: "frame_1 -> frame_2",
			args: args{
				prev: "test_data/capture_frame_1.png",
				next: "test_data/capture_frame_2.png",
			},
			wantErr:  false,
			expected: "test_data/capture_frame_2.png",
		},
		{
			name: "frame_2 -> frame_3",
			args: args{
				prev: "test_data/capture_frame_2.png",
				next: "test_data/capture_frame_3.png",
			},
			wantErr:  false,
			expected: "test_data/capture_frame_3.png",
		},
		{
			name: "frame_3 -> frame_4",
			args: args{
				prev: "test_data/capture_frame_3.png",
				next: "test_data/capture_frame_4.png",
			},
			wantErr:  false,
			expected: "test_data/capture_frame_4.png",
		},
	}
	for _, tt := range tests {
		prev := loadPng(t, tt.args.prev)
		next := loadPng(t, tt.args.next)

		destImg := loadPng(t, tt.args.prev)

		buffer := NewBufferWithSize(len(prev.Pix))

		t.Run(tt.name, func(t *testing.T) {
			if err := CalculateBitmask(prev, next, buffer); (err != nil) != tt.wantErr {
				t.Errorf("CalculateBitmask() error = %v, wantErr %v", err, tt.wantErr)
			}

			applyMaskToImage(t, destImg, buffer)

			assertRgbaEqual(t, loadPng(t, tt.expected), destImg)

		})
	}
}

func TestCalculateBitmaskParallel(t *testing.T) {
	type args struct {
		prev string
		next string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected string
	}{
		{
			name: "frame_1 -> frame_2",
			args: args{
				prev: "test_data/capture_frame_1.png",
				next: "test_data/capture_frame_2.png",
			},
			wantErr:  false,
			expected: "test_data/capture_frame_2.png",
		},
		{
			name: "frame_2 -> frame_3",
			args: args{
				prev: "test_data/capture_frame_2.png",
				next: "test_data/capture_frame_3.png",
			},
			wantErr:  false,
			expected: "test_data/capture_frame_3.png",
		},
		{
			name: "frame_3 -> frame_4",
			args: args{
				prev: "test_data/capture_frame_3.png",
				next: "test_data/capture_frame_4.png",
			},
			wantErr:  false,
			expected: "test_data/capture_frame_4.png",
		},
	}
	for _, tt := range tests {
		prev := loadPng(t, tt.args.prev)
		next := loadPng(t, tt.args.next)

		destImg := loadPng(t, tt.args.prev)

		buffer := NewBufferWithSize(len(prev.Pix))

		t.Run(tt.name, func(t *testing.T) {
			if err := CalculateBitmaskProcPerRow(prev, next, buffer); (err != nil) != tt.wantErr {
				t.Errorf("CalculateBitmaskProcPerRow() error = %v, wantErr %v", err, tt.wantErr)
			}

			applyMaskToImage(t, destImg, buffer)

			assertRgbaEqual(t, loadPng(t, tt.expected), destImg)

		})
	}
}
func BenchmarkCalculateBitmask(b *testing.B) {
	prevImg := loadPng(b, "test_data/xor_prev.png")
	nextImg := loadPng(b, "test_data/xor_next.png")
	buffer := NewBuffer(prevImg.Rect.Dx(), prevImg.Rect.Dy())
	b.Run("calculate bitmask", func(b *testing.B) {
		if err := CalculateBitmask(prevImg, nextImg, buffer); err != nil {
			b.Fatal(err)
		}
	})
}

func BenchmarkCalculateBitmaskParallel(b *testing.B) {
	prevImg := loadPng(b, "test_data/xor_prev.png")
	nextImg := loadPng(b, "test_data/xor_next.png")
	buffer := NewBuffer(prevImg.Rect.Dx(), prevImg.Rect.Dy())
	b.Run("calculate bitmask", func(b *testing.B) {
		if err := CalculateBitmaskProcPerRow(prevImg, nextImg, buffer); err != nil {
			b.Fatal(err)
		}
	})
}
func BenchmarkCalculateBitmaskParallel4(b *testing.B) {
	prevImg := loadPng(b, "test_data/xor_prev.png")
	nextImg := loadPng(b, "test_data/xor_next.png")
	buffer := NewBuffer(prevImg.Rect.Dx(), prevImg.Rect.Dy())
	b.Run("calculate bitmask", func(b *testing.B) {
		if err := CalculateBitmask4ProcPerRow(prevImg, nextImg, buffer); err != nil {
			b.Fatal(err)
		}
	})
}
