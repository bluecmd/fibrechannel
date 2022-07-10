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
			"SOFf":  {Value: 0x1, Comment: "TODO"},
			"SOFi4": {Value: 0x1, Comment: "TODO"},
			"SOFi2": {Value: 0x1, Comment: "TODO"},
			"SOFi3": {Value: 0x1, Comment: "TODO"},
			"SOFn4": {Value: 0x1, Comment: "TODO"},
			"SOFn2": {Value: 0x1, Comment: "TODO"},
			"SOFn3": {Value: 0x1, Comment: "TODO"},
			"SOFc4": {Value: 0x1, Comment: "TODO"},
		}}

	eof := &e.Enum{
		Name: "EOF",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"EOFn":   {Value: 0x1, Comment: "TODO"},
			"EOFt":   {Value: 0x1, Comment: "TODO"},
			"EOFrt":  {Value: 0x1, Comment: "TODO"},
			"EOFdt":  {Value: 0x1, Comment: "TODO"},
			"EOFni":  {Value: 0x1, Comment: "TODO"},
			"EOFdti": {Value: 0x1, Comment: "TODO"},
			"EOFrti": {Value: 0x1, Comment: "TODO"},
			"EOFa":   {Value: 0x1, Comment: "TODO"},
		}}

	t := &e.Enum{
		Name: "Type",
		Size: 1 * e.Bytes,
		Values: map[string]e.Constant{
			"TypeBLS":      {Value: 0x1, Comment: "TODO"},
			"TypeELS":      {Value: 0x1, Comment: "TODO"},
			"TypeLLCSNAP":  {Value: 0x1, Comment: "TODO"},
			"TypeIP":       {Value: 0x1, Comment: "TODO"},
			"TypeFCP":      {Value: 0x1, Comment: "TODO"},
			"TypeGPP":      {Value: 0x1, Comment: "TODO"},
			"TypeSBToCU":   {Value: 0x1B, Comment: "FICON / FC-SB-3: Control Unit -> Channel"},
			"TypeSBFromCU": {Value: 0x1C, Comment: "FICON / FC-SB-3: Channel -> Control Unit"},
			"TypeFCCT":     {Value: 0x1, Comment: "TODO"},
			"TypeSWILS":    {Value: 0x1, Comment: "TODO"},
			"TypeAL":       {Value: 0x1, Comment: "TODO"},
			"TypeSNMP":     {Value: 0x1, Comment: "TODO"},
			"TypeNVME":     {Value: 0x1, Comment: "TODO"},
			"TypeSPINFAB":  {Value: 0x1, Comment: "TODO"},
			"TypeDIAG":     {Value: 0x1, Comment: "TODO"},
		}}

	fctl := e.NewBitStruct("FrameControl")
	fctl.IntField("TODO1", 6)
	prioen := fctl.BoolBit("PriorityEnable")
	fctl.IntField("TODO2", 17)

	fc.Field("RCtl", e.Uint8)
	// Address for source/destination Nx_Ports
	// Each Nx_Port shall have a native N_Port_ID that is unique within the
	// address domain of a Fabric. An N_Port_ID of binary zeros indicates that
	// an Nx_Port is unidentified. When a PN_Port completes Link Initialization,
	// it shall be unidentified (i.e., it shall have a single Nx_Port for which
	// the N_Port_ID is 00 00 00h).
	fc.Field("DestinationID", &e.ByteArray{Count: 3})

	csctl := e.NewStruct("CSCtl")
	csctl.Field("Data", e.Uint8)
	prio := e.NewStruct("Prio")
	prio.Field("Data", e.Uint8)

	csctlPrio := &e.SwitchedType{
		Name:       "CsctlPriority",
		Size:       1 * e.Bytes,
		SwitchedOn: prioen,
		Cases: map[string]e.Type{
			"false": csctl,
			"true":  prio,
		}}
	fc.Field("CsctlPriority", csctlPrio)

	fc.Field("SourceID", &e.ByteArray{Count: 3})

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

	fc.Field("Parameters", &e.ByteArray{Count: 4})

	payload := &e.SwitchedType{
		Name:       "Payload",
		Size:       e.RemainingBytes,
		SwitchedOn: ftype,
		Cases: map[string]e.Type{
			"TypeELS": &e.Object{Class: "els.Frame"},
		}}
	fc.Field("Payload", payload)

	imports := []string{
		"github.com/bluecmd/fibrechannel/els",
	}
	b, err := e.Generate("fibrechannel", imports, fc, sof, eof)
	if err != nil {
		log.Fatalf("Generate failed: %v", err)
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
