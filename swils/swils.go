package swils

import (
	"bytes"
	"io"

	"github.com/bluecmd/fibrechannel/encoding"
)

type Command uint8

const (
	CmdSWRJT  = 0x01 // Switch Fabric Internal Link Service Reject
	CmdSWACC  = 0x02 // Switch Fabric Internal Link Service Accept
	CmdELP    = 0x10 // Exchange Link Parameters
	CmdEFP    = 0x11 // Exchange Fabric Parameters
	CmdDIA    = 0x12 // Domain Identifier Assigned
	CmdRDI    = 0x13 // Request Domain_ID
	CmdHLO    = 0x14 // Hello
	CmdLSU    = 0x15 // Link State Update
	CmdLSA    = 0x16 // Link State Acknowledgement
	CmdBF     = 0x17 // Build Fabric
	CmdRCF    = 0x18 // Reconfigure Fabric
	CmdSWRSCN = 0x1B // Inter-Switch Registered State Change Notification
	CmdDRLIR  = 0x1E // Distribute Registered Link Incident Records
	CmdDSCN   = 0x20 // Obsoleted in FC-SW-5
	CmdLOOPD  = 0x21 // Obsoleted in FC-SW-3
	CmdMR     = 0x22 // Merge Request
	CmdACA    = 0x23 // Acquire Change Authorization
	CmdRCA    = 0x24 // Release Change Authorization
	CmdSFC    = 0x25 // Stage Fabric Configuration
	CmdUFC    = 0x26 // Update Fabric Configuration
	CmdCEC    = 0x29 // Check E_Port Connectivity
	// 2A has loads of meanings, skipping for now
	CmdESC  = 0x30 // Exchange Switch Capabilities
	CmdESS  = 0x31 // Exchange Switch Support
	CmdMRRA = 0x34 // Merge Request Resource Allocation
	CmdSTR  = 0x35 // Switch Trace Route
	CmdEVFP = 0x36 // Exchange Virtual Fabrics Parameters
	CmdFFI  = 0x50 // Fast Fabric Initialization for the Avionics Environment
)

type Frame struct {
	Command    Command `fc:"@0"`
	RawPayload []byte  `fc:"@4"`
	Payload    interface{}
}

func (f *Frame) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFromAndPost(r, f)
}

func (f *Frame) PostUnmarshal() error {
	var sf io.ReaderFrom
	switch f.Command {
	case CmdELP:
		sf = &ELP{}
	}

	if sf == nil {
		return nil
	}

	_, err := sf.ReadFrom(bytes.NewReader(f.RawPayload))
	if err != nil {
		return err
	}
	f.Payload = sf
	f.RawPayload = nil
	return nil
}

func (f *Frame) PreMarshal() error {
	if f.Payload == nil {
		return nil
	}
	b := bytes.NewBuffer(f.RawPayload)
	_, err := f.Payload.(io.WriterTo).WriteTo(b)
	f.RawPayload = b.Bytes()
	return err
}

func (f *Frame) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteToAndPre(w, f)
}
