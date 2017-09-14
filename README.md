# Decoding IS-IS Link State PDUs (LSP)

From [Package json](https://golang.org/pkg/encoding/json/) we learn _array and slice values encode as JSON arrays_, however _[]byte encodes as a **base64-encoded** string_. Therefore you might end up with a base64-encoded string in your JSON outputs. That was our case when reading IS-IS PDU's, therefore we used this code to translate it to a human readable format.

Disclaimer: At this point in time, it does not cover all the exiting TLV's, just those we found in our network. We will keep adding more on-demand.

## Use

`tlvdecode` reads either from a base64 encoded [file](input/data64) or Telemetry IOS XR [message](input/full1.json) and outputs the LSP details.

### Raw base64 LSP data

- **Input**

```console
$cat input/data64
AVECUAACAAAAAAANb0kDChE2qzGHORPLqjGVYAZIkTzGGQEGBUkAAAFigQGO5QIAAokWbXJzdG4tNTUwMi0yLmNpc2NvLmNvb
<snip>
```

- **Output**

```console
$ ./tlvdecode data64 
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
Type222,  L013: Neighbor System ID: 0151.0250.0001.00, Metric: 10
Type237,  L118: MT ID: IPv6 Unicast
Prefixes:
2001:558:2::2/128, Metric:1
2001:f00:ba::/64, Metric:10
2001:f00:bb::/64, Metric:10
2001:f00:bc::/64, Metric:10
2001:f00:bd::/64, Metric:65000
2001:f00:be::/64, Metric:65000
```

### Base64 LSP data inside a Telemetry message

- **Input**

```json
{ "Source": "[2001:420:2cff:1204::5502:1]:36597",
  "Telemetry": {...},
    "Rows": [
        {
            "Timestamp": 1495073881691,
            "Keys": {...},
            "Content": {
                "lsp_header_data": {...},
                "lsp_body": "AVECUAABA..."
            }
        },
        {
            "Timestamp": 1495073881691,
            "Keys": {...},
            "Content": {
                "lsp_header_data": {...},
                "lsp_body": "AVECUAACAA..."
            }
        }
    ]
}
```

- **Output**

```console
$./tlvdecode full1.json
===== LSP Details (lenght: 190) ====
LSPID:      0151.0250.0001.0000-0000
Seq Num:    0x01a1
Checksum:   0x985f
Type Block: 0x03
===== TLV Details (total: 008) ====
Type010,  L017: 0x365175ce39d3977ec4d48de0f6569df972
Type001,  L006: 49.0000.0162
<snip>
===== LSP Details (lenght: 226) ====
LSPID:      0151.0250.0002.0000-0000
Seq Num:    0x000d
Checksum:   0x6f49
Type Block: 0x03
===== TLV Details (total: 008) ====
Type010,  L017: 0x36ab31873913cbaa3195600648913cc619
Type001,  L006: 49.0000.0162
<snip>
```

## IS-IS Theory

### TLV 1 (Area Address)

Includes the Area Addresses to which the Intermediate System is connected.

### TLV 10 (Authentication)

The information that is used to authenticate the PDU.

### TLV 22 (Extended IS Reachability)

The original IS reachability (TLV type 2, defined in [ISO 10589](https://www.iso.org/standard/30932.html) contains information about a series of IS neighbors. The extended IS reachability TLV proposed on [RFC 3784](https://tools.ietf.org/html/rfc3784) contains a new data structure, consisting of:

```
7 octets of system Id and pseudonode number
3 octets of default metric
1 octet of length of sub-TLVs
0-244 octets of sub-TLVs,
   where each sub-TLV consists of a sequence of
		1 octet of sub-type
		1 octet of length of the value field of the sub-TLV
		0-242 octets of value
```

### TLV 129 (Protocols Supported)

Carries the Network Layer Protocol Identifiers (NLPID) for Network Layer protocols that the IS (Intermediate System) is capable. It refers to the Data Protocols that are supported. 
- IPv6 NLPID is 142 (0x8E)
- IPv4 NLPID is 204 (0xCC)
- CLNS NLPID is 129 (0x81)

### TLV 137 (Dynamic Hostname)

Identifies the symbolic name of the router originating the Link State PDU.

### TLV 140 (IPv6 TE Router ID)

The IPv6 TE Router ID TLV contains a 16-octet IPv6 address. A stable global IPv6 address MUST be used, so that the router ID provides a routable address, regardless of the state of a node's interfaces.

### TLV 222 (MT Intermediate Systems)

It is aligned with extended IS reachability TLV type 22 beside an additional two bytes in front at the beginning of the TLV. After the 2-byte MT membership format, the MT IS content is in the same format as extended IS TLV, type 22

```
 +--------------------------------+
 |R |R |R |R |        MT ID       |      2
 +--------------------------------+
 | extended IS TLV format         |    11 - 253
 +--------------------------------+
 .                                .
 .                                .
 +--------------------------------+
 | extended IS TLV format         |    11 - 253
 +--------------------------------+
```

### TLV 229 (Multi-Topology)

It contains one or more MTs the router is participating.

```
  +--------------------------------+
  |O |A |R |R |        MT ID       |      2
  +--------------------------------+

Bit O represents the OVERLOAD bit for the MT
Bit A represents the ATTACH bit for the MT
Bits R are reserved
```

#### Reserved MT ID Values

-  MT ID #0: Equivalent to the "standard" topology.
-  MT ID #1: Reserved for IPv4 in-band management purposes.
-  MT ID #2: Reserved for IPv6 routing topology.
-  MT ID #3: Reserved for IPv4 multicast routing topology.
-  MT ID #4: Reserved for IPv6 multicast routing topology.
-  MT ID #5: Reserved for IPv6 in-band management purposes.
-  MT ID #6-#4095: Reserved.

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

### TLV 237 (Multi-Topology Reachable IPv6 Prefixes)

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

```
   0 1 2 3 4 5 6 7...
  +-+-+-+-+-+-+-+-+...
  |X|R|N|          ...
  +-+-+-+-+-+-+-+-+...

X-Flag:  External Prefix Flag (Bit 0)
R-Flag:  Re-advertisement Flag (Bit 1)
N-flag:  Node Flag (Bit 2)
```

### OSI Sytem ID in IOS XR 

Described as a string in [Cisco-IOS-XR-types](https://github.com/YangModels/yang/blob/master/vendor/cisco/xr/621/Cisco-IOS-XR-types.yang#L261)

```
  typedef Osi-area-address {
    type string {
      pattern '[a-fA-F0-9]{2}(\.[a-fA-F0-9]{4}){0,6}';
    }
    description "An OSI area address should consist of an odd number
                 of octets, and be of the form 01 or 01.2345 etc up
                 to 01.2345.6789.abcd.ef01.2345.6789. This data type
                 restricts each character to a hex character.";
  }
```


## Links

- [Intermediate System-to-Intermediate System (IS-IS) TLVs](http://www.cisco.com/c/en/us/support/docs/ip/integrated-intermediate-system-to-intermediate-system-is-is/5739-tlvs-5739.html)
- [IS-IS TLV Codepoints](https://www.iana.org/assignments/isis-tlv-codepoints/isis-tlv-codepoints.xhtml)
- RFC 1195: [Use of OSI IS-IS for Routing in TCP/IP and Dual Environments](https://tools.ietf.org/html/rfc1195)
- RFC 3784: [IS-IS Extensions for Traffic Engineering (TE)](https://tools.ietf.org/html/rfc3784)
- RFC 5120: [Multi Topology (MT) Routing in IS-ISs](https://tools.ietf.org/html/rfc5120)
- RFC 5308: [Routing IPv6 with IS-IS](https://tools.ietf.org/html/rfc5308)
- RFC 6119: [IPv6 Traffic Engineering in IS-IS](https://tools.ietf.org/html/rfc6119)