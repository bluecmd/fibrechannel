package els

import (
	"io"
)

type ClassSvcParams struct {
	// TODO
}

type CommonSvcParams struct {
	// TODO
}

type FLOGI struct {
	CSP  CommonSvcParams
	WWPN WWN
	WWNN WWN
	CSSP [4]ClassSvcParams
}

func (s *FLOGI) UnmarshalBinary(b []byte) error {
	if len(b) < 112 {
		return io.ErrUnexpectedEOF
	}

	// TODO CSP / CSSP

	copy(s.WWPN[:], b[16:24])
	copy(s.WWNN[:], b[24:32])
	return nil
}
