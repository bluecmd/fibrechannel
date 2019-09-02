package swils

import (
	"encoding/binary"
	"io"

	fc "github.com/bluecmd/fibrechannel"
)

type ELP struct {
	Revision              uint8
	Flags                 uint16
	BBSCN                 uint8
	RATOV                 uint32
	EDTOV                 uint32
	Port                  fc.WWN   // RequesterPortName
	Switch                fc.WWN   // RequesterSwitchName
	ClassFParameters      [16]byte // TODO
	Class2Parameters      [4]byte  // TODO
	Class3Parameters      [4]byte  // TODO
	ISLFlowControlMode    uint16
	FlowControlParameters []byte
}

func (s *ELP) UnmarshalBinary(b []byte) error {
	if len(b) < 80 {
		return io.ErrUnexpectedEOF
	}

	s.Revision = uint8(b[0])
	s.Flags = binary.BigEndian.Uint16(b[1:3])
	s.BBSCN = uint8(b[3])
	s.RATOV = binary.BigEndian.Uint32(b[4:8])
	s.EDTOV = binary.BigEndian.Uint32(b[8:12])
	copy(s.Port[:], b[12:20])
	copy(s.Switch[:], b[20:28])

	//ClassFParameters      [16]byte
	//_                     [4]byte // Obsolete in FC-SW-5
	//Class2Parameters      [4]byte
	//Class3Parameters      [4]byte
	//_                     [20]byte // Reserved
	//ISLFlowControlMode    uint16
	//FlowControlParameters []byte

	n := binary.BigEndian.Uint16(b[78:80])
	s.FlowControlParameters = make([]byte, n)
	copy(s.FlowControlParameters, b[80:80+n])

	return nil
}
