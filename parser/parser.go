package parser

import (
	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/scanner"
	"github.com/yuya-takeyama/petitgo/token"
)

// パーサー構造体
type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.TokenInfo
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.lexer.NextToken()
}

func (p *Parser) ParseStatement() Statement {
	switch p.currentToken.Type {
	case lexer.IDENT:
		if p.currentToken.Literal == "var" {
			return p.parseVarStatement()
		}
		// 次のトークンが ASSIGN なら代入文、そうでなければ式文
		nextLexer := lexer.NewLexer(p.lexer.Input()[p.lexer.Position():])
		nextToken := nextLexer.NextToken()
		if nextToken.Type == lexer.ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.BREAK:
		return p.parseBreakStatement()
	case lexer.CONTINUE:
		return p.parseContinueStatement()
	case lexer.LBRACE:
		return p.parseBlockStatement()
	default:
		// 式文として処理
		return p.parseExpressionStatement()
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

func (p *Parser) parseIfStatement() Statement {
	// if
	p.nextToken()

	// condition
	condition := p.ParseExpression()

	// {
	if p.currentToken.Type != lexer.LBRACE {
		// エラー: とりあえずダミーを返す
		return &IfStatement{}
	}

	// then block
	thenBlock := p.parseBlockStatement()

	var elseBlock *BlockStatement = nil

	// else があるかチェック
	if p.currentToken.Type == lexer.ELSE {
		p.nextToken()
		if p.currentToken.Type == lexer.LBRACE {
			elseBlock = p.parseBlockStatement()
		}
	}

	return &IfStatement{
		Condition: condition,
		ThenBlock: thenBlock,
		ElseBlock: elseBlock,
	}
}

func (p *Parser) parseForStatement() Statement {
	// for
	p.nextToken()

	// for文の形式を判定
	// 1. for condition { ... } (condition-only)
	// 2. for init; condition; update { ... } (full form)

	// とりあえず condition-only form で実装
	condition := p.ParseExpression()

	// {
	if p.currentToken.Type != lexer.LBRACE {
		// エラー: とりあえずダミーを返す
		return &ForStatement{}
	}

	// body
	body := p.parseBlockStatement()

	return &ForStatement{
		Init:      nil,
		Condition: condition,
		Update:    nil,
		Body:      body,
	}
}

func (p *Parser) parseBreakStatement() Statement {
	// break
	p.nextToken()
	return &BreakStatement{}
}

func (p *Parser) parseContinueStatement() Statement {
	// continue
	p.nextToken()
	return &ContinueStatement{}
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	// {
	p.nextToken()

	statements := []Statement{}

	for p.currentToken.Type != lexer.RBRACE && p.currentToken.Type != lexer.EOF {
		stmt := p.ParseStatement()
		statements = append(statements, stmt)
	}

	// }
	if p.currentToken.Type == lexer.RBRACE {
		p.nextToken()
	}

	return &BlockStatement{Statements: statements}
}

func (p *Parser) parseExpressionStatement() Statement {
	expression := p.ParseExpression()
	return &ExpressionStatement{Expression: expression}
}

func (p *Parser) ParseExpression() ASTNode {
	return p.parseComparison()
}

// 比較演算子の解析 (==, !=, <, >, <=, >=)
func (p *Parser) parseComparison() ASTNode {
	left := p.parseTerm()

	for p.currentToken.Type == lexer.EQL || p.currentToken.Type == lexer.NEQ ||
		p.currentToken.Type == lexer.LSS || p.currentToken.Type == lexer.GTR ||
		p.currentToken.Type == lexer.LEQ || p.currentToken.Type == lexer.GEQ {
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseTerm()
		left = &BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseTerm() ASTNode {
	left := p.parseMultiplyDivide()

	for p.currentToken.Type == lexer.ADD || p.currentToken.Type == lexer.SUB {
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

	for p.currentToken.Type == lexer.MUL || p.currentToken.Type == lexer.QUO {
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
	if p.currentToken.Type == lexer.INT {
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

	if p.currentToken.Type == lexer.IDENT {
		name := p.currentToken.Literal
		p.nextToken()

		// 関数呼び出しかチェック
		if p.currentToken.Type == lexer.LPAREN {
			p.nextToken() // '(' を消費

			arguments := []ASTNode{}

			// 引数をパース
			for p.currentToken.Type != lexer.RPAREN && p.currentToken.Type != lexer.EOF {
				arg := p.ParseExpression()
				arguments = append(arguments, arg)
				// カンマは今回スキップ（簡単のため）
			}

			if p.currentToken.Type == lexer.RPAREN {
				p.nextToken() // ')' を消費
			}

			return &CallNode{Function: name, Arguments: arguments}
		}

		return &VariableNode{Name: name}
	}

	if p.currentToken.Type == lexer.LPAREN {
		p.nextToken() // '(' を消費
		expr := p.ParseExpression()
		if p.currentToken.Type == lexer.RPAREN {
			p.nextToken() // ')' を消費
		}
		return expr
	}

	// エラーケース: とりあえず 0 を返す
	return &NumberNode{Value: 0}
}
