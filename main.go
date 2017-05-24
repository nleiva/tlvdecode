// main.go
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/nleiva/tlv"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func read237(t []byte) (mtid string, err error) {
	if len(t) < 2 {
		return mtid, fmt.Errorf("Not a valid TLV, lenght: %v", len(t))
	}
	switch binary.BigEndian.Uint16(t[0:2]) {
	case 2:
		mtid = "IPv6"
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
	var mask uint8
	var flags uint8
	var metric uint32
	// subTLV can have different lenght!, might improve this in the future
	subtlv := make([]byte, 4)

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

func main() {
	fileb64 := "data64"
	var enc = base64.NewEncoding(encodeStd)

	// Read the the whole file at once and put it an a byte array
	src, err := ioutil.ReadFile(fileb64)
	check(err, "Error opening file")

	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	check(err, "Error decoding file")

	// Read Header
	r, err := readHeader(dst, n)
	check(err, "Error reading header")

	// Read the TLV's from the Reader and put them on a slice
	tlvs, err := tlv.Read(r)
	check(err, "Failed to read TLVs: ")
	ts := tlvs.GetThemAll()

	fmt.Printf("===== TLV Details (total: %03d) ====\n", tlvs.Length())
	for _, tl := range ts {
		switch tl.Type() {
		case 1:
			fmt.Printf("Type%03d,  L%03d: %x.%x.%x\n", tl.Type(), tl.Length(), tl.Value()[1:2], tl.Value()[2:4], tl.Value()[4:6])
		case 137:
			fmt.Printf("Type%03d,  L%03d: %s\n", tl.Type(), tl.Length(), tl.Value())
		case 232:
			fmt.Printf("Type%03d,  L%03d: %v\n", tl.Type(), tl.Length(), net.IP(tl.Value()))
		case 237:
			fmt.Printf("Type%03d,  L%03d: ", tl.Type(), tl.Length())
			_, err := read237(tl.Value()[:tl.Length()])
			check(err, "Failed to read TLV237: ")
		default:
			fmt.Printf("Type%03d,  L%03d: %#x\n", tl.Type(), tl.Length(), tl.Value())
		}
	}

}
