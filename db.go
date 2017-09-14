package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type entries struct {
	list []*entry
}

type entry struct {
	localID   string
	area      string
	hostname  string
	neighbors []neighbor
	prefixes  []prefix
}

type neighbor struct {
	remoteID string
	metric   uint32
}

type prefix struct {
	ip     string
	mask   uint8
	metric uint32
}

func displayTable(d *entries) {
	var data [][]string
	for _, s := range d.list {
		for _, n := range s.neighbors {
			data = append(data, []string{s.localID, s.hostname, n.remoteID, fmt.Sprint(n.metric)})
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Local ID", "Hostname", "Remote ID", "Metric"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
