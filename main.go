package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yuya-takeyama/petitgo/asmgen"
	"github.com/yuya-takeyama/petitgo/ast"
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
		case "help", "-h", "--help":
			showHelp()
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Available commands: build, run, ast, asm, help")
			os.Exit(1)
		}
	}

	// REPLを起動
	repl.StartREPL()
}

// showHelp displays usage information and available commands
func showHelp() {
	fmt.Println("petitgo - A minimal Go implementation aiming for self-hosting capability")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  petitgo [command] [arguments]")
	fmt.Println("")
	fmt.Println("COMMANDS:")
	fmt.Println("  build <file.pg>    Compile a petitgo program to native binary")
	fmt.Println("  run <file.pg>      Compile and run a petitgo program")
	fmt.Println("  ast <file.pg>      Display the Abstract Syntax Tree as JSON")
	fmt.Println("  asm <file.pg>      Generate ARM64 assembly code")
	fmt.Println("  help, -h, --help   Show this help message")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  petitgo                    # Start interactive REPL")
	fmt.Println("  petitgo run examples/fibonacci.pg")
	fmt.Println("  petitgo build hello.pg     # Creates 'hello' executable")
	fmt.Println("  petitgo ast program.pg     # View AST structure")
	fmt.Println("  petitgo asm program.pg     # View generated assembly")
	fmt.Println("")
	fmt.Println("PETITGO LANGUAGE FEATURES:")
	fmt.Println("  - Arithmetic operations (+, -, *, /)")
	fmt.Println("  - Variables and assignments (x := 10, x = 20)")
	fmt.Println("  - Control flow (if/else, for, switch/case)")
	fmt.Println("  - Functions with parameters and return values")
	fmt.Println("  - Basic types (int, string, bool)")
	fmt.Println("  - Struct definitions and field access")
	fmt.Println("  - Comments (// and /* */)")
	fmt.Println("  - Built-in functions: println(), len(), append()")
	fmt.Println("")
	fmt.Println("For more information, visit: https://github.com/yuya-takeyama/petitgo")
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

	// Generate ARM64 assembly directly (no more Go codegen)
	generator := asmgen.NewAsmGenerator()
	assembly := generator.Generate(statements)
	runtimeCode := generator.GenerateRuntime()

	// Write assembly to temporary file
	asmFile := "/tmp/petitgo_temp.s"
	fullAsm := assembly + runtimeCode
	err = os.WriteFile(asmFile, []byte(fullAsm), 0644)
	if err != nil {
		fmt.Printf("Error writing assembly file: %v\n", err)
		os.Exit(1)
	}

	// Get output filename
	outputName := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Assemble and link (platform-specific commands)
	var cmd *exec.Cmd
	objFile := "/tmp/petitgo_temp.o"

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// macOS ARM64
		cmd = exec.Command("as", "-arch", "arm64", "-o", objFile, asmFile)
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		// Linux x86_64
		cmd = exec.Command("as", "--64", "-o", objFile, asmFile)
	} else {
		fmt.Printf("Unsupported platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("Supported platforms: darwin/arm64, linux/amd64")
		os.Exit(1)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error assembling: %v\n", err)
		os.Exit(1)
	}

	// Platform-specific linker
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// macOS ARM64
		cmd = exec.Command("clang", "-arch", "arm64", "-o", outputName, objFile)
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		// Linux x86_64
		cmd = exec.Command("ld", "-o", outputName, objFile)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error linking: %v\n", err)
		os.Exit(1)
	}

	// Clean up
	os.Remove(asmFile)
	os.Remove("/tmp/petitgo_temp.o")

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

	// Generate ARM64 assembly directly
	generator := asmgen.NewAsmGenerator()
	assembly := generator.Generate(statements)
	runtimeCode := generator.GenerateRuntime()

	// Create temporary files properly
	tempDir := os.TempDir()

	// Create temp assembly file
	asmFile, err := os.CreateTemp(tempDir, "petitgo_run_*.s")
	if err != nil {
		fmt.Printf("Error creating temp assembly file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(asmFile.Name())

	// Write assembly content
	fullAsm := assembly + runtimeCode
	_, err = asmFile.WriteString(fullAsm)
	asmFile.Close()
	if err != nil {
		fmt.Printf("Error writing assembly file: %v\n", err)
		os.Exit(1)
	}

	// Create temp object file
	objFile, err := os.CreateTemp(tempDir, "petitgo_run_*.o")
	if err != nil {
		fmt.Printf("Error creating temp object file: %v\n", err)
		os.Exit(1)
	}
	objFile.Close()
	defer os.Remove(objFile.Name())

	// Create temp executable file
	execFile, err := os.CreateTemp(tempDir, "petitgo_run_*")
	if err != nil {
		fmt.Printf("Error creating temp executable file: %v\n", err)
		os.Exit(1)
	}
	execFile.Close()
	defer os.Remove(execFile.Name())

	// Assemble and link (platform-specific commands)
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// macOS ARM64
		cmd = exec.Command("as", "-arch", "arm64", "-o", objFile.Name(), asmFile.Name())
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		// Linux x86_64
		cmd = exec.Command("as", "--64", "-o", objFile.Name(), asmFile.Name())
	} else {
		fmt.Printf("Unsupported platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("Supported platforms: darwin/arm64, linux/amd64")
		os.Exit(1)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error assembling: %v\n", err)
		os.Exit(1)
	}

	// Platform-specific linker
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// macOS ARM64
		cmd = exec.Command("clang", "-arch", "arm64", "-o", execFile.Name(), objFile.Name())
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		// Linux x86_64
		cmd = exec.Command("ld", "-o", execFile.Name(), objFile.Name())
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error linking: %v\n", err)
		os.Exit(1)
	}

	// Execute the compiled binary
	cmd = exec.Command(execFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
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
