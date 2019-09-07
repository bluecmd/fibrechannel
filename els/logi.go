package els

import (
	"io"

	"github.com/bluecmd/fibrechannel/common"
	"github.com/bluecmd/fibrechannel/encoding"
)

type CommonFeatures uint32

type CommonSvcParams struct {
	FCPHVersion uint16
	// The buffer-to-buffer credits is the number of buffers available for holding
	// Class 2/3 frames received
	B2BCredits       uint16
	Features         CommonFeatures // TODO(bluecmd): get/set for bits
	MaxConcurrentSeq uint16         // TODO(bluecmd): Verify size
	RelOffsetInfoCat uint16
	// E_D_TOV is "Error Detect TimeOut Value"
	// If "E_D_TOV Resolution" feature is set to zero, the E_D_TOV value
	// is in milliseconds. If the bit is one it is instead nanoseconds.
	EDTOV uint32
}

type ClassSvcParams struct {
	Service              uint16
	Initiator            uint16
	Recipient            uint16
	ReceiveDataFieldSize uint16 // TODO(bluecmd): Really only 12 bits
	_                    uint8
	ConcurrentSeq        uint8
	E2ECredits           uint16
	_                    uint8
	OpenSeqPerExch       uint8
	_                    uint16
}

type LOGI struct {
	CommonSvcParams CommonSvcParams   `fc:"@3"`
	PortName        common.WWN        `fc:"@19"`
	NodeName        common.WWN        `fc:"@27"`
	ClassSvcParams  [3]ClassSvcParams `fc:"@35"`
	AuxSvcParams    ClassSvcParams    `fc:"@83"`
	VendorVersion   [16]byte          `fc:"@99"`
}

type PLOGI struct {
	LOGI
}

type FLOGI struct {
	LOGI
}

func (s *LOGI) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFromAndPost(r, s)
}

func (s *LOGI) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteToAndPre(w, s)
}

func (s *LOGI) PostUnmarshal() error {
	// if CommonSvcParams.HasPayload {
	return nil
}

func (s *LOGI) PreMarshal() error {
	// if CommonSvcParams.HasPayload {
	return nil
}
