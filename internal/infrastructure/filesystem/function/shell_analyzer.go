package function

import (
	"fmt"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

type FunctionInfo struct {
	Name          string
	StartLine     uint
	EndLine       uint
	HasDocs       bool
	Documentation []string
}

type CodeOutsideFunction struct {
	Line    uint
	Content string
}

type ShellAnalyzer struct {
	Functions   []FunctionInfo
	CodeOutside []CodeOutsideFunction
	sourceLines []string
}

func NewShellAnalyzer() *ShellAnalyzer {
	return &ShellAnalyzer{
		Functions:   make([]FunctionInfo, 0),
		CodeOutside: make([]CodeOutsideFunction, 0),
	}
}

func (sa *ShellAnalyzer) AnalyzeScript(script string) error {
	sa.sourceLines = strings.Split(script, "\n")

	parser := syntax.NewParser(syntax.KeepComments(true))
	file, err := parser.Parse(strings.NewReader(script), "")
	if err != nil {
		return fmt.Errorf("failed to parse script: %w", err)
	}

	sa.walkAST(file)
	sa.findCodeOutsideFunctions(file)

	return nil
}

func (sa *ShellAnalyzer) walkAST(file *syntax.File) {
	var comments []*syntax.Comment

	// Collect all comments first
	syntax.Walk(file, func(node syntax.Node) bool {
		if comment, ok := node.(*syntax.Comment); ok {
			comments = append(comments, comment)
		}
		return true
	})

	// Find functions and their documentation
	syntax.Walk(file, func(node syntax.Node) bool {
		if funcDecl, ok := node.(*syntax.FuncDecl); ok {
			sa.analyzeFunctionWithDocs(funcDecl, comments)
		}
		return true
	})
}

func (sa *ShellAnalyzer) analyzeFunctionWithDocs(funcDecl *syntax.FuncDecl, comments []*syntax.Comment) {
	funcInfo := FunctionInfo{
		Name:          funcDecl.Name.Value,
		StartLine:     funcDecl.Pos().Line(),
		EndLine:       funcDecl.End().Line(),
		HasDocs:       false,
		Documentation: make([]string, 0),
	}

	// Find comments directly above the function (consecutive comment block only)
	funcStartLine := funcDecl.Pos().Line()
	var docComments []string

	// Find the consecutive comment block immediately before the function
	// Only comments that are directly adjacent (no gaps) should be considered
	for lineNum := int(funcStartLine) - 1; lineNum >= 1; lineNum-- {
		commentFound := false
		for _, comment := range comments {
			if comment.Pos().Line() == uint(lineNum) {
				// This is a comment line directly above the function
				commentText := strings.TrimSpace(strings.TrimPrefix(comment.Text, "#"))
				// Skip shebang lines
				if !strings.HasPrefix(commentText, "!") && commentText != "" {
					docComments = append([]string{commentText}, docComments...) // prepend
				}
				commentFound = true
				break
			}
		}

		// If we didn't find a comment on this line, stop looking
		// This ensures only consecutive comment blocks are considered
		if !commentFound {
			break
		}
	}

	if len(docComments) > 0 {
		funcInfo.HasDocs = true
		funcInfo.Documentation = docComments
	}

	sa.Functions = append(sa.Functions, funcInfo)
}

func (sa *ShellAnalyzer) findCodeOutsideFunctions(file *syntax.File) {
	// Get all function ranges
	funcRanges := make([][2]uint, len(sa.Functions))
	for i, fn := range sa.Functions {
		funcRanges[i] = [2]uint{fn.StartLine, fn.EndLine}
	}

	// Track processed lines to avoid duplicates
	processedLines := make(map[uint]bool)

	// Walk through all statements
	syntax.Walk(file, func(node syntax.Node) bool {
		// Skip comments and function declarations themselves
		if _, isComment := node.(*syntax.Comment); isComment {
			return true
		}
		if _, isFunc := node.(*syntax.FuncDecl); isFunc {
			return true
		}

		// Check for executable statements - handle different types
		var stmtLine uint
		switch stmt := node.(type) {
		case *syntax.CallExpr:
			stmtLine = stmt.Pos().Line()
		case *syntax.DeclClause:
			stmtLine = stmt.Pos().Line()
		case *syntax.Assign:
			stmtLine = stmt.Pos().Line()
		case *syntax.IfClause:
			stmtLine = stmt.Pos().Line()
		case *syntax.ForClause:
			stmtLine = stmt.Pos().Line()
		case *syntax.WhileClause:
			stmtLine = stmt.Pos().Line()
		case *syntax.CaseClause:
			stmtLine = stmt.Pos().Line()
		case *syntax.Block:
			stmtLine = stmt.Pos().Line()
		default:
			return true
		}

		// Skip if already processed
		if processedLines[stmtLine] {
			return true
		}

		// Skip if it's a shebang line or empty
		if stmtLine == 1 && len(sa.sourceLines) > 0 && strings.HasPrefix(sa.sourceLines[0], "#!") {
			return true
		}

		// Check if this statement is inside any function
		insideFunction := false
		for _, funcRange := range funcRanges {
			if stmtLine >= funcRange[0] && stmtLine <= funcRange[1] {
				insideFunction = true
				break
			}
		}

		if !insideFunction {
			content := ""
			if int(stmtLine-1) < len(sa.sourceLines) {
				content = strings.TrimSpace(sa.sourceLines[stmtLine-1])
			}

			// Skip empty lines and comment-only lines
			if content != "" && !strings.HasPrefix(content, "#") {
				sa.CodeOutside = append(sa.CodeOutside, CodeOutsideFunction{
					Line:    stmtLine,
					Content: content,
				})
				processedLines[stmtLine] = true
			}
		}
		return true
	})
}

func (sa *ShellAnalyzer) PrintReport() {
	fmt.Println("=== Shell Script Analysis Report ===")
	fmt.Println()

	// Report functions
	fmt.Printf("Found %d function(s):\n", len(sa.Functions))
	for _, fn := range sa.Functions {
		fmt.Printf("  📋 %s (lines %d-%d)\n", fn.Name, fn.StartLine, fn.EndLine)
		if fn.HasDocs {
			fmt.Println("    ✅ Has documentation:")
			for _, doc := range fn.Documentation {
				fmt.Printf("       %s\n", doc)
			}
		} else {
			fmt.Println("    ❌ Missing documentation")
		}
		fmt.Println()
	}

	// Report code outside functions
	if len(sa.CodeOutside) > 0 {
		fmt.Printf("⚠️  Found %d line(s) with code outside functions:\n", len(sa.CodeOutside))
		for _, code := range sa.CodeOutside {
			fmt.Printf("  Line %d: %s\n", code.Line, code.Content)
		}
		fmt.Println()
		fmt.Println("🚨 Warning: Code outside functions may have side effects on your injections. Please move this code into functions.")
	} else {
		fmt.Println("✅ No code found outside functions - good practice!")
	}
}
