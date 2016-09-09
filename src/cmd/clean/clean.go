package clean

import (
	"github.com/iocage/libiocage/clean/cleanJails"
	"github.com/iocage/libiocage/list/baseDatasets"
	"github.com/iocage/libiocage/list/jailDatasets"
	"github.com/iocage/libiocage/list/templateDatasets"
)

func Args(Pool *string, force bool, ctype ...bool) {
	jails := ctype[0]
	base := ctype[1]
	template := ctype[2]
	var datasets []string

	switch {
	case base:
		datasets = baseDatasets.Args(Pool)
	case template:
		datasets = templateDatasets.Args(Pool)
	case jails:
		datasets = jailDatasets.Args(Pool)
		cleanJails.Args(Pool, force, datasets)
	}
}
