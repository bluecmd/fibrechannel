package els

import (
	"fmt"
	"io"

	"github.com/bluecmd/fibrechannel/encoding"
)

type ServiceParameter struct {
	Type byte     `fc:"@0"`
	Todo [15]byte `fc:"@1"`
}

type PRLI struct {
	PagesLength       uint8  `fc:"@0"`
	PayloadLength     uint16 `fc:"@1"`
	ServiceParameters []ServiceParameter
}

func (s *PRLI) ReadFrom(r io.Reader) (int64, error) {
	n, err := encoding.ReadFrom(r, s)
	if err != nil {
		return n, err
	}
	for i := 0; i < int(s.PagesLength)/16; i++ {
		x := ServiceParameter{}
		m, err := encoding.ReadFrom(r, &x)
		if err != nil {
			return n + m, err
		}
		n += m

		s.ServiceParameters = append(s.ServiceParameters, x)
	}
	return n, nil
}

func (s *PRLI) WriteTo(w io.Writer) (int64, error) {
	if len(s.ServiceParameters)*16 > 255 {
		return 0, fmt.Errorf("Too many ServiceParameters")
	}
	s.PagesLength = uint8(len(s.ServiceParameters) * 16)
	n, err := encoding.WriteTo(w, s)
	if err != nil {
		return n, err
	}
	for _, p := range s.ServiceParameters {
		m, err := encoding.WriteTo(w, &p)
		if err != nil {
			return n + m, err
		}
		n += m
	}
	return n, nil

}
