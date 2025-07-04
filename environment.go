package main

// Environment は変数を管理する
type Environment struct {
	variables map[string]int
}

func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]int),
	}
}

func (env *Environment) Set(name string, value int) {
	env.variables[name] = value
}

func (env *Environment) Get(name string) (int, bool) {
	value, exists := env.variables[name]
	return value, exists
}

// Statement を評価する
func EvalStatement(stmt Statement, env *Environment) {
	switch s := stmt.(type) {
	case *VarStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	case *AssignStatement:
		value := EvalWithEnvironment(s.Value, env)
		env.Set(s.Name, value)
	}
}

// Environment を使って式を評価する
func EvalWithEnvironment(node ASTNode, env *Environment) int {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value
	case *VariableNode:
		value, exists := env.Get(n.Name)
		if !exists {
			// とりあえず 0 を返す（エラーハンドリングは後で）
			return 0
		}
		return value
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
		}
	}

	// エラーケース: とりあえず 0 を返す
	return 0
}
