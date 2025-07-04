package codegen

import (
	"fmt"
	"strings"

	"github.com/yuya-takeyama/petitgo/ast"
	"github.com/yuya-takeyama/petitgo/token"
)

// Generator converts AST to Go source code
type Generator struct {
	output strings.Builder
	indent int
}

// NewGenerator creates a new code generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate converts an AST node to Go source code
func (g *Generator) Generate(node ast.ASTNode) string {
	g.output.Reset()
	g.generateNode(node)
	return g.output.String()
}

// GenerateProgram converts a list of statements to a complete Go program
func (g *Generator) GenerateProgram(statements []ast.Statement) string {
	g.output.Reset()

	// Add package declaration
	g.writeLine("package main")
	g.writeLine("")

	// Add imports - none needed for println as it's built-in
	g.writeLine("")

	// Check if there's already a main function
	hasMainFunc := false
	for _, stmt := range statements {
		if funcStmt, ok := stmt.(*ast.FuncStatement); ok && funcStmt.Name == "main" {
			hasMainFunc = true
			break
		}
	}

	// Add main function wrapper if no main function exists
	if !hasMainFunc {
		g.writeLine("func main() {")
		g.indent++
	}

	// Generate statements
	for _, stmt := range statements {
		g.generateStatement(stmt)
		g.writeLine("")
	}

	// Close main function if we added it
	if !hasMainFunc {
		g.indent--
		g.writeLine("}")
	}

	return g.output.String()
}

func (g *Generator) generateNode(node ast.ASTNode) {
	switch n := node.(type) {
	case *ast.NumberNode:
		g.write(fmt.Sprintf("%d", n.Value))
	case *ast.StringNode:
		g.write(fmt.Sprintf("\"%s\"", n.Value))
	case *ast.BooleanNode:
		if n.Value {
			g.write("true")
		} else {
			g.write("false")
		}
	case *ast.VariableNode:
		g.write(n.Name)
	case *ast.FieldAccessNode:
		g.generateNode(n.Object)
		g.write(".")
		g.write(n.Field)
	case *ast.BinaryOpNode:
		// Always wrap binary operations in parentheses to preserve exact precedence
		g.write("(")
		g.generateNode(n.Left)
		g.write(" ")
		g.generateOperator(n.Operator)
		g.write(" ")
		g.generateNode(n.Right)
		g.write(")")
	case *ast.CallNode:
		if n.Function == "print" {
			g.write("println(")
		} else {
			g.write(n.Function + "(")
		}
		for i, arg := range n.Arguments {
			if i > 0 {
				g.write(", ")
			}
			g.generateNode(arg)
		}
		g.write(")")
	case *ast.StructLiteral:
		g.write(n.TypeName + "{")
		first := true
		for fieldName, fieldValue := range n.Fields {
			if !first {
				g.write(", ")
			}
			first = false
			g.write(fieldName + ": ")
			g.generateNode(fieldValue)
		}
		g.write("}")
	default:
		g.write("/* unsupported node */")
	}
}

