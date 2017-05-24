// main.go
package main

import (
	"bytes"
	"encoding/base64"
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

func main() {
	fileb64 := "data64"
	var enc = base64.NewEncoding(encodeStd)

	// Read the the whole file at once and put it an a byte array
	src, err := ioutil.ReadFile(fileb64)
	check(err, "Error opening file")

	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	check(err, "Error decoding file")

	fmt.Printf("===== LSP Details (lenght: %v) ====\n", n)

	fmt.Printf("LSPID:      %X.%X.%X.%X-%X\n", dst[0:2], dst[2:4], dst[4:6], dst[6:8], dst[8:10])
	fmt.Printf("Seq Num:    %#x\n", dst[10:12])
	fmt.Printf("Checksum:   %#x\n", dst[12:14])
	fmt.Printf("Type Block: %#x\n", dst[14:15])

	// Get a io.Reader from a []byte slice
	r := bytes.NewReader(dst[15:n])

	// Read the TLV's from the Reader and put them on a slice
	tlvs, err := tlv.Read(r)
	check(err, "Failed to read TLVs")
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
			fmt.Printf("Type%03d,  L%03d:\n", tl.Type(), tl.Length())
		default:
			fmt.Printf("Type%03d,  L%03d: %#x\n", tl.Type(), tl.Length(), tl.Value())
		}
	}

}
