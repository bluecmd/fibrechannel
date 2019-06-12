package fibrechannel

import (
	"encoding/binary"
	"errors"
	"fmt"
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

	TypeBLS      = 0x00
	TypeELS      = 0x01
	TypeLLCSNAP  = 0x04
	TypeIP       = 0x05
	TypeFCP      = 0x08
	TypeGPP      = 0x09
	TypeSBToCU   = 0x1B // FICON / FC-SB-3: Control Unit -> Channel
	TypeSBFromCU = 0x1C // FICON / FC-SB-3: Channel -> Control Unit
	TypeFCCT     = 0x20
	TypeSWILS    = 0x22
	TypeAL       = 0x23
	TypeSNMP     = 0x24
	TypeSPINFAB  = 0xEE
	TypeDIAG     = 0xEF
)

var (
	ErrInvalidEOF = errors.New("invalid EOF")
	ErrInvalidSOF = errors.New("invalid SOF")
)

type RoutingControl byte
type ClassControl byte
type Priority byte
type Address [3]byte
type Type byte
type FrameControl uint32
type DataFieldControl byte
type SequenceID byte
type SequenceCount uint16
type ExchangeID uint16

type Frame struct {
	RCtl RoutingControl

	// Which field of these two are in use is controlled by the
	// "CS_CTL/Priority Enable" bit in F_CTL.
	CSCtl    *ClassControl
	Priority *Priority

	// Address for source/destination Nx_Ports
	// Each Nx_Port shall have a native N_Port_ID that is unique within the
	// address domain of a Fabric.  An N_Port_ID of binary zeros indicates that
	// an Nx_Port is unidentified. When a PN_Port completes Link Initialization,
	// it shall be unidentified (i.e., it shall have a single Nx_Port for which
	// the N_Port_ID is 00 00 00h).
	Destination Address
	Source      Address

	Type Type

	FCtl FrameControl

	// Sequence ID
	// The sequence count (SEQ_CNT) is a two-byte field (Word 3, Bits 15-0)
	// that shall indicate the sequential order of Data frame transmission within
	// a single Sequence or multiple consecutive Sequences for the same Exchange.
	// The SEQ_CNT of the first Data frame of the first Sequence of the Exchange
	// transmitted by either the Originator or Responder shall be binary zero.
	// The SEQ_CNT of each subsequent Data frame in the Sequence shall be
	// incremented by one.
	SeqID SequenceID

	DFCtl DataFieldControl

	// Sequence count
	SeqCount SequenceCount

	// Originator Exchange_ID
	// If the Originator is enforcing uniqueness via the OX_ID mechanism, it
	// shall set a unique value for OX_ID other than FF FFh in the first Data
	// frame of the first Sequence of an Exchange. An OX_ID of FF FFh indicates
	// that the OX_ID is unassigned and that the Originator is not enforcing
	// uniqueness via the OX_ID mechanism. If an Originator uses the unassigned
	// value of FF FFh to identify the Exchange, it shall have only one Exchange
	// (OX_ID set to FF FFh) with a given Responder.
	OXID ExchangeID

	// Responder Exchange_ID
	// An RX_ID of FF FFh shall indicate that the RX_ID is unassigned. If the
	// Responder does not assign an RX_ID other than FF FFh by the end of the
	// first Sequence, then the Responder is not enforcing uniqueness via the
	// RX_ID mechanism.
	RXID ExchangeID

	// TODO(bluecmd): Optional fields

	Payload []byte
}

func (f *Frame) UnmarshalBinary(sof SOF, b []byte, eof EOF) error {
	if len(b) < 20 {
		return io.ErrUnexpectedEOF
	}
	// FC Frames are always 4-byte aligned
	if len(b)%4 != 0 {
		return io.ErrUnexpectedEOF
	}

	if err := f.RCtl.UnmarshalBinary(b[0]); err != nil {
		return err
	}
	if err := f.Destination.UnmarshalBinary(b[1:4]); err != nil {
		return err
	}
	if err := f.Source.UnmarshalBinary(b[5:8]); err != nil {
		return err
	}
	if err := f.Type.UnmarshalBinary(b[8]); err != nil {
		return err
	}
	if err := f.FCtl.UnmarshalBinary(b[9:12]); err != nil {
		return err
	}
	if err := f.SeqID.UnmarshalBinary(b[12]); err != nil {
		return err
	}
	if err := f.DFCtl.UnmarshalBinary(b[13]); err != nil {
		return err
	}
	if err := f.SeqCount.UnmarshalBinary(b[14:16]); err != nil {
		return err
	}
	if err := f.OXID.UnmarshalBinary(b[16:18]); err != nil {
		return err
	}
	if err := f.RXID.UnmarshalBinary(b[18:20]); err != nil {
		return err
	}

	if f.FCtl.CSCtlEnabled() {
		var cc ClassControl
		if err := cc.UnmarshalBinary(b[4]); err != nil {
			return err
		}
		f.CSCtl = &cc
	} else {
		var p Priority
		if err := p.UnmarshalBinary(b[4]); err != nil {
			return err
		}
		f.Priority = &p
	}

	b = b[20:]

	if f.DFCtl.HasESP() {
		return fmt.Errorf("ESP is not implemented")
	}
	if f.DFCtl.HasNetworkHeader() {
		return fmt.Errorf("Network header is not implemented")
	}

	f.Payload = make([]byte, len(b))
	copy(f.Payload, b[:])
	return nil
}

func (s *RoutingControl) UnmarshalBinary(b byte) error {
	*s = RoutingControl(b)
	return nil
}

func (s *Address) UnmarshalBinary(b []byte) error {
	if len(b) != 3 {
		return io.ErrUnexpectedEOF
	}
	copy(s[:], b)
	return nil
}

func (s *FrameControl) UnmarshalBinary(b []byte) error {
	if len(b) != 3 {
		return io.ErrUnexpectedEOF
	}
	*s = FrameControl(binary.BigEndian.Uint32(append([]byte{0}, b...)))
	return nil
}

func (s *FrameControl) CSCtlEnabled() bool {
	return *s&0x20000 == 0

}

func (s *Type) UnmarshalBinary(b byte) error {
	*s = Type(b)
	return nil
}

func (s *ClassControl) UnmarshalBinary(b byte) error {
	*s = ClassControl(b)
	return nil
}

func (s *Priority) UnmarshalBinary(b byte) error {
	*s = Priority(b)
	return nil
}

func (s *DataFieldControl) UnmarshalBinary(b byte) error {
	*s = DataFieldControl(b)
	return nil
}

func (s *DataFieldControl) HasESP() bool {
	return *s&0x40 != 0
}

func (s *DataFieldControl) HasNetworkHeader() bool {
	return *s&0x20 != 0
}

func (s *DataFieldControl) DeviceHeaderSize() int {
	return int(*s & 0x03 << 4)
}

func (s *SequenceID) UnmarshalBinary(b byte) error {
	*s = SequenceID(b)
	return nil
}

func (s *SequenceCount) UnmarshalBinary(b []byte) error {
	if len(b) != 2 {
		return io.ErrUnexpectedEOF
	}
	*s = SequenceCount(binary.BigEndian.Uint16(b))
	return nil
}

func (s *ExchangeID) UnmarshalBinary(b []byte) error {
	if len(b) != 2 {
		return io.ErrUnexpectedEOF
	}
	*s = ExchangeID(binary.BigEndian.Uint16(b))
	return nil
}
