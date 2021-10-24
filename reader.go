package bitslice

import "encoding/binary"

type BitReader struct {
	bs        BitSlice
	position  uint
	byteOrder binary.ByteOrder
}

func NewBitReader(bs BitSlice, bo binary.ByteOrder) *BitReader {
	return &BitReader{
		bs:        bs,
		byteOrder: bo,
	}
}

func (br *BitReader) ReadUint64(length uint) uint64 {
	var out uint64
	for i, b := range br.bs.Slice[br.position : br.position+length] {
		if !b {
			continue
		}
		out |= 1 << (int(length) - i - 1)
	}
	return out
}
