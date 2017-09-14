// tlvdecode decodes IS-IS base64 encoded TLVs.
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	file := flag.String("f", "input/full1.json", "Input file")
	table := flag.Bool("t", false, "Display Link Table")
	flag.Parse()

	info := new(entries)
	if strings.Contains(*file, ".json") {
		data := new(Data)
		err := decodeTelemetry(data, *file)
		check(err, "Error decoding JSON file")
		for _, b := range data.Rows {
			system := new(entry)
			err = readBytes([]byte(b.Content.LspBody), system, table)
			check(err, "Error reading bytes")
			info.list = append(info.list, system)
		}
		if *table {
			displayTable(info)
		}
		return
	}
	// Read the the whole file at once and put it an a byte array
	src, err := ioutil.ReadFile(*file)
	check(err, "Error opening file")
	system := new(entry)
	err = readBytes(src, system, table)
	check(err, "Error reading bytes")
	if *table {
		displayTable(info)
	}
}
