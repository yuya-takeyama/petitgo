package main

func Eval(node ASTNode) int {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value
	case *BinaryOpNode:
		left := Eval(n.Left)
		right := Eval(n.Right)

		switch n.Operator {
		case ADD:
			return left + right
		case SUB:
			return left - right
		case MUL:
			return left * right
		case QUO:
			return left / right
		}
	}

	// エラーケース: とりあえず 0 を返す
	return 0
}
