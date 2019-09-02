package swils

import (
	"io"
)

type Command int

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
	Command Command
	Payload []byte
}

func (f *Frame) UnmarshalBinary(b []byte) error {
	if len(b) < 4 {
		return io.ErrUnexpectedEOF
	}

	f.Command = Command(b[0])
	f.Payload = b[4:]
	return nil
}
