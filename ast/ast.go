package ast

import "github.com/yuya-takeyama/petitgo/token"

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

// BooleanNode represents a boolean literal
type BooleanNode struct {
	Value bool
}

func (n *BooleanNode) String() string {
	return "BooleanNode"
}

// StringNode represents a string literal
type StringNode struct {
	Value string
}

func (n *StringNode) String() string {
	return "StringNode"
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

// VariableNode represents a variable reference
type VariableNode struct {
	Name string
}

func (n *VariableNode) String() string {
	return "VariableNode"
}

// CallNode represents a function call (function(args...))
type CallNode struct {
	Function  string
	Arguments []ASTNode
}

func (n *CallNode) String() string {
	return "CallNode"
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

// AssignStatement represents an assignment (x := 42 or x = 42)
type AssignStatement struct {
	Name  string
	Value ASTNode
}

func (n *AssignStatement) String() string {
	return "AssignStatement"
}

// BlockStatement represents a block of statements ({ ... })
type BlockStatement struct {
	Statements []Statement
}

func (n *BlockStatement) String() string {
	return "BlockStatement"
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

// BreakStatement represents a break statement
type BreakStatement struct{}

func (n *BreakStatement) String() string {
	return "BreakStatement"
}

// ContinueStatement represents a continue statement
type ContinueStatement struct{}

func (n *ContinueStatement) String() string {
	return "ContinueStatement"
}

// ExpressionStatement represents an expression as a statement
type ExpressionStatement struct {
	Expression ASTNode
}

func (n *ExpressionStatement) String() string {
	return "ExpressionStatement"
}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type string
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

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value ASTNode // nil for empty return
}

func (n *ReturnStatement) String() string {
	return "ReturnStatement"
}

// StructField represents a field in a struct definition
type StructField struct {
	Name string
	Type string
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

// StructLiteral represents a struct literal (Person{Name: "Alice", Age: 25})
type StructLiteral struct {
	TypeName string
	Fields   map[string]ASTNode // field name -> value expression
}

func (n *StructLiteral) String() string {
	return "StructLiteral"
}

// FieldAccess represents field access (person.Name)
type FieldAccess struct {
	Object ASTNode // the struct instance
	Field  string  // field name
}

func (n *FieldAccess) String() string {
	return "FieldAccess"
}

// SliceType represents a slice type ([]int, []string, etc.)
type SliceType struct {
	ElementType string
}

// SliceLiteral represents a slice literal ([]int{1, 2, 3})
type SliceLiteral struct {
	ElementType string
	Elements    []ASTNode
}

func (n *SliceLiteral) String() string {
	return "SliceLiteral"
}

// IndexAccess represents slice/array indexing (slice[0])
type IndexAccess struct {
	Object ASTNode // the slice/array
	Index  ASTNode // index expression
}

func (n *IndexAccess) String() string {
	return "IndexAccess"
}

// PackageStatement represents a package declaration (package main)
type PackageStatement struct {
	Name string // package name
}

func (n *PackageStatement) String() string {
	return "PackageStatement"
}

func (n *PackageStatement) statement() {}

// ImportStatement represents an import declaration (import "fmt")
type ImportStatement struct {
	Path string // import path like "fmt", "os"
}

func (n *ImportStatement) String() string {
	return "ImportStatement"
}

func (n *ImportStatement) statement() {}
