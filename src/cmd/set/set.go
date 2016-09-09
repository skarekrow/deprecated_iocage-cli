package set

import (
	"fmt"
	"github.com/iocage/libiocage/get/uuidPathTag"
	"github.com/iocage/libiocage/list/jailDatasets"
	"github.com/iocage/libiocage/set/jailProp"
	"os"
	"strings"
)

// set accepts a property and a jail. It will exit silently for success.
func Args(property, jail string, Pool, Iocroot *string) error {
	jails := jailDatasets.Args(Pool)
	p := strings.Split(property, "=")
	prop, val := p[0], p[1]

	switch {
	// notes will annoyingly get mangled if it has an '=' in it.
	case property[0:5] == "notes":
		val = property[6:]
	}

	if jail == "default" {
		err := jailProp.Args(*Iocroot, prop, val, true)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		_, path, _, err := uuidPathTag.Args(Pool, jails, jail)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = jailProp.Args(path, prop, val, false)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return nil
}
