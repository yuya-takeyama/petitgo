package eval

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
