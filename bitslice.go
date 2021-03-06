package bitslice

import (
	"bytes"
	"encoding/binary"
	"io"
)

//BitSlice slice with bits
type BitSlice struct {
	Slice     []bool
	ByteOrder binary.ByteOrder
}

//NewEmptyBitSlice returns empty BitSlice
func NewEmptyBitSlice(bo binary.ByteOrder) *BitSlice {
	return &BitSlice{
		ByteOrder: bo,
	}
}

//NewBitSliceFromAllReader returns BitSlice exported from all io.Reader
func NewBitSliceFromAllReader(r io.Reader, bo binary.ByteOrder) (*BitSlice, error) {
	bs := &BitSlice{
		ByteOrder: bo,
	}
	var b byte
	cnt := 0
	for {
		err := binary.Read(r, bo, &b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		newByte := make([]bool, 8)
		bs.Slice = append(bs.Slice, newByte...)
		for j := range [8]struct{}{} {
			bs.Slice[8*cnt+j] = (b<<7)>>7 == 1
			b = b >> 1
		}
		cnt++
	}
	return bs, nil
}

//NewBitSliceFromBool returns BitSlice exported from []bool
func NewBitSliceFromBool(b []bool, bo binary.ByteOrder) *BitSlice {
	return &BitSlice{
		b, bo,
	}
}

//NewBitSliceFromBytes returns BitSlice exported from []byte
func NewBitSliceFromBytes(b []byte, bo binary.ByteOrder) (*BitSlice, error) {
	buf := bytes.NewBuffer(b)
	return NewBitSliceFromReader(buf, bo, uint(len(b)))
}

//NewBitSliceFromReader returns BitSlice exported from io.Reader
func NewBitSliceFromReader(r io.Reader, bo binary.ByteOrder, bytesSize uint) (*BitSlice, error) {
	var inBytes = make([]byte, bytesSize)
	err := binary.Read(r, bo, &inBytes)
	if err != nil {
		return nil, err
	}
	bs := &BitSlice{
		Slice:     make([]bool, bytesSize*8),
		ByteOrder: bo,
	}

	lFunc := func(bPtr *byte, index int) {
		bs.Slice[index] = (*bPtr<<7)>>7 == 1
		*bPtr = *bPtr >> 1
	}

	bFunc := func(bPtr *byte, index int) {
		bs.Slice[index] = *bPtr>>7 == 1
		*bPtr = *bPtr << 1
	}
	nowFunc := lFunc
	if bo.String() == binary.BigEndian.String() {
		nowFunc = bFunc
	}

	for i, b := range inBytes {
		for j := range [8]struct{}{} {
			nowFunc(&b, 8*i+j)
		}
	}
	return bs, nil
}

//ToBytes returns []byte from BitSlice
func (s BitSlice) ToBytes() []byte {
	var packed []byte
	var flagTrue byte
	var flagTrueDefault byte = 1

	lFunc := func() {
		flagTrue = flagTrue << 1
	}

	bFunc := func() {
		flagTrue = flagTrue >> 1
	}
	nowFunc := lFunc
	if s.ByteOrder.String() == binary.BigEndian.String() {
		nowFunc = bFunc
		flagTrueDefault = flagTrueDefault << 7
	}
	flagTrue = flagTrueDefault

	for i, flag := range s.Slice {
		if i%8 == 0 {
			packed = append(packed, 0)
			flagTrue = flagTrueDefault
		}
		if flag {
			packed[i/8] |= flagTrue
		}
		nowFunc()
	}
	return packed
}

//ToBuffer write bytes from BitSlice to io.Writer
func (s BitSlice) ToBuffer(w io.Writer) error {
	err := binary.Write(w, s.ByteOrder, s.ToBytes())
	if err != nil {
		return err
	}
	return nil
}

//Len return length of slice
func (s BitSlice) Len() int {
	return len(s.Slice)
}

//LenBytes return length of slice
func (s BitSlice) LenBytes() int {
	lenBits := s.Len()
	lenBytes := lenBits / 8
	if lenBits%8 > 0 {
		lenBytes += 1
	}
	return lenBytes
}

//ShiftLeft returns shifted BitSlice, like << operation
func (s BitSlice) ShiftLeft(val int) BitSlice {
	if val == 0 {
		return s
	}
	if val < 0 {
		return s.ShiftRight(-val)
	}
	newSlice := s
	newSlice.Slice = make([]bool, s.Len())
	for i, bit := range s.Slice[val:] {
		newSlice.Slice[i] = bit
	}
	return newSlice
}

//ShiftRight returns shifted BitSlice, like >> operation
func (s BitSlice) ShiftRight(val int) BitSlice {
	if val == 0 {
		return s
	}
	if val < 0 {
		return s.ShiftLeft(-val)
	}
	newSlice := s
	sLen := s.Len()
	newSlice.Slice = make([]bool, sLen)
	ns := s.Slice[:sLen-val]
	nsLen := len(ns)
	for i, bit := range ns {
		newSlice.Slice[sLen-nsLen+i] = bit
	}
	return newSlice
}

//Inverse returns inverted BitSlice
func (s BitSlice) Inverse() BitSlice {
	newSlice := s
	for i, bit := range newSlice.Slice {
		newSlice.Slice[i] = !bit
	}
	return newSlice
}

//Or returns combine of 2 BitSlice, like | operation
func (s BitSlice) Or(bs BitSlice) BitSlice {
	smallSlice := bs
	bigSlice := s
	if bsLen := bs.Len(); bsLen > s.Len() {
		smallSlice = s
		bigSlice = bs
	}
	for i := range smallSlice.Slice {
		smallSlice.Slice[i] = smallSlice.Slice[i] || bigSlice.Slice[i]
	}
	smallSlice.ByteOrder = s.ByteOrder
	return smallSlice
}

//And returns combine of 2 BitSlice, like & operation
func (s BitSlice) And(bs BitSlice) BitSlice {
	sliceLen := s.Len()
	smallSlice := bs
	bigSlice := s
	if bsLen := bs.Len(); bsLen > s.Len() {
		sliceLen = bsLen
		smallSlice = s
		bigSlice = bs
	}
	newSlice := BitSlice{make([]bool, sliceLen), s.ByteOrder}
	for i := range smallSlice.Slice {
		newSlice.Slice[i] = smallSlice.Slice[i] && bigSlice.Slice[i]
	}
	return newSlice
}
