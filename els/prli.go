package els

import (
	"io"

	"github.com/bluecmd/fibrechannel/encoding"
)

type ServiceParameter struct {
	Type byte `fc:"@0"`
	_    [15]byte
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
