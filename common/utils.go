package common

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/davecgh/go-spew/spew"
)

type SerDes interface {
	io.ReaderFrom
	io.WriterTo
}

func TestFrameFiles(t *testing.T, f func() SerDes) {
	s := spew.NewDefaultConfig()
	s.DisablePointerAddresses = true

	err := filepath.Walk("testdata", func(path string, o os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if o.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".fc") {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			d, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			frm := f()
			n, err := frm.ReadFrom(bytes.NewReader(d))
			if err != nil {
				t.Errorf("ReadFrom: %v", err)
			}
			delta := n - int64(len(d))
			if delta > 0 {
				t.Errorf("ReadFrom left %d bytes", delta)
			}

			dc, err := ioutil.ReadFile(path + ".golden")
			if err != nil {
				t.Fatalf("open golden: %v", err)
			}
			want := string(dc)
			got := s.Sdump(frm)
			if want != got {
				t.Errorf("Result differ from golden:\n%v", diff.LineDiff(want, got))
			}

			dd := new(bytes.Buffer)
			_, err = frm.WriteTo(dd)
			if err != nil {
				t.Fatalf("Failed to re-serialize: %v", err)
			}
			if want, got := d, dd.Bytes(); !bytes.Equal(want, got) {
				t.Fatalf("unexpected de-serialized output:\n- want: %v\n-  got: %v", want, got)
			}
		})
		return nil
	})

	if err != nil {
		t.Fatalf("err from walk: %v", err)
	}
}
