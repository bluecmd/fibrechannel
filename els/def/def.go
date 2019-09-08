package main

import (
	"log"
	"os"

	e "github.com/bluecmd/fibrechannel/encoding"
)

func defPLOGI() e.Type {
	plogi := e.NewStruct("PLOGI")
	plogi.Field("", &e.Skip{3 * e.Bytes})

	features := e.NewStruct("PLOGIFeatures")
	features.Field("Data", e.Uint32)

	common := e.NewStruct("PLOGICommonSvcParams")
	common.Field("FCPHVersion", e.Uint16)
	common.Field("B2BCredits", e.Uint16)
	common.Field("Features", features)
	common.Field("MaxConcurrentSeq", e.Uint16)
	common.Field("RelOffsetInfoCat", e.Uint16)
	common.Field("EDTOV", e.Uint32)

	class := e.NewStruct("PLOGIClassSvcParams")
	class.Field("Service", e.Uint16)
	class.Field("Initiator", e.Uint16)
	class.Field("Recipient", e.Uint16)
	class.Field("ReceiveDataFieldSize", &e.Unsigned{12 * e.Bits})
	class.Field("", &e.Skip{1 * e.Bytes})
	class.Field("ConcurrentSeq", e.Uint8)
	class.Field("E2ECredits", e.Uint16)
	class.Field("", &e.Skip{1 * e.Bytes})
	class.Field("OpenSeqPerExch", e.Uint8)
	class.Field("", &e.Skip{2 * e.Bytes})

	plogi.Field("CommonSvcParams", common)
	plogi.Field("PortName", &e.Object{"common.WWN"})
	plogi.Field("NodeName", &e.Object{"common.WWN"})
	plogi.Field("ClassSvcParams", &e.Array{3, class})
	plogi.Field("AuxSvcParams", class)
	plogi.Field("VendorVersion", &e.ByteArray{16})

	return plogi
}

func main() {
	els := e.NewStruct("Frame")

	rctl := &e.Enum{
		Name: "Route",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"RouteSolicited": e.Constant{0x21, "Solicited ELS"},
			"RouteRequest":   e.Constant{0x22, "ELS Request"},
			"RouteReply":     e.Constant{0x23, "ELS Reply"},
		}}

	cmd := &e.Enum{
		Name: "Command",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"CmdLSRJT":     e.Constant{0x01, "ESL reject"},
			"CmdLSACC":     e.Constant{0x02, "ESL Accept"},
			"CmdPLOGI":     e.Constant{0x03, "N_Port login"},
			"CmdFLOGI":     e.Constant{0x04, "F_Port login"},
			"CmdLOGO":      e.Constant{0x05, "Logout"},
			"CmdABTX":      e.Constant{0x06, "Abort exchange - obsolete"},
			"CmdRCS":       e.Constant{0x07, "read connection status"},
			"CmdRES":       e.Constant{0x08, "read exchange status block"},
			"CmdRSS":       e.Constant{0x09, "read sequence status block"},
			"CmdRSI":       e.Constant{0x0a, "read sequence initiative"},
			"CmdESTS":      e.Constant{0x0b, "establish streaming"},
			"CmdESTC":      e.Constant{0x0c, "estimate credit"},
			"CmdADVC":      e.Constant{0x0d, "advise credit"},
			"CmdRTV":       e.Constant{0x0e, "read timeout value"},
			"CmdRLS":       e.Constant{0x0f, "read link error status block"},
			"CmdEcho":      e.Constant{0x10, "echo"},
			"CmdTest":      e.Constant{0x11, "test"},
			"CmdRRQ":       e.Constant{0x12, "reinstate recovery qualifier"},
			"CmdREC":       e.Constant{0x13, "read exchange concise"},
			"CmdSRR":       e.Constant{0x14, "sequence retransmission request"},
			"CmdPRLI":      e.Constant{0x20, "process login"},
			"CmdPRLO":      e.Constant{0x21, "process logout"},
			"CmdSCN":       e.Constant{0x22, "state change notification"},
			"CmdTPLS":      e.Constant{0x23, "test process login state"},
			"CmdTPRLO":     e.Constant{0x24, "third party process logout"},
			"CmdLCLM":      e.Constant{0x25, "login control list mgmt (obs)"},
			"CmdGAID":      e.Constant{0x30, "get alias_ID"},
			"CmdFACT":      e.Constant{0x31, "fabric activate alias_id"},
			"CmdFDACDT":    e.Constant{0x32, "fabric deactivate alias_id"},
			"CmdNACT":      e.Constant{0x33, "N-port activate alias_id"},
			"CmdNDACT":     e.Constant{0x34, "N-port deactivate alias_id"},
			"CmdQOSR":      e.Constant{0x40, "quality of service request"},
			"CmdRVCS":      e.Constant{0x41, "read virtual circuit status"},
			"CmdPDISC":     e.Constant{0x50, "discover N_port service params"},
			"CmdFDISC":     e.Constant{0x51, "discover F_port service params"},
			"CmdADISC":     e.Constant{0x52, "discover address"},
			"CmdRNC":       e.Constant{0x53, "report node cap (obs)"},
			"CmdFARPReq":   e.Constant{0x54, "FC ARP request"},
			"CmdFARPReply": e.Constant{0x55, "FC ARP reply"},
			"CmdRPS":       e.Constant{0x56, "read port status block"},
			"CmdRPL":       e.Constant{0x57, "read port list"},
			"CmdRPBC":      e.Constant{0x58, "read port buffer condition"},
			"CmdFAN":       e.Constant{0x60, "fabric address notification"},
			"CmdRSCN":      e.Constant{0x61, "registered state change notification"},
			"CmdSCR":       e.Constant{0x62, "state change registration"},
			"CmdRNFT":      e.Constant{0x63, "report node FC-4 types"},
			"CmdCSR":       e.Constant{0x68, "clock synch. request"},
			"CmdCSU":       e.Constant{0x69, "clock synch. update"},
			"CmdLInit":     e.Constant{0x70, "loop initialize"},
			"CmdLSTS":      e.Constant{0x72, "loop status"},
			"CmdRNID":      e.Constant{0x78, "request node ID data"},
			"CmdRLIR":      e.Constant{0x79, "registered link incident report"},
			"CmdLIRR":      e.Constant{0x7a, "link incident record registration"},
			"CmdSRL":       e.Constant{0x7b, "scan remote loop"},
			"CmdSBRP":      e.Constant{0x7c, "set bit-error reporting params"},
			"CmdRPSC":      e.Constant{0x7d, "report speed capabilities"},
			"CmdQSA":       e.Constant{0x7e, "query security attributes"},
			"CmdEVFP":      e.Constant{0x7f, "exchange virt. fabrics params"},
			"CmdLKA":       e.Constant{0x80, "link keep-alive"},
			"CmdAuthELS":   e.Constant{0x90, "authentication ELS"},
		}}

	plogi := defPLOGI()

	fcmd := els.Field("cmd", cmd)

	payload := &e.SwitchedType{"Payload", fcmd, map[string]e.Type{
		"CmdPLOGI": plogi,
	}}
	els.Field("Payload", payload)

	imports := []string{
		"github.com/bluecmd/fibrechannel/common",
	}
	b, err := e.Generate("els", imports, els, rctl, plogi)
	if err != nil {
		log.Fatalf("Generate failed: %v", err)
	}
	os.Stdout.Write(b)
}
