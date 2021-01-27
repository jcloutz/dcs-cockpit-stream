package cockpit_stream

import (
	"testing"
)

//func TestCalculateBitmask(t *testing.T) {
//	type args struct {
//		prev   *image.RGBA
//		next   *image.RGBA
//		buffer *Buffer
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := CalculateBitmask(tt.args.prev, tt.args.next, tt.args.buffer); (err != nil) != tt.wantErr {
//				t.Errorf("CalculateBitmask() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

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
