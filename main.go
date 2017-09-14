// tlvdecode decodes IS-IS base64 encoded TLVs.
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	file := "input/" + os.Args[1]
	if strings.Contains(file, ".json") {
		data := new(Data)
		err := decodeTelemetry(data, file)
		check(err, "Error decoding JSON file")
		for _, b := range data.Rows {
			err = readBytes([]byte(b.Content.LspBody))
			check(err, "Error reading bytes")
		}
		return
	}
	// Read the the whole file at once and put it an a byte array
	src, err := ioutil.ReadFile(file)
	check(err, "Error opening file")
	err = readBytes(src)
	check(err, "Error reading bytes")
}
