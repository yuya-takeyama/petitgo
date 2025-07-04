package eval

import (
	"os"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

// ControlFlowException for break and continue
type ControlFlowException struct {
	Type string // "break" or "continue"
}

// printInt converts an integer to string and outputs it (without fmt package)
func printInt(n int) {
	if n == 0 {
		os.Stdout.Write([]byte("0"))
		return
	}

	if n < 0 {
		os.Stdout.Write([]byte("-"))
		n = -n
	}

	// get digits in reverse order
	digits := []byte{}
	for n > 0 {
		digit := n % 10
		digits = append(digits, byte('0'+digit))
		n /= 10
	}

	// output digits in correct order
	for i := len(digits) - 1; i >= 0; i-- {
		os.Stdout.Write([]byte{digits[i]})
	}
}

func Eval(node ast.ASTNode) int {
	env := NewEnvironment()
	return EvalWithEnvironment(node, env)
}

func EvalWithEnvironment(node ast.ASTNode, env *Environment) int {
	switch n := node.(type) {
	case *ast.NumberNode:
		return n.Value
	case *ast.VariableNode:
		if value, exists := env.Get(n.Name); exists {
			return value
		}
		return 0 // undefined variables return 0
	case *ast.BinaryOpNode:
		left := EvalWithEnvironment(n.Left, env)
		right := EvalWithEnvironment(n.Right, env)

		switch n.Operator {
		case token.ADD:
			return left + right
		case token.SUB:
			return left - right
		case token.MUL:
			return left * right
		case token.QUO:
			return left / right
		case token.EQL:
			if left == right {
				return 1
			}
			return 0
		case token.NEQ:
			if left != right {
				return 1
			}
			return 0
		case token.LSS:
			if left < right {
				return 1
			}
			return 0
		case token.GTR:
			if left > right {
				return 1
			}
			return 0
		case token.LEQ:
			if left <= right {
				return 1
			}
			return 0
		case token.GEQ:
			if left >= right {
				return 1
			}
			return 0
		}
	case *ast.CallNode:
		if n.Function == "print" && len(n.Arguments) > 0 {
			value := EvalWithEnvironment(n.Arguments[0], env)
			printInt(value)
			return value
		}
	}

	// エラーケース: とりあえず 0 を返す
	return 0
}

func EvalStatement(stmt ast.Statement, env *Environment) {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *ast.AssignStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *ast.ExpressionStatement:
		EvalWithEnvironment(s.Expression, env)
	case *ast.IfStatement:
		condition := EvalWithEnvironment(s.Condition, env)
		if condition != 0 { // 0 is false, anything else is true
			EvalBlockStatement(s.ThenBlock, env)
		} else if s.ElseBlock != nil {
			EvalBlockStatement(s.ElseBlock, env)
		}
	case *ast.ForStatement:
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
	case *ast.BlockStatement:
		EvalBlockStatement(s, env)
	case *ast.BreakStatement:
		// TODO: implement proper break handling
	case *ast.ContinueStatement:
		// TODO: implement proper continue handling
	}
}

func EvalBlockStatement(block *ast.BlockStatement, env *Environment) {
	for _, stmt := range block.Statements {
		EvalStatement(stmt, env)
	}
}
