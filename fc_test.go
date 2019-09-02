package fibrechannel

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

//func TestFrameMarshalBinary(t *testing.T) {
//	var tests = []struct {
//		desc string
//		f    *Frame
//		b    []byte
//		err  error
//	}{
//		{
//			desc: "Normal frame",
//			f:    &Frame{CSCtl: new(ClassControl)},
//			b:    bytes.Repeat([]byte{0}, 24),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.desc, func(t *testing.T) {
//			b, err := tt.f.MarshalBinary()
//			if err != nil {
//				if want, got := tt.err, err; want != got {
//					t.Fatalf("unexpected error: %v != %v", want, got)
//				}
//				return
//			}
//
//			if want, got := tt.b, b; !bytes.Equal(want, got) {
//				t.Fatalf("unexpected Frame bytes:\n- want: %v\n-  got: %v", want, got)
//			}
//		})
//	}
//}

func TestFrameUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		s    SOF
		f    *Frame
		e    EOF
		b    []byte
		err  error
	}{
		{
			desc: "nil buffer",
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "short buffer",
			b:    bytes.Repeat([]byte{0}, 12),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "normal frame - csctl enabled",
			b:    bytes.Repeat([]byte{0}, 24),
			s:    SOFf,
			e:    EOFn,
			f:    &Frame{SOF: SOFf, EOF: EOFn, CSCtl: new(ClassControl)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			f := new(Frame)
			if err := ReadFrame(tt.s, bytes.NewBuffer(tt.b), tt.e, f); err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}
				return
			}

			if want, got := tt.f, f; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Frame:\n- want: %v\n-  got: %v", spew.Sdump(want), spew.Sdump(got))
			}
		})
	}
}
