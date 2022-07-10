package main

import (
	"log"
	"os"

	. "github.com/bluecmd/fibrechannel/encoding"
)

func defPLOGI() Type {
	plogi := NewStruct("PLOGI")
	plogi.Field("", &Skip{Size: 3 * Bytes})

	common := NewBitStruct("PLOGICommonSvcParams")

	// Word 0
	common.IntField("FCPHVersion", 16) // 31-16
	common.IntField("B2BCredits", 16)  // 15-0

	// Word 1
	common.BoolBit("ContIncrRelOffset")       // 31
	common.BoolBit("RandomRelOffset")         // 30
	common.BoolBit("ValidVendorVersionLevel") // 29
	// N_Port/F_Port=0 for an N_Port, and N_Port/F_Port=1 for an F_Port
	common.BoolBit("NorFPort") // 28
	// BB_Credit Management=0 for an N_Port or F_Port, BB_Credit_Management=1 for an L_Port
	common.BoolBit("BBCreditMgmt")               // 27
	common.BoolBit("EDTOVResolution")            // 26
	common.BoolBit("EnergyEffLPIModeSupported")  // 25
	common.SkipBit(1)                            // 24
	common.BoolBit("PriorityTaggingSupported")   // 23
	common.BoolBit("QueryDataBufferCond")        // 22
	common.BoolBit("SecurityBit")                // 21
	common.BoolBit("ClockSyncPrimitiveCapable")  // 20
	common.BoolBit("RTTOVValue")                 // 19
	common.BoolBit("DynamicHalfDuplexSupported") // 18
	common.BoolBit("SeqCntVendorSpec")           // 17
	common.BoolBit("PayloadBit")                 // 16
	common.IntField("BBSCN", 4)                  // 15-12
	common.IntField("B2BRecvDataFieldSize", 12)  // 11-0

	// Word 2
	common.SkipBit(5)                              // 31-27
	common.BoolBit("AppHdrSupport")                // 27
	common.SkipBit(2)                              // 25-24
	common.IntField("NxPortTotalConcurrentSeq", 8) // 23-16
	common.IntField("RelOffsetInfoCat", 16)        // 15-0

	// Word 3
	common.IntField("EDTOV", 32)

	class := NewStruct("PLOGIClassSvcParams")
	class.Field("Service", Uint16)
	class.Field("Initiator", Uint16)
	class.Field("Recipient", Uint16)
	class.Field("ReceiveDataFieldSize", &Unsigned{Size: 12 * Bits})
	class.Field("", &Skip{Size: 1 * Bytes})
	class.Field("ConcurrentSeq", Uint8)
	class.Field("E2ECredits", Uint16)
	class.Field("", &Skip{Size: 1 * Bytes})
	class.Field("OpenSeqPerExch", Uint8)
	class.Field("", &Skip{Size: 2 * Bytes})

	plogi.Field("CommonSvcParams", common)
	plogi.Field("PortName", &Object{Class: "common.WWN"})
	plogi.Field("NodeName", &Object{Class: "common.WWN"})
	plogi.Field("ClassSvcParams", &Array{Count: 3, Type: class})
	plogi.Field("AuxSvcParams", class)
	plogi.Field("VendorVersion", &ByteArray{Count: 16})

	return plogi
}

