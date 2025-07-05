package ast

import (
	"encoding/json"
	"testing"

	"github.com/yuya-takeyama/petitgo/token"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		node interface{}
		want map[string]interface{}
	}{
		{
			name: "NumberNode",
			node: &NumberNode{Value: 42},
			want: map[string]interface{}{
				"type":  "NumberNode",
				"value": float64(42), // JSON numbers are float64
			},
		},
		{
			name: "VarStatement",
			node: &VarStatement{
				Name:     "x",
				TypeName: "int",
				Value:    &NumberNode{Value: 42},
			},
			want: map[string]interface{}{
				"type":     "VarStatement",
				"name":     "x",
				"typeName": "int",
				"value": map[string]interface{}{
					"type":  "NumberNode",
					"value": float64(42),
				},
			},
		},
		{
			name: "AssignStatement",
			node: &AssignStatement{
				Name:  "x",
				Value: &NumberNode{Value: 10},
			},
			want: map[string]interface{}{
				"type": "AssignStatement",
				"name": "x",
				"value": map[string]interface{}{
					"type":  "NumberNode",
					"value": float64(10),
				},
			},
		},
		{
			name: "FuncStatement",
			node: &FuncStatement{
				Name: "add",
				Parameters: []Parameter{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "int"},
				},
				ReturnType: "int",
				Body: &BlockStatement{
					Statements: []Statement{},
				},
			},
			want: map[string]interface{}{
				"type": "FuncStatement",
				"name": "add",
				"parameters": []interface{}{
					map[string]interface{}{"name": "a", "type": "int"},
					map[string]interface{}{"name": "b", "type": "int"},
				},
				"returnType": "int",
				"body": map[string]interface{}{
					"type":       "BlockStatement",
					"statements": []interface{}{},
				},
			},
		},
		{
			name: "PackageStatement",
			node: &PackageStatement{Name: "main"},
			want: map[string]interface{}{
				"type": "PackageStatement",
				"name": "main",
			},
		},
		{
			name: "ImportStatement",
			node: &ImportStatement{Path: "fmt"},
			want: map[string]interface{}{
				"type": "ImportStatement",
				"path": "fmt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.node)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			// Unmarshal back to map for comparison
			var got map[string]interface{}
			if err := json.Unmarshal(jsonData, &got); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			// Compare JSON structure
			if !deepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryOpNodeMarshalJSON(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 1},
		Operator: token.ADD,
		Right:    &NumberNode{Value: 2},
	}

	jsonData, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(jsonData, &got); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	want := map[string]interface{}{
		"type":     "BinaryOpNode",
		"operator": "+",
		"left": map[string]interface{}{
			"type":  "NumberNode",
			"value": float64(1),
		},
		"right": map[string]interface{}{
			"type":  "NumberNode",
			"value": float64(2),
		},
	}

	if !deepEqual(got, want) {
		t.Errorf("BinaryOpNode MarshalJSON() = %v, want %v", got, want)
	}
}

// Simple deep equality check for maps
func deepEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !valueEqual(v, bv) {
			return false
		}
	}
	return true
}

func valueEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case map[string]interface{}:
		if bv, ok := b.(map[string]interface{}); ok {
			return deepEqual(av, bv)
		}
	case []interface{}:
		if bv, ok := b.([]interface{}); ok {
			if len(av) != len(bv) {
				return false
			}
			for i, item := range av {
				if !valueEqual(item, bv[i]) {
					return false
				}
			}
			return true
		}
	default:
		return a == b
	}
	return false
}

