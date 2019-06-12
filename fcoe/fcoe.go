package fcoe

import (
	"encoding/binary"
	"io"

	fc "github.com/bluecmd/fibrechannel"
)

const (
	EtherType uint16 = 0x8906
)

var (
	eofMap = map[uint8]fc.EOF{
		0x41: fc.EOFn,
		0x42: fc.EOFt,
		0x44: fc.EOFrt,
		0x46: fc.EOFdt,
		0x49: fc.EOFni,
		0x4E: fc.EOFdti,
		0x4F: fc.EOFrti,
		0x50: fc.EOFa,
	}
	sofMap = map[uint8]fc.SOF{
		0x28: fc.SOFf,
		0x29: fc.SOFi4,
		0x2D: fc.SOFi2,
		0x2E: fc.SOFi3,
		0x31: fc.SOFn4,
		0x35: fc.SOFn2,
		0x36: fc.SOFn3,
		0x39: fc.SOFc4,
	}
)

type Frame struct {
	Version int
	SOF     fc.SOF

	Payload []byte

	EOF fc.EOF

	CRC32 uint32
}

func (f *Frame) UnmarshalBinary(b []byte) error {
	var ok bool
	if len(b) < 18 {
		return io.ErrUnexpectedEOF
	}

	f.Version = int(b[0])
	f.SOF, ok = sofMap[b[13]]
	if !ok {
		return fc.ErrInvalidSOF
	}
	f.Payload = b[14 : len(b)-8]
	f.EOF, ok = eofMap[b[len(b)-4]]
	if !ok {
		return fc.ErrInvalidEOF
	}
	f.CRC32 = binary.BigEndian.Uint32(b[len(b)-8 : len(b)-4])

	return nil
}
