package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/nleiva/tlv"
	"github.com/pkg/errors"
)

func readBytes(src []byte, d *entry, b *bool) error {
	// Decode base64 data into bytes
	var enc = base64.NewEncoding(encodeStd)
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	if err != nil {
		return errors.Wrap(err, "error decoding file")
	}
	// Read PDU info
	r, err := readHeader(dst, n, d, b)
	if err != nil {
		return errors.Wrap(err, "error reading header")
	}
	// DEBUG
	// fmt.Printf("%X", r)
	// Read the TLV's from the Reader and put them on a slice
	tlvs, err := tlv.Read(r)
	if err != nil {
		return errors.Wrap(err, "failed to read TLVs: ")
	}
	// Print the TLV details
	if !*b {
		fmt.Printf("===== TLV Details (total: %03d) ====\n", tlvs.Length())
	}
	err = exploreTLV(tlvs.GetThemAll(), d, b)
	if err != nil {
		return errors.Wrap(err, "failed to read TLV details: ")
	}
	return nil
}

func read222(t []byte, d *entry, b *bool) (mtid string, err error) {
	if len(t) < 13 {
		return mtid, fmt.Errorf("not a valid TLV, lenght: %v", len(t))
	}
	var mt uint16
	var subtlv uint8
	var nsel uint8
	s := make([]byte, 6)
	m := make([]byte, 3)

	buf := bytes.NewReader(t)
	err = binary.Read(buf, binary.BigEndian, &mt)
	check(err, "failed to read MT ID: ")
	err = binary.Read(buf, binary.BigEndian, &s)
	check(err, "failed to read System ID: ")
	err = binary.Read(buf, binary.BigEndian, &nsel)
	check(err, "failed to read NSAP selector: ")
	err = binary.Read(buf, binary.BigEndian, &m)
	check(err, "failed to read Metric: ")
	err = binary.Read(buf, binary.BigEndian, &subtlv)
	check(err, "failed to read subTLV: ")

	switch mt {
	case 2:
		mtid = "IPv6 Unicast"
	default:
		mtid = "Unknown"
	}
	// Metric has three bytes
	met := uint32(m[0])*65536 + uint32(m[1])*256 + uint32(m[2])
	// Format the System ID
	sysid := fmt.Sprintf("%x", (s[0:2])) + "." + fmt.Sprintf("%x", (s[2:4])) + "." + fmt.Sprintf("%x", (s[4:6]))

	n := neighbor{
		remoteID: sysid,
		metric:   met,
	}
	d.neighbors = append(d.neighbors, n)
	if !*b {
		fmt.Printf("Neighbor System ID: %v.%02d, Metric: %v\n", sysid, nsel, met)
	}
	if subtlv != 0 {
		// TODO Process the sub-TLV
		// return mtid, fmt.Errorf("missed a subTLV")
	}
	return mtid, err
}

func read237(t []byte, d *entry, b *bool) (mtid string, err error) {
	if len(t) < 2 {
		return mtid, fmt.Errorf("not a valid TLV, lenght: %v", len(t))
	}
	switch binary.BigEndian.Uint16(t[0:2]) {
	case 2:
		mtid = "IPv6 Unicast"
	default:
		mtid = "Unknown"
	}
	if !*b {
		fmt.Printf("MT ID: %v\nPrefixes:\n", mtid)
	}

	err = readPrefix(bytes.NewReader(t[2:]), b)
	return mtid, err
}

func readPrefix(buf *bytes.Reader, b *bool) (err error) {
	if buf.Len() == 0 {
		return err
	}
	if buf.Len() <= 6 {
		return fmt.Errorf("not a valid Prefix, lenght: %v", buf.Len())
	}
	var mask, flags, slen uint8
	var metric uint32

	err = binary.Read(buf, binary.BigEndian, &metric)
	check(err, "failed to read Metric: ")
	err = binary.Read(buf, binary.BigEndian, &flags)
	check(err, "failed to read SubTLV: ")
	err = binary.Read(buf, binary.BigEndian, &mask)
	check(err, "failed to read Mask: ")
	prefix := make([]byte, mask/8)
	err = binary.Read(buf, binary.BigEndian, &prefix)
	check(err, "failed to read Prefix: ")
	// Pad with additional bytes for IPv6 address compliance
	pad := make([]byte, 16-mask/8)
	prefix = append(prefix, pad...)

	// Check if subtlv present flag is on
	if flags&(1<<5) != 0 {
		// TODO Process the sub-TLV
		err = binary.Read(buf, binary.BigEndian, &slen)
		subtlv := make([]byte, int(slen))
		err = binary.Read(buf, binary.BigEndian, &subtlv)
		check(err, "failed to read subTLV: ")
	}
	if !*b {
		fmt.Printf("%v/%v, Metric:%v\n", net.IP(prefix), mask, metric)
	}
	err = readPrefix(buf, b)
	return err
}

func readHeader(h []byte, n int, d *entry, b *bool) (buf *bytes.Reader, err error) {
	if len(h) < 15 {
		return buf, fmt.Errorf("not a valid Header, lenght: %v", len(h))
	}
	sysid := fmt.Sprintf("%x", (h[0:2])) + "." + fmt.Sprintf("%x", (h[2:4])) + "." + fmt.Sprintf("%x", (h[4:6]))
	d.localID = sysid
	if !*b {
		fmt.Printf("===== LSP Details (lenght: %v) ====\n", n)
		fmt.Printf("LSPID:      %s.%x-%x\n", sysid, h[6:8], h[8:10])
		fmt.Printf("Seq Num:    %#x\n", h[10:12])
		fmt.Printf("Checksum:   %#x\n", h[12:14])
		fmt.Printf("Type Block: %#x\n", h[14:15])
	}
	// Get a io.Reader from a []byte slice
	buf = bytes.NewReader(h[15:])
	return buf, err
}

func exploreTLV(ts []tlv.TLV, d *entry, b *bool) error {
	for _, tl := range ts {
		switch tl.Type() {
		case 1:
			a := fmt.Sprintf("%x.%x.%x", tl.Value()[1:2], tl.Value()[2:4], tl.Value()[4:6])
			d.area = a
			if !*b {
				fmt.Printf("Type%03d,  L%03d: %s\n", tl.Type(), tl.Length(), a)
			}
		case 137:
			d.hostname = string(tl.Value())
			if !*b {
				fmt.Printf("Type%03d,  L%03d: %s\n", tl.Type(), tl.Length(), tl.Value())
			}
		case 140, 232:
			if !*b {
				fmt.Printf("Type%03d,  L%03d: %v\n", tl.Type(), tl.Length(), net.IP(tl.Value()))
			}
		case 222:
			if !*b {
				fmt.Printf("Type%03d,  L%03d: ", tl.Type(), tl.Length())
			}
			_, err := read222(tl.Value()[:tl.Length()], d, b)
			if err != nil {
				return errors.Wrap(err, "failed to read TLV 222")
			}
		case 237:
			if !*b {
				fmt.Printf("Type%03d,  L%03d: ", tl.Type(), tl.Length())
			}
			_, err := read237(tl.Value()[:tl.Length()], d, b)
			if err != nil {
				return errors.Wrap(err, "failed to read TLV 237")
			}
		default:
			if !*b {
				fmt.Printf("Type%03d,  L%03d: %#x\n", tl.Type(), tl.Length(), tl.Value())
			}
		}
	}
	return nil
}

func decodeTelemetry(v interface{}, f string) error {
	file, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("could not open the file: %s; %v", f, err)
	}
	return json.NewDecoder(file).Decode(v)
}
