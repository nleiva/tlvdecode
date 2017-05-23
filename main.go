// main.go
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

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

	fmt.Printf("LSPID: %X\n", dst[:10])
	fmt.Printf("Seq Num: %#x\n", dst[10:12])
	fmt.Printf("Checksum: %#x\n", dst[12:14])
	fmt.Printf("Type Block: %#x\n", dst[14:15])

	// Read individual TLV from byte array
	//tmpTLV, err := tlv.FromBytes(dst[15:n])
	//check(err, "Failed to read TLV")
	//fmt.Printf("T%v,  L%v: %#x\n", tmpTLV.Type(), tmpTLV.Length(), tmpTLV.Value())

	// Get a io.Reader from a []byte slice
	r := bytes.NewReader(dst[15:n])

	// Read the TLV's from the Reader and put them on a slice
	rtlvs, err := tlv.Read(r)
	check(err, "Failed to read TLVs")
	ts := rtlvs.GetThemAll()

	fmt.Printf("===== TLV Details (total: %03d) ====\n", rtlvs.Length())
	for _, tl := range ts {
		fmt.Printf("Type%03d,  L%03d: %#x\n", tl.Type(), tl.Length(), tl.Value())
	}

	// Manual way to look at TLV's
	//fmt.Printf("TLV: %#x, %s\n", dst[49:51], dst[51:73])
}
