package fcoe

import (
	"encoding/binary"
	"hash/crc32"
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
	if len(b) < 22 {
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

func (f *Frame) Checksum() uint32 {
	h := crc32.ChecksumIEEE(f.Payload)
	// Change from big-endian to host encoding
	b := [4]byte{
		uint8(h >> 0),
		uint8(h >> 8),
		uint8(h >> 16),
		uint8(h >> 24),
	}
	return binary.BigEndian.Uint32(b[:])
}

func (f *Frame) MarshalBinary() ([]byte, error) {
	b := make([]byte, f.length())
	err := f.read(b)
	return b, err
}

func (f *Frame) read(b []byte) error {
	// Version
	b[0] = byte(f.Version)

	// TODO(bluecmd): This can be made faster, but the map is so small so
	// should be fine
	for k, v := range sofMap {
		if v == f.SOF {
			b[13] = k
			break
		}
	}
	copy(b[14:], f.Payload)

	binary.BigEndian.PutUint32(b[len(b)-8:len(b)-4], f.CRC32)
	for k, v := range eofMap {
		if v == f.EOF {
			b[len(b)-4] = k
			break
		}
	}

	return nil
}

func (f *Frame) length() int {
	return len(f.Payload) + 22
}
