meta:
  id: fibrechannel
  file-extension: fc
  license: LGPL-3.0-or-later
  xref:
    doc: FC-FS-4
    incits: 488
  endian: be
  imports:
    - /fc_bls
    - /fc_els
seq:
  - id: r_ctl
    type: routing_control
    size: 1
    doc: |
      Example Documentation of R_CTL
  - id: destination_id
    type: address
    size: 3
  - id: csctl_or_priority
    type: u1
    doc: |
      This field is either CS_Ctl or Priority depending on the value of
      f_ctl.priority_enable
  - id: source_id
    type: address
    size: 3
  - id: type
    type: u1
    enum: fctype_enum
  - id: f_ctl
    type: frame_control
    size: 3
  - id: sequence_id
    type: u1
  - id: df_ctl
    type: data_field_control
  - id: sequence_count
    type: u2
  - id: ox_id
    type: u2
  - id: rx_id
    type: u2
  - id: parameter
    type: parameter
    size: 4
  - id: bls
    size-eos: true
    if: type == fctype_enum::bls
    type: fc_bls
  - id: els
    size-eos: true
    if: type == fctype_enum::els
    type: fc_els
    #  - id: payload
    #    size-eos: true
    #    type:
    #      switch-on: type
    #      cases:
    #        'fctype_enum::bls': fc_bls
    #        'fctype_enum::els': fc_els

enums:
  fctype_enum:
    0x00: bls
    0x01: els
    0x04: llcsnap
    0x05: ip
    0x08: fcp
    0x09: gpp
    0x1b:
      id: sb_to_cu
      doc: "FICON / FC-SB-3: Control Unit -> Channel"
    0x1c:
      id: sb_from_cu
      doc: "FICON / FC-SB-3: Channel -> Control Unit"
    0x20: fcct
    0x22: swils
    0x23: al
    0x24: snmp
    0x28: nvme
    0xee: spinfab
    0xef: diag

types:
  wwn:
    seq:
      - id: data
        size: 8
  address:
    seq:
      - id: data
        size: 3
  parameter:
    seq:
      - id: placeholder
        size: 4
  routing_control:
    seq:
      - id: placeholder
        size: 1
  data_field_control:
    seq:
      - id: placeholder
        size: 1
  frame_control:
    seq:
      - id: exchange_context
        type: b1
      - id: sequence_context
        type: b1
      - id: first_sequence
        type: b1
      - id: last_sequence
        type: b1
      - id: end_sequence
        type: b1
      - id: reserved5
        type: b1
      - id: priority_enable
        type: b1
      - id: sequence_initative
        type: b1
      - id: reserved4
        type: b1
      - id: reserved3
        type: b1
      - id: ack_form
        type: b2
      - id: reserved2
        type: b1
      - id: reserved1
        type: b1
      - doc: Obsolete
        id: retransmitted_sequence
        type: b1
      - doc: Obsolete
        id: unidirectional_transmit
        type: b1
      - doc: Obsolete
        id: continue_sequence_condition
        type: b2
      - id: abort_sequence_condition
        type: b2
      - id: relative_offset_present
        type: b1
      - id: exchange_reassembly
        type: b1
      - id: fill_bytes
        type: b2

# vim: set ft=yaml:ts=2:sw=2
