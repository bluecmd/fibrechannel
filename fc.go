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
	TypeNVME     = 0x28
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
	// Start-of-frame, not part of the wire-format
	// The SOF/EOF is an important part of FC so this is a convenience field
	// for implementations to keep track of which marker is used for the frame
	SOF SOF

	// End-of-frame, see "start-of-frame"
	EOF EOF

	RCtl RoutingControl `fc:"@0"`

	// Address for source/destination Nx_Ports
	// Each Nx_Port shall have a native N_Port_ID that is unique within the
	// address domain of a Fabric.  An N_Port_ID of binary zeros indicates that
	// an Nx_Port is unidentified. When a PN_Port completes Link Initialization,
	// it shall be unidentified (i.e., it shall have a single Nx_Port for which
	// the N_Port_ID is 00 00 00h).
	Destination Address `fc:"@1"`

	// Which field of these two are in use is controlled by the
	// "CS_CTL/Priority Enable" bit in F_CTL.
	CsctlPriorityRaw byte          `fc:"@4"`
	CSCtl            *ClassControl // Handled by PostUnmarshal
	Priority         *Priority     // Handled by PostUnmarshal

	// Same as Destination
	Source Address `fc:"@5"`

	Type Type `fc:"@8"`

	FCtl FrameControl `fc:"@9"`

	// Sequence ID
	// The sequence count (SEQ_CNT) is a two-byte field (Word 3, Bits 15-0)
	// that shall indicate the sequential order of Data frame transmission within
	// a single Sequence or multiple consecutive Sequences for the same Exchange.
	// The SEQ_CNT of the first Data frame of the first Sequence of the Exchange
	// transmitted by either the Originator or Responder shall be binary zero.
	// The SEQ_CNT of each subsequent Data frame in the Sequence shall be
	// incremented by one.
	SeqID SequenceID `fc:"@12"`

	DFCtl DataFieldControl `fc:"@13"`

	// Sequence count
	SeqCount SequenceCount `fc:"@14"`

	// Originator Exchange_ID
	// If the Originator is enforcing uniqueness via the OX_ID mechanism, it
	// shall set a unique value for OX_ID other than FF FFh in the first Data
	// frame of the first Sequence of an Exchange. An OX_ID of FF FFh indicates
	// that the OX_ID is unassigned and that the Originator is not enforcing
	// uniqueness via the OX_ID mechanism. If an Originator uses the unassigned
	// value of FF FFh to identify the Exchange, it shall have only one Exchange
	// (OX_ID set to FF FFh) with a given Responder.
	OXID ExchangeID `fc:"@16"`

	// Responder Exchange_ID
	// An RX_ID of FF FFh shall indicate that the RX_ID is unassigned. If the
	// Responder does not assign an RX_ID other than FF FFh by the end of the
	// first Sequence, then the Responder is not enforcing uniqueness via the
	// RX_ID mechanism.
	RXID ExchangeID `fc:"@18"`

	// TODO(bluecmd): Optional fields
	// TODO(bluecmd): Parameters

	RawPayload []byte      `fc:"@24"`
	Payload    interface{} // Handled by Rebuild
}

func (f *Frame) PostUnmarshal() error {
	if f.FCtl.CSCtlEnabled() {
		cc := ClassControl(f.CsctlPriorityRaw)
		f.CSCtl = &cc
	} else {
		p := Priority(f.CsctlPriorityRaw)
		f.Priority = &p
	}
	if f.DFCtl.HasESP() {
		return fmt.Errorf("ESP is not implemented")
	}
	if f.DFCtl.HasNetworkHeader() {
		return fmt.Errorf("Network header is not implemented")
	}
	return nil
}

func (f *Frame) PreMarshal() error {
	if f.FCtl.CSCtlEnabled() {
		if f.CSCtl == nil {
			return fmt.Errorf("CSCtl is missing but enabled")
		}
		f.CsctlPriorityRaw = byte(*f.CSCtl)
	} else {
		if f.Priority == nil {
			return fmt.Errorf("Priority is missing but enabled")
		}
		f.CsctlPriorityRaw = byte(*f.Priority)
	}

	// TODO: Parameters
	if f.DFCtl.HasESP() {
		return fmt.Errorf("ESP is not implemented")
	}
	if f.DFCtl.HasNetworkHeader() {
		return fmt.Errorf("Network header is not implemented")
	}
	return nil
}

func (s *FrameControl) ReadFrom(r io.Reader) (int64, error) {
	b := [3]byte{}
	n, err := r.Read(b[:])
	if err != nil {
		return int64(n), err
	}
	*s = FrameControl(binary.BigEndian.Uint32(append([]byte{0}, b[:]...)))
	return 3, nil
}

func (s *FrameControl) WriteTo(r io.Writer) (int64, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(*s))
	n, err := r.Write(buf[1:3])
	return int64(n), err
}

func (s *FrameControl) CSCtlEnabled() bool {
	return *s&0x20000 == 0
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

// Helper to save the SOF and EOF markers in the Fibre Channel frame
func ReadFrame(sof SOF, r io.Reader, eof EOF, f *Frame) (int64, error) {
	f.SOF = sof
	f.EOF = eof
	return ReadFrom(r, f)
}
