package fcoe

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	fc "github.com/bluecmd/fibrechannel"
)

func TestFrameMarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		f    *Frame
		b    []byte
		err  error
	}{
		{
			desc: "Normal frame",
			f:    &Frame{SOF: fc.SOFf, EOF: fc.EOFn, CRC32: 0x10},
			b:    append(bytes.Repeat([]byte{0}, 13), 0x28, 0, 0, 0, 0x10, 0x41, 0, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b, err := tt.f.MarshalBinary()
			if err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}
				return
			}

			if want, got := tt.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %v\n-  got: %v", want, got)
			}
		})
	}
}

func TestFrameUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		f    *Frame
		b    []byte
		err  error
	}{
		{
			desc: "nil buffer",
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "short buffer",
			b:    bytes.Repeat([]byte{0}, 13),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "standard frame",
			b: []byte{
				0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x28,
				0x12, 0x34, 0x56, 0x78,
				0x41, 0x00, 0x00, 0x00,
			},
			f: &Frame{
				CRC32:   0x12345678,
				SOF:     fc.SOFf,
				Payload: []byte{},
				EOF:     fc.EOFn,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			f := new(Frame)
			if err := f.UnmarshalBinary(tt.b); err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}
				return
			}

			if want, got := tt.f, f; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Frame:\n- want: %v\n-  got: %v", want, got)
			}
		})
	}
}

func TestCRC32(t *testing.T) {
	p := []byte{
		0x22, 0xff, 0xff, 0xfd, 0x00, 0xed, 0x01, 0x00, 0x01, 0x29, 0x00, 0x00,
		0xf0, 0x00, 0x00, 0x00, 0x03, 0xf8, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x62, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}
	f := &Frame{Payload: p}
	f.CRC32 = f.Checksum()
	want := uint32(0x1b726f69)
	got := f.CRC32
	if want != got {
		t.Fatalf("CRC calculation failed: wanted %08x, got %08x", want, got)
	}
}
