package els

import (
	"io"

	fc "github.com/bluecmd/fibrechannel"
)

type ClassSvcParams struct {
	// TODO
}

type CommonSvcParams struct {
	// TODO
}

type FLOGI struct {
	CSP  CommonSvcParams
	WWPN fc.WWN `fc:"@16"`
	WWNN fc.WWN `fc:"@24"`
	CSSP [4]ClassSvcParams
}

func (s *FLOGI) ReadFrom(r io.Reader) (int64, error) {
	return fc.ReadFrom(r, s)
}
