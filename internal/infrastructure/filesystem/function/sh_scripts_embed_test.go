package function

import (
	"strings"
	"testing"
)

func TestGetScriptAnalysis(t *testing.T) {
	analyzer, err := GetScriptAnalysis(RequireShScript)
	if err != nil {
		t.Fatalf("GetRequireScriptAnalysis() failed: %v", err)
	}

	if analyzer == nil {
		t.Fatal("Analyzer should not be nil")
	}

	// Should have functions (assuming require.sh has functions)
	if len(analyzer.Functions) == 0 {
		t.Error("require.sh should contain functions")
	}

	// Check for proper structure - no global code
	if len(analyzer.CodeOutside) > 0 {
		t.Errorf("require.sh should not have code outside functions, found %d lines", len(analyzer.CodeOutside))
		for _, code := range analyzer.CodeOutside {
			t.Logf("  Line %d: %s", code.Line, code.Content)
		}
	}

	// All functions should have documentation
	undocumentedFunctions := 0
	for _, fn := range analyzer.Functions {
		if !fn.HasDocs {
			undocumentedFunctions++
			t.Errorf("Function %s is missing documentation", fn.Name)
		}
	}

	if undocumentedFunctions > 0 {
		t.Errorf("Found %d undocumented functions in require.sh", undocumentedFunctions)
	}
}

func TestPrintScriptDocs(t *testing.T) {
	// This test mainly ensures the function doesn't crash
	err := Example__PrintRequireScriptDocs()
	if err != nil {
		t.Errorf("PrintRequireScriptDocs() failed: %v", err)
	}
}

func TestEmbeddedScriptConsistency(t *testing.T) {
	// Test that the embedded script is consistent with analysis
	script := RequireShScript

	analyzer := newShellAnalyzer()
	err := analyzer.analyzeScript(script)
	if err != nil {
		t.Fatalf("Failed to analyze embedded script: %v", err)
	}

	analyzer2, err := GetScriptAnalysis(RequireShScript)
	if err != nil {
		t.Fatalf("GetRequireScriptAnalysis failed: %v", err)
	}

	// Should have same number of functions
	if len(analyzer.Functions) != len(analyzer2.Functions) {
		t.Errorf("Inconsistent function count: direct=%d, helper=%d",
			len(analyzer.Functions), len(analyzer2.Functions))
	}

	// Should have same number of issues
	if len(analyzer.CodeOutside) != len(analyzer2.CodeOutside) {
		t.Errorf("Inconsistent code outside functions: direct=%d, helper=%d",
			len(analyzer.CodeOutside), len(analyzer2.CodeOutside))
	}
}

func TestRequireScriptValidShell(t *testing.T) {
	script := RequireShScript
	lines := strings.Split(script, "\n")

	if len(lines) == 0 {
		t.Fatal("Script should not be empty")
	}

	// First line should be shebang
	if !strings.HasPrefix(lines[0], "#!") {
		t.Error("Script should start with shebang")
	}

	// Should contain function keywords
	hasFunction := false
	for _, line := range lines {
		if strings.Contains(line, "function ") ||
			(strings.Contains(line, "()") && strings.Contains(line, "{")) {
			hasFunction = true
			break
		}
	}

	if !hasFunction {
		t.Error("Script should contain at least one function")
	}
}

// Benchmark test to ensure embedding is efficient
func BenchmarkGetRequireScript(b *testing.B) {
	for i := 0; i < b.N; i++ {
		script := RequireShScript
		if script == "" {
			b.Error("Script should not be empty")
		}
	}
}

func BenchmarkGetRequireScriptAnalysis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		analyzer, err := GetScriptAnalysis(RequireShScript)
		if err != nil {
			b.Errorf("Analysis failed: %v", err)
		}
		if analyzer == nil {
			b.Error("Analyzer should not be nil")
		}
	}
}
