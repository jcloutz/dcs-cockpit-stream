package cockpit_stream

import (
	"github.com/pierrec/lz4"
)

type payloadDecoder struct {
	Bytes           *Buffer
	controlBytes    []byte
	compressionType PayloadType
	size            uint32
	width           uint32
	height          uint32
	posX            uint32
	posY            uint32
}

func NewPayloadDecoder(buffer []uint8) (*payloadDecoder, error) {
	size := decodeUintAtOffset(buffer[RangeSizeStart:RangeSizeEnd])

	newBuffer := NewBufferWithSize(int(size))
	decodeSize, err := lz4.UncompressBlock(buffer[RangeDataStart:], newBuffer.Bytes)
	if err != nil {
		return nil, err
	}
	newBuffer.Size = decodeSize

	return &payloadDecoder{
		Bytes:           newBuffer,
		compressionType: PayloadType(buffer[RangeTypeStart]),
		size:            size,
		width:           decodeUintAtOffset(buffer[RangeWidthStart:RangeWidthEnd]),
		height:          decodeUintAtOffset(buffer[RangeHeightStart:RangeHeightEnd]),
		posX:            decodeUintAtOffset(buffer[RangePosXStart:RangePosXEnd]),
		posY:            decodeUintAtOffset(buffer[RangePosYStart:RangePosYEnd]),
	}, nil
}

func (p *payloadDecoder) GetSize() int {
	return int(p.size)
}

func (p *payloadDecoder) GetWidth() int {
	return int(p.width)
}

func (p *payloadDecoder) GetHeight() int {
	return int(p.height)
}

func (p *payloadDecoder) GetPosX() int {
	return int(p.posX)
}

func (p *payloadDecoder) GetPosY() int {
	return int(p.posY)
}
