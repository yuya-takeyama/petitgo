# petitgo Examples

This directory contains example programs written in petitgo to demonstrate various language features.

## Running Examples

### Using the build command

```bash
# Compile to executable
./petitgo build examples/hello.pg
./hello

# Compile and run more complex examples
./petitgo build examples/calculator.pg
./calculator
```

### Using the run command (compile and run directly)

```bash
./petitgo run examples/hello.pg
./petitgo run examples/calculator.pg
./petitgo run examples/variables.pg
./petitgo run examples/conditionals.pg
./petitgo run examples/functions.pg
./petitgo run examples/fibonacci.pg
```

## Example Programs

### hello.pg

Simple variable assignment and printing.

### calculator.pg

Demonstrates arithmetic operations and operator precedence.

### variables.pg

Shows variable declarations, type annotations, and reassignments.

### conditionals.pg

If/else statements and comparison operators.

### functions.pg

Function definitions, parameters, return values, and function calls.

### fibonacci.pg

Recursive fibonacci calculation with for loops.

## Working Examples

✅ **hello.pg**: `42`  
✅ **calculator.pg**: Outputs `15`, `5`, `50`, `2`, `20`, `30` (each on new line)  
✅ **variables.pg**: Outputs `52`, `10`, `100` (each on new line)  
✅ **simple_if.pg**: `1`  
✅ **functions.pg**: Outputs `8`, `15`, `10` (each on new line)

## Partially Working

⚠️ **conditionals.pg**: Complex if/else chains may cause infinite loops  
⚠️ **fibonacci.pg**: For loops not fully implemented yet

## Notes

- Uses Go's built-in `println()` function for output with automatic newlines
- Variable reassignment (`x = y`) not implemented yet, use new variables
- Comments (`//`) not implemented yet
- For loops need more work
