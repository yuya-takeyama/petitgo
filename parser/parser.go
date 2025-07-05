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
	for {
		p.currentToken = p.scanner.NextToken()
		// Skip comments
		if p.currentToken.Type != token.COMMENT {
			break
		}
	}
}

func (p *Parser) ParseStatement() ast.Statement {
	// Check for EOF first
	if p.currentToken.Type == token.EOF {
		return nil
	}

	switch p.currentToken.Type {
	case token.IDENT:
		if p.currentToken.Literal == "var" {
			return p.parseVarStatement()
		}
		// Check what follows the identifier
		nextScanner := scanner.NewScanner(p.scanner.Input()[p.scanner.Position():])
		nextToken := nextScanner.NextToken()

		if nextToken.Type == token.ASSIGN {
			if nextToken.Literal == ":=" {
				return p.parseAssignStatement()
			} else {
				return p.parseReassignStatement()
			}
		} else if nextToken.Type == token.INC {
			return p.parseIncStatement()
		} else if nextToken.Type == token.DEC {
			return p.parseDecStatement()
		} else if nextToken.Type == token.ADD_ASSIGN || nextToken.Type == token.SUB_ASSIGN ||
			nextToken.Type == token.MUL_ASSIGN || nextToken.Type == token.QUO_ASSIGN {
			return p.parseCompoundAssignStatement()
		}
		return p.parseExpressionStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.SWITCH:
		return p.parseSwitchStatement()
	case token.TYPE:
		return p.parseTypeStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.FUNC:
		return p.parseFuncStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.PACKAGE:
		return p.parsePackageStatement()
	case token.IMPORT:
		return p.parseImportStatement()
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

func (p *Parser) parseReassignStatement() ast.Statement {
	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// =
	p.nextToken()

	// 式
	value := p.ParseExpression()

	return &ast.ReassignStatement{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseIncStatement() ast.Statement {
	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// ++
	p.nextToken()

	return &ast.IncStatement{
		Name: name,
	}
}

func (p *Parser) parseDecStatement() ast.Statement {
	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// --
	p.nextToken()

	return &ast.DecStatement{
		Name: name,
	}
}

func (p *Parser) parseCompoundAssignStatement() ast.Statement {
	// 変数名
	name := p.currentToken.Literal
	p.nextToken()

	// compound operator (+=, -=, etc.)
	operator := p.currentToken.Type
	p.nextToken()

	// 式
	value := p.ParseExpression()

	return &ast.CompoundAssignStatement{
		Name:     name,
		Operator: operator,
		Value:    value,
	}
}

func (p *Parser) parseSwitchStatement() ast.Statement {
	// switch
	p.nextToken()

	// value to switch on
	value := p.ParseExpression()

	// {
	if p.currentToken.Type != token.LBRACE {
		return &ast.SwitchStatement{}
	}
	p.nextToken()

	var cases []*ast.CaseStatement
	var defaultCase *ast.BlockStatement

	// Parse cases
	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		if p.currentToken.Type == token.CASE {
			// case value:
			p.nextToken()
			caseValue := p.ParseExpression()

			// :
			if p.currentToken.Type == token.COLON {
				p.nextToken()
			}

			// Parse statements until next case/default/}
			var statements []ast.Statement
			for p.currentToken.Type != token.CASE &&
				p.currentToken.Type != token.DEFAULT &&
				p.currentToken.Type != token.RBRACE &&
				p.currentToken.Type != token.EOF {
				if stmt := p.ParseStatement(); stmt != nil {
					statements = append(statements, stmt)
				} else {
					// If ParseStatement returns nil, advance token to prevent infinite loop
					p.nextToken()
				}
			}

			cases = append(cases, &ast.CaseStatement{
				Value: caseValue,
				Body:  &ast.BlockStatement{Statements: statements},
			})
		} else if p.currentToken.Type == token.DEFAULT {
			// default:
			p.nextToken()

			// :
			if p.currentToken.Type == token.COLON {
				p.nextToken()
			}

			// Parse statements until }
			var statements []ast.Statement
			for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
				if stmt := p.ParseStatement(); stmt != nil {
					statements = append(statements, stmt)
				} else {
					// If ParseStatement returns nil, advance token to prevent infinite loop
					p.nextToken()
				}
			}

			defaultCase = &ast.BlockStatement{Statements: statements}
		} else {
			p.nextToken()
		}
	}

	// }
	if p.currentToken.Type == token.RBRACE {
		p.nextToken()
	}

	return &ast.SwitchStatement{
		Value:   value,
		Cases:   cases,
		Default: defaultCase,
	}
}

func (p *Parser) parseTypeStatement() ast.Statement {
	// type
	p.nextToken()

	// type name
	if p.currentToken.Type != token.IDENT {
		return &ast.TypeStatement{}
	}
	typeName := p.currentToken.Literal
	p.nextToken()

	// struct
	if p.currentToken.Literal != "struct" {
		return &ast.TypeStatement{}
	}
	p.nextToken()

	// {
	if p.currentToken.Type != token.LBRACE {
		return &ast.TypeStatement{}
	}
	p.nextToken()

	var fields []*ast.FieldDef

	// Parse fields
	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		if p.currentToken.Type == token.IDENT {
			fieldName := p.currentToken.Literal
			p.nextToken()

			if p.currentToken.Type == token.IDENT {
				fieldType := p.currentToken.Literal
				p.nextToken()

				fields = append(fields, &ast.FieldDef{
					Name: fieldName,
					Type: fieldType,
				})
			}
		} else {
			p.nextToken()
		}
	}

	// }
	if p.currentToken.Type == token.RBRACE {
		p.nextToken()
	}

	return &ast.TypeStatement{
		Name:   typeName,
		Fields: fields,
	}
}

func (p *Parser) parseIfStatement() ast.Statement {
	// if
	p.nextToken()

	// condition - special handling for if statement context
	condition := p.parseIfCondition()

	// {
	if p.currentToken.Type != token.LBRACE {
		// エラー: とりあえずダミーを返す（conditionは設定する）
		return &ast.IfStatement{Condition: condition}
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

// parseIfCondition parses conditions in if statements (prevents struct literal confusion)
func (p *Parser) parseIfCondition() ast.ASTNode {
	// For if conditions, we need to be careful about struct literal vs variable references
	// If we see IDENT followed by LBRACE, treat it as a variable, not a struct literal
	if p.currentToken.Type == token.IDENT {
		name := p.currentToken.Literal
		p.nextToken()

		// If the next token is LBRACE, this is definitely a variable reference
		// because the LBRACE belongs to the if statement block
		if p.currentToken.Type == token.LBRACE {
			return &ast.VariableNode{Name: name}
		}

		// If not LBRACE, parse as normal expression (might be comparison, etc.)
		// We need to "put back" the IDENT token and parse normally
		// For now, handle simple cases
		if p.currentToken.Type == token.GTR || p.currentToken.Type == token.LSS ||
			p.currentToken.Type == token.EQL || p.currentToken.Type == token.NEQ ||
			p.currentToken.Type == token.LEQ || p.currentToken.Type == token.GEQ {
			// Binary comparison: IDENT op EXPR
			operator := p.currentToken.Type
			p.nextToken()
			right := p.parseTerm()
			return &ast.BinaryOpNode{
				Left:     &ast.VariableNode{Name: name},
				Operator: operator,
				Right:    right,
			}
		}

		// Just a variable reference
		return &ast.VariableNode{Name: name}
	}

	// For other cases, use normal expression parsing
	// This handles cases like "1 > 0", "x + y > 5", etc.
	return p.ParseExpression()
}

func (p *Parser) parseForStatement() ast.Statement {
	// for
	p.nextToken()

	// for文の形式を判定
	// 1. for condition { ... } (condition-only)
	// 2. for init; condition; update { ... } (full form)

	// Check if this is full form (has semicolons)
	// We need to look ahead to see if there are semicolons
	if p.hasForLoopSemicolons() {
		return p.parseFullForStatement()
	} else {
		return p.parseConditionOnlyForStatement()
	}
}

func (p *Parser) hasForLoopSemicolons() bool {
	// Simple check: scan ahead for semicolons before opening brace
	scanner := p.scanner
	pos := scanner.Position()
	input := scanner.Input()

	for pos < len(input) && input[pos] != '{' && input[pos] != '\n' {
		if input[pos] == ';' {
			return true
		}
		pos++
	}
	return false
}

func (p *Parser) parseConditionOnlyForStatement() ast.Statement {
	// for condition { ... }
	condition := p.ParseExpression()

	// {
	if p.currentToken.Type != token.LBRACE {
		return &ast.ForStatement{}
	}

	body := p.parseBlockStatement()

	return &ast.ForStatement{
		Init:      nil,
		Condition: condition,
		Update:    nil,
		Body:      body,
	}
}

func (p *Parser) parseFullForStatement() ast.Statement {
	// for init; condition; update { ... }

	// init statement
	var init ast.Statement
	if p.currentToken.Type != token.SEMICOLON {
		init = p.ParseStatement()
	}

	// skip semicolon
	if p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	// condition
	var condition ast.ASTNode
	if p.currentToken.Type != token.SEMICOLON {
		condition = p.ParseExpression()
	}

	// skip semicolon
	if p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	// update statement
	var update ast.Statement
	if p.currentToken.Type != token.LBRACE {
		update = p.ParseStatement()
	}

	// {
	if p.currentToken.Type != token.LBRACE {
		return &ast.ForStatement{}
	}

	body := p.parseBlockStatement()

	return &ast.ForStatement{
		Init:      init,
		Condition: condition,
		Update:    update,
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

	if p.currentToken.Type == token.TRUE {
		p.nextToken()
		return &ast.BooleanNode{Value: true}
	}

	if p.currentToken.Type == token.FALSE {
		p.nextToken()
		return &ast.BooleanNode{Value: false}
	}

	if p.currentToken.Type == token.STRING {
		value := p.currentToken.Literal
		p.nextToken()
		return &ast.StringNode{Value: value}
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

				// カンマをスキップ
				if p.currentToken.Type == token.COMMA {
					p.nextToken()
				}
			}

			if p.currentToken.Type == token.RPAREN {
				p.nextToken() // ')' を消費
			}

			return &ast.CallNode{Function: name, Arguments: arguments}
		}

		// struct literal かチェック (Person{...})
		if p.currentToken.Type == token.LBRACE {
			return p.parseStructLiteral(name)
		}

		// 変数または struct instance として処理
		var result ast.ASTNode = &ast.VariableNode{Name: name}

		// field access チェック (.field) or index access ([index])
		for p.currentToken.Type == token.PERIOD || p.currentToken.Type == token.LBRACK {
			if p.currentToken.Type == token.PERIOD {
				p.nextToken() // '.' を消費
				if p.currentToken.Type == token.IDENT {
					fieldName := p.currentToken.Literal
					p.nextToken()
					result = &ast.FieldAccessNode{
						Object: result,
						Field:  fieldName,
					}
				}
			} else if p.currentToken.Type == token.LBRACK {
				p.nextToken() // '[' を消費
				index := p.ParseExpression()
				if p.currentToken.Type == token.RBRACK {
					p.nextToken() // ']' を消費
				}
				result = &ast.IndexAccess{
					Object: result,
					Index:  index,
				}
			}
		}

		return result
	}

	if p.currentToken.Type == token.LPAREN {
		p.nextToken() // '(' を消費
		expr := p.ParseExpression()
		if p.currentToken.Type == token.RPAREN {
			p.nextToken() // ')' を消費
		}
		return expr
	}

	// Slice literal: []type{elements...}
	if p.currentToken.Type == token.LBRACK {
		p.nextToken() // '[' を消費
		if p.currentToken.Type == token.RBRACK {
			p.nextToken() // ']' を消費
			if p.currentToken.Type == token.IDENT {
				elementType := p.currentToken.Literal
				p.nextToken()
				if p.currentToken.Type == token.LBRACE {
					return p.parseSliceLiteral(elementType)
				}
			}
		}
	}

	// エラーケース: とりあえず 0 を返す
	return &ast.NumberNode{Value: 0}
}

// parseFuncStatement parses function definitions: func name(param type, ...) returnType { body }
func (p *Parser) parseFuncStatement() ast.Statement {
	// consume 'func'
	p.nextToken()

	// function name
	name := p.currentToken.Literal
	p.nextToken()

	// consume '('
	if p.currentToken.Type != token.LPAREN {
		// Error: expected '('
		return nil
	}
	p.nextToken()

	// parse parameters
	parameters := []ast.Parameter{}
	for p.currentToken.Type != token.RPAREN && p.currentToken.Type != token.EOF {
		// parameter name
		paramName := p.currentToken.Literal
		p.nextToken()

		// parameter type
		paramType := p.currentToken.Literal
		p.nextToken()

		parameters = append(parameters, ast.Parameter{Name: paramName, Type: paramType})

		// skip comma if present
		if p.currentToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	// consume ')'
	if p.currentToken.Type == token.RPAREN {
		p.nextToken()
	}

	// return type (optional for now, default to empty)
	returnType := ""
	if p.currentToken.Type == token.IDENT {
		returnType = p.currentToken.Literal
		p.nextToken()
	}

	// function body
	var body *ast.BlockStatement
	if p.currentToken.Type == token.LBRACE {
		body = p.parseBlockStatement()
	}

	return &ast.FuncStatement{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
}

// parseReturnStatement parses return statements: return [expression]
func (p *Parser) parseReturnStatement() ast.Statement {
	// consume 'return'
	p.nextToken()

	// parse optional return value
	var value ast.ASTNode
	if p.currentToken.Type != token.EOF && p.currentToken.Type != token.RBRACE {
		value = p.ParseExpression()
	}

	return &ast.ReturnStatement{Value: value}
}

// parseStructDefinition parses struct definitions: type Person struct { Name string; Age int }

// parseStructLiteral parses struct literals: Person{Name: "Alice", Age: 25}
func (p *Parser) parseStructLiteral(typeName string) ast.ASTNode {
	// '{' は既に確認済み
	p.nextToken() // '{' を消費

	fields := make(map[string]ast.ASTNode)

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		// field name
		if p.currentToken.Type != token.IDENT {
			break
		}
		fieldName := p.currentToken.Literal
		p.nextToken()

		// expect ':'
		if p.currentToken.Type != token.COLON {
			break // error: expected ':'
		}
		p.nextToken() // ':' を消費

		// field value
		value := p.ParseExpression()
		fields[fieldName] = value

		// skip comma if present
		if p.currentToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	// consume '}'
	if p.currentToken.Type == token.RBRACE {
		p.nextToken()
	}

	return &ast.StructLiteral{
		TypeName: typeName,
		Fields:   fields,
	}
}

// parseSliceLiteral parses slice literals: []int{1, 2, 3}
func (p *Parser) parseSliceLiteral(elementType string) ast.ASTNode {
	// '{' は既に確認済み
	p.nextToken() // '{' を消費

	var elements []ast.ASTNode

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		element := p.ParseExpression()
		elements = append(elements, element)

		// skip comma if present
		if p.currentToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	// consume '}'
	if p.currentToken.Type == token.RBRACE {
		p.nextToken()
	}

	return &ast.SliceLiteral{
		ElementType: elementType,
		Elements:    elements,
	}
}

// parsePackageStatement parses package declarations: package main
func (p *Parser) parsePackageStatement() ast.Statement {
	// consume 'package'
	p.nextToken()

	// package name
	if p.currentToken.Type != token.IDENT {
		// error: expected package name
		return nil
	}
	name := p.currentToken.Literal
	p.nextToken()

	return &ast.PackageStatement{
		Name: name,
	}
}

// parseImportStatement parses import declarations: import "fmt"
func (p *Parser) parseImportStatement() ast.Statement {
	// consume 'import'
	p.nextToken()

	// import path (string literal)
	if p.currentToken.Type != token.STRING {
		// error: expected string literal
		return nil
	}
	path := p.currentToken.Literal
	p.nextToken()

	return &ast.ImportStatement{
		Path: path,
	}
}
