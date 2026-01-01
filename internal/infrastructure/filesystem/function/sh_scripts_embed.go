package function

/*
This file is declaring all the scripts from sh_scripts/ as embedded strings,
using Go's embed package.

It's the best pattern to have the scripts embedded in the binary
and still be able to run manually like regular shell scripts.
*/

import (
	_ "embed"
	"fmt"
)

//go:embed sh_scripts/require.sh
var RequireShScript string

// This function could be removed later, it's just an example of how to print docs
// Example__PrintRequireScriptDocs prints the documentation for functions in require.sh
func Example__PrintRequireScriptDocs() error {
	analyzer, err := GetScriptAnalysis(RequireShScript)
	if err != nil {
		return err
	}

	fmt.Println("=== require.sh Function Documentation ===")
	for _, fn := range analyzer.Functions {
		fmt.Printf("📋 %s:\n", fn.Name)
		if fn.HasDocs {
			for _, doc := range fn.Documentation {
				fmt.Printf("   %s\n", doc)
			}
		} else {
			fmt.Println("   ❌ No documentation")
		}
		fmt.Println()
	}

	if len(analyzer.CodeOutside) > 0 {
		fmt.Println("⚠️ Warning: Script has code outside functions")
	}

	return nil
}
