package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sync"
)

type postUnmarshaller interface {
	PostUnmarshal() error
}

type preMarshaller interface {
	PreMarshal() error
}

type fcTag struct {
	off int64
}

type fcTagCache struct {
	tags []*fcTag
}

var (
	tagCacheMap  = map[reflect.Type]*fcTagCache{}
	tagCacheLock = sync.RWMutex{}
)

func parseTag(f *reflect.StructField) (*fcTag, error) {
	fc := f.Tag.Get("fc")
	if fc == "" {
		return nil, nil
	}
	tag := fcTag{}
	fmt.Sscanf(fc, "@%d", &tag.off)
	return &tag, nil
}

func tagcache(t reflect.Type) (*fcTagCache, error) {
	tagCacheLock.RLock()
	tc, ok := tagCacheMap[t]
	tagCacheLock.RUnlock()
	if !ok {
		tc = &fcTagCache{
			tags: []*fcTag{},
		}
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			tag, err := parseTag(&ft)
			if err != nil {
				return nil, fmt.Errorf("%v tag error: %v", ft, err)
			}
			tc.tags = append(tc.tags, tag)
		}
		tagCacheLock.Lock()
		tagCacheMap[t] = tc
		tagCacheLock.Unlock()
	}
	return tc, nil
}

// Reads a fibrechannel package annotated structure from the io.Reader
// If fed a structure that does not have any `fc` tags it will do nothing
func ReadFrom(r io.Reader, f interface{}) (int64, error) {
	ptr := reflect.ValueOf(f)
	if ptr.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("Expected pointer to fibrechannel Frame, got %v", ptr)
	}
	frm := ptr.Elem()
	frmt := frm.Type()

	tc, err := tagcache(frmt)
	if err != nil {
		return 0, err
	}

	var pos int64
	for i, tag := range tc.tags {
		if tag == nil {
			continue
		}
		if pos > tag.off {
			return pos, fmt.Errorf("would have gone backwards on field %v", frmt.Field(i))
		}
		// Skip ahead to the new position if needed
		if tag.off > pos {
			skip := int64(tag.off - pos)
			n, err := io.CopyN(ioutil.Discard, r, skip)
			if err != nil {
				return pos + n, err
			}
			if n != skip {
				return pos + n, io.ErrUnexpectedEOF
			}
		}
		fi := frm.Field(i).Addr().Interface()
		// If the field has a ReadFrom method, use it - otherwise use binary.Read
		if rdr, ok := fi.(io.ReaderFrom); ok {
			n, err := rdr.ReadFrom(r)
			if err != nil {
				if err == io.EOF {
					return pos + n, io.ErrUnexpectedEOF
				} else {
					return pos + n, err
				}
			}
			pos = tag.off + n
			continue
		}

		// If the field is a []byte assume it's the payload and shove the rest in there
		if slice, ok := fi.(*[]byte); ok {
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(r)
			if err != nil {
				return pos, err
			}
			*slice = buf.Bytes()
			// We're done, there cannot be anything left
			break
		}

		// Otherwise default to binary.Read
		if err := binary.Read(r, binary.BigEndian, fi); err != nil {
			if err == io.EOF {
				return pos, io.ErrUnexpectedEOF
			} else {
				return pos, err
			}
		}
		pos = tag.off + int64(frmt.Field(i).Type.Size())
	}

	// If there is a PostUnmarshal, call it with the byte array
	if pm, ok := f.(postUnmarshaller); ok {
		return pos, pm.PostUnmarshal()
	}
	return pos, nil
}

// Writes a fibrechannel package annotated structure to the io.Writer
// If fed a structure that does not have any `fc` tags it will do nothing
func WriteTo(w io.Writer, f interface{}) (int64, error) {
	// If there is a PostUnmarshal, call it with the byte array
	if pm, ok := f.(preMarshaller); ok {
		if err := pm.PreMarshal(); err != nil {
			return 0, err
		}
	}

	ptr := reflect.ValueOf(f)
	if ptr.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("Expected pointer to fibrechannel Frame, got %v", ptr)
	}
	frm := ptr.Elem()
	frmt := frm.Type()

	tc, err := tagcache(frmt)
	if err != nil {
		return 0, err
	}

	var pos int64
	for i, tag := range tc.tags {
		if tag == nil {
			continue
		}
		// We require the struct to be defined in order to not have to jump around
		if pos > tag.off {
			return pos, fmt.Errorf("would have gone backwards on field %v", frmt.Field(i))
		}
		// Skip ahead to the new position if needed
		if tag.off > pos {
			skip := tag.off - pos
			// TODO(bluecmd): A lot of wasted allocations?
			n, err := w.Write(make([]byte, skip))
			if err != nil {
				return pos + int64(n), err
			}
			if int64(n) != skip {
				return pos + int64(n), io.ErrUnexpectedEOF
			}
		}
		// If the field has a WriteTo method, use it - otherwise use binary.Write
		fi := frm.Field(i).Addr().Interface()
		if wrt, ok := fi.(io.WriterTo); ok {
			n, err := wrt.WriteTo(w)
			if err != nil {
				return pos + n, err
			}
			pos = tag.off + n
		} else {
			if err := binary.Write(w, binary.BigEndian, fi); err != nil {
				return pos, err
			}
			pos = tag.off + int64(frmt.Field(i).Type.Size())
		}
	}

	return pos, nil
}
