package deactivate

import (
	"fmt"
	"github.com/iocage/libiocage/deactivate/deactivatePool"
	"os"
)

func Args(pool string) string {
	cmd, err := deactivatePool.Args(pool)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cmd
}
