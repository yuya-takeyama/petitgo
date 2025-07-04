package main

import "os"

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

func Eval(node ASTNode) int {
	env := NewEnvironment()
	return EvalWithEnvironment(node, env)
}

func EvalWithEnvironment(node ASTNode, env *Environment) int {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value
	case *VariableNode:
		if value, exists := env.Get(n.Name); exists {
			return value
		}
		return 0 // 未定義変数は 0
	case *BinaryOpNode:
		left := EvalWithEnvironment(n.Left, env)
		right := EvalWithEnvironment(n.Right, env)

		switch n.Operator {
		case ADD:
			return left + right
		case SUB:
			return left - right
		case MUL:
			return left * right
		case QUO:
			return left / right
		case EQL:
			if left == right {
				return 1
			}
			return 0
		case NEQ:
			if left != right {
				return 1
			}
			return 0
		case LSS:
			if left < right {
				return 1
			}
			return 0
		case GTR:
			if left > right {
				return 1
			}
			return 0
		case LEQ:
			if left <= right {
				return 1
			}
			return 0
		case GEQ:
			if left >= right {
				return 1
			}
			return 0
		}
	case *CallNode:
		if n.Function == "print" && len(n.Arguments) > 0 {
			value := EvalWithEnvironment(n.Arguments[0], env)
			printInt(value)
			return value
		}
	}

	// エラーケース: とりあえず 0 を返す
	return 0
}

func EvalStatement(stmt Statement, env *Environment) {
	switch s := stmt.(type) {
	case *VarStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *AssignStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *ExpressionStatement:
		EvalWithEnvironment(s.Expression, env)
	case *IfStatement:
		condition := EvalWithEnvironment(s.Condition, env)
		if condition != 0 { // 0 は false、それ以外は true
			EvalBlockStatement(s.ThenBlock, env)
		} else if s.ElseBlock != nil {
			EvalBlockStatement(s.ElseBlock, env)
		}
	case *ForStatement:
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
	case *BlockStatement:
		EvalBlockStatement(s, env)
	case *BreakStatement:
		// TODO: implement proper break handling
	case *ContinueStatement:
		// TODO: implement proper continue handling
	}
}

func EvalBlockStatement(block *BlockStatement, env *Environment) {
	for _, stmt := range block.Statements {
		EvalStatement(stmt, env)
	}
}
