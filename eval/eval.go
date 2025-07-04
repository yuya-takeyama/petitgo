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

// ReturnException for return statements
type ReturnException struct {
	Value int
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

// EvalValue evaluates an expression and returns a proper Value
func EvalValue(node ast.ASTNode) Value {
	env := NewEnvironment()
	return EvalValueWithEnvironment(node, env)
}

// EvalValueWithEnvironment evaluates an expression with proper type system
func EvalValueWithEnvironment(node ast.ASTNode, env *Environment) Value {
	switch n := node.(type) {
	case *ast.NumberNode:
		return &IntValue{Value: n.Value}
	case *ast.BooleanNode:
		return &BoolValue{Value: n.Value}
	case *ast.StringNode:
		return &StringValue{Value: n.Value}
	case *ast.VariableNode:
		if value, exists := env.Get(n.Name); exists {
			return value
		}
		return &IntValue{Value: 0}
	case *ast.BinaryOpNode:
		return evalBinaryOpWithTypes(n, env)
	case *ast.CallNode:
		return evalCallWithTypes(n, env)
	case *ast.StructLiteral:
		return evalStructLiteral(n, env)
	case *ast.FieldAccessNode:
		return evalFieldAccess(n, env)
	case *ast.SliceLiteral:
		return evalSliceLiteral(n, env)
	case *ast.IndexAccess:
		return evalIndexAccess(n, env)
	}

	// Default case: convert old evaluation to IntValue
	oldResult := EvalWithEnvironment(node, env)
	return &IntValue{Value: oldResult}
}

// evalBinaryOpWithTypes evaluates binary operations with proper type checking
func evalBinaryOpWithTypes(node *ast.BinaryOpNode, env *Environment) Value {
	left := EvalValueWithEnvironment(node.Left, env)
	right := EvalValueWithEnvironment(node.Right, env)

	// Type checking and operation dispatch
	switch node.Operator {
	case token.ADD:
		return addValues(left, right)
	case token.SUB, token.MUL, token.QUO:
		return arithmeticOp(left, right, node.Operator)
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		return compareValues(left, right, node.Operator)
	}

	// Default: return false for unsupported operations
	return &BoolValue{Value: false}
}

// addValues handles addition with type checking (int + int, string + string)
func addValues(left, right Value) Value {
	leftInt, leftIsInt := left.(*IntValue)
	rightInt, rightIsInt := right.(*IntValue)

	// int + int
	if leftIsInt && rightIsInt {
		return &IntValue{Value: leftInt.Value + rightInt.Value}
	}

	leftStr, leftIsStr := left.(*StringValue)
	rightStr, rightIsStr := right.(*StringValue)

	// string + string
	if leftIsStr && rightIsStr {
		return &StringValue{Value: leftStr.Value + rightStr.Value}
	}

	// Type mismatch: convert to string and concatenate
	return &StringValue{Value: left.String() + right.String()}
}

// arithmeticOp handles arithmetic operations (only for int types)
func arithmeticOp(left, right Value, op token.Token) Value {
	leftInt, leftIsInt := left.(*IntValue)
	rightInt, rightIsInt := right.(*IntValue)

	if !leftIsInt || !rightIsInt {
		// Type error: return 0 for now
		return &IntValue{Value: 0}
	}

	switch op {
	case token.SUB:
		return &IntValue{Value: leftInt.Value - rightInt.Value}
	case token.MUL:
		return &IntValue{Value: leftInt.Value * rightInt.Value}
	case token.QUO:
		if rightInt.Value == 0 {
			// Division by zero: return 0 for now
			return &IntValue{Value: 0}
		}
		return &IntValue{Value: leftInt.Value / rightInt.Value}
	}

	return &IntValue{Value: 0}
}

// compareValues handles comparison operations
func compareValues(left, right Value, op token.Token) Value {
	// For now, only support int comparisons
	leftInt, leftIsInt := left.(*IntValue)
	rightInt, rightIsInt := right.(*IntValue)

	if !leftIsInt || !rightIsInt {
		// Type error: compare as strings for now
		leftStr := left.String()
		rightStr := right.String()
		return compareStrings(leftStr, rightStr, op)
	}

	return compareInts(leftInt.Value, rightInt.Value, op)
}

// compareInts compares two integers
func compareInts(left, right int, op token.Token) Value {
	switch op {
	case token.EQL:
		return &BoolValue{Value: left == right}
	case token.NEQ:
		return &BoolValue{Value: left != right}
	case token.LSS:
		return &BoolValue{Value: left < right}
	case token.GTR:
		return &BoolValue{Value: left > right}
	case token.LEQ:
		return &BoolValue{Value: left <= right}
	case token.GEQ:
		return &BoolValue{Value: left >= right}
	}
	return &BoolValue{Value: false}
}

// compareStrings compares two strings
func compareStrings(left, right string, op token.Token) Value {
	switch op {
	case token.EQL:
		return &BoolValue{Value: left == right}
	case token.NEQ:
		return &BoolValue{Value: left != right}
	case token.LSS:
		return &BoolValue{Value: left < right}
	case token.GTR:
		return &BoolValue{Value: left > right}
	case token.LEQ:
		return &BoolValue{Value: left <= right}
	case token.GEQ:
		return &BoolValue{Value: left >= right}
	}
	return &BoolValue{Value: false}
}

// evalCallWithTypes evaluates function calls with proper type system
func evalCallWithTypes(node *ast.CallNode, env *Environment) Value {
	// Built-in function: print
	if node.Function == "print" && len(node.Arguments) > 0 {
		value := EvalValueWithEnvironment(node.Arguments[0], env)
		// Print the value using its String() method
		os.Stdout.Write([]byte(value.String()))
		return value
	}

	// Built-in function: len
	if node.Function == "len" && len(node.Arguments) == 1 {
		value := EvalValueWithEnvironment(node.Arguments[0], env)
		switch v := value.(type) {
		case *SliceValue:
			return &IntValue{Value: len(v.Elements)}
		case *StringValue:
			return &IntValue{Value: len(v.Value)}
		default:
			return &IntValue{Value: 0}
		}
	}

	// Built-in function: append
	if node.Function == "append" && len(node.Arguments) >= 2 {
		sliceVal := EvalValueWithEnvironment(node.Arguments[0], env)
		if slice, ok := sliceVal.(*SliceValue); ok {
			// Create new slice with appended elements
			newElements := make([]Value, len(slice.Elements))
			copy(newElements, slice.Elements)

			// Append all remaining arguments
			for i := 1; i < len(node.Arguments); i++ {
				elem := EvalValueWithEnvironment(node.Arguments[i], env)
				newElements = append(newElements, elem)
			}

			return &SliceValue{
				ElementType: slice.ElementType,
				Elements:    newElements,
			}
		}
		return sliceVal // Return original if not a slice
	}

	// User-defined function: for now, fallback to old implementation
	if function, exists := env.GetFunction(node.Function); exists {
		result := callUserFunction(function, node.Arguments, env)
		return &IntValue{Value: result}
	}

	return &IntValue{Value: 0}
}

func EvalWithEnvironment(node ast.ASTNode, env *Environment) int {
	switch n := node.(type) {
	case *ast.NumberNode:
		return n.Value
	case *ast.BooleanNode:
		if n.Value {
			return 1 // true = 1
		}
		return 0 // false = 0
	case *ast.StringNode:
		// For now, return string length as evaluation result
		// In the future, we'll need proper type system
		return len(n.Value)
	case *ast.VariableNode:
		if value, exists := env.GetInt(n.Name); exists {
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
		// Built-in function: print
		if n.Function == "print" && len(n.Arguments) > 0 {
			value := EvalWithEnvironment(n.Arguments[0], env)
			printInt(value)
			return value
		}

		// User-defined function
		if function, exists := env.GetFunction(n.Function); exists {
			return callUserFunction(function, n.Arguments, env)
		}
	}

	// エラーケース: とりあえず 0 を返す
	return 0
}

func EvalStatement(stmt ast.Statement, env *Environment) {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		// Use type-aware evaluation
		value := EvalValueWithEnvironment(s.Value, env)

		// Type checking: verify that the value matches the declared type
		expectedType := s.TypeName
		actualType := value.Type()

		if expectedType != actualType {
			// Type mismatch - for now, we'll create a zero value of the expected type
			// In a more sophisticated implementation, this would be a compile-time error
			switch expectedType {
			case "int":
				value = &IntValue{Value: 0}
			case "string":
				value = &StringValue{Value: ""}
			case "bool":
				value = &BoolValue{Value: false}
			default:
				// Unknown type, keep original value
			}
		}

		env.Set(s.Name, value)
	case *ast.AssignStatement:
		// Use type-aware evaluation with type inference
		value := EvalValueWithEnvironment(s.Value, env)

		// Check if variable already exists
		if existingValue, exists := env.Get(s.Name); exists {
			// Variable exists - check type compatibility
			existingType := existingValue.Type()
			newType := value.Type()

			if existingType != newType {
				// Type mismatch - use zero value of existing type for type safety
				// This maintains Go's type safety principles
				switch existingType {
				case "int":
					value = &IntValue{Value: 0}
				case "string":
					value = &StringValue{Value: ""}
				case "bool":
					value = &BoolValue{Value: false}
				default:
					// Unknown type, keep new value (fallback)
				}
			}
		}
		// If variable doesn't exist, infer type from value (type inference)

		env.Set(s.Name, value)
	case *ast.ReassignStatement:
		// Variable reassignment - must exist already
		value := EvalValueWithEnvironment(s.Value, env)

		// Check if variable exists
		if existingValue, exists := env.Get(s.Name); exists {
			// Variable exists - check type compatibility
			existingType := existingValue.Type()
			newType := value.Type()

			if existingType != newType {
				// Type mismatch - use zero value of existing type for type safety
				switch existingType {
				case "int":
					value = &IntValue{Value: 0}
				case "string":
					value = &StringValue{Value: ""}
				case "bool":
					value = &BoolValue{Value: false}
				default:
					// Unknown type, keep new value (fallback)
				}
			}

			env.Set(s.Name, value)
		} else {
			// Variable doesn't exist - this would be a compile error in real Go
			// For now, we'll just ignore it
		}
	case *ast.IncStatement:
		// Increment statement - variable must exist
		if existingValue, exists := env.Get(s.Name); exists {
			if intVal, ok := existingValue.(*IntValue); ok {
				env.Set(s.Name, &IntValue{Value: intVal.Value + 1})
			}
		}
	case *ast.DecStatement:
		// Decrement statement - variable must exist
		if existingValue, exists := env.Get(s.Name); exists {
			if intVal, ok := existingValue.(*IntValue); ok {
				env.Set(s.Name, &IntValue{Value: intVal.Value - 1})
			}
		}
	case *ast.CompoundAssignStatement:
		// Compound assignment - variable must exist
		if existingValue, exists := env.Get(s.Name); exists {
			if intVal, ok := existingValue.(*IntValue); ok {
				newValue := EvalValueWithEnvironment(s.Value, env)
				if newIntVal, ok := newValue.(*IntValue); ok {
					var result int
					switch s.Operator {
					case token.ADD_ASSIGN:
						result = intVal.Value + newIntVal.Value
					case token.SUB_ASSIGN:
						result = intVal.Value - newIntVal.Value
					case token.MUL_ASSIGN:
						result = intVal.Value * newIntVal.Value
					case token.QUO_ASSIGN:
						if newIntVal.Value != 0 {
							result = intVal.Value / newIntVal.Value
						} else {
							result = 0 // Handle division by zero
						}
					}
					env.Set(s.Name, &IntValue{Value: result})
				}
			}
		}
	case *ast.ExpressionStatement:
		// Use type-aware evaluation for expressions
		EvalValueWithEnvironment(s.Expression, env)
	case *ast.IfStatement:
		// Use type-aware evaluation for conditions
		condition := EvalValueWithEnvironment(s.Condition, env)
		if condition.IsTruthy() {
			EvalBlockStatement(s.ThenBlock, env)
		} else if s.ElseBlock != nil {
			EvalBlockStatement(s.ElseBlock, env)
		}
	case *ast.ForStatement:
		// Execute init statement if present
		if s.Init != nil {
			EvalStatement(s.Init, env)
		}

		for {
			// condition check with type-aware evaluation
			if s.Condition != nil {
				condition := EvalValueWithEnvironment(s.Condition, env)
				if !condition.IsTruthy() {
					break
				}
			}

			// body execution
			EvalBlockStatement(s.Body, env)

			// execute update statement if present
			if s.Update != nil {
				EvalStatement(s.Update, env)
			}
		}
	case *ast.SwitchStatement:
		// Evaluate the switch value
		switchValue := EvalValueWithEnvironment(s.Value, env)

		// Try to match each case
		matched := false
		for _, caseStmt := range s.Cases {
			caseValue := EvalValueWithEnvironment(caseStmt.Value, env)
			// Simple equality check - could be enhanced for type-aware comparison
			if switchValue.String() == caseValue.String() {
				EvalBlockStatement(caseStmt.Body, env)
				matched = true
				break // Go switch has implicit break
			}
		}

		// If no case matched, execute default block
		if !matched && s.Default != nil {
			EvalBlockStatement(s.Default, env)
		}
	case *ast.BlockStatement:
		EvalBlockStatement(s, env)
	case *ast.BreakStatement:
		// TODO: implement proper break handling
	case *ast.ContinueStatement:
		// TODO: implement proper continue handling
	case *ast.FuncStatement:
		// Register function in environment
		function := &Function{
			Name:       s.Name,
			Parameters: s.Parameters,
			ReturnType: s.ReturnType,
			Body:       s.Body,
		}
		env.SetFunction(s.Name, function)
	case *ast.ReturnStatement:
		// Handle return statement with exception (type-aware)
		var value Value = &IntValue{Value: 0}
		if s.Value != nil {
			value = EvalValueWithEnvironment(s.Value, env)
		}
		// For now, still use int for return values (backward compatibility)
		if intVal, ok := value.(*IntValue); ok {
			panic(&ReturnException{Value: intVal.Value})
		} else {
			panic(&ReturnException{Value: 0})
		}
	case *ast.StructDefinition:
		// Register struct definition in environment
		env.SetStruct(s.Name, s)
	case *ast.PackageStatement:
		// Handle package declaration
		// For now, just store the package name in environment
		env.SetPackage(s.Name)
	case *ast.ImportStatement:
		// Handle import declaration
		// For now, just store the import path in environment
		env.AddImport(s.Path)
	}
}

func EvalBlockStatement(block *ast.BlockStatement, env *Environment) {
	if block != nil && block.Statements != nil {
		for _, stmt := range block.Statements {
			EvalStatement(stmt, env)
		}
	}
}

// callUserFunction calls a user-defined function with arguments
func callUserFunction(function *Function, args []ast.ASTNode, env *Environment) int {
	// Create new scope for function execution
	localEnv := NewEnvironment()

	// Copy current environment's functions to local scope
	for name, fn := range env.functions {
		localEnv.SetFunction(name, fn)
	}

	// Bind arguments to parameters (type-aware with type checking)
	for i, param := range function.Parameters {
		if i < len(args) {
			value := EvalValueWithEnvironment(args[i], env)

			// Type checking: verify argument type matches parameter type
			expectedType := param.Type
			actualType := value.Type()

			if expectedType != actualType {
				// Type mismatch - create zero value of expected type
				switch expectedType {
				case "int":
					value = &IntValue{Value: 0}
				case "string":
					value = &StringValue{Value: ""}
				case "bool":
					value = &BoolValue{Value: false}
				default:
					// Unknown type, keep original value
				}
			}

			localEnv.Set(param.Name, value)
		} else {
			// Missing argument - set zero value of parameter type
			switch param.Type {
			case "int":
				localEnv.Set(param.Name, &IntValue{Value: 0})
			case "string":
				localEnv.Set(param.Name, &StringValue{Value: ""})
			case "bool":
				localEnv.Set(param.Name, &BoolValue{Value: false})
			default:
				localEnv.Set(param.Name, &IntValue{Value: 0})
			}
		}
	}

	// Execute function body with return handling
	returnValue := 0
	if function.Body != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if returnEx, ok := r.(*ReturnException); ok {
						// Capture return value
						returnValue = returnEx.Value
					} else {
						// Re-panic for other exceptions
						panic(r)
					}
				}
			}()

			EvalBlockStatement(function.Body, localEnv)
		}()
	}

	return returnValue
}

