package main

import (
	"bufio"
	"os"
)

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	env := NewEnvironment()

	print("petitgo calculator with variables and control flow\n")
	print("Try: var x int = 42, x := 10, x = 20, expressions like x + 5\n")
	print("Control flow: if x > 5 { print(x) }, comparisons like 5 > 3\n")
	print("Type 'exit' to quit\n")

	for {
		print("> ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		// 入力を評価
		result := evaluateInput(input, env)
		printInt(result)
		print("\n")
	}
}

func evaluateInput(input string, env *Environment) int {
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	// Statement か Expression かを判定
	if isStatement(input) {
		stmt := parser.ParseStatement()
		EvalStatement(stmt, env)
		// Statement の場合は結果を返さない（0 を返す）
		return 0
	} else {
		expr := parser.ParseExpression()
		return EvalWithEnvironment(expr, env)
	}
}

func isStatement(input string) bool {
	// 簡単な判定: "var", "if", "for", "break", "continue" で始まるか、":=" や " = " を含むか、{ を含むかどうか
	if len(input) >= 3 && input[:3] == "var" {
		return true
	}
	if len(input) >= 2 && input[:2] == "if" {
		return true
	}
	if len(input) >= 3 && input[:3] == "for" {
		return true
	}
	if len(input) >= 5 && input[:5] == "break" {
		return true
	}
	if len(input) >= 8 && input[:8] == "continue" {
		return true
	}

	// { を含む（ブロック文）
	for i := 0; i < len(input); i++ {
		if input[i] == '{' {
			return true
		}
	}

	// := を含む
	for i := 0; i < len(input)-1; i++ {
		if input[i] == ':' && input[i+1] == '=' {
			return true
		}
	}

	// = を含む（ただし == ではない）
	for i := 0; i < len(input); i++ {
		if input[i] == '=' {
			// == でないことを確認
			if i > 0 && input[i-1] == '=' {
				continue
			}
			if i < len(input)-1 && input[i+1] == '=' {
				continue
			}
			return true
		}
	}

	return false
}

func evaluateExpression(input string) int {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	ast := parser.ParseExpression()
	return Eval(ast)
}

// fmt を使わずに文字列を出力
func print(s string) {
	os.Stdout.WriteString(s)
}