func main() {
	els := NewStruct("Frame")

	rctl := &Enum{
		Name: "Route",
		Size: 1 * Bytes,
		Values: map[string]Constant{
			"RouteSolicited": {Value: 0x1, Comment: "Solicited ELS"},
			"RouteRequest":   {Value: 0x1, Comment: "ELS Request"},
			"RouteReply":     {Value: 0x1, Comment: "ELS Reply"},
		}}

	cmd := &Enum{
		Name: "Command",
		Size: 1 * Bytes,
		Values: map[string]Constant{
			"CmdLSRJT":     {Value: 0x1, Comment: "ESL reject"},
			"CmdLSACC":     {Value: 0x1, Comment: "ESL Accept"},
			"CmdPLOGI":     {Value: 0x1, Comment: "N_Port login"},
			"CmdFLOGI":     {Value: 0x1, Comment: "F_Port login"},
			"CmdLOGO":      {Value: 0x1, Comment: "Logout"},
			"CmdABTX":      {Value: 0x1, Comment: "Abort exchange - obsolete"},
			"CmdRCS":       {Value: 0x1, Comment: "read connection status"},
			"CmdRES":       {Value: 0x1, Comment: "read exchange status block"},
			"CmdRSS":       {Value: 0x1, Comment: "read sequence status block"},
			"CmdRSI":       {Value: 0x1, Comment: "read sequence initiative"},
			"CmdESTS":      {Value: 0x1, Comment: "establish streaming"},
			"CmdESTC":      {Value: 0x1, Comment: "estimate credit"},
			"CmdADVC":      {Value: 0x1, Comment: "advise credit"},
			"CmdRTV":       {Value: 0x1, Comment: "read timeout value"},
			"CmdRLS":       {Value: 0x1, Comment: "read link error status block"},
			"CmdEcho":      {Value: 0x1, Comment: "echo"},
			"CmdTest":      {Value: 0x1, Comment: "test"},
			"CmdRRQ":       {Value: 0x1, Comment: "reinstate recovery qualifier"},
			"CmdREC":       {Value: 0x1, Comment: "read exchange concise"},
			"CmdSRR":       {Value: 0x1, Comment: "sequence retransmission request"},
			"CmdPRLI":      {Value: 0x1, Comment: "process login"},
			"CmdPRLO":      {Value: 0x1, Comment: "process logout"},
			"CmdSCN":       {Value: 0x1, Comment: "state change notification"},
			"CmdTPLS":      {Value: 0x1, Comment: "test process login state"},
			"CmdTPRLO":     {Value: 0x1, Comment: "third party process logout"},
			"CmdLCLM":      {Value: 0x1, Comment: "login control list mgmt (obs)"},
			"CmdGAID":      {Value: 0x1, Comment: "get alias_ID"},
			"CmdFACT":      {Value: 0x1, Comment: "fabric activate alias_id"},
			"CmdFDACDT":    {Value: 0x1, Comment: "fabric deactivate alias_id"},
			"CmdNACT":      {Value: 0x1, Comment: "N-port activate alias_id"},
			"CmdNDACT":     {Value: 0x1, Comment: "N-port deactivate alias_id"},
			"CmdQOSR":      {Value: 0x1, Comment: "quality of service request"},
			"CmdRVCS":      {Value: 0x1, Comment: "read virtual circuit status"},
			"CmdPDISC":     {Value: 0x1, Comment: "discover N_port service params"},
			"CmdFDISC":     {Value: 0x1, Comment: "discover F_port service params"},
			"CmdADISC":     {Value: 0x1, Comment: "discover address"},
			"CmdRNC":       {Value: 0x1, Comment: "report node cap (obs)"},
			"CmdFARPReq":   {Value: 0x1, Comment: "FC ARP request"},
			"CmdFARPReply": {Value: 0x1, Comment: "FC ARP reply"},
			"CmdRPS":       {Value: 0x1, Comment: "read port status block"},
			"CmdRPL":       {Value: 0x1, Comment: "read port list"},
			"CmdRPBC":      {Value: 0x1, Comment: "read port buffer condition"},
			"CmdFAN":       {Value: 0x1, Comment: "fabric address notification"},
			"CmdRSCN":      {Value: 0x1, Comment: "registered state change notification"},
			"CmdSCR":       {Value: 0x1, Comment: "state change registration"},
			"CmdRNFT":      {Value: 0x1, Comment: "report node FC-4 types"},
			"CmdCSR":       {Value: 0x1, Comment: "clock synch. request"},
			"CmdCSU":       {Value: 0x1, Comment: "clock synch. update"},
			"CmdLInit":     {Value: 0x1, Comment: "loop initialize"},
			"CmdLSTS":      {Value: 0x1, Comment: "loop status"},
			"CmdRNID":      {Value: 0x1, Comment: "request node ID data"},
			"CmdRLIR":      {Value: 0x1, Comment: "registered link incident report"},
			"CmdLIRR":      {Value: 0x1, Comment: "link incident record registration"},
			"CmdSRL":       {Value: 0x1, Comment: "scan remote loop"},
			"CmdSBRP":      {Value: 0x1, Comment: "set bit-error reporting params"},
			"CmdRPSC":      {Value: 0x1, Comment: "report speed capabilities"},
			"CmdQSA":       {Value: 0x1, Comment: "query security attributes"},
			"CmdEVFP":      {Value: 0x1, Comment: "exchange virt. fabrics params"},
			"CmdLKA":       {Value: 0x1, Comment: "link keep-alive"},
			"CmdAuthELS":   {Value: 0x1, Comment: "authentication ELS"},
		}}

	plogi := defPLOGI()

	fcmd := els.Field("cmd", cmd)

	var payload = &SwitchedType{
		Name:       "Payload",
		Size:       RemainingBytes,
		SwitchedOn: fcmd,
		Cases: map[string]Type{
			"CmdPLOGI": plogi,
		},
	}
	els.Field("Payload", payload)

	imports := []string{
		"github.com/bluecmd/fibrechannel/common",
	}
	b, err := Generate("els", imports, els, rctl, plogi)
	if err != nil {
		log.Fatalf("Generate failed: %v", err)
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
