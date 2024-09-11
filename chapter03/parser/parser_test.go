package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func astNum(i int) *ast.BasicLit {
	return &ast.BasicLit{Value: strconv.Itoa(i), Kind: token.INT}
}

func astExpr(a any) ast.Expr {
	switch t := a.(type) {
	case int:
		return astNum(t)
	case ast.Expr:
		return t
	}
	return nil
}

func astBin(lhs any, op string, rhs any) *ast.BinaryExpr {
	return &ast.BinaryExpr{X: astExpr(lhs), Op: astOp(op), Y: astExpr(rhs)}
}

func astUn(op string, value any) *ast.UnaryExpr {
	return &ast.UnaryExpr{X: astExpr(value), Op: astOp(op)}
}

func astOp(op string) token.Token {
	switch op {
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	}
	return 0
}

func (suite *ParserTestSuite) equalAST(expected ast.Node, given ast.Node, failMsg string) {
	suite.IsType(expected, given, failMsg)
	switch e := expected.(type) {
	case *ast.BinaryExpr:
		g := given.(*ast.BinaryExpr)
		suite.equalAST(e.X, g.X, failMsg)
		suite.equalAST(e.Y, g.Y, failMsg)
		suite.Equal(e.Op, g.Op, failMsg)
		suite.NotZero(g.OpPos, failMsg)
	case *ast.UnaryExpr:
		g := given.(*ast.UnaryExpr)
		suite.equalAST(e.X, g.X, failMsg)
		suite.Equal(e.Op, g.Op, failMsg)
		suite.NotZero(g.OpPos, failMsg)
	case *ast.BasicLit:
		g := given.(*ast.BasicLit)
		suite.Equal(e.Value, g.Value, failMsg)
		suite.Equal(e.Kind, g.Kind, failMsg)
		suite.NotZero(g.ValuePos, failMsg)
	default:
		suite.FailNow("Unexpected type: %T", e)
	}
}

func (suite *ParserTestSuite) TestPrecedence() {
	subTests := []struct {
		name     string
		input    string
		expected ast.Node
	}{
		{
			name:     "mul->add",
			input:    "1*2+3",
			expected: astBin(astBin(1, "*", 2), "+", 3),
		},
		{
			name:     "add->mul",
			input:    "1+2*3",
			expected: astBin(1, "+", astBin(2, "*", 3)),
		},
		{
			name:     "div->add",
			input:    "1/2+3",
			expected: astBin(astBin(1, "/", 2), "+", 3),
		},
		{
			name:     "add->div",
			input:    "1+2/3",
			expected: astBin(1, "+", astBin(2, "/", 3)),
		},
		{
			name:     "mul->sub",
			input:    "1*2-3",
			expected: astBin(astBin(1, "*", 2), "-", 3),
		},
		{
			name:     "sub->mul",
			input:    "1-2*3",
			expected: astBin(1, "-", astBin(2, "*", 3)),
		},
		{
			name:     "div->sub",
			input:    "1/2-3",
			expected: astBin(astBin(1, "/", 2), "-", 3),
		},
		{
			name:     "sub->div",
			input:    "1-2/3",
			expected: astBin(1, "-", astBin(2, "/", 3)),
		},
		{
			name:     "minus->add",
			input:    "-1+2",
			expected: astBin(astUn("-", 1), "+", 2),
		},
		{
			name:     "add->minus",
			input:    "1+-2",
			expected: astBin(1, "+", astUn("-", 2)),
		},
		{
			name:     "minus->mul",
			input:    "-1*2",
			expected: astBin(astUn("-", 1), "*", 2),
		},
		{
			name:     "mul->minus",
			input:    "1*-2",
			expected: astBin(1, "*", astUn("-", 2)),
		},
		{
			name:     "sub->add",
			input:    "1-2+3",
			expected: astBin(astBin(1, "-", 2), "+", 3),
		},
		{
			name:     "add->sub",
			input:    "1+2-3",
			expected: astBin(astBin(1, "+", 2), "-", 3),
		},
		{
			name:     "div->mul",
			input:    "1/2*3",
			expected: astBin(astBin(1, "/", 2), "*", 3),
		},
		{
			name:     "mul->div",
			input:    "1*2/3",
			expected: astBin(astBin(1, "*", 2), "/", 3),
		},
		{
			name:     "all",
			input:    "1*2+3/-4-5",
			expected: astBin(astBin(astBin(1, "*", 2), "+", astBin(3, "/", astUn("-", 4))), "-", 5),
		},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			p := NewParser(test.input)
			n, err := p.parseBinary()
			suite.NoError(err)
			failMsg := fmt.Sprintf("input:\n\t%s\nparsed:\n\t%s\n", test.input, p.string(n))
			suite.equalAST(test.expected, n, failMsg)
		})
	}
}

func (suite *ParserTestSuite) TestMissingSemicolon() {
	subTests := []struct {
		name  string
		input string
	}{
		{name: "var decl", input: "int a = 2"},
		{name: "return", input: "return 3"},
		{name: "in block", input: "{ int a = 2 }"},
		{name: "multiple statements", input: "int a = 2 int b = 3"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			p := NewParser("int a = 2")
			_, err := p.Parse()
			suite.Errorf(err, "expected ;")
		})
	}
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
