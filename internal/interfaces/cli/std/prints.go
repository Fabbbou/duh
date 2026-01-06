package std

import (
	"fmt"
	"os"
)

func F(format string, a ...any) (int, error) {
	return fmt.Printf(format, a...)
}

func Ln(a ...any) (int, error) {
	return fmt.Println(a...)
}

func Lnf(format string, a ...any) (int, error) {
	return fmt.Printf(format+"\n", a...)
}

func Errf(format string, a ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

func ErrLnf(format string, a ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, format+"\n", a...)
}
