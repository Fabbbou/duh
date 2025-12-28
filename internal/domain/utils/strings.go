package utils

import "strings"

func EscapeShellString(input string) string {
	// Characters that need escaping in shell contexts
	replacer := strings.NewReplacer(
		`\`, `\\`, // backslash must be first
		`$`, `\$`, // variable expansion
		`"`, `\"`, // double quote
	)
	return replacer.Replace(input)
}