// Test all String() methods for 100% coverage
func TestStringMethods(t *testing.T) {
	tests := []struct {
		name string
		node interface{ String() string }
		want string
	}{
		{"NumberNode", &NumberNode{Value: 42}, "NumberNode"},
		{"StringNode", &StringNode{Value: "hello"}, "StringNode"},
		{"BooleanNode", &BooleanNode{Value: true}, "BooleanNode"},
		{"VariableNode", &VariableNode{Name: "x"}, "VariableNode"},
		{"BinaryOpNode", &BinaryOpNode{}, "BinaryOpNode"},
		{"CallNode", &CallNode{Function: "test"}, "CallNode"},
		{"VarStatement", &VarStatement{Name: "x"}, "VarStatement"},
		{"AssignStatement", &AssignStatement{Name: "x"}, "AssignStatement"},
		{"ReassignStatement", &ReassignStatement{Name: "x"}, "ReassignStatement"},
		{"CompoundAssignStatement", &CompoundAssignStatement{Name: "x"}, "CompoundAssignStatement"},
		{"IncStatement", &IncStatement{Name: "x"}, "IncStatement"},
		{"DecStatement", &DecStatement{Name: "x"}, "DecStatement"},
		{"SwitchStatement", &SwitchStatement{}, "SwitchStatement"},
		{"CaseStatement", &CaseStatement{}, "CaseStatement"},
		{"IfStatement", &IfStatement{}, "IfStatement"},
		{"ForStatement", &ForStatement{}, "ForStatement"},
		{"BreakStatement", &BreakStatement{}, "BreakStatement"},
		{"ContinueStatement", &ContinueStatement{}, "ContinueStatement"},
		{"BlockStatement", &BlockStatement{}, "BlockStatement"},
		{"ExpressionStatement", &ExpressionStatement{}, "ExpressionStatement"},
		{"FuncStatement", &FuncStatement{Name: "test"}, "FuncStatement"},
		{"ReturnStatement", &ReturnStatement{}, "ReturnStatement"},
		{"StructLiteral", &StructLiteral{TypeName: "Person"}, "StructLiteral"},
		{"SliceLiteral", &SliceLiteral{ElementType: "int"}, "SliceLiteral"},
		{"FieldAccessNode", &FieldAccessNode{Field: "name"}, "FieldAccessNode"},
		{"IndexAccess", &IndexAccess{}, "IndexAccess"},
		{"StructDefinition", &StructDefinition{Name: "Person"}, "StructDefinition"},
		{"TypeStatement", &TypeStatement{Name: "Person"}, "TypeStatement"},
		{"PackageStatement", &PackageStatement{Name: "main"}, "PackageStatement"},
		{"ImportStatement", &ImportStatement{Path: "fmt"}, "ImportStatement"},
		{"ArrayLiteral", &ArrayLiteral{ElementType: "int"}, "ArrayLiteral"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.String()
			if result != tt.want {
				t.Errorf("%s.String() = %v, want %v", tt.name, result, tt.want)
			}
		})
	}
}

// Test Statement interface compliance
func TestStatementInterface(t *testing.T) {
	statements := []Statement{
		&VarStatement{},
		&AssignStatement{},
		&ReassignStatement{},
		&CompoundAssignStatement{},
		&IncStatement{},
		&DecStatement{},
		&SwitchStatement{},
		&IfStatement{},
		&ForStatement{},
		&BreakStatement{},
		&ContinueStatement{},
		&BlockStatement{},
		&ExpressionStatement{},
		&FuncStatement{},
		&ReturnStatement{},
		&StructDefinition{},
		&TypeStatement{},
		&PackageStatement{},
		&ImportStatement{},
	}

	// Just ensure all statements implement the Statement interface
	for i, stmt := range statements {
		if stmt == nil {
			t.Errorf("Statement at index %d is nil", i)
		}
	}
}

