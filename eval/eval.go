package eval

import (
	"os"

	"github.com/yuya-takeyama/petitgo/lexer"
	"github.com/yuya-takeyama/petitgo/parser"
)

// ControlFlowException for break and continue
type ControlFlowException struct {
	Type string // "break" or "continue"
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

func Eval(node parser.ASTNode) int {
	env := NewEnvironment()
	return EvalWithEnvironment(node, env)
}

func EvalWithEnvironment(node parser.ASTNode, env *Environment) int {
	switch n := node.(type) {
	case *parser.NumberNode:
		return n.Value
	case *parser.VariableNode:
		if value, exists := env.Get(n.Name); exists {
			return value
		}
		return 0 // 未定義変数は 0
	case *parser.BinaryOpNode:
		left := EvalWithEnvironment(n.Left, env)
		right := EvalWithEnvironment(n.Right, env)

		switch n.Operator {
		case lexer.ADD:
			return left + right
		case lexer.SUB:
			return left - right
		case lexer.MUL:
			return left * right
		case lexer.QUO:
			return left / right
		case lexer.EQL:
			if left == right {
				return 1
			}
			return 0
		case lexer.NEQ:
			if left != right {
				return 1
			}
			return 0
		case lexer.LSS:
			if left < right {
				return 1
			}
			return 0
		case lexer.GTR:
			if left > right {
				return 1
			}
			return 0
		case lexer.LEQ:
			if left <= right {
				return 1
			}
			return 0
		case lexer.GEQ:
			if left >= right {
				return 1
			}
			return 0
		}
	case *parser.CallNode:
		if n.Function == "print" && len(n.Arguments) > 0 {
			value := EvalWithEnvironment(n.Arguments[0], env)
			printInt(value)
			return value
		}
	}

	// エラーケース: とりあえず 0 を返す
	return 0
}

func EvalStatement(stmt parser.Statement, env *Environment) {
	switch s := stmt.(type) {
	case *parser.VarStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *parser.AssignStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *parser.ExpressionStatement:
		EvalWithEnvironment(s.Expression, env)
	case *parser.IfStatement:
		condition := EvalWithEnvironment(s.Condition, env)
		if condition != 0 { // 0 は false、それ以外は true
			EvalBlockStatement(s.ThenBlock, env)
		} else if s.ElseBlock != nil {
			EvalBlockStatement(s.ElseBlock, env)
		}
	case *parser.ForStatement:
		for {
			// condition check
			if s.Condition != nil {
				condition := EvalWithEnvironment(s.Condition, env)
				if condition == 0 { // false
					break
				}
			}

			// body execution
			EvalBlockStatement(s.Body, env)

			// update (not implemented for condition-only form)
		}
	case *parser.BlockStatement:
		EvalBlockStatement(s, env)
	case *parser.BreakStatement:
		// TODO: implement proper break handling
	case *parser.ContinueStatement:
		// TODO: implement proper continue handling
	}
}

func EvalBlockStatement(block *parser.BlockStatement, env *Environment) {
	for _, stmt := range block.Statements {
		EvalStatement(stmt, env)
	}
}
