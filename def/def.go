package main

import (
	"log"
	"os"

	e "github.com/bluecmd/fibrechannel/encoding"
)

func main() {
	fc := e.NewStruct("Frame")

	sof := &e.Enum{
		Name: "SOF",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"SOFf":  e.Constant{0x1, "TODO"},
			"SOFi4": e.Constant{0x2, "TODO"},
			"SOFi2": e.Constant{0x3, "TODO"},
			"SOFi3": e.Constant{0x4, "TODO"},
			"SOFn4": e.Constant{0x5, "TODO"},
			"SOFn2": e.Constant{0x6, "TODO"},
			"SOFn3": e.Constant{0x7, "TODO"},
			"SOFc4": e.Constant{0x8, "TODO"},
		}}

	eof := &e.Enum{
		Name: "EOF",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"EOFn":   e.Constant{0x1, "TODO"},
			"EOFt":   e.Constant{0x2, "TODO"},
			"EOFrt":  e.Constant{0x3, "TODO"},
			"EOFdt":  e.Constant{0x4, "TODO"},
			"EOFni":  e.Constant{0x5, "TODO"},
			"EOFdti": e.Constant{0x6, "TODO"},
			"EOFrti": e.Constant{0x7, "TODO"},
			"EOFa":   e.Constant{0x8, "TODO"},
		}}

	t := &e.Enum{
		Name: "Type",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"TypeBLS      ": e.Constant{0x00, "TODO"},
			"TypeELS      ": e.Constant{0x01, "TODO"},
			"TypeLLCSNAP  ": e.Constant{0x04, "TODO"},
			"TypeIP       ": e.Constant{0x05, "TODO"},
			"TypeFCP      ": e.Constant{0x08, "TODO"},
			"TypeGPP      ": e.Constant{0x09, "TODO"},
			"TypeSBToCU   ": e.Constant{0x1B, "FICON / FC-SB-3: Control Unit -> Channel"},
			"TypeSBFromCU ": e.Constant{0x1C, "FICON / FC-SB-3: Channel -> Control Unit"},
			"TypeFCCT     ": e.Constant{0x20, "TODO"},
			"TypeSWILS    ": e.Constant{0x22, "TODO"},
			"TypeAL       ": e.Constant{0x23, "TODO"},
			"TypeSNMP     ": e.Constant{0x24, "TODO"},
			"TypeNVME     ": e.Constant{0x28, "TODO"},
			"TypeSPINFAB  ": e.Constant{0xEE, "TODO"},
			"TypeDIAG     ": e.Constant{0xEF, "TODO"},
		}}

	fctl := e.NewBitStruct("FrameControl")
	fctl.IntField("TODO1", 6)
	prioen := fctl.BoolBit("PriorityEnable")
	fctl.IntField("TODO2", 17)

	fc.Field("RCtl", e.Uint8)
	// Address for source/destination Nx_Ports
	// Each Nx_Port shall have a native N_Port_ID that is unique within the
	// address domain of a Fabric.  An N_Port_ID of binary zeros indicates that
	// an Nx_Port is unidentified. When a PN_Port completes Link Initialization,
	// it shall be unidentified (i.e., it shall have a single Nx_Port for which
	// the N_Port_ID is 00 00 00h).
	fc.Field("DestinationID", &e.ByteArray{3})

	csctl := e.NewStruct("CSCtl")
	csctl.Field("Data", e.Uint8)
	prio := e.NewStruct("Prio")
	prio.Field("Data", e.Uint8)

	csctlPrio := &e.SwitchedType{"CsctlPriority", 1 * e.Bytes, prioen, map[string]e.Type{
		"false": csctl,
		"true":  prio,
	}}
	fc.Field("CsctlPriority", csctlPrio)

	fc.Field("SourceID", &e.ByteArray{3})

	ftype := fc.Field("fcType", t)

	fc.Field("FCtl", fctl)

	fc.Field("SeqID", e.Uint8)
	fc.Field("DFCtl", e.Uint8)
	fc.Field("SeqCount", e.Uint16)

	// Originator Exchange_ID
	// If the Originator is enforcing uniqueness via the OX_ID mechanism, it
	// shall set a unique value for OX_ID other than FF FFh in the first Data
	// frame of the first Sequence of an Exchange. An OX_ID of FF FFh indicates
	// that the OX_ID is unassigned and that the Originator is not enforcing
	// uniqueness via the OX_ID mechanism. If an Originator uses the unassigned
	// value of FF FFh to identify the Exchange, it shall have only one Exchange
	// (OX_ID set to FF FFh) with a given Responder.
	fc.Field("OXID", e.Uint16)

	// Responder Exchange_ID
	// An RX_ID of FF FFh shall indicate that the RX_ID is unassigned. If the
	// Responder does not assign an RX_ID other than FF FFh by the end of the
	// first Sequence, then the Responder is not enforcing uniqueness via the
	// RX_ID mechanism.
	fc.Field("RXID", e.Uint16)

	fc.Field("Parameters", &e.ByteArray{4})

	payload := &e.SwitchedType{"Payload", e.RemainingBytes, ftype, map[string]e.Type{
		"TypeELS": &e.Object{"els.Frame"},
	}}
	fc.Field("Payload", payload)

	imports := []string{
		"github.com/bluecmd/fibrechannel/els",
	}
	b, err := e.Generate("fibrechannel", imports, fc, sof, eof)
	if err != nil {
		os.Stdout.Write(b)
		log.Fatalf("Generate failed: %v", err)
	}
	os.Stdout.Write(b)
}
