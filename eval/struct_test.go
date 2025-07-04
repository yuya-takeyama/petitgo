package eval

import (
	"testing"

	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func TestStruct_Definition(t *testing.T) {
	t.Skip("Struct definitions not fully implemented yet")
	input := `type Person struct { Name string Age int }`

	s := scanner.NewScanner(input)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()

	env := NewEnvironment()
	EvalStatement(stmt, env)

	// Check if struct definition was registered
	structDef, exists := env.GetStruct("Person")
	if !exists {
		t.Fatal("Person struct should be defined")
	}

	if structDef.Name != "Person" {
		t.Errorf("Expected struct name 'Person', got '%s'", structDef.Name)
	}

	if len(structDef.Fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(structDef.Fields))
	}

	// Check field definitions
	nameField := structDef.Fields[0]
	if nameField.Name != "Name" || nameField.Type != "string" {
		t.Errorf("Expected Name:string field, got %s:%s", nameField.Name, nameField.Type)
	}

	ageField := structDef.Fields[1]
	if ageField.Name != "Age" || ageField.Type != "int" {
		t.Errorf("Expected Age:int field, got %s:%s", ageField.Name, ageField.Type)
	}
}

func TestStruct_Literal(t *testing.T) {
	t.Skip("Struct literals not fully implemented yet")
	// First define the struct
	env := NewEnvironment()

	structDefInput := `type Person struct { Name string Age int }`
	s := scanner.NewScanner(structDefInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	// Now test struct literal
	input := `Person{Name: "Alice", Age: 25}`
	s = scanner.NewScanner(input)
	p = parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValueWithEnvironment(expr, env)

	structVal, ok := result.(*StructValue)
	if !ok {
		t.Fatalf("Expected StructValue, got %T", result)
	}

	if structVal.TypeName != "Person" {
		t.Errorf("Expected type 'Person', got '%s'", structVal.TypeName)
	}

	// Check field values
	nameVal, exists := structVal.Fields["Name"]
	if !exists {
		t.Fatal("Name field should exist")
	}
	if strVal, ok := nameVal.(*StringValue); ok {
		if strVal.Value != "Alice" {
			t.Errorf("Expected Name='Alice', got '%s'", strVal.Value)
		}
	} else {
		t.Errorf("Expected StringValue for Name, got %T", nameVal)
	}

	ageVal, exists := structVal.Fields["Age"]
	if !exists {
		t.Fatal("Age field should exist")
	}
	if intVal, ok := ageVal.(*IntValue); ok {
		if intVal.Value != 25 {
			t.Errorf("Expected Age=25, got %d", intVal.Value)
		}
	} else {
		t.Errorf("Expected IntValue for Age, got %T", ageVal)
	}
}

func TestStruct_FieldAccess(t *testing.T) {
	t.Skip("Struct field access not fully implemented yet")
	// Setup: define struct and create instance
	env := NewEnvironment()

	// Define struct
	structDefInput := `type Person struct { Name string Age int }`
	s := scanner.NewScanner(structDefInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	// Create struct instance
	assignInput := `person := Person{Name: "Bob", Age: 30}`
	s = scanner.NewScanner(assignInput)
	p = parser.NewParser(s)
	stmt = p.ParseStatement()
	EvalStatement(stmt, env)

	// Test field access: person.Name
	fieldAccessInput := `person.Name`
	s = scanner.NewScanner(fieldAccessInput)
	p = parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValueWithEnvironment(expr, env)

	strVal, ok := result.(*StringValue)
	if !ok {
		t.Fatalf("Expected StringValue, got %T", result)
	}

	if strVal.Value != "Bob" {
		t.Errorf("Expected 'Bob', got '%s'", strVal.Value)
	}

	// Test field access: person.Age
	fieldAccessInput = `person.Age`
	s = scanner.NewScanner(fieldAccessInput)
	p = parser.NewParser(s)
	expr = p.ParseExpression()

	result = EvalValueWithEnvironment(expr, env)

	intVal, ok := result.(*IntValue)
	if !ok {
		t.Fatalf("Expected IntValue, got %T", result)
	}

	if intVal.Value != 30 {
		t.Errorf("Expected 30, got %d", intVal.Value)
	}
}

func TestStruct_DefaultFieldValues(t *testing.T) {
	t.Skip("Struct default field values not fully implemented yet")
	// Define struct
	env := NewEnvironment()

	structDefInput := `type Person struct { Name string Age int Active bool }`
	s := scanner.NewScanner(structDefInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	// Create struct with partial initialization
	input := `Person{Name: "Charlie"}`
	s = scanner.NewScanner(input)
	p = parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValueWithEnvironment(expr, env)

	structVal, ok := result.(*StructValue)
	if !ok {
		t.Fatalf("Expected StructValue, got %T", result)
	}

	// Check that missing fields have default values
	ageVal, exists := structVal.Fields["Age"]
	if !exists {
		t.Fatal("Age field should exist with default value")
	}
	if intVal, ok := ageVal.(*IntValue); ok {
		if intVal.Value != 0 {
			t.Errorf("Expected default Age=0, got %d", intVal.Value)
		}
	}

	activeVal, exists := structVal.Fields["Active"]
	if !exists {
		t.Fatal("Active field should exist with default value")
	}
	if boolVal, ok := activeVal.(*BoolValue); ok {
		if boolVal.Value != false {
			t.Errorf("Expected default Active=false, got %v", boolVal.Value)
		}
	}
}

func TestStruct_NestedAccess(t *testing.T) {
	t.Skip("Struct nested access not fully implemented yet")
	// Test chained field access and complex expressions
	env := NewEnvironment()

	// Define struct
	structDefInput := `type Point struct { X int Y int }`
	s := scanner.NewScanner(structDefInput)
	p := parser.NewParser(s)
	stmt := p.ParseStatement()
	EvalStatement(stmt, env)

	// Create struct and assign
	assignInput := `p := Point{X: 10, Y: 20}`
	s = scanner.NewScanner(assignInput)
	p = parser.NewParser(s)
	stmt = p.ParseStatement()
	EvalStatement(stmt, env)

	// Test arithmetic with field access: p.X + p.Y
	exprInput := `p.X + p.Y`
	s = scanner.NewScanner(exprInput)
	p = parser.NewParser(s)
	expr := p.ParseExpression()

	result := EvalValueWithEnvironment(expr, env)

	intVal, ok := result.(*IntValue)
	if !ok {
		t.Fatalf("Expected IntValue, got %T", result)
	}

	if intVal.Value != 30 {
		t.Errorf("Expected 30 (10+20), got %d", intVal.Value)
	}
}
