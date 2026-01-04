package handler

import (
	"fmt"
	"os"
)

func stdPrintln(s string) {
	fmt.Fprint(os.Stdout, s+"\n")
}
