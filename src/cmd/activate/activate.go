package activate

import (
	"fmt"
	"github.com/iocage/libiocage/activate/activatePool"
	"os"
)

func Args(pool string) string {
	cmd, err := activatePool.Args(pool)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cmd
}