func (g *Generator) generateStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		g.writeIndent()
		g.write(fmt.Sprintf("var %s %s = ", s.Name, s.TypeName))
		g.generateNode(s.Value)
	case *ast.AssignStatement:
		g.writeIndent()
		g.write(s.Name + " := ")
		g.generateNode(s.Value)
	case *ast.ReassignStatement:
		g.writeIndent()
		g.write(s.Name + " = ")
		g.generateNode(s.Value)
	case *ast.IncStatement:
		g.writeIndent()
		g.write(s.Name + "++")
	case *ast.DecStatement:
		g.writeIndent()
		g.write(s.Name + "--")
	case *ast.CompoundAssignStatement:
		g.writeIndent()
		operatorStr := ""
		switch s.Operator {
		case token.ADD_ASSIGN:
			operatorStr = "+="
		case token.SUB_ASSIGN:
			operatorStr = "-="
		case token.MUL_ASSIGN:
			operatorStr = "*="
		case token.QUO_ASSIGN:
			operatorStr = "/="
		}
		g.write(s.Name + " " + operatorStr + " ")
		g.generateNode(s.Value)
	case *ast.ExpressionStatement:
		g.writeIndent()
		g.generateNode(s.Expression)
	case *ast.FuncStatement:
		g.writeIndent()
		g.write(fmt.Sprintf("func %s(", s.Name))
		for i, param := range s.Parameters {
			if i > 0 {
				g.write(", ")
			}
			g.write(fmt.Sprintf("%s %s", param.Name, param.Type))
		}
		g.write(")")
		if s.ReturnType != "" {
			g.write(" " + s.ReturnType)
		}
		g.write(" ")
		g.generateBlockStatement(s.Body)
	case *ast.ReturnStatement:
		g.writeIndent()
		g.write("return")
		if s.Value != nil {
			g.write(" ")
			g.generateNode(s.Value)
		}
	case *ast.IfStatement:
		g.writeIndent()
		g.write("if ")
		g.generateNode(s.Condition)
		g.write(" ")
		g.generateBlockStatement(s.ThenBlock)
		if s.ElseBlock != nil {
			g.write(" else ")
			g.generateBlockStatement(s.ElseBlock)
		}
	case *ast.ForStatement:
		g.writeIndent()
		g.write("for ")

		// Generate init, condition, update if present
		if s.Init != nil || s.Condition != nil || s.Update != nil {
			if s.Init != nil {
				g.generateStatement(s.Init)
				g.write("; ")
			} else {
				g.write("; ")
			}

			if s.Condition != nil {
				g.generateNode(s.Condition)
			}
			g.write("; ")

			if s.Update != nil {
				g.generateStatement(s.Update)
			}
		} else if s.Condition != nil {
			// condition-only form
			g.generateNode(s.Condition)
		}

		g.write(" ")
		g.generateBlockStatement(s.Body)
	case *ast.SwitchStatement:
		g.writeIndent()
		g.write("switch ")
		g.generateNode(s.Value)
		g.write(" {\n")
		g.indent++

		// Generate cases
		for _, caseStmt := range s.Cases {
			g.writeIndent()
			g.write("case ")
			g.generateNode(caseStmt.Value)
			g.write(":\n")
			g.indent++
			for _, stmt := range caseStmt.Body.Statements {
				g.generateStatement(stmt)
				g.writeLine("")
			}
			g.indent--
		}

		// Generate default case if present
		if s.Default != nil {
			g.writeIndent()
			g.write("default:\n")
			g.indent++
			for _, stmt := range s.Default.Statements {
				g.generateStatement(stmt)
				g.writeLine("")
			}
			g.indent--
		}

		g.indent--
		g.writeIndent()
		g.write("}")
	case *ast.TypeStatement:
		g.writeIndent()
		g.write("type " + s.Name + " struct {\n")
		g.indent++
		for _, field := range s.Fields {
			g.writeIndent()
			g.write(field.Name + " " + field.Type + "\n")
		}
		g.indent--
		g.writeIndent()
		g.write("}")
	case *ast.InterfaceStatement:
		g.writeIndent()
		g.write("type " + s.Name + " interface {\n")
		g.indent++
		for _, method := range s.Methods {
			g.writeIndent()
			g.write(method.Name + "(")
			for i, param := range method.Parameters {
				if i > 0 {
					g.write(", ")
				}
				g.write(param.Name + " " + param.Type)
			}
			g.write(")")
			if method.ReturnType != "" {
				g.write(" " + method.ReturnType)
			}
			g.write("\n")
		}
		g.indent--
		g.writeIndent()
		g.write("}")
	case *ast.BlockStatement:
		g.generateBlockStatement(s)
	default:
		g.writeIndent()
		g.write("/* unsupported statement */")
	}
}

func (g *Generator) generateBlockStatement(block *ast.BlockStatement) {
	g.write("{\n")
	g.indent++

	if block != nil && block.Statements != nil {
		for _, stmt := range block.Statements {
			g.generateStatement(stmt)
			g.writeLine("")
		}
	}

	g.indent--
	g.writeIndent()
	g.write("}")
}

func (g *Generator) generateOperator(op token.Token) {
	// Convert token to Go operator
	switch op {
	case token.ADD:
		g.write("+")
	case token.SUB:
		g.write("-")
	case token.MUL:
		g.write("*")
	case token.QUO:
		g.write("/")
	case token.EQL:
		g.write("==")
	case token.NEQ:
		g.write("!=")
	case token.LSS:
		g.write("<")
	case token.GTR:
		g.write(">")
	case token.LEQ:
		g.write("<=")
	case token.GEQ:
		g.write(">=")
	default:
		g.write("/*unknown op*/")
	}
}

func (g *Generator) write(s string) {
	g.output.WriteString(s)
}

func (g *Generator) writeLine(s string) {
	g.output.WriteString(s + "\n")
}

func (g *Generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		g.output.WriteString("\t")
	}
}
