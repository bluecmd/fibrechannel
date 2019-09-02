package fibrechannel

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

type postUnmarshal interface {
	PostUnmarshal() error
}

type fcTag struct {
	off int
}

func parseTag(f *reflect.StructField) (*fcTag, error) {
	fc := f.Tag.Get("fc")
	if fc == "" {
		return nil, nil
	}
	tag := fcTag{}
	fmt.Sscanf(fc, "@%d", &tag.off)
	return &tag, nil
}

// Reads a fibrechannel package annotated structure from the io.Reader
// If fed a structure that does not have any `fc` tags it will do nothing
func Read(r io.Reader, f interface{}) error {
	// TODO(bluecmd): I'd be interested in annotating this to see how
	// much speed we lose on doing reflection vs. hardcoding all the copies.
	ptr := reflect.ValueOf(f)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("Expected pointer to fibrechannel Frame, got %v", ptr)
	}
	frm := ptr.Elem()

	pos := 0
	for i := 0; i < frm.NumField(); i++ {
		ft := frm.Type().Field(i)
		tag, err := parseTag(&ft)
		if err != nil {
			return fmt.Errorf("%v tag error: %v", ft, err)
		}
		if tag == nil {
			continue
		}
		// We require the struct to be defined in order to not have to jump around
		if pos > tag.off {
			return fmt.Errorf("would have gone backwards on field %v", ft)
		}
		// Skip ahead to the new position if needed
		if tag.off > pos {
			skip := int64(tag.off - pos)
			n, err := io.CopyN(ioutil.Discard, r, skip)
			if err != nil {
				return err
			}
			if n != skip {
				return io.ErrUnexpectedEOF
			}
		}
		// If the field has a ReadFrom method, use it - otherwise use binary.Read
		fi := frm.Field(i).Addr().Interface()
		if rdr, ok := fi.(io.ReaderFrom); ok {
			n, err := rdr.ReadFrom(r);
			if err != nil {
				if err == io.EOF {
					return io.ErrUnexpectedEOF
				} else {
					return err
				}
			}
			pos = tag.off + int(n)
		} else {
			if err := binary.Read(r, binary.BigEndian, fi); err != nil {
				if err == io.EOF {
					return io.ErrUnexpectedEOF
				} else {
					return err
				}
			}
			pos = tag.off + int(ft.Type.Size())
		}
	}

	// If there is a PostUnmarshal, call it with the byte array
	if pm, ok := f.(postUnmarshal); ok {
		return pm.PostUnmarshal()
	}
	return nil
}
