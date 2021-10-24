package bitslice

import "encoding/binary"

type BitWriter struct {
	bs        *BitSlice
	byteOrder binary.ByteOrder
}

func NewBitWriter(bo binary.ByteOrder) *BitWriter {
	return &BitWriter{
		bs:        NewEmptyBitSlice(binary.BigEndian),
		byteOrder: bo,
	}
}

func (bw *BitWriter) WriteUint64(val uint64, length uint) {
	out := make([]bool, length)
	for i := range out {
		if val&1 == 1 {
			out[int(length)-1-i] = true
		}
		val = val >> 1
	}
	if bw.byteOrder == binary.BigEndian {
		bw.bs.Slice = append(bw.bs.Slice, out...)
		return
	}
	var newSlice []bool
	newSlice = append(newSlice, out...)
	newSlice = append(newSlice, bw.bs.Slice...)
}
