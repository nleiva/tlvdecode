# tlvdecode
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
