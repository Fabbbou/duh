package function

import (
	"fmt"
	"os"

	"mvdan.cc/sh/v3/syntax"

	"strings"
)

func GetScript(path string) (string, error) {
	scriptBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	parser := syntax.NewParser(syntax.KeepComments(true))
	file, err := parser.Parse(strings.NewReader(string(scriptBytes)), "")
	if err != nil {
		return "", err
	}

	// Walk through the AST to find comments
	syntax.Walk(file, func(node syntax.Node) bool {
		if comment, ok := node.(*syntax.Comment); ok {
			fmt.Println("Found comment:", comment.Text)
		}
		return true
	})
	return string(scriptBytes), nil
}
