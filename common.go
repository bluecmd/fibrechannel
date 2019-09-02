package fibrechannel

import (
	"fmt"
)

type WWN [8]byte

func (s *WWN) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x",
		s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7])
}
