package main

import (
	"bytes"
	"log"
	"os"

	"github.com/google/gopacket/layers"
	pcap "github.com/google/gopacket/pcapgo"
	gp "github.com/google/gopacket"
)

func main() {
	f, err := os.Create("somefile.pcapng")
	if err != nil {
		log.Fatalf("Create: %v", err)
	}
	defer f.Close()

	// TODO: Change to 225 w/ SOF/EOF
	w, err := pcap.NewNgWriter(f, layers.LinkType(224))
	if err != nil {
		log.Fatalf("NewNgWriter: %v", err)
	}
	defer w.Flush()

	for _, path := range os.Args[1:] {
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			log.Fatalf("Open: %v", err)
		}
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(fr); err != nil {
			log.Fatalf("Read: %v", err)
		}
		ci := gp.CaptureInfo{
			CaptureLength: buf.Len(),
			Length: buf.Len(),
		}
		if err := w.WritePacket(ci, buf.Bytes()); err != nil {
			log.Fatalf("WritePacket: %v", err)
		}
	}
}
