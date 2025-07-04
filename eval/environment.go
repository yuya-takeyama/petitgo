package eval

import "github.com/yuya-takeyama/petitgo/ast"

// Function represents a user-defined function
type Function struct {
	Name       string
	Parameters []ast.Parameter
	ReturnType string
	Body       *ast.BlockStatement
}

// Environment は変数と関数を管理する
type Environment struct {
	variables map[string]int
	functions map[string]*Function
}

func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]int),
		functions: make(map[string]*Function),
	}
}

func (env *Environment) Set(name string, value int) {
	env.variables[name] = value
}

func (env *Environment) Get(name string) (int, bool) {
	value, exists := env.variables[name]
	return value, exists
}

func (env *Environment) SetFunction(name string, function *Function) {
	env.functions[name] = function
}

func (env *Environment) GetFunction(name string) (*Function, bool) {
	function, exists := env.functions[name]
	return function, exists
}
