package els

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	fc "github.com/bluecmd/fibrechannel"
)

func TestFrameUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		c    io.ReaderFrom
		f    io.ReaderFrom
		b    []byte
		err  error
	}{
		{
			desc: "nil buffer",
			c:    &Frame{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "els frame",
			b:    append([]byte{0x04, 0, 0, 0}, bytes.Repeat([]byte{1}, 1)...),
			c:    &Frame{},
			f:    &Frame{Command: CmdFLOGI, Payload: []byte{1}},
		},
		{
			desc: "flogi",
			b: append([]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
				0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
				bytes.Repeat([]byte{0}, 84)...),
			c: &FLOGI{},
			f: &FLOGI{
				WWPN: fc.WWN([8]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}),
				WWNN: fc.WWN([8]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if _, err := tt.c.ReadFrom(bytes.NewReader(tt.b)); err != nil {
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
