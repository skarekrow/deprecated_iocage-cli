package destroy

import (
	"fmt"
	"github.com/iocage/libiocage/destroy/destroyJails"
	"github.com/iocage/libiocage/get/uuidPathTag"
	"github.com/iocage/libiocage/list/jailDatasets"
	"os"
)

func Args(Pool *string, force, clean bool, jail []string) {
	var uuid, path, tag string
	var err error

	for _, j := range jail {
		if !clean {
			jails := jailDatasets.Args(Pool)
			uuid, path, tag, err = uuidPathTag.Args(Pool, jails, j)
		} else {
			uuid, path, tag, err = uuidPathTag.Args(Pool, []string{j}, j)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		destroyJails.Args(Pool, force, j, uuid, path, tag)
	}
}
