package cockpit_stream

import (
	"encoding/binary"

	"github.com/pierrec/lz4"
)

type PayloadType uint

const (
	PayloadTypeImage PayloadType = 1
	PayloadTypeMask  PayloadType = 2
)

type PayloadRangeOffset int

const (
	RangeTypeStart   PayloadRangeOffset = 0
	RangeSizeStart   PayloadRangeOffset = 1
	RangeSizeEnd     PayloadRangeOffset = 5
	RangeWidthStart  PayloadRangeOffset = 5
	RangeWidthEnd    PayloadRangeOffset = 9
	RangeHeightStart PayloadRangeOffset = 9
	RangeHeightEnd   PayloadRangeOffset = 13
	RangePosXStart   PayloadRangeOffset = 13
	RangePosXEnd     PayloadRangeOffset = 17
	RangePosYStart   PayloadRangeOffset = 17
	RangePosYEnd     PayloadRangeOffset = 21
	RangeDataStart   PayloadRangeOffset = 21
)

type payloadEncoder struct {
	Bytes           []byte
	controlBytes    []byte
	compressionType PayloadType
	size            uint32
	width           uint32
	height          uint32
	posX            uint32
	posY            uint32
}

func NewPayloadEncoder() *payloadEncoder {
	return &payloadEncoder{
		controlBytes:    make([]byte, 21),
		Bytes:           []byte{},
		compressionType: 0,
		size:            0,
		width:           0,
		height:          0,
		posX:            0,
		posY:            0,
	}
}

//func DecodeCompressedMask(buffer []uint8) {
//	cType := buffer[RangeTypeStart:RangeTypeEnd]
//	size := decodeUintAtOffset(buffer[RangeSizeStart:RangeSizeEnd])
//	width := decodeUintAtOffset(buffer[RangeWidthStart:RangeWidthEnd])
//	height := decodeUintAtOffset(buffer[RangeHeightStart:RangeHeightEnd])
//	posX := decodeUintAtOffset(buffer[RangePosXStart:RangePosXEnd])
//	posY := decodeUintAtOffset(buffer[RangePosYStart:RangePosYEnd])
//
//	fmt.Println(cType, size, width, height, posX, posY)
//}

func (c *payloadEncoder) encodeUintAtOffset(i uint32, start PayloadRangeOffset, end PayloadRangeOffset) {
	binary.BigEndian.PutUint32(c.controlBytes[start:end], i)
}

func decodeUintAtOffset(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func (c *payloadEncoder) SetType(t PayloadType) *payloadEncoder {
	c.compressionType = t
	c.controlBytes[0] = byte(t)

	return c
}

func (c *payloadEncoder) setSize(size uint32) *payloadEncoder {
	c.size = size
	c.encodeUintAtOffset(size, RangeSizeStart, RangeSizeEnd)

	return c
}

func (c *payloadEncoder) SetWidth(width uint32) *payloadEncoder {
	c.encodeUintAtOffset(width, RangeWidthStart, RangeWidthEnd)

	return c
}
func (c *payloadEncoder) SetHeight(height uint32) *payloadEncoder {
	c.encodeUintAtOffset(height, RangeHeightStart, RangeHeightEnd)

	return c
}

func (c *payloadEncoder) SetPosX(posX uint32) *payloadEncoder {
	c.posX = posX
	c.encodeUintAtOffset(posX, RangePosXStart, RangePosXEnd)

	return c
}
func (c *payloadEncoder) SetPosY(posY uint32) *payloadEncoder {
	c.posY = posY
	c.encodeUintAtOffset(posY, RangePosYStart, RangePosYEnd)

	return c
}

func (c *payloadEncoder) SetBytes(buffer *Buffer) error {
	_, err := lz4.CompressBlock(buffer.Bytes, buffer.Bytes, nil)
	if err != nil {
		return err
	}
	c.setSize(uint32(buffer.Size))
	c.Bytes = append(c.controlBytes, buffer.Bytes[:buffer.Size]...)

	return nil
}

func (c *payloadEncoder) GetSize() PayloadType {
	return c.compressionType
}

func (c *payloadEncoder) GetWidth() uint32 {
	return c.width
}

func (c *payloadEncoder) GetHeight() uint32 {
	return c.height
}

func (c *payloadEncoder) GetStride() uint32 {
	return c.width * 4
}

func (c *payloadEncoder) GetPosX() uint32 {
	return c.posX
}

func (c *payloadEncoder) GetPosY() uint32 {
	return c.posY
}
