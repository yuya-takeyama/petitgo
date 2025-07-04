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

	// Add imports
	g.writeLine("import (")
	g.writeLine("\t\"os\"")
	g.writeLine(")")
	g.writeLine("")

	// Add custom print function
	g.addPrintFunction()

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
	case *ast.BinaryOpNode:
		g.generateNode(n.Left)
		g.write(" ")
		g.generateOperator(n.Operator)
		g.write(" ")
		g.generateNode(n.Right)
	case *ast.CallNode:
		if n.Function == "print" {
			g.write("printInt(")
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

func (g *Generator) addPrintFunction() {
	g.writeLine("// Custom print function for petitgo compatibility")
	g.writeLine("func printInt(n int) {")
	g.writeLine("\tif n == 0 {")
	g.writeLine("\t\tos.Stdout.Write([]byte(\"0\"))")
	g.writeLine("\t\treturn")
	g.writeLine("\t}")
	g.writeLine("")
	g.writeLine("\tif n < 0 {")
	g.writeLine("\t\tos.Stdout.Write([]byte(\"-\"))")
	g.writeLine("\t\tn = -n")
	g.writeLine("\t}")
	g.writeLine("")
	g.writeLine("\tdigits := []byte{}")
	g.writeLine("\tfor n > 0 {")
	g.writeLine("\t\tdigits = append([]byte{byte('0' + n%10)}, digits...)")
	g.writeLine("\t\tn /= 10")
	g.writeLine("\t}")
	g.writeLine("\tos.Stdout.Write(digits)")
	g.writeLine("}")
	g.writeLine("")
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
