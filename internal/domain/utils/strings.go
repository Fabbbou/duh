package utils

import "strings"

func EscapeDoubleQuotes(input string) string {
	return strings.ReplaceAll(input, `"`, `\"`)
}
