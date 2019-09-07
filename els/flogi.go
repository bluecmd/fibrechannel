package els

import (
	"io"

	"github.com/bluecmd/fibrechannel/common"
	"github.com/bluecmd/fibrechannel/encoding"
)

type ClassSvcParams struct {
	// TODO
	X [16]byte
}

type CommonSvcParams struct {
	// TODO
	X [16]byte
}

type FLOGI struct {
	CSP           CommonSvcParams   `fc:"@3"`
	WWPN          common.WWN        `fc:"@19"`
	WWNN          common.WWN        `fc:"@27"`
	CSSP          [4]ClassSvcParams `fc:"@35"`
	VendorVersion [16]byte          `fc:"@99"`
}

func (s *FLOGI) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, s)
}

func (s *FLOGI) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, s)
}
