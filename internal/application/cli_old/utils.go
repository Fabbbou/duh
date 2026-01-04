package cli_old

import (
	"fmt"
	"os"
)

func stdPrint(s string) {
	fmt.Fprint(os.Stdout, s)
}
