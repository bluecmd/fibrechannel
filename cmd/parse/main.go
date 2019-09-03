package main

import (
	"log"
	"os"

	fc "github.com/bluecmd/fibrechannel"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	frm := fc.Frame{}
	_, err := fc.ReadFrom(os.Stdin, &frm)
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}
	spew.Dump(frm)
}
