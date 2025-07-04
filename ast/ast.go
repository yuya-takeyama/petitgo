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
