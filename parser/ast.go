package parser

import "github.com/yuya-takeyama/petitgo/lexer"

// AST ノードのインターフェース
type ASTNode interface {
	String() string
}

// 数値ノード
type NumberNode struct {
	Value int
}

func (n *NumberNode) String() string {
	return "NumberNode"
}

// 二項演算ノード
type BinaryOpNode struct {
	Left     ASTNode
	Operator lexer.Token
	Right    ASTNode
}

func (n *BinaryOpNode) String() string {
	return "BinaryOpNode"
}

// 変数参照ノード
type VariableNode struct {
	Name string
}

func (n *VariableNode) String() string {
	return "VariableNode"
}

// 関数呼び出しノード (function(args...))
type CallNode struct {
	Function  string
	Arguments []ASTNode
}

func (n *CallNode) String() string {
	return "CallNode"
}

// Statement インターフェース
type Statement interface {
	String() string
}

// 変数宣言文ノード (var x int = 42)
type VarStatement struct {
	Name     string
	TypeName string
	Value    ASTNode
}

func (n *VarStatement) String() string {
	return "VarStatement"
}

// 代入文ノード (x := 42 or x = 42)
type AssignStatement struct {
	Name  string
	Value ASTNode
}

func (n *AssignStatement) String() string {
	return "AssignStatement"
}

// ブロック文ノード ({ ... })
type BlockStatement struct {
	Statements []Statement
}

func (n *BlockStatement) String() string {
	return "BlockStatement"
}

// If文ノード (if condition { block } else { block })
type IfStatement struct {
	Condition ASTNode
	ThenBlock *BlockStatement
	ElseBlock *BlockStatement // nil if no else
}

func (n *IfStatement) String() string {
	return "IfStatement"
}

// For文ノード (for init; condition; update { block })
type ForStatement struct {
	Init      Statement // nil for condition-only for loops
	Condition ASTNode   // nil for infinite loops
	Update    Statement // nil for condition-only for loops
	Body      *BlockStatement
}

func (n *ForStatement) String() string {
	return "ForStatement"
}

// Break文ノード
type BreakStatement struct{}

func (n *BreakStatement) String() string {
	return "BreakStatement"
}

// Continue文ノード
type ContinueStatement struct{}

func (n *ContinueStatement) String() string {
	return "ContinueStatement"
}

// 式文ノード (expression;)
type ExpressionStatement struct {
	Expression ASTNode
}

func (n *ExpressionStatement) String() string {
	return "ExpressionStatement"
}
