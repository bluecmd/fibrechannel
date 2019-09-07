package els

import (
	"io"

	"github.com/bluecmd/fibrechannel/encoding"
)

// FLOGI and PLOGI shares data format
type PLOGI FLOGI

func (s *PLOGI) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFrom(r, s)
}

func (s *PLOGI) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, s)
}
