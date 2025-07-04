package main

import (
	"bufio"
	"os"
)

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	print("petitgo calculator\n")
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

		// 式を評価
		result := evaluateExpression(input)
		printInt(result)
		print("\n")
	}
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

// fmt を使わずに整数を出力
func printInt(n int) {
	if n == 0 {
		print("0")
		return
	}

	if n < 0 {
		print("-")
		n = -n
	}

	// 数値を文字列に変換
	digits := ""
	for n > 0 {
		digit := n % 10
		digits = string(rune('0'+digit)) + digits
		n = n / 10
	}

	print(digits)
}
