package ast

import (
	"encoding/json"
	"fmt"

	"github.com/yuya-takeyama/petitgo/token"
)

// ASTNode interface for all AST nodes
type ASTNode interface {
	String() string
}

// NumberNode represents a numeric literal
type NumberNode struct {
	Value int
}

func (n *NumberNode) String() string {
	return "NumberNode"
}

func (n *NumberNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "NumberNode",
		"value": n.Value,
	})
}

// BooleanNode represents a boolean literal
type BooleanNode struct {
	Value bool
}

func (n *BooleanNode) String() string {
	return "BooleanNode"
}

func (n *BooleanNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "BooleanNode",
		"value": n.Value,
	})
}

// StringNode represents a string literal
type StringNode struct {
	Value string
}

func (n *StringNode) String() string {
	return "StringNode"
}

func (n *StringNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "StringNode",
		"value": n.Value,
	})
}

// BinaryOpNode represents a binary operation
type BinaryOpNode struct {
	Left     ASTNode
	Operator token.Token
	Right    ASTNode
}

func (n *BinaryOpNode) String() string {
	return "BinaryOpNode"
}

func (n *BinaryOpNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     "BinaryOpNode",
		"left":     n.Left,
		"operator": tokenToString(n.Operator),
		"right":    n.Right,
	})
}

// VariableNode represents a variable reference
type VariableNode struct {
	Name string
}

func (n *VariableNode) String() string {
	return "VariableNode"
}

func (n *VariableNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "VariableNode",
		"name": n.Name,
	})
}

// CallNode represents a function call (function(args...))
type CallNode struct {
	Function  string
	Arguments []ASTNode
}

func (n *CallNode) String() string {
	return "CallNode"
}

func (n *CallNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":      "CallNode",
		"function":  n.Function,
		"arguments": n.Arguments,
	})
}

// Statement interface for all statement nodes
type Statement interface {
	String() string
}

// VarStatement represents a variable declaration (var x int = 42)
type VarStatement struct {
	Name     string
	TypeName string
	Value    ASTNode
}

func (n *VarStatement) String() string {
	return "VarStatement"
}

func (n *VarStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     "VarStatement",
		"name":     n.Name,
		"typeName": n.TypeName,
		"value":    n.Value,
	})
}

// AssignStatement represents an assignment (x := 42 or x = 42)
type AssignStatement struct {
	Name  string
	Value ASTNode
}

func (n *AssignStatement) String() string {
	return "AssignStatement"
}

func (n *AssignStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "AssignStatement",
		"name":  n.Name,
		"value": n.Value,
	})
}

// ReassignStatement represents a variable reassignment (x = 42)
type ReassignStatement struct {
	Name  string
	Value ASTNode
}

func (n *ReassignStatement) String() string {
	return "ReassignStatement"
}

func (n *ReassignStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "ReassignStatement",
		"name":  n.Name,
		"value": n.Value,
	})
}

// CompoundAssignStatement represents compound assignment (x += y, x -= y, etc.)
type CompoundAssignStatement struct {
	Name     string
	Operator token.Token // ADD_ASSIGN, SUB_ASSIGN, etc.
	Value    ASTNode
}

func (n *CompoundAssignStatement) String() string {
	return "CompoundAssignStatement"
}

func (n *CompoundAssignStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     "CompoundAssignStatement",
		"name":     n.Name,
		"operator": tokenToString(n.Operator),
		"value":    n.Value,
	})
}

// IncStatement represents an increment statement (x++)
type IncStatement struct {
	Name string
}

func (n *IncStatement) String() string {
	return "IncStatement"
}

func (n *IncStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "IncStatement",
		"name": n.Name,
	})
}

// DecStatement represents a decrement statement (x--)
type DecStatement struct {
	Name string
}

func (n *DecStatement) String() string {
	return "DecStatement"
}

func (n *DecStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "DecStatement",
		"name": n.Name,
	})
}

// SwitchStatement represents a switch statement
type SwitchStatement struct {
	Value   ASTNode
	Cases   []*CaseStatement
	Default *BlockStatement
}

func (n *SwitchStatement) String() string {
	return "SwitchStatement"
}

func (n *SwitchStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "SwitchStatement",
		"value":   n.Value,
		"cases":   n.Cases,
		"default": n.Default,
	})
}

// CaseStatement represents a case in a switch
type CaseStatement struct {
	Value ASTNode
	Body  *BlockStatement
}

func (n *CaseStatement) String() string {
	return "CaseStatement"
}

func (n *CaseStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "CaseStatement",
		"value": n.Value,
		"body":  n.Body,
	})
}

// BlockStatement represents a block of statements ({ ... })
type BlockStatement struct {
	Statements []Statement
}

func (n *BlockStatement) String() string {
	return "BlockStatement"
}

func (n *BlockStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       "BlockStatement",
		"statements": n.Statements,
	})
}

// IfStatement represents an if statement (if condition { block } else { block })
type IfStatement struct {
	Condition ASTNode
	ThenBlock *BlockStatement
	ElseBlock *BlockStatement // nil if no else
}

func (n *IfStatement) String() string {
	return "IfStatement"
}

func (n *IfStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":      "IfStatement",
		"condition": n.Condition,
		"thenBlock": n.ThenBlock,
		"elseBlock": n.ElseBlock,
	})
}

// ForStatement represents a for loop (for init; condition; update { block })
type ForStatement struct {
	Init      Statement // nil for condition-only for loops
	Condition ASTNode   // nil for infinite loops
	Update    Statement // nil for condition-only for loops
	Body      *BlockStatement
}

