package ir

import (
	"testing"

	goast "go/ast"

	"github.com/SeaOfNodes/Simple-Go/chapter03/utils/ast"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
}

func (suite *GeneratorTestSuite) TestPrint() {
	DisablePeephole = true
	ret := ast.Ret(ast.Bin(ast.Bin(1, "+", ast.Bin(2, "*", 3)), "+", ast.Un("-", 5)))
	retNode, err := NewGenerator().Generate(ast.Block(ret))
	suite.NoError(err)
	suite.Equal("return ((1+(2*3))+(-5));", ToString(retNode))
	DisablePeephole = false
}

func (suite *GeneratorTestSuite) TestOperations() {
	subTests := []struct {
		name     string
		input    *goast.ReturnStmt
		expected string
	}{
		{name: "add", input: ast.Ret(ast.Bin(1, "+", 2)), expected: "return 3;"},
		{name: "sub", input: ast.Ret(ast.Bin(1, "-", 2)), expected: "return -1;"},
		{name: "mul", input: ast.Ret(ast.Bin(2, "*", 3)), expected: "return 6;"},
		{name: "div", input: ast.Ret(ast.Bin(6, "/", 2)), expected: "return 3;"},
		{name: "minus", input: ast.Ret(ast.Bin(1, "-", ast.Un("-", 2))), expected: "return 3;"},
		{name: "arithmetic", input: ast.Ret(ast.Bin(1, "+", ast.Bin(ast.Bin(2, "*", 3), "+", ast.Un("-", 5)))), expected: "return 2;"},
		{name: "paren", input: ast.Ret(ast.Paren(1)), expected: "return 1;"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			retNode, err := NewGenerator().Generate(ast.Block(test.input))
			suite.NoError(err)
			suite.Equal(test.expected, ToString(retNode))
		})
	}
}

func (suite *GeneratorTestSuite) TestVars() {
	subTests := []struct {
		name     string
		input    *goast.BlockStmt
		expected string
	}{
		{name: "ret", input: ast.Block(ast.Decl("a", 1), ast.Ret("a")), expected: "return 1;"},
		{name: "assign to int", input: ast.Block(ast.Decl("a", 1), ast.Assign("a", 2), ast.Ret("a")), expected: "return 2;"},
		{name: "assign to var", input: ast.Block(ast.Decl("a", 1), ast.Decl("b", 2), ast.Assign("a", "b"), ast.Ret("a")), expected: "return 2;"},
		{name: "assign in block", input: ast.Block(ast.Decl("a", 1), ast.Block(ast.Assign("a", 2)), ast.Ret("a")), expected: "return 2;"},
		{name: "arithmetic", input: ast.Block(ast.Decl("a", 1), ast.Decl("b", 3), ast.Ret(ast.Bin("a", "+", "b"))), expected: "return 4;"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			retNode, err := NewGenerator().Generate(ast.Block(test.input))
			suite.NoError(err)
			suite.Equal(test.expected, ToString(retNode))
		})
	}
}

func (suite *GeneratorTestSuite) TestUnknownIdent() {
	subTests := []struct {
		name  string
		input goast.Stmt
	}{
		{name: "ret", input: ast.Ret("a")},
		{name: "assign", input: ast.Assign("a", 1)},
		{name: "decl", input: ast.Decl("b", "a")},
		{name: "decl in block", input: ast.Block(ast.Block(ast.Decl("a", 1)), ast.Ret("a"))},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			_, err := NewGenerator().Generate(ast.Block(test.input))
			suite.ErrorContains(err, "unknown identifier")
		})
	}
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
