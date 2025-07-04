package eval

import "fmt"

// Value represents a value with a specific type in petitgo
type Value interface {
	Type() string
	String() string
	IsTruthy() bool
}

// IntValue represents an integer value
type IntValue struct {
	Value int
}

func (v *IntValue) Type() string   { return "int" }
func (v *IntValue) String() string { return fmt.Sprintf("%d", v.Value) }
func (v *IntValue) IsTruthy() bool { return v.Value != 0 }

// StringValue represents a string value
type StringValue struct {
	Value string
}

func (v *StringValue) Type() string   { return "string" }
func (v *StringValue) String() string { return v.Value }
func (v *StringValue) IsTruthy() bool { return len(v.Value) > 0 }

// BoolValue represents a boolean value
type BoolValue struct {
	Value bool
}

func (v *BoolValue) Type() string { return "bool" }
func (v *BoolValue) String() string {
	if v.Value {
		return "true"
	}
	return "false"
}
func (v *BoolValue) IsTruthy() bool { return v.Value }

// StructValue represents a struct instance
type StructValue struct {
	TypeName string
	Fields   map[string]Value
}

func (v *StructValue) Type() string { return v.TypeName }
func (v *StructValue) String() string {
	// For now, simple representation
	return fmt.Sprintf("%s{...}", v.TypeName)
}
func (v *StructValue) IsTruthy() bool { return v.Fields != nil }

// SliceValue represents a slice
type SliceValue struct {
	ElementType string
	Elements    []Value
}

func (v *SliceValue) Type() string { return "[]" + v.ElementType }
func (v *SliceValue) String() string {
	// For now, simple representation
	return fmt.Sprintf("[%d elements]", len(v.Elements))
}
func (v *SliceValue) IsTruthy() bool { return len(v.Elements) > 0 }
