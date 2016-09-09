package get

import (
	"fmt"
	"github.com/iocage/libiocage/get/uclProp"
	"github.com/iocage/libiocage/get/uuidPathTag"
	"github.com/iocage/libiocage/list/jailDatasets"
	"github.com/olekukonko/tablewriter"
	"os"
)

// get accepts a property and a jail. It accepts a recursive boolean which
// governs recursing through every jail and printing that property.
// Otherwise it will return the requested property of the supplied jail.
func Args(property, jail string, Pool, Iocroot *string, flags ...bool) string {
	var prop string
	var tbl []string
	recursive := flags[0]
	header := flags[1]

	jails := jailDatasets.Args(Pool)

	if jail == "default" {
		prop, _ = uclProp.Args(*Iocroot, property, true)
		fmt.Println(prop)
		os.Exit(0)
	}

	if recursive {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeader([]string{"UUID", "TAG", "PROP - " + property})
		for _, j := range jails {
			uuid, path, tag, _ := uuidPathTag.Args(Pool, []string{j}, j)
			prop, _ = uclProp.Args(path, property, false)
			tbl = []string{uuid, tag, prop}
			table.Append(tbl)
			if header {
				fmt.Printf("%s %s %s\n", uuid, tag, prop)
			}
		}

		if !header {
			table.Render()
		}

		os.Exit(0)
	}

	// Check to see if the jail exists.
	_, path, _, err := uuidPathTag.Args(Pool, jails, jail)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prop, err = uclProp.Args(path, property, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return prop
}
