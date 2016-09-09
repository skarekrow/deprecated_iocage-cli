package list

import (
	"fmt"
	"github.com/iocage/libiocage/list/baseDatasets"
	"github.com/iocage/libiocage/list/jailDatasets"
	"github.com/iocage/libiocage/list/prepare"
	"github.com/iocage/libiocage/list/templateDatasets"
	"github.com/olekukonko/tablewriter"
	"os"
)

// list accepts a list type and then runs the appropriate list function.
func Args(ltype string, Pool *string, flags ...bool) {
	header := flags[0]
	base := flags[1]
	template := flags[2]
	var list [][]string
	var datasets []string
	var t, action string

	switch {
	case base:
		ltype = "base"
		t = "bases"
		action = "fetched"
		datasets = baseDatasets.Args(Pool)
	case template:
		ltype = "template"
		t = "templates"
		action = "made"
		datasets = templateDatasets.Args(Pool)
	default:
		ltype = "all"
		t = "jails"
		action = "created"
		datasets = jailDatasets.Args(Pool)
	}

	if !header {
		list = prepare.Args(Pool, datasets, ltype, false)
		if len(list) == 0 {
			fmt.Printf("No %s have been %s.\n", t, action)
			os.Exit(0)
		}
	} else {
		list = prepare.Args(Pool, datasets, ltype, true)
		if len(list) == 0 {
			fmt.Printf("No %s have been %s.\n", t, action)
			os.Exit(0)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	switch ltype {
	case "all":
		table.SetHeader([]string{"JID", "UUID", "BOOT", "STATE", "JAIL TAG",
			"IP4", "BASE", "JAIL #"})
		table.SetFooter([]string{"JID", "UUID", "BOOT", "STATE", "JAIL TAG",
			"IP4", "BASE", "JAIL #"})
	case "template":
		table.SetHeader([]string{"JID", "UUID", "BOOT", "STATE",
			"TEMPLATE TAG", "IP4", "BASE"})
		table.SetFooter([]string{"JID", "UUID", "BOOT", "STATE",
			"TEMPLATE TAG", "IP4", "BASE"})
	case "base":
		table.SetHeader([]string{"FETCHED BASES"})
		table.SetFooter([]string{"FETCHED BASES"})
	default:
		fmt.Printf("%s is not a valid type to list.\n", ltype)
		os.Exit(1)
	}

	if header {
		for _, d := range list {
			fmt.Println(d[0])
		}
	} else {
		table.AppendBulk(list)
		table.Render()
	}
}
