package swils

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

func TestFrameFiles(t *testing.T) {
	common.TestFrameFiles(t, func() common.SerDes { return &Frame{} })
}