func (n *ForStatement) String() string {
	return "ForStatement"
}

func (n *ForStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":      "ForStatement",
		"init":      n.Init,
		"condition": n.Condition,
		"update":    n.Update,
		"body":      n.Body,
	})
}

// BreakStatement represents a break statement
type BreakStatement struct{}

func (n *BreakStatement) String() string {
	return "BreakStatement"
}

func (n *BreakStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "BreakStatement",
	})
}

// ContinueStatement represents a continue statement
type ContinueStatement struct{}

func (n *ContinueStatement) String() string {
	return "ContinueStatement"
}

func (n *ContinueStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "ContinueStatement",
	})
}

// ExpressionStatement represents an expression as a statement
type ExpressionStatement struct {
	Expression ASTNode
}

func (n *ExpressionStatement) String() string {
	return "ExpressionStatement"
}

func (n *ExpressionStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       "ExpressionStatement",
		"expression": n.Expression,
	})
}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type string
}

func (p *Parameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name": p.Name,
		"type": p.Type,
	})
}

// FuncStatement represents a function definition
type FuncStatement struct {
	Name       string
	Parameters []Parameter
	ReturnType string
	Body       *BlockStatement
}

func (n *FuncStatement) String() string {
	return "FuncStatement"
}

func (n *FuncStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       "FuncStatement",
		"name":       n.Name,
		"parameters": n.Parameters,
		"returnType": n.ReturnType,
		"body":       n.Body,
	})
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value ASTNode // nil for empty return
}

func (n *ReturnStatement) String() string {
	return "ReturnStatement"
}

func (n *ReturnStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "ReturnStatement",
		"value": n.Value,
	})
}

// StructField represents a field in a struct definition
type StructField struct {
	Name string
	Type string
}

func (sf *StructField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name": sf.Name,
		"type": sf.Type,
	})
}

// StructDefinition represents a struct type definition (type Person struct { ... })
type StructDefinition struct {
	Name   string
	Fields []StructField
}

func (n *StructDefinition) String() string {
	return "StructDefinition"
}

func (n *StructDefinition) statement() {}

func (n *StructDefinition) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   "StructDefinition",
		"name":   n.Name,
		"fields": n.Fields,
	})
}

// StructLiteral represents a struct literal (Person{Name: "Alice", Age: 25})
type StructLiteral struct {
	TypeName string
	Fields   map[string]ASTNode // field name -> value expression
}

func (n *StructLiteral) String() string {
	return "StructLiteral"
}

func (n *StructLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     "StructLiteral",
		"typeName": n.TypeName,
		"fields":   n.Fields,
	})
}

// FieldAccess represents field access (person.Name)
type FieldAccess struct {
	Object ASTNode // the struct instance
	Field  string  // field name
}

func (n *FieldAccess) String() string {
	return "FieldAccess"
}

func (n *FieldAccess) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   "FieldAccess",
		"object": n.Object,
		"field":  n.Field,
	})
}

// SliceType represents a slice type ([]int, []string, etc.)
type SliceType struct {
	ElementType string
}

func (st *SliceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "SliceType",
		"elementType": st.ElementType,
	})
}

// SliceLiteral represents a slice literal ([]int{1, 2, 3})
type SliceLiteral struct {
	ElementType string
	Elements    []ASTNode
}

func (n *SliceLiteral) String() string {
	return "SliceLiteral"
}

func (n *SliceLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "SliceLiteral",
		"elementType": n.ElementType,
		"elements":    n.Elements,
	})
}

// IndexAccess represents slice/array indexing (slice[0])
type IndexAccess struct {
	Object ASTNode // the slice/array
	Index  ASTNode // index expression
}

func (n *IndexAccess) String() string {
	return "IndexAccess"
}

func (n *IndexAccess) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   "IndexAccess",
		"object": n.Object,
		"index":  n.Index,
	})
}

// PackageStatement represents a package declaration (package main)
type PackageStatement struct {
	Name string // package name
}

func (n *PackageStatement) String() string {
	return "PackageStatement"
}

func (n *PackageStatement) statement() {}

func (n *PackageStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "PackageStatement",
		"name": n.Name,
	})
}

// ImportStatement represents an import declaration (import "fmt")
type ImportStatement struct {
	Path string // import path like "fmt", "os"
}

func (n *ImportStatement) String() string {
	return "ImportStatement"
}

func (n *ImportStatement) statement() {}

func (n *ImportStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "ImportStatement",
		"path": n.Path,
	})
}

// tokenToString converts a token to its string representation
func tokenToString(tok token.Token) string {
	switch tok {
	case token.ADD:
		return "+"
	case token.SUB:
		return "-"
	case token.MUL:
		return "*"
	case token.QUO:
		return "/"
	case token.EQL:
		return "=="
	case token.NEQ:
		return "!="
	case token.LSS:
		return "<"
	case token.GTR:
		return ">"
	case token.LEQ:
		return "<="
	case token.GEQ:
		return ">="
	case token.LAND:
		return "&&"
	case token.LOR:
		return "||"
	case token.NOT:
		return "!"
	case token.ADD_ASSIGN:
		return "+="
	case token.SUB_ASSIGN:
		return "-="
	case token.MUL_ASSIGN:
		return "*="
	case token.QUO_ASSIGN:
		return "/="
	default:
		return fmt.Sprintf("TOKEN_%d", int(tok))
	}
}
