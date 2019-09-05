package els

import (
	"io"

	"github.com/bluecmd/fibrechannel/common"
	"github.com/bluecmd/fibrechannel/encoding"
)

type ClassSvcParams struct {
	// TODO
}

type CommonSvcParams struct {
	// TODO
}

type FLOGI struct {
	CSP  CommonSvcParams
	WWPN common.WWN `fc:"@16"`
	WWNN common.WWN `fc:"@24"`
	CSSP [4]ClassSvcParams
}

func (s *FLOGI) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, s)
}
