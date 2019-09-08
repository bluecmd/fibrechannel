package els

import (
	"bytes"
	"io"
	"testing"

	"github.com/bluecmd/fibrechannel/common"
)

func TestNilBuffer(t *testing.T) {
	c := &Frame{}
	_, err := c.ReadFrom(bytes.NewReader([]byte{}))
	if err != io.EOF {
		t.Fatalf("got unexpected error %v, wanted io.EOF", err)
	}
}

func TestUnsupportedCmd(t *testing.T) {
	c := &Frame{}
	_, err := c.ReadFrom(bytes.NewReader([]byte{0xff}))
	if err != nil {
		t.Fatalf("got error %v, expected no error", err)
	}
	_, err = c.WriteTo(new(bytes.Buffer))
	if err == nil {
		t.Fatalf("got no error, expected one")
	}
}

func TestFrameFiles(t *testing.T) {
	common.TestFrameFiles(t, func() common.SerDes { return &Frame{} })
}
