package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestSlice_Literal(t *testing.T) {
	tests := []struct {
		input        string
		elementType  string
		elementCount int
		firstValue   int
		lastValue    int
	}{
		{"[]int{1, 2, 3}", "int", 3, 1, 3},
		{"[]int{42}", "int", 1, 42, 42},
		{"[]int{}", "int", 0, 0, 0},
		{"[]int{10, 20, 30, 40, 50}", "int", 5, 10, 50},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValue(expr)

		sliceVal, ok := result.(*SliceValue)
		if !ok {
			t.Fatalf("Input: %s, expected SliceValue, got %T", tt.input, result)
		}

		if sliceVal.ElementType != tt.elementType {
			t.Errorf("Input: %s, expected element type %s, got %s", tt.input, tt.elementType, sliceVal.ElementType)
		}

		if len(sliceVal.Elements) != tt.elementCount {
			t.Errorf("Input: %s, expected %d elements, got %d", tt.input, tt.elementCount, len(sliceVal.Elements))
		}

		// Check first and last values if slice is not empty
		if tt.elementCount > 0 {
			if firstInt, ok := sliceVal.Elements[0].(*IntValue); ok {
				if firstInt.Value != tt.firstValue {
					t.Errorf("Input: %s, expected first value %d, got %d", tt.input, tt.firstValue, firstInt.Value)
				}
			}

			if lastInt, ok := sliceVal.Elements[len(sliceVal.Elements)-1].(*IntValue); ok {
				if lastInt.Value != tt.lastValue {
					t.Errorf("Input: %s, expected last value %d, got %d", tt.input, tt.lastValue, lastInt.Value)
				}
			}
		}
	}
}

func TestSlice_Indexing(t *testing.T) {
	env := NewEnvironment()

	// Create a slice
	sliceInput := `numbers := []int{10, 20, 30, 40, 50}`
	s := scanner.NewScanner(sliceInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	tests := []struct {
		input    string
		expected int
	}{
		{"numbers[0]", 10},
		{"numbers[1]", 20},
		{"numbers[2]", 30},
		{"numbers[3]", 40},
		{"numbers[4]", 50},
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValueWithEnvironment(expr, env)

		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatalf("Input: %s, expected IntValue, got %T", tt.input, result)
		}

		if intVal.Value != tt.expected {
			t.Errorf("Input: %s, expected %d, got %d", tt.input, tt.expected, intVal.Value)
		}
	}
}

func TestSlice_OutOfBounds(t *testing.T) {
	env := NewEnvironment()

	// Create a slice
	sliceInput := `numbers := []int{10, 20, 30}`
	s := scanner.NewScanner(sliceInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	tests := []string{
		"numbers[-1]", // negative index
		"numbers[3]",  // index == length
		"numbers[10]", // index > length
	}

	for _, input := range tests {
		s := scanner.NewScanner(input)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValueWithEnvironment(expr, env)

		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatalf("Input: %s, expected IntValue for out of bounds, got %T", input, result)
		}

		if intVal.Value != 0 {
			t.Errorf("Input: %s, expected 0 for out of bounds, got %d", input, intVal.Value)
		}
	}
}

func TestSlice_LenFunction(t *testing.T) {
	tests := []struct {
		setup    string
		lenExpr  string
		expected int
	}{
		{`nums := []int{1, 2, 3}`, `len(nums)`, 3},
		{`empty := []int{}`, `len(empty)`, 0},
		{`str := "hello"`, `len(str)`, 5},
		{`big := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}`, `len(big)`, 10},
	}

	for _, tt := range tests {
		env := NewEnvironment()

		// Setup
		s := scanner.NewScanner(tt.setup)
		p := parser.NewParser(s)
		stmt := p.ParseStatement()
		EvalStatement(stmt, env)

		// Test len()
		s = scanner.NewScanner(tt.lenExpr)
		p = parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValueWithEnvironment(expr, env)

		intVal, ok := result.(*IntValue)
		if !ok {
			t.Fatalf("Setup: %s, Input: %s, expected IntValue, got %T", tt.setup, tt.lenExpr, result)
		}

		if intVal.Value != tt.expected {
			t.Errorf("Setup: %s, Input: %s, expected %d, got %d", tt.setup, tt.lenExpr, tt.expected, intVal.Value)
		}
	}
}

