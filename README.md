# Decoding IS-IS Link State PDUs (LSP)

Decode IS-IS base64 encoded TLV's


## Code Examples

`tlvdecode` reads from a base64 encoded [file](data64) and outputs the LSP details

```console
$ ./tlvdecode 
===== LSP Details (lenght: 226) ====
LSPID:      0151.0250.0002.0000-0000
Seq Num:    0x000d
Checksum:   0x6f49
Type Block: 0x03
===== TLV Details (total: 008) ====
Type010,  L017: 0x36ab31873913cbaa3195600648913cc619
Type001,  L006: 49.0000.0162
Type129,  L001: 0x8e
Type229,  L002: 0x0002
Type137,  L022: mrstn-5502-2.cisco.com
Type232,  L016: 2001:558:2::2
Type222,  L013: 0x00020151025000010000000a00
Type237,  L118:
```

## IS-IS Theory


### TLV 1 (Area Address)

Includes the Area Addresses to which the Intermediate System is connected.

### TLV 10 (Authentication)

The information that is used to authenticate the PDU.

### TLV 129 (Protocols Supported)

Carries the Network Layer Protocol Identifiers (NLPID) for Network Layer protocols that the IS (Intermediate System) is capable. It refers to the Data Protocols that are supported. 
- IPv6 NLPID is 142 (0x8E)
- IPv4 NLPID is 204 (0xCC)
- CLNS NLPID is 129 (0x81)

### TLV 137 (Dynamic Hostname)

Identifies the symbolic name of the router originating the link-state packet (LSP).

### TLV 229 (Multi-Topology)

It contains one or more MTs the router is participating.

```
  +--------------------------------+
  |O |A |R |R |        MT ID       |      2
  +--------------------------------+
```

### TLV 232 (IPv6 Interface Address)

Maps directly to "IP Interface Address" TLV in [RFC1195](https://tools.ietf.org/html/rfc1195).

```
   0                   1                   2                   3
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Type = 232   |    Length     |   Interface Address 1(*) ..   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                  .. Interface Address 1(*) ..                 |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                  .. Interface Address 1(*) ..                 |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                  .. Interface Address 1(*) ..                 |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |   Interface Address 1(*) ..   |   Interface Address 2(*) ..
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   * - if present
```

### TLV 236 (IPv6 Reachability)

The "IPv6 Reachability" TLV describes network reachability through the specification of a routing prefix, metric information, a bit to indicate if the prefix is being advertised down from a higher level, a bit to indicate if the prefix is being distributed from another routing protocol, and OPTIONALLY the existence of Sub-TLVs to allow for later extension

```
   0                   1                   2                   3
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Type = 236   |    Length     |          Metric ..            |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          .. Metric            |U|X|S| Reserve |  Prefix Len   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Prefix ...
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Sub-TLV Len(*) | Sub-TLVs(*) ...
   * - if present

   U - up/down bit
   X - external original bit
   S - subtlv present bit
```

### TLV 237 (Multi-Topology Reachable IPv6 Prefixes TLV)

It is aligned with IPv6 Reachability TLV type 236 beside an additional two bytes in front

```
  +--------------------------------+
  |R |R |R |R |        MT ID       |      2
  +--------------------------------+
  | IPv6 Reachability format       |    6 - 253
  +--------------------------------+
  .                                .
  +--------------------------------+
  | IPv6 Reachability format       |    6 - 253
  +--------------------------------+
```

### IPv4/IPv6 Extended Reachability Attribute Flags

This sub-TLV supports the advertisement of additional flags associated with a given prefix advertisement

```console
   0 1 2 3 4 5 6 7...
  +-+-+-+-+-+-+-+-+...
  |X|R|N|          ...
  +-+-+-+-+-+-+-+-+...

X-Flag:  External Prefix Flag (Bit 0)
R-Flag:  Re-advertisement Flag (Bit 1)
N-flag:  Node Flag (Bit 2)
```


## Links

- [Intermediate System-to-Intermediate System (IS-IS) TLVs](http://www.cisco.com/c/en/us/support/docs/ip/integrated-intermediate-system-to-intermediate-system-is-is/5739-tlvs-5739.html)
- [IS-IS TLV Codepoints](https://www.iana.org/assignments/isis-tlv-codepoints/isis-tlv-codepoints.xhtml)
- [Routing IPv6 with IS-IS](https://tools.ietf.org/html/rfc5308)