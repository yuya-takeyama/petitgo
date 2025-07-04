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

// 関数呼び出しノード (function(args...))
type CallNode struct {
	Function  string
	Arguments []ASTNode
}

func (n *CallNode) String() string {
	return "CallNode"
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
		// 次のトークンが ASSIGN なら代入文、そうでなければ式文
		nextLexer := NewLexer(p.lexer.input[p.lexer.position:])
		nextToken := nextLexer.NextToken()
		if nextToken.Type == ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	case IF:
		return p.parseIfStatement()
	case FOR:
		return p.parseForStatement()
	case BREAK:
		return p.parseBreakStatement()
	case CONTINUE:
		return p.parseContinueStatement()
	case LBRACE:
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
	if p.currentToken.Type != LBRACE {
		// エラー: とりあえずダミーを返す
		return &IfStatement{}
	}

	// then block
	thenBlock := p.parseBlockStatement()

	var elseBlock *BlockStatement = nil

	// else があるかチェック
	if p.currentToken.Type == ELSE {
		p.nextToken()
		if p.currentToken.Type == LBRACE {
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
	if p.currentToken.Type != LBRACE {
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

	for p.currentToken.Type != RBRACE && p.currentToken.Type != EOF {
		stmt := p.ParseStatement()
		statements = append(statements, stmt)
	}

	// }
	if p.currentToken.Type == RBRACE {
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

	for p.currentToken.Type == EQL || p.currentToken.Type == NEQ ||
		p.currentToken.Type == LSS || p.currentToken.Type == GTR ||
		p.currentToken.Type == LEQ || p.currentToken.Type == GEQ {
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

		// 関数呼び出しかチェック
		if p.currentToken.Type == LPAREN {
			p.nextToken() // '(' を消費

			arguments := []ASTNode{}

			// 引数をパース
			for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
				arg := p.ParseExpression()
				arguments = append(arguments, arg)
				// カンマは今回スキップ（簡単のため）
			}

			if p.currentToken.Type == RPAREN {
				p.nextToken() // ')' を消費
			}

			return &CallNode{Function: name, Arguments: arguments}
		}

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