func TestSlice_AppendFunction(t *testing.T) {
	env := NewEnvironment()

	// Create initial slice
	setupInput := `nums := []int{1, 2, 3}`
	s := scanner.NewScanner(setupInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	// Append single element
	appendInput := `nums = append(nums, 4)`
	s = scanner.NewScanner(appendInput)
	p = parser.NewParser(s)
	stmt = p.ParseStatement()
	EvalStatement(stmt, env)

	// Check length
	lenInput := `len(nums)`
	s = scanner.NewScanner(lenInput)
	p = parser.NewParser(s)
	expr := p.ParseExpression()
	result := EvalValueWithEnvironment(expr, env)

	if intVal, ok := result.(*IntValue); ok {
		if intVal.Value != 4 {
			t.Errorf("Expected length 4 after append, got %d", intVal.Value)
		}
	}

	// Check last element
	lastInput := `nums[3]`
	s = scanner.NewScanner(lastInput)
	p = parser.NewParser(s)
	expr = p.ParseExpression()
	result = EvalValueWithEnvironment(expr, env)

	if intVal, ok := result.(*IntValue); ok {
		if intVal.Value != 4 {
			t.Errorf("Expected last element 4, got %d", intVal.Value)
		}
	}

	// Append multiple elements
	multiAppendInput := `nums = append(nums, 5, 6, 7)`
	s = scanner.NewScanner(multiAppendInput)
	p = parser.NewParser(s)
	stmt = p.ParseStatement()
	EvalStatement(stmt, env)

	// Check new length
	s = scanner.NewScanner(lenInput)
	p = parser.NewParser(s)
	expr = p.ParseExpression()
	result = EvalValueWithEnvironment(expr, env)

	if intVal, ok := result.(*IntValue); ok {
		if intVal.Value != 7 {
			t.Errorf("Expected length 7 after multiple append, got %d", intVal.Value)
		}
	}
}

func TestSlice_ComplexOperations(t *testing.T) {
	env := NewEnvironment()

	// Test slice with expressions
	input := `numbers := []int{1 + 1, 2 * 2, 3 * 3}`
	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	// Check values
	tests := []struct {
		expr     string
		expected int
	}{
		{"numbers[0]", 2}, // 1 + 1
		{"numbers[1]", 4}, // 2 * 2
		{"numbers[2]", 9}, // 3 * 3
	}

	for _, tt := range tests {
		s := scanner.NewScanner(tt.expr)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValueWithEnvironment(expr, env)

		if intVal, ok := result.(*IntValue); ok {
			if intVal.Value != tt.expected {
				t.Errorf("Expression %s: expected %d, got %d", tt.expr, tt.expected, intVal.Value)
			}
		}
	}

	// Test using variables in slice
	lines := []string{"x := 10", "y := 20", "z := 30", "vars := []int{x, y, z}"}
	for _, line := range lines {
		s = scanner.NewScanner(line)
		p = parser.NewParser(s)
		stmt = p.ParseStatement()
		EvalStatement(stmt, env)
	}

	// Check variable values in slice
	varTests := []struct {
		expr     string
		expected int
	}{
		{"vars[0]", 10},
		{"vars[1]", 20},
		{"vars[2]", 30},
	}

	for _, tt := range varTests {
		s := scanner.NewScanner(tt.expr)
		p := parser.NewParser(s)
		expr := p.ParseExpression()

		result := EvalValueWithEnvironment(expr, env)

		if intVal, ok := result.(*IntValue); ok {
			if intVal.Value != tt.expected {
				t.Errorf("Expression %s: expected %d, got %d", tt.expr, tt.expected, intVal.Value)
			}
		}
	}
}
