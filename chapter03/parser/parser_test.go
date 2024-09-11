package parser

import (
	"fmt"
	goast "go/ast"
	"testing"

	"github.com/SeaOfNodes/Simple-Go/chapter03/utils/ast"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func (suite *ParserTestSuite) equalAST(expected goast.Node, given goast.Node, failMsg string) {
	suite.IsType(expected, given, failMsg)
	switch e := expected.(type) {
	case *goast.BinaryExpr:
		g := given.(*goast.BinaryExpr)
		suite.equalAST(e.X, g.X, failMsg)
		suite.equalAST(e.Y, g.Y, failMsg)
		suite.Equal(e.Op, g.Op, failMsg)
		suite.NotZero(g.OpPos, failMsg)
	case *goast.UnaryExpr:
		g := given.(*goast.UnaryExpr)
		suite.equalAST(e.X, g.X, failMsg)
		suite.Equal(e.Op, g.Op, failMsg)
		suite.NotZero(g.OpPos, failMsg)
	case *goast.BasicLit:
		g := given.(*goast.BasicLit)
		suite.Equal(e.Value, g.Value, failMsg)
		suite.Equal(e.Kind, g.Kind, failMsg)
		suite.NotZero(g.ValuePos, failMsg)
	case *goast.ParenExpr:
		g := given.(*goast.ParenExpr)
		suite.equalAST(e.X, g.X, failMsg)
	default:
		suite.FailNow("Unexpected type", "Type: %T", e)
	}
}

func (suite *ParserTestSuite) TestPrecedence() {
	subTests := []struct {
		name     string
		input    string
		expected goast.Node
	}{
		{
			name:     "mul->add",
			input:    "1*2+3",
			expected: ast.Bin(ast.Bin(1, "*", 2), "+", 3),
		},
		{
			name:     "add->mul",
			input:    "1+2*3",
			expected: ast.Bin(1, "+", ast.Bin(2, "*", 3)),
		},
		{
			name:     "div->add",
			input:    "1/2+3",
			expected: ast.Bin(ast.Bin(1, "/", 2), "+", 3),
		},
		{
			name:     "add->div",
			input:    "1+2/3",
			expected: ast.Bin(1, "+", ast.Bin(2, "/", 3)),
		},
		{
			name:     "mul->sub",
			input:    "1*2-3",
			expected: ast.Bin(ast.Bin(1, "*", 2), "-", 3),
		},
		{
			name:     "sub->mul",
			input:    "1-2*3",
			expected: ast.Bin(1, "-", ast.Bin(2, "*", 3)),
		},
		{
			name:     "div->sub",
			input:    "1/2-3",
			expected: ast.Bin(ast.Bin(1, "/", 2), "-", 3),
		},
		{
			name:     "sub->div",
			input:    "1-2/3",
			expected: ast.Bin(1, "-", ast.Bin(2, "/", 3)),
		},
		{
			name:     "minus->add",
			input:    "-1+2",
			expected: ast.Bin(ast.Un("-", 1), "+", 2),
		},
		{
			name:     "add->minus",
			input:    "1+-2",
			expected: ast.Bin(1, "+", ast.Un("-", 2)),
		},
		{
			name:     "minus->mul",
			input:    "-1*2",
			expected: ast.Bin(ast.Un("-", 1), "*", 2),
		},
		{
			name:     "mul->minus",
			input:    "1*-2",
			expected: ast.Bin(1, "*", ast.Un("-", 2)),
		},
		{
			name:     "sub->add",
			input:    "1-2+3",
			expected: ast.Bin(ast.Bin(1, "-", 2), "+", 3),
		},
		{
			name:     "add->sub",
			input:    "1+2-3",
			expected: ast.Bin(ast.Bin(1, "+", 2), "-", 3),
		},
		{
			name:     "div->mul",
			input:    "1/2*3",
			expected: ast.Bin(ast.Bin(1, "/", 2), "*", 3),
		},
		{
			name:     "mul->div",
			input:    "1*2/3",
			expected: ast.Bin(ast.Bin(1, "*", 2), "/", 3),
		},
		{
			name:     "all",
			input:    "1*2+3/-4-5",
			expected: ast.Bin(ast.Bin(ast.Bin(1, "*", 2), "+", ast.Bin(3, "/", ast.Un("-", 4))), "-", 5),
		},
		{
			name:     "paren",
			input:    "1*(2+3)",
			expected: ast.Bin(1, "*", ast.Paren(ast.Bin(2, "+", 3))),
		},
		{
			name:     "manyParen",
			input:    "1*(2+(-3)*(1+(1-1)))",
			expected: ast.Bin(1, "*", ast.Paren(ast.Bin(2, "+", ast.Bin(ast.Paren(ast.Un("-", 3)), "*", ast.Paren(ast.Bin(1, "+", ast.Paren(ast.Bin(1, "-", 1)))))))),
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
			p := NewParser(test.input)
			_, err := p.Parse()
			suite.ErrorContains(err, "expected ;")
		})
	}
}

func (suite *ParserTestSuite) TestBlockNotClosed() {
	subTests := []struct {
		name  string
		input string
	}{
		{name: "start", input: "{ return 1;"},
		{name: "end", input: "return 1; {"},
		{name: "in block", input: "{ int a = 2; { }"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			p := NewParser(test.input)
			_, err := p.Parse()
			suite.ErrorContains(err, "expected a statement got")
		})
	}

}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
