package cockpit_stream

type Buffer struct {
	Bytes []uint8
	Size  int
}

func NewBuffer(width int, height int) *Buffer {
	return &Buffer{
		Bytes: make([]uint8, width*height*4),
		Size:  width * height * 4,
	}
}

func NewBufferWithSize(size int) *Buffer {
	return &Buffer{
		Bytes: make([]uint8, size),
		Size:  size,
	}
}
