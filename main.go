// main.go
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	//"github.com/Akagi201/tlv"
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

	// Read the the whole file at once and pute it an a byte array
	src, err := ioutil.ReadFile(fileb64)
	check(err, "Error opening file")

	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	check(err, "Error decoding file")

	fmt.Printf("Lenght: %v\n", n)

	fmt.Printf("LSPID: %X\n", dst[:10])
	fmt.Printf("Seq Num: %#x\n", dst[10:12])
	fmt.Printf("Checksum: %#x\n", dst[12:14])
	fmt.Printf("Type Block: %#x\n", dst[14:15])
	fmt.Printf("T10,  L17: %#x\n", dst[15:34])
	fmt.Printf("T01,  L06: %#x\n", dst[34:42])
	fmt.Printf("T129, L01: %#x\n", dst[42:45])
	fmt.Printf("T229, L02: %#x\n", dst[45:49])
	fmt.Printf("T137, L22: %#x, %s\n", dst[49:51], dst[51:73])
	fmt.Printf("T232, L16: %#x\n", dst[73:91])
	fmt.Printf("T222, L13: %#x\n", dst[91:106])
	fmt.Printf("T237, L118: %#x\n", dst[106:n])

}
