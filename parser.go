package main

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
	Operator Token
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

// パーサー構造体
type Parser struct {
	lexer        *Lexer
	currentToken TokenInfo
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.lexer.NextToken()
}

func (p *Parser) ParseStatement() Statement {
	switch p.currentToken.Type {
	case IDENT:
		if p.currentToken.Literal == "var" {
			return p.parseVarStatement()
		}
		return p.parseAssignStatement()
	default:
		// とりあえずダミーを返す
		return &VarStatement{}
	}
}

func (p *Parser) parseVarStatement() Statement {
	// var
	p.nextToken()

	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// 型名
	typeName := p.currentToken.Literal
	p.nextToken()

	// =
	p.nextToken()

	// 式
	value := p.ParseExpression()

	return &VarStatement{
		Name:     name,
		TypeName: typeName,
		Value:    value,
	}
}

func (p *Parser) parseAssignStatement() Statement {
	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// := or =
	p.nextToken()

	// 式
	value := p.ParseExpression()

	return &AssignStatement{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) ParseExpression() ASTNode {
	return p.parseTerm()
}

func (p *Parser) parseTerm() ASTNode {
	left := p.parseMultiplyDivide()

	for p.currentToken.Type == ADD || p.currentToken.Type == SUB {
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseMultiplyDivide()
		left = &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseMultiplyDivide() ASTNode {
	left := p.parseFactor()

	for p.currentToken.Type == MUL || p.currentToken.Type == QUO {
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseFactor()
		left = &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseFactor() ASTNode {
	if p.currentToken.Type == INT {
		value := 0
		literal := p.currentToken.Literal

		// 手動で文字列を数値に変換（strconv なしで）
		for i := 0; i < len(literal); i++ {
			digit := int(literal[i] - '0')
			value = value*10 + digit
		}

		p.nextToken()
		return &NumberNode{Value: value}
	}

	if p.currentToken.Type == IDENT {
		name := p.currentToken.Literal
		p.nextToken()
		return &VariableNode{Name: name}
	}

	if p.currentToken.Type == LPAREN {
		p.nextToken() // '(' を消費
		expr := p.ParseExpression()
		if p.currentToken.Type == RPAREN {
			p.nextToken() // ')' を消費
		}
		return expr
	}

	// エラーケース: とりあえず 0 を返す
	return &NumberNode{Value: 0}
}
