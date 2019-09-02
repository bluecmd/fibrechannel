package fibrechannel

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
			f:    &Frame{CSCtl: new(ClassControl), RawPayload: []byte{1, 2, 3, 4}},
			b:    append(bytes.Repeat([]byte{0}, 24), 1, 2, 3, 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := Write(buf, tt.f)
			if err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}
				return
			}

			if want, got := tt.b, buf.Bytes(); !bytes.Equal(want, got) {
				t.Fatalf("unexpected Frame bytes:\n- want: %v\n-  got: %v", want, got)
			}
		})
	}
}

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
			b:    bytes.Repeat([]byte{0}, 28),
			s:    SOFf,
			e:    EOFn,
			f:    &Frame{SOF: SOFf, EOF: EOFn, CSCtl: new(ClassControl), RawPayload: []byte{0, 0, 0, 0}},
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

func BenchmarkUnmarshal(b *testing.B) {
	buf := bytes.NewReader(bytes.Repeat([]byte{0}, 28))
	f := new(Frame)
	for n := 0; n < b.N; n++ {
		buf.Seek(0, io.SeekStart)
		Read(buf, f)
	}
}

func BenchmarkMarshal(b *testing.B) {
	buf := new(bytes.Buffer)
	buf.Grow(10000)
	f := &Frame{CSCtl: new(ClassControl), RawPayload: []byte{1, 2, 3, 4}}
	for n := 0; n < b.N; n++ {
		buf.Reset()
		Write(buf, f)
	}
}