// evalStructLiteral evaluates struct literal expressions
func evalStructLiteral(node *ast.StructLiteral, env *Environment) Value {
	// Get struct definition
	structDef, exists := env.GetStruct(node.TypeName)
	if !exists {
		// Unknown struct type: return empty struct for now
		return &StructValue{
			TypeName: node.TypeName,
			Fields:   make(map[string]Value),
		}
	}

	// Evaluate field values
	fields := make(map[string]Value)
	for fieldName, fieldExpr := range node.Fields {
		value := EvalValueWithEnvironment(fieldExpr, env)
		fields[fieldName] = value
	}

	// Initialize missing fields with zero values
	for _, field := range structDef.Fields {
		if _, exists := fields[field.Name]; !exists {
			// Set default zero values based on type
			switch field.Type {
			case "int":
				fields[field.Name] = &IntValue{Value: 0}
			case "string":
				fields[field.Name] = &StringValue{Value: ""}
			case "bool":
				fields[field.Name] = &BoolValue{Value: false}
			default:
				// Unknown type: set to int 0
				fields[field.Name] = &IntValue{Value: 0}
			}
		}
	}

	return &StructValue{
		TypeName: node.TypeName,
		Fields:   fields,
	}
}

// evalFieldAccess evaluates field access expressions (obj.field)
func evalFieldAccess(node *ast.FieldAccessNode, env *Environment) Value {
	obj := EvalValueWithEnvironment(node.Object, env)

	// Check if object is a struct
	structVal, ok := obj.(*StructValue)
	if !ok {
		// Not a struct: return zero value
		return &IntValue{Value: 0}
	}

	// Get field value
	if fieldValue, exists := structVal.Fields[node.Field]; exists {
		return fieldValue
	}

	// Field not found: return zero value
	return &IntValue{Value: 0}
}

