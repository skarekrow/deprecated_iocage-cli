package fetch

import (
	"github.com/iocage/libiocage/fetch/fetchBase"
	// "fmt"
)

func Args(Pool, Iocroot *string, props []string) {
	fetchBase.Args(Pool, Iocroot, props)
	// fmt.Printf("%s\n", results)
}
