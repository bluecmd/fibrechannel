package fibrechannel

import (
	//"encoding/binary"
	"errors"
	"io"
)

type SOF int
type EOF int

const (
	EOFn   = 0x1
	EOFt   = 0x2
	EOFrt  = 0x3
	EOFdt  = 0x4
	EOFni  = 0x5
	EOFdti = 0x6
	EOFrti = 0x7
	EOFa   = 0x8

	SOFf  = 0x1
	SOFi4 = 0x2
	SOFi2 = 0x3
	SOFi3 = 0x4
	SOFn4 = 0x5
	SOFn2 = 0x6
	SOFn3 = 0x7
	SOFc4 = 0x8
)

var (
	ErrInvalidEOF = errors.New("invalid EOF")
	ErrInvalidSOF = errors.New("invalid SOF")
)

type Frame struct {
}

func (f *Frame) UnmarshalBinary(sof SOF, b []byte, eof EOF) error {
	if len(b) < 28 {
		return io.ErrUnexpectedEOF
	}

	return nil
}
