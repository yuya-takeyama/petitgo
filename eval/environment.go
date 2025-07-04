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
	variables map[string]Value
	functions map[string]*Function
	structs   map[string]*ast.StructDefinition
}

func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]Value),
		functions: make(map[string]*Function),
		structs:   make(map[string]*ast.StructDefinition),
	}
}

func (env *Environment) Set(name string, value Value) {
	env.variables[name] = value
}

func (env *Environment) Get(name string) (Value, bool) {
	value, exists := env.variables[name]
	return value, exists
}

// SetInt is a helper function for backward compatibility
func (env *Environment) SetInt(name string, value int) {
	env.variables[name] = &IntValue{Value: value}
}

// GetInt is a helper function for backward compatibility
func (env *Environment) GetInt(name string) (int, bool) {
	if value, exists := env.variables[name]; exists {
		if intVal, ok := value.(*IntValue); ok {
			return intVal.Value, true
		}
	}
	return 0, false
}

func (env *Environment) SetFunction(name string, function *Function) {
	env.functions[name] = function
}

func (env *Environment) GetFunction(name string) (*Function, bool) {
	function, exists := env.functions[name]
	return function, exists
}

func (env *Environment) SetStruct(name string, definition *ast.StructDefinition) {
	env.structs[name] = definition
}

func (env *Environment) GetStruct(name string) (*ast.StructDefinition, bool) {
	definition, exists := env.structs[name]
	return definition, exists
}
