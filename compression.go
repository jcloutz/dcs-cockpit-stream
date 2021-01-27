package cockpit_stream

import (
	"github.com/pierrec/lz4"
)

func NewCompressionBuffer(length int) *Buffer {
	maxSize := lz4.CompressBlockBound(length)
	return &Buffer{
		Bytes: make([]byte, maxSize),
		Size:  maxSize,
	}
}

func CompressBuffer(buffer *Buffer) (int, error) {
	size, err := lz4.CompressBlock(buffer.Bytes, buffer.Bytes, nil)
	if err != nil {
		return 0, nil
	}
	buffer.Size = size

	return size, nil
}

func DecompressBuffer(buffer *Buffer, size int) (*Buffer, int, error) {
	decodeBuffer := NewBufferWithSize(size)

	length, err := lz4.UncompressBlock(buffer.Bytes, decodeBuffer.Bytes)
	if err != nil {
		return nil, 0, err
	}

	return decodeBuffer, length, err
}
