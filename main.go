// main.go
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/nleiva/tlv"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	var enc = base64.NewEncoding(encodeStd)

	// Read the the whole file at once and put it an a byte array
	src, err := ioutil.ReadFile("input/" + os.Args[1])
	check(err, "Error opening file")

	// Decode base64 data into bytes
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	check(err, "Error decoding file")

	// Read PDU info
	r, err := readHeader(dst, n)
	check(err, "Error reading header")

	// Read the TLV's from the Reader and put them on a slice
	tlvs, err := tlv.Read(r)
	check(err, "Failed to read TLVs: ")
	ts := tlvs.GetThemAll()

	// Print the TLV details
	fmt.Printf("===== TLV Details (total: %03d) ====\n", tlvs.Length())
	for _, tl := range ts {
		switch tl.Type() {
		case 1:
			fmt.Printf("Type%03d,  L%03d: %x.%x.%x\n", tl.Type(), tl.Length(), tl.Value()[1:2], tl.Value()[2:4], tl.Value()[4:6])
		case 137:
			fmt.Printf("Type%03d,  L%03d: %s\n", tl.Type(), tl.Length(), tl.Value())
		case 140, 232:
			fmt.Printf("Type%03d,  L%03d: %v\n", tl.Type(), tl.Length(), net.IP(tl.Value()))
		case 222:
			fmt.Printf("Type%03d,  L%03d: ", tl.Type(), tl.Length())
			_, err := read222(tl.Value()[:tl.Length()])
			check(err, "Failed to read TLV 222: ")
		case 237:
			fmt.Printf("Type%03d,  L%03d: ", tl.Type(), tl.Length())
			_, err := read237(tl.Value()[:tl.Length()])
			check(err, "Failed to read TLV 237: ")
		default:
			fmt.Printf("Type%03d,  L%03d: %#x\n", tl.Type(), tl.Length(), tl.Value())
		}
	}

}
