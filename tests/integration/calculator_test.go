package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCalculatorExample(t *testing.T) {
	// Skip native compilation tests on non-ARM64 platforms
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("Native compilation only supported on macOS ARM64")
	}

	// Create a temporary file with calculator code
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "calc_test.pg")

	code := `func main() {
    x := 10
    y := 5
    
    // Basic operations
    sum := x + y
    diff := x - y
    prod := x * y
    quot := x / y
    
    println(sum)   // 15
    println(diff)  // 5
    println(prod)  // 50
    println(quot)  // 2
    
    // Operator precedence
    result := x + y * 2
    println(result) // 20 (10 + 10)
    
    // Parentheses should override precedence
    result2 := (x + y) * 2
    println(result2) // 30 ((10 + 5) * 2)
}`

	err := os.WriteFile(testFile, []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Run petitgo
	cmd := exec.Command("go", "run", "../../main.go", "run", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run petitgo: %v\nOutput: %s", err, output)
	}

	// Check output
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	expected := []string{"15", "5", "50", "2", "20", "30"}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d\nOutput:\n%s", len(expected), len(lines), output)
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %s, got %s", i+1, expected[i], line)
		}
	}
}

func TestParenthesesPrecedence(t *testing.T) {
	// Skip native compilation tests on non-ARM64 platforms
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("Native compilation only supported on macOS ARM64")
	}

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "basic parentheses",
			code: `func main() {
    result := (2 + 3) * 4
    println(result)
}`,
			expected: "20",
		},
		{
			name: "nested parentheses",
			code: `func main() {
    result := ((1 + 2) * 3) + 4
    println(result)
}`,
			expected: "13",
		},
		{
			name: "parentheses with variables",
			code: `func main() {
    x := 10
    y := 5
    result := (x + y) * 2
    println(result)
}`,
			expected: "30",
		},
		{
			name: "complex expression",
			code: `func main() {
    a := 2
    b := 3
    c := 4
    d := 5
    result := (a + b) * (c + d)
    println(result)
}`,
			expected: "45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.pg")

			err := os.WriteFile(testFile, []byte(tt.code), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Run petitgo
			cmd := exec.Command("go", "run", "../../main.go", "run", testFile)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to run petitgo: %v\nOutput: %s", err, output)
			}

			// Check output
			result := strings.TrimSpace(string(output))
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
