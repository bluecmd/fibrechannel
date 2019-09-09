package encoding

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	R     io.Reader
	Error error
	Pos   int64
}

type Writer struct {
	W     io.Writer
	Error error
	Pos   int64
}

func (r *Reader) Read(b []byte) (int, error) {
	if n, err := r.R.Read(b); err != nil {
		r.Pos += int64(n)
		r.Error = err
		return n, err
	}
	r.Pos += int64(len(b))
	return len(b), nil
}

func (r *Reader) Skip(n int) {
	n, err := r.R.Read(make([]byte, n))
	if err != nil {
		r.Error = err
	}
	r.Pos += int64(n)
}

func (r *Reader) ReadObject(v interface{}) {
	if err := binary.Read(r, binary.BigEndian, v); err != nil {
		r.Error = err
	}
}

func (w *Writer) Write(b []byte) (int, error) {
	if n, err := w.W.Write(b); err != nil {
		w.Pos += int64(n)
		w.Error = err
		return n, err
	}
	w.Pos += int64(len(b))
	return len(b), nil
}

func (w *Writer) Skip(n int) {
	n, err := w.W.Write(make([]byte, n))
	if err != nil {
		w.Error = err
	}
	w.Pos += int64(n)
}

func (w *Writer) WriteObject(v interface{}) {
	if err := binary.Write(w, binary.BigEndian, v); err != nil {
		w.Error = err
	}
}
