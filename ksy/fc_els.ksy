meta:
  id: fc_els
  file-extension: fc
  license: LGPL-3.0-or-later
  xref:
    doc: FC-LS-3
    incits: 487
  endian: be
seq:
  - id: command
    type: u1
    enum: elscmd_enum
  - id: payload
    size-eos: true
    type:
      switch-on: command
      cases:
        'elscmd_enum::prli': els_prli
enums:
  elscmd_enum:
    0x01:
      id: lsrjt
      doc: "ESL reject"
    0x02:
      id: lsacc
      doc: "ESL Accept"
    0x03:
      id: plogi
      doc: "N_Port login"
    0x04:
      id: flogi
      doc: "F_Port login"
    0x05:
      id: logo
      doc: "Logout"
    0x06:
      id: abtx
      doc: "Abort exchange - obsolete"
    0x07:
      id: rcs
      doc: "read connection status"
    0x08:
      id: res
      doc: "read exchange status block"
    0x09:
      id: rss
      doc: "read sequence status block"
    0x0a:
      id: rsi
      doc: "read sequence initiative"
    0x0b:
      id: ests
      doc: "establish streaming"
    0x0c:
      id: estc
      doc: "estimate credit"
    0x0d:
      id: advc
      doc: "advise credit"
    0x0e:
      id: rtv
      doc: "read timeout value"
    0x0f:
      id: rls
      doc: "read link error status block"
    0x10:
      id: echo
      doc: "echo"
    0x11:
      id: test
      doc: "test"
    0x12:
      id: rrq
      doc: "reinstate recovery qualifier"
    0x13:
      id: rec
      doc: "read exchange concise"
    0x14:
      id: srr
      doc: "sequence retransmission request"
    0x20:
      id: prli
      doc: "process login"
    0x21:
      id: prlo
      doc: "process logout"
    0x22:
      id: scn
      doc: "state change notification"
    0x23:
      id: tpls
      doc: "test process login state"
    0x24:
      id: tprlo
      doc: "third party process logout"
    0x25:
      id: lclm
      doc: "login control list mgmt (obs)"
    0x30:
      id: gaid
      doc: "get alias_ID"
    0x31:
      id: fact
      doc: "fabric activate alias_id"
    0x32:
      id: fdacdt
      doc: "fabric deactivate alias_id"
    0x33:
      id: nact
      doc: "N-port activate alias_id"
    0x34:
      id: ndact
      doc: "N-port deactivate alias_id"
    0x40:
      id: qosr
      doc: "quality of service request"
    0x41:
      id: rvcs
      doc: "read virtual circuit status"
    0x50:
      id: pdisc
      doc: "discover N_port service params"
    0x51:
      id: fdisc
      doc: "discover F_port service params"
    0x52:
      id: adisc
      doc: "discover address"
    0x53:
      id: rnc
      doc: "report node cap (obs)"
    0x54:
      id: farpreq
      doc: "FC ARP request"
    0x55:
      id: farpreply
      doc: "FC ARP reply"
    0x56:
      id: rps
      doc: "read port status block"
    0x57:
      id: rpl
      doc: "read port list"
    0x58:
      id: rpbc
      doc: "read port buffer condition"
    0x60:
      id: fan
      doc: "fabric address notification"
    0x61:
      id: rscn
      doc: "registered state change notification"
    0x62:
      id: scr
      doc: "state change registration"
    0x63:
      id: rnft
      doc: "report node FC-4 types"
    0x68:
      id: csr
      doc: "clock synch. request"
    0x69:
      id: csu
      doc: "clock synch. update"
    0x70:
      id: linit
      doc: "loop initialize"
    0x72:
      id: lsts
      doc: "loop status"
    0x78:
      id: rnid
      doc: "request node ID data"
    0x79:
      id: rlir
      doc: "registered link incident report"
    0x7a:
      id: lirr
      doc: "link incident record registration"
    0x7b:
      id: srl
      doc: "scan remote loop"
    0x7c:
      id: sbrp
      doc: "set bit-error reporting params"
    0x7d:
      id: rpsc
      doc: "report speed capabilities"
    0x7e:
      id: qsa
      doc: "query security attributes"
    0x7f:
      id: evfp
      doc: "exchange virt. fabrics params"
    0x80:
      id: lka
      doc: "link keep-alive"
    0x90:
      id: authels
      doc: "authentication ELS"
types:
  els_prli:
    seq:
      - id: page_length
        type: u1
      - id: payload_length
        type: u2
      - id: service_parameters
        size: 4
        repeat: expr 
        repeat-expr: page_length / 4