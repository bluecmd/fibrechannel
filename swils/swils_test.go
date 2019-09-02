package swils

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	fc "github.com/bluecmd/fibrechannel"
)

type swilsCmd interface {
	UnmarshalBinary(b []byte) error
}

func TestFrameUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		c    swilsCmd
		f    swilsCmd
		b    []byte
		err  error
	}{
		{
			desc: "nil buffer",
			c:    &Frame{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "swils frame",
			b:    append([]byte{0x10, 0, 0, 0}, bytes.Repeat([]byte{1}, 1)...),
			c:    &Frame{},
			f:    &Frame{Command: CmdELP, Payload: []byte{1}},
		},
		{
			desc: "elp",
			b: append(append([]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
				0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
				bytes.Repeat([]byte{0}, 50)...), 0, 0x03, 1, 2, 3, 0, 0, 0),
			c: &ELP{},
			f: &ELP{
				Port:                  fc.WWN([8]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}),
				Switch:                fc.WWN([8]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}),
				FlowControlParameters: []byte{1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if err := tt.c.UnmarshalBinary(tt.b); err != nil {
				if want, got := tt.err, err; want != got {
					t.Fatalf("unexpected error: %v != %v", want, got)
				}
				return
			}

			if want, got := tt.f, tt.c; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Frame:\n- want: %v\n-  got: %v", want, got)
			}
		})
	}
}