// Test remaining MarshalJSON methods
func TestRemainingMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		node interface{}
	}{
		{"StringNode", &StringNode{Value: "test"}},
		{"BooleanNode", &BooleanNode{Value: true}},
		{"VariableNode", &VariableNode{Name: "x"}},
		{"BinaryOpNode", &BinaryOpNode{Left: &NumberNode{Value: 1}, Operator: token.ADD, Right: &NumberNode{Value: 2}}},
		{"CallNode", &CallNode{Function: "test", Arguments: []ASTNode{}}},
		{"ReassignStatement", &ReassignStatement{Name: "x", Value: &NumberNode{Value: 1}}},
		{"CompoundAssignStatement", &CompoundAssignStatement{Name: "x", Operator: token.ADD, Value: &NumberNode{Value: 1}}},
		{"IncStatement", &IncStatement{Name: "x"}},
		{"DecStatement", &DecStatement{Name: "x"}},
		{"SwitchStatement", &SwitchStatement{Value: &VariableNode{Name: "x"}, Cases: []*CaseStatement{}}},
		{"CaseStatement", &CaseStatement{Value: &NumberNode{Value: 1}, Body: &BlockStatement{}}},
		{"IfStatement", &IfStatement{Condition: &BooleanNode{Value: true}, ThenBlock: &BlockStatement{}}},
		{"ForStatement", &ForStatement{Condition: &BooleanNode{Value: true}, Body: &BlockStatement{}}},
		{"BreakStatement", &BreakStatement{}},
		{"ContinueStatement", &ContinueStatement{}},
		{"ExpressionStatement", &ExpressionStatement{Expression: &NumberNode{Value: 1}}},
		{"ReturnStatement", &ReturnStatement{Value: &NumberNode{Value: 1}}},
		{"StructLiteral", &StructLiteral{TypeName: "Person", Fields: map[string]ASTNode{}}},
		{"SliceLiteral", &SliceLiteral{ElementType: "int", Elements: []ASTNode{}}},
		{"FieldAccessNode", &FieldAccessNode{Object: &VariableNode{Name: "x"}, Field: "name"}},
		{"IndexAccess", &IndexAccess{Object: &VariableNode{Name: "arr"}, Index: &NumberNode{Value: 0}}},
		{"StructDefinition", &StructDefinition{Name: "Person", Fields: []StructField{}}},
		{"TypeStatement", &TypeStatement{Name: "Person", Fields: []*FieldDef{}}},
		{"ArrayLiteral", &ArrayLiteral{ElementType: "int", Size: 10, Elements: []ASTNode{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.node)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", tt.name, err)
			}

			// Unmarshal to ensure valid JSON
			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", tt.name, err)
			}

			// Check that type field is correct
			if typeField, ok := result["type"]; !ok || typeField != tt.name {
				t.Errorf("%s MarshalJSON type field = %v, want %v", tt.name, typeField, tt.name)
			}
		})
	}
}

// Test tokenToString function
func TestTokenToString(t *testing.T) {
	tests := []struct {
		name  string
		token token.Token
		want  string
	}{
		{"ADD", token.ADD, "+"},
		{"SUB", token.SUB, "-"},
		{"MUL", token.MUL, "*"},
		{"QUO", token.QUO, "/"},
		{"EQL", token.EQL, "=="},
		{"NEQ", token.NEQ, "!="},
		{"LSS", token.LSS, "<"},
		{"GTR", token.GTR, ">"},
		{"LEQ", token.LEQ, "<="},
		{"GEQ", token.GEQ, ">="},
		{"ADD_ASSIGN", token.ADD_ASSIGN, "+="},
		{"SUB_ASSIGN", token.SUB_ASSIGN, "-="},
		{"MUL_ASSIGN", token.MUL_ASSIGN, "*="},
		{"QUO_ASSIGN", token.QUO_ASSIGN, "/="},
		{"LAND", token.LAND, "&&"},
		{"LOR", token.LOR, "||"},
		{"NOT", token.NOT, "!"},
		{"EOF", token.EOF, "TOKEN_1"}, // EOF is handled by default case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenToString(tt.token)
			if result != tt.want {
				t.Errorf("tokenToString(%v) = %v, want %v", tt.token, result, tt.want)
			}
		})
	}
}
