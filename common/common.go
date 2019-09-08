package common

import (
	"encoding/binary"
	"fmt"
	"io"
)

type WWN [8]byte

func (s *WWN) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x",
		s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7])
}

func (s *WWN) ReadFrom(r io.Reader) (int64, error) {
	if err := binary.Read(r, binary.BigEndian, s); err != nil {
		return 0, err
	}
	return 8, nil
}

func (s *WWN) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.BigEndian, s); err != nil {
		return 0, err
	}
	return 8, nil
}

type Uint16SizedByteArray []byte

func (p *Uint16SizedByteArray) ReadFrom(r io.Reader) (int64, error) {
	var cnt uint16
	if err := binary.Read(r, binary.BigEndian, &cnt); err != nil {
		return 0, err
	}
	*p = make([]byte, cnt)
	n, err := r.Read(*p)
	if err != nil {
		return int64(n), err
	}
	if n != int(cnt) {
		return int64(n), io.EOF
	}
	return int64(n), nil
}

func (p *Uint16SizedByteArray) WriteTo(w io.Writer) (int64, error) {
	s := uint16(len(*p))
	if err := binary.Write(w, binary.BigEndian, &s); err != nil {
		return 0, err
	}
	n, err := w.Write(*p)
	if err != nil {
		return int64(n), err
	}
	if n != int(s) {
		return int64(n), io.EOF
	}
	return int64(n), nil
}
