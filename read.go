package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func read222(t []byte) (mtid string, err error) {
	if len(t) < 13 {
		return mtid, fmt.Errorf("Not a valid TLV, lenght: %v", len(t))
	}
	var mt uint16
	var subtlv uint8
	var nsel uint8
	s := make([]byte, 6)
	m := make([]byte, 3)

	buf := bytes.NewReader(t)
	err = binary.Read(buf, binary.BigEndian, &mt)
	check(err, "Failed to read MT ID: ")
	err = binary.Read(buf, binary.BigEndian, &s)
	check(err, "Failed to read System ID: ")
	err = binary.Read(buf, binary.BigEndian, &nsel)
	check(err, "Failed to read NSAP selector: ")
	err = binary.Read(buf, binary.BigEndian, &m)
	check(err, "Failed to read Metric: ")
	err = binary.Read(buf, binary.BigEndian, &subtlv)
	check(err, "Failed to read subTLV: ")

	switch mt {
	case 2:
		mtid = "IPv6 Unicast"
	default:
		mtid = "Unknown"
	}
	// Metric has three bytes
	metric := uint32(m[0])*65536 + uint32(m[1])*256 + uint32(m[2])
	// Format the System ID
	sysid := fmt.Sprintf("%x", (s[0:2])) + "." + fmt.Sprintf("%x", (s[2:4])) + "." + fmt.Sprintf("%x", (s[4:6]))

	// This is just temporary
	fmt.Printf("Neighbor System ID: %v.%02d, Metric: %v\n", sysid, nsel, metric)
	if subtlv != 0 {
		return mtid, fmt.Errorf("Missed a subTLV")
	}
	return mtid, err
}

func read237(t []byte) (mtid string, err error) {
	if len(t) < 2 {
		return mtid, fmt.Errorf("Not a valid TLV, lenght: %v", len(t))
	}
	switch binary.BigEndian.Uint16(t[0:2]) {
	case 2:
		mtid = "IPv6 Unicast"
	default:
		mtid = "Unknown"
	}
	// This is just temporary
	fmt.Printf("MT ID: %v\nPrefixes:\n", mtid)

	err = readPrefix(bytes.NewReader(t[2:]))
	return mtid, err
}

func readPrefix(buf *bytes.Reader) (err error) {
	if buf.Len() == 0 {
		return err
	}
	if buf.Len() <= 6 {
		return fmt.Errorf("Not a valid Prefix, lenght: %v", buf.Len())
	}
	var mask, flags, slen uint8
	var metric uint32

	err = binary.Read(buf, binary.BigEndian, &metric)
	check(err, "Failed to read Metric: ")
	err = binary.Read(buf, binary.BigEndian, &flags)
	check(err, "Failed to read SubTLV: ")
	err = binary.Read(buf, binary.BigEndian, &mask)
	check(err, "Failed to read Mask: ")
	prefix := make([]byte, mask/8)
	err = binary.Read(buf, binary.BigEndian, &prefix)
	check(err, "Failed to read Prefix: ")
	// Pad with additional bytes for IPv6 address compliance
	pad := make([]byte, 16-mask/8)
	prefix = append(prefix, pad...)

	// Check if subtlv present flag is on
	if flags&(1<<5) != 0 {
		err = binary.Read(buf, binary.BigEndian, &slen)
		subtlv := make([]byte, int(slen))
		err = binary.Read(buf, binary.BigEndian, &subtlv)
		check(err, "Failed to read subTLV: ")
	}

	fmt.Printf("%v/%v, Metric:%v\n", net.IP(prefix), mask, metric)
	err = readPrefix(buf)

	return err
}

func readHeader(h []byte, n int) (buf *bytes.Reader, err error) {
	if len(h) < 15 {
		return buf, fmt.Errorf("Not a valid Header, lenght: %v", len(h))
	}
	fmt.Printf("===== LSP Details (lenght: %v) ====\n", n)
	fmt.Printf("LSPID:      %X.%X.%X.%X-%X\n", h[0:2], h[2:4], h[4:6], h[6:8], h[8:10])
	fmt.Printf("Seq Num:    %#x\n", h[10:12])
	fmt.Printf("Checksum:   %#x\n", h[12:14])
	fmt.Printf("Type Block: %#x\n", h[14:15])

	// Get a io.Reader from a []byte slice
	buf = bytes.NewReader(h[15:])
	return buf, err
}
