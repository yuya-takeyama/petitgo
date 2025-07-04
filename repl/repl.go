package repl

import (
	"bufio"
	"os"

	"github.com/yuya-takeyama/petitgo/eval"
	"github.com/yuya-takeyama/petitgo/parser"
	"github.com/yuya-takeyama/petitgo/scanner"
)

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	env := eval.NewEnvironment()

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
		printValue(result)
		print("\n")
	}
}

func evaluateInput(input string, env *eval.Environment) eval.Value {
	sc := scanner.NewScanner(input)
	parser := parser.NewParser(sc)

	// Statement か Expression かを判定
	if isStatement(input) {
		stmt := parser.ParseStatement()
		eval.EvalStatement(stmt, env)
		// Statement の場合は結果を返さない（空文字列を返す）
		return &eval.StringValue{Value: ""}
	} else {
		expr := parser.ParseExpression()
		return eval.EvalValueWithEnvironment(expr, env)
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

// fmt を使わずに文字列を出力
func print(s string) {
	os.Stdout.WriteString(s)
}

// printValue outputs a Value using its String() method
func printValue(v eval.Value) {
	if v != nil && v.String() != "" {
		print(v.String())
	}
}

// 数値を文字列に変換して出力する（fmt なしで）
func printInt(n int) {
	if n == 0 {
		os.Stdout.Write([]byte("0"))
		return
	}

	if n < 0 {
		os.Stdout.Write([]byte("-"))
		n = -n
	}

	// 数字を逆順で取得
	digits := []byte{}
	for n > 0 {
		digit := n % 10
		digits = append(digits, byte('0'+digit))
		n /= 10
	}

	// 逆順で出力
	for i := len(digits) - 1; i >= 0; i-- {
		os.Stdout.Write([]byte{digits[i]})
	}
}
