package main

import (
	"log"
	"os"

	fc "github.com/bluecmd/fibrechannel"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	frm := fc.Frame{}
	_, err := frm.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}
	s := spew.NewDefaultConfig()
	s.DisablePointerAddresses = true
	s.Dump(&frm)
}
