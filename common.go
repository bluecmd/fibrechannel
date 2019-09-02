package fibrechannel

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
