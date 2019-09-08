package fibrechannel

import (
	"bytes"
	"io"
	"testing"

	"github.com/bluecmd/fibrechannel/common"
)

func TestNilBuffer(t *testing.T) {
	c := &Frame{}
	_, err := c.ReadFrom(bytes.NewReader([]byte{}))
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("got unexpected error %v, wanted io.ErrUnexpectedEOF", err)
	}
}

func TestShortBuffer(t *testing.T) {
	c := &Frame{}
	_, err := c.ReadFrom(bytes.NewReader([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("got unexpected error %v, wanted io.ErrUnexpectedEOF", err)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	buf := bytes.NewReader(bytes.Repeat([]byte{0}, 28))
	f := new(Frame)
	for n := 0; n < b.N; n++ {
		buf.Seek(0, io.SeekStart)
		f.ReadFrom(buf)
	}
}

func BenchmarkMarshal(b *testing.B) {
	buf := new(bytes.Buffer)
	buf.Grow(10000)
	f := &Frame{CSCtl: new(ClassControl), RawPayload: []byte{1, 2, 3, 4}}
	for n := 0; n < b.N; n++ {
		buf.Reset()
		f.WriteTo(buf)
	}
}

func TestFrameFiles(t *testing.T) {
	common.TestFrameFiles(t, func() common.SerDes { return &Frame{} })
}