// evalSliceLiteral evaluates slice literal expressions
func evalSliceLiteral(node *ast.SliceLiteral, env *Environment) Value {
	var elements []Value

	// Evaluate each element
	for _, elem := range node.Elements {
		value := EvalValueWithEnvironment(elem, env)
		elements = append(elements, value)
	}

	return &SliceValue{
		ElementType: node.ElementType,
		Elements:    elements,
	}
}

// evalIndexAccess evaluates index access expressions (slice[index])
func evalIndexAccess(node *ast.IndexAccess, env *Environment) Value {
	obj := EvalValueWithEnvironment(node.Object, env)
	index := EvalValueWithEnvironment(node.Index, env)

	// Check if object is a slice
	sliceVal, ok := obj.(*SliceValue)
	if !ok {
		// Not a slice: return zero value
		return &IntValue{Value: 0}
	}

	// Check if index is an integer
	indexVal, ok := index.(*IntValue)
	if !ok {
		// Invalid index type: return zero value
		return &IntValue{Value: 0}
	}

	// Bounds checking
	if indexVal.Value < 0 || indexVal.Value >= len(sliceVal.Elements) {
		// Out of bounds: return zero value
		return &IntValue{Value: 0}
	}

	return sliceVal.Elements[indexVal.Value]
}
