package swils

import (
	"io"

	"github.com/bluecmd/fibrechannel/common"
	"github.com/bluecmd/fibrechannel/encoding"
)

type ELP struct {
	Revision           uint8                       `fc:"@0"`
	Flags              uint16                      `fc:"@1"`
	BBSCN              uint8                       `fc:"@3"`
	RATOV              uint32                      `fc:"@4"`
	EDTOV              uint32                      `fc:"@8"`
	Port               common.WWN                  `fc:"@12"` // RequesterPortName
	Switch             common.WWN                  `fc:"@20"` // RequesterSwitchName
	ClassFParameters   [16]byte                    `fc:"@28"`
	Class2Parameters   [4]byte                     `fc:"@48"`
	Class3Parameters   [4]byte                     `fc:"@52"`
	ISLFlowControlMode uint16                      `fc:"@76"`
	FCParam            common.Uint16SizedByteArray `fc:"@78"`
}

func (s *ELP) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, s)
}

func (s *ELP) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, s)
}
