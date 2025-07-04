package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yuya-takeyama/petitgo/asmgen"
	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/codegen"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/repl"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]

		switch command {
		case "build":
			if len(os.Args) < 3 {
				fmt.Println("Usage: petitgo build <file.go>")
				os.Exit(1)
			}
			buildFile(os.Args[2])
			return
		case "run":
			if len(os.Args) < 3 {
				fmt.Println("Usage: petitgo run <file.go>")
				os.Exit(1)
			}
			runFile(os.Args[2])
			return
		case "ast":
			if len(os.Args) < 3 {
				fmt.Println("Usage: petitgo ast <file.go>")
				os.Exit(1)
			}
			astFile(os.Args[2])
			return
		case "asm":
			if len(os.Args) < 3 {
				fmt.Println("Usage: petitgo asm <file.pg>")
				os.Exit(1)
			}
			asmFile(os.Args[2])
			return
		case "native":
			if len(os.Args) < 3 {
				fmt.Println("Usage: petitgo native <file.pg>")
				os.Exit(1)
			}
			nativeBuild(os.Args[2])
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Available commands: build, run, ast, asm, native")
			os.Exit(1)
		}
	}

	// REPLを起動
	repl.StartREPL()
}

// buildFile compiles a petitgo file to a Go executable
func buildFile(filename string) {
	// Read the petitgo source file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Parse the petitgo code
	statements := parseProgram(string(content))

	// Generate Go source code
	gen := codegen.NewGenerator()
	goSource := gen.GenerateProgram(statements)

	// Create temporary Go file
	tempDir := os.TempDir()
	tempGoFile := filepath.Join(tempDir, "petitgo_output.go")

	err = os.WriteFile(tempGoFile, []byte(goSource), 0644)
	if err != nil {
		fmt.Printf("Error writing temporary Go file: %v\n", err)
		os.Exit(1)
	}

	// Compile with go build
	outputName := strings.TrimSuffix(filename, filepath.Ext(filename))
	cmd := exec.Command("go", "build", "-o", outputName, tempGoFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error compiling Go code: %v\n", err)
		fmt.Printf("Go compiler output:\n%s\n", output)
		os.Exit(1)
	}

	// Clean up temporary file
	os.Remove(tempGoFile)

	fmt.Printf("Successfully compiled %s to %s\n", filename, outputName)
}

// runFile compiles and runs a petitgo file
func runFile(filename string) {
	// Read and parse like buildFile
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	statements := parseProgram(string(content))
	gen := codegen.NewGenerator()
	goSource := gen.GenerateProgram(statements)

	// Create temporary Go file
	tempDir := os.TempDir()
	tempGoFile := filepath.Join(tempDir, "petitgo_run.go")

	err = os.WriteFile(tempGoFile, []byte(goSource), 0644)
	if err != nil {
		fmt.Printf("Error writing temporary Go file: %v\n", err)
		os.Exit(1)
	}

	// Run with go run
	cmd := exec.Command("go", "run", tempGoFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Clean up
	os.Remove(tempGoFile)
}

// parseProgram parses a petitgo program and returns statements
func parseProgram(source string) []ast.Statement {
	s := scanner.NewScanner(source)
	p := parser.NewParser(s)

	var statements []ast.Statement

	for {
		stmt := p.ParseStatement()
		if stmt == nil {
			break
		}
		statements = append(statements, stmt)
	}

	return statements
}

// astFile parses a petitgo file and outputs the AST as JSON
func astFile(filename string) {
	// Read the petitgo source file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Parse the petitgo code
	statements := parseProgram(string(content))

	// Convert AST to JSON using MarshalJSON methods
	program := map[string]interface{}{
		"type":       "Program",
		"statements": statements,
	}

	// Pretty print JSON
	jsonBytes, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling AST to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}

// asmFile generates ARM64 assembly from petitgo source
func asmFile(filename string) {
	// Read the petitgo source file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Parse the petitgo code
	statements := parseProgram(string(content))

	// Generate ARM64 assembly
	generator := asmgen.NewAsmGenerator()
	assembly := generator.Generate(statements)
	runtime := generator.GenerateRuntime()

	fmt.Print(assembly)
	fmt.Print(runtime)
}

// nativeBuild compiles petitgo source to native ARM64 binary
func nativeBuild(filename string) {
	// Read the petitgo source file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Parse the petitgo code
	statements := parseProgram(string(content))

	// Generate ARM64 assembly
	generator := asmgen.NewAsmGenerator()
	assembly := generator.Generate(statements)
	runtime := generator.GenerateRuntime()

	// Write assembly to temporary file
	asmFile := "/tmp/petitgo_temp.s"
	fullAsm := assembly + runtime
	err = os.WriteFile(asmFile, []byte(fullAsm), 0644)
	if err != nil {
		fmt.Printf("Error writing assembly file: %v\n", err)
		os.Exit(1)
	}

	// Get output filename
	outputFile := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Assemble and link
	cmd := exec.Command("as", "-arch", "arm64", "-o", "/tmp/petitgo_temp.o", asmFile)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error assembling: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("clang", "-arch", "arm64", "-o", outputFile, "/tmp/petitgo_temp.o")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error linking: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully compiled %s to native binary %s\n", filename, outputFile)
}
