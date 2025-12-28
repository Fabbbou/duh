package cli

import (
	"fmt"
	"os"
)

func stdPrint(s string) {
	fmt.Fprint(os.Stdout, s)
}
