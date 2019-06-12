package fcoe

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	fc "github.com/bluecmd/fibrechannel"
)

func TestFrameUnmarshalBinary(t *testing.T) {
	fcfb := bytes.Repeat([]byte{0}, 28)
	fcf := fc.Frame{}

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
			b: append([]byte{
				0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x28,
			}, append(fcfb, []byte{
				0x12, 0x34, 0x56, 0x78,
				0x41, 0x00, 0x00, 0x00,
			}...)...),
			f: &Frame{
				CRC32:   0x12345678,
				SOF:     fc.SOFf,
				Payload: fcf,
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
