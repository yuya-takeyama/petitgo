package parser

import (
	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/scanner"
	"github.com/yuya-takeyama/petitgo/token"
)

// Parser struct for parsing tokens into AST
type Parser struct {
	scanner      *scanner.Scanner
	currentToken token.TokenInfo
}

func NewParser(s *scanner.Scanner) *Parser {
	p := &Parser{scanner: s}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.scanner.NextToken()
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.IDENT:
		if p.currentToken.Literal == "var" {
			return p.parseVarStatement()
		}
		// 次のトークンが ASSIGN なら代入文、そうでなければ式文
		nextScanner := scanner.NewScanner(p.scanner.Input()[p.scanner.Position():])
		nextToken := nextScanner.NextToken()
		if nextToken.Type == token.ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		// 式文として処理
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() ast.Statement {
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

	return &ast.VarStatement{
		Name:     name,
		TypeName: typeName,
		Value:    value,
	}
}

func (p *Parser) parseAssignStatement() ast.Statement {
	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// := or =
	p.nextToken()

	// 式
	value := p.ParseExpression()

	return &ast.AssignStatement{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseIfStatement() ast.Statement {
	// if
	p.nextToken()

	// condition
	condition := p.ParseExpression()

	// {
	if p.currentToken.Type != token.LBRACE {
		// エラー: とりあえずダミーを返す
		return &ast.IfStatement{}
	}

	// then block
	thenBlock := p.parseBlockStatement()

	var elseBlock *ast.BlockStatement = nil

	// else があるかチェック
	if p.currentToken.Type == token.ELSE {
		p.nextToken()
		if p.currentToken.Type == token.LBRACE {
			elseBlock = p.parseBlockStatement()
		}
	}

	return &ast.IfStatement{
		Condition: condition,
		ThenBlock: thenBlock,
		ElseBlock: elseBlock,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	// for
	p.nextToken()

	// for文の形式を判定
	// 1. for condition { ... } (condition-only)
	// 2. for init; condition; update { ... } (full form)

	// とりあえず condition-only form で実装
	condition := p.ParseExpression()

	// {
	if p.currentToken.Type != token.LBRACE {
		// エラー: とりあえずダミーを返す
		return &ast.ForStatement{}
	}

	// body
	body := p.parseBlockStatement()

	return &ast.ForStatement{
		Init:      nil,
		Condition: condition,
		Update:    nil,
		Body:      body,
	}
}

func (p *Parser) parseBreakStatement() ast.Statement {
	// break
	p.nextToken()
	return &ast.BreakStatement{}
}

func (p *Parser) parseContinueStatement() ast.Statement {
	// continue
	p.nextToken()
	return &ast.ContinueStatement{}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// {
	p.nextToken()

	statements := []ast.Statement{}

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		stmt := p.ParseStatement()
		statements = append(statements, stmt)
	}

	// }
	if p.currentToken.Type == token.RBRACE {
		p.nextToken()
	}

	return &ast.BlockStatement{Statements: statements}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expression := p.ParseExpression()
	return &ast.ExpressionStatement{Expression: expression}
}

func (p *Parser) ParseExpression() ast.ASTNode {
	return p.parseComparison()
}

// 比較演算子の解析 (==, !=, <, >, <=, >=)
func (p *Parser) parseComparison() ast.ASTNode {
	left := p.parseTerm()

	for p.currentToken.Type == token.EQL || p.currentToken.Type == token.NEQ ||
		p.currentToken.Type == token.LSS || p.currentToken.Type == token.GTR ||
		p.currentToken.Type == token.LEQ || p.currentToken.Type == token.GEQ {
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseTerm()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseTerm() ast.ASTNode {
	left := p.parseMultiplyDivide()

	for p.currentToken.Type == token.ADD || p.currentToken.Type == token.SUB {
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseMultiplyDivide()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseMultiplyDivide() ast.ASTNode {
	left := p.parseFactor()

	for p.currentToken.Type == token.MUL || p.currentToken.Type == token.QUO {
		operator := p.currentToken.Type
		p.nextToken()
		right := p.parseFactor()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseFactor() ast.ASTNode {
	if p.currentToken.Type == token.INT {
		value := 0
		literal := p.currentToken.Literal

		// 手動で文字列を数値に変換（strconv なしで）
		for i := 0; i < len(literal); i++ {
			digit := int(literal[i] - '0')
			value = value*10 + digit
		}

		p.nextToken()
		return &ast.NumberNode{Value: value}
	}

	if p.currentToken.Type == token.IDENT {
		name := p.currentToken.Literal
		p.nextToken()

		// 関数呼び出しかチェック
		if p.currentToken.Type == token.LPAREN {
			p.nextToken() // '(' を消費

			arguments := []ast.ASTNode{}

			// 引数をパース
			for p.currentToken.Type != token.RPAREN && p.currentToken.Type != token.EOF {
				arg := p.ParseExpression()
				arguments = append(arguments, arg)
				// カンマは今回スキップ（簡単のため）
			}

			if p.currentToken.Type == token.RPAREN {
				p.nextToken() // ')' を消費
			}

			return &ast.CallNode{Function: name, Arguments: arguments}
		}

		return &ast.VariableNode{Name: name}
	}

	if p.currentToken.Type == token.LPAREN {
		p.nextToken() // '(' を消費
		expr := p.ParseExpression()
		if p.currentToken.Type == token.RPAREN {
			p.nextToken() // ')' を消費
		}
		return expr
	}

	// エラーケース: とりあえず 0 を返す
	return &ast.NumberNode{Value: 0}
}
