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
