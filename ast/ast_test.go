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
