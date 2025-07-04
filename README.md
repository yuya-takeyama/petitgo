# petitgo

A minimal Go implementation aiming for self-hosting capability.

## Overview

petitgo is a small-scale Go language implementation that aims to eventually compile itself (self-hosting). Starting with minimal features, we incrementally add functionality following a test-driven development approach.

## Features

### ✅ Implemented

- **Phase 1: Minimal Calculator** - Four arithmetic operations with operator precedence
- **Phase 2: Variables and Statements** - Variable declarations, assignments, and basic statements
- **Phase 3: Control Flow** - if/else statements, for loops, break/continue
- **Phase 4: Functions** - Function definitions, calls, parameters, and return values
- **Phase 5: Type System** - Basic types (int, string, bool), type checking, and type inference
- **Phase 8: Native Compiler** - Direct ARM64 assembly generation (no Go dependency)
- **Advanced Features**:
  - Switch statements with case matching
  - Struct definitions and field access
  - Basic slice operations (len, append, indexing)
  - Comments (line and block comments)
  - Increment/decrement operators (++, --)
  - Compound assignment operators (+=, -=, *=, /=)
  - Complete for loops (init; condition; update)

### 🚧 In Progress

- **Phase 9: Self-hosting** - Compiling petitgo with petitgo itself

## Installation

```bash
git clone https://github.com/yuya-takeyama/petitgo.git
cd petitgo
go build -o petitgo
```

## Usage

### Interactive REPL

```bash
./petitgo
```

Example REPL session:
```
> 2 + 3 * 4
14
> x := 10
> x + 5
15
> if x > 5 { print("big") }
big
```

### Compile and Run Programs

Create a `.pg` file:
```go
// fibonacci.pg
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    for i := 0; i < 10; i++ {
        println(fibonacci(i))
    }
}
```

Compile and run:
```bash
# Direct execution
./petitgo run fibonacci.pg

# Compile to native binary
./petitgo build fibonacci.pg
./fibonacci
```

### Other Commands

```bash
# View AST as JSON
./petitgo ast fibonacci.pg

# Generate ARM64 assembly
./petitgo asm fibonacci.pg
```

## Examples

See the [examples/](examples/) directory for various sample programs:

- `hello.pg` - Simple hello world
- `calculator.pg` - Arithmetic operations
- `fibonacci.pg` - Recursive fibonacci implementation
- `switch_demo.pg` - Switch statement usage
- `struct_field_test.pg` - Struct field access
- And more...

## Architecture

petitgo consists of several packages:

- `token/` - Token definitions for lexical analysis
- `scanner/` - Lexical analyzer (tokenizer)
- `ast/` - Abstract Syntax Tree node definitions
- `parser/` - Syntax analyzer (parser)
- `eval/` - Expression evaluator with type system
- `asmgen/` - ARM64 assembly code generator
- `repl/` - Read-Eval-Print Loop implementation

## Compilation Process

1. **Source Code** (.pg files)
2. **Scanner** → Tokens
3. **Parser** → Abstract Syntax Tree (AST)
4. **ASM Generator** → ARM64 Assembly
5. **System Assembler** (as) → Object File
6. **System Linker** (clang) → Native Executable

No Go compiler dependency for the final binary!

## Supported Language Features

### Types
- `int` - Integer numbers
- `string` - String literals with escape sequences
- `bool` - Boolean values (true/false)
- `struct` - Structure types with fields

### Operators
- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Assignment: `:=`, `=`
- Compound: `+=`, `-=`, `*=`, `/=`
- Increment/Decrement: `++`, `--`

### Control Flow
- `if`/`else` statements
- `for` loops (condition-only and full form)
- `switch`/`case` statements
- `break`/`continue`
- `return` statements

### Functions
- Function definitions with parameters and return types
- Recursive function calls
- Built-in functions: `println()`, `len()`, `append()`

### Advanced Features
- Struct field access (`obj.field`)
- Basic slice operations
- Line comments (`//`) and block comments (`/* */`)
- Variable reassignment and compound operators

## Development

This project follows Test-Driven Development (TDD) principles and uses a Pull Request based workflow.

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./parser
go test ./eval
```

### Development Workflow

1. Create feature branch: `yuya-takeyama/feat/FEATURE_NAME`
2. Implement with TDD approach
3. Create draft Pull Request: `gh pr create --draft`
4. Complete implementation and tests
5. Remove draft status for review

## Roadmap

- **Phase 6**: Enhanced standard library functions
- **Phase 7**: Package system and imports
- **Phase 9**: Self-hosting capability

See [docs/ROADMAP.md](docs/ROADMAP.md) for detailed development phases.

## Contributing

Contributions are welcome! Please read our development guide in [CLAUDE.md](CLAUDE.md) for detailed instructions on:

- Test-Driven Development approach
- Code style guidelines
- Pull Request workflow
- Commit message conventions

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Inspired by the Go programming language
- Built with Test-Driven Development principles
- Aims for eventual self-hosting capability