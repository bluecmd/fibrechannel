package els

import (
	"bytes"
	"io"

	"github.com/bluecmd/fibrechannel/encoding"
)

type Command uint8

const (
	CmdLSRJT     = 0x01 // ESL reject
	CmdLSACC     = 0x02 // ESL Accept
	CmdPLOGI     = 0x03 // N_Port login
	CmdFLOGI     = 0x04 // F_Port login
	CmdLOGO      = 0x05 // Logout
	CmdABTX      = 0x06 // Abort exchange - obsolete
	CmdRCS       = 0x07 // read connection status
	CmdRES       = 0x08 // read exchange status block
	CmdRSS       = 0x09 // read sequence status block
	CmdRSI       = 0x0a // read sequence initiative
	CmdESTS      = 0x0b // establish streaming
	CmdESTC      = 0x0c // estimate credit
	CmdADVC      = 0x0d // advise credit
	CmdRTV       = 0x0e // read timeout value
	CmdRLS       = 0x0f // read link error status block
	CmdEcho      = 0x10 // echo
	CmdTest      = 0x11 // test
	CmdRRQ       = 0x12 // reinstate recovery qualifier
	CmdREC       = 0x13 // read exchange concise
	CmdSRR       = 0x14 // sequence retransmission request
	CmdPRLI      = 0x20 // process login
	CmdPRLO      = 0x21 // process logout
	CmdSCN       = 0x22 // state change notification
	CmdTPLS      = 0x23 // test process login state
	CmdTPRLO     = 0x24 // third party process logout
	CmdLCLM      = 0x25 // login control list mgmt (obs)
	CmdGAID      = 0x30 // get alias_ID
	CmdFACT      = 0x31 // fabric activate alias_id
	CmdFDACDT    = 0x32 // fabric deactivate alias_id
	CmdNACT      = 0x33 // N-port activate alias_id
	CmdNDACT     = 0x34 // N-port deactivate alias_id
	CmdQOSR      = 0x40 // quality of service request
	CmdRVCS      = 0x41 // read virtual circuit status
	CmdPDISC     = 0x50 // discover N_port service params
	CmdFDISC     = 0x51 // discover F_port service params
	CmdADISC     = 0x52 // discover address
	CmdRNC       = 0x53 // report node cap (obs)
	CmdFARPReq   = 0x54 // FC ARP request
	CmdFARPReply = 0x55 // FC ARP reply
	CmdRPS       = 0x56 // read port status block
	CmdRPL       = 0x57 // read port list
	CmdRPBC      = 0x58 // read port buffer condition
	CmdFAN       = 0x60 // fabric address notification
	CmdRSCN      = 0x61 // registered state change notification
	CmdSCR       = 0x62 // state change registration
	CmdRNFT      = 0x63 // report node FC-4 types
	CmdCSR       = 0x68 // clock synch. request
	CmdCSU       = 0x69 // clock synch. update
	CmdLInit     = 0x70 // loop initialize
	CmdLSTS      = 0x72 // loop status
	CmdRNID      = 0x78 // request node ID data
	CmdRLIR      = 0x79 // registered link incident report
	CmdLIRR      = 0x7a // link incident record registration
	CmdSRL       = 0x7b // scan remote loop
	CmdSBRP      = 0x7c // set bit-error reporting params
	CmdRPSC      = 0x7d // report speed capabilities
	CmdQSA       = 0x7e // query security attributes
	CmdEVFP      = 0x7f // exchange virt. fabrics params
	CmdLKA       = 0x80 // link keep-alive
	CmdAuthELS   = 0x90 // authentication ELS
)

type Frame struct {
	Command    Command `fc:"@0"`
	RawPayload []byte  `fc:"@1"`
	Payload    interface{}
}

func (f *Frame) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadFromAndPost(r, f)
}

func (f *Frame) PostUnmarshal() error {
	var sf io.ReaderFrom
	switch f.Command {
	case CmdFLOGI:
		sf = &FLOGI{}
	case CmdPLOGI:
		sf = &PLOGI{}
	case CmdPRLI:
		sf = &PRLI{}
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
