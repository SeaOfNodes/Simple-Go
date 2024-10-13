package ir

import (
	"testing"

	goast "go/ast"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
	"github.com/SeaOfNodes/Simple-Go/chapter04/utils/ast"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
}

func (suite *GeneratorTestSuite) TestPrint() {
	DisablePeephole = true
	ret := ast.Ret(ast.Bin(ast.Bin(1, "+", ast.Bin(2, "*", 3)), "+", ast.Un("-", 5)))
	retNode, err := NewGenerator(types.Bottom).Generate(ast.Block(ret))
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
			retNode, err := NewGenerator(types.Bottom).Generate(ast.Block(test.input))
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
			retNode, err := NewGenerator(types.Bottom).Generate(ast.Block(test.input))
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
			_, err := NewGenerator(types.Bottom).Generate(ast.Block(test.input))
			suite.ErrorContains(err, "unknown identifier")
		})
	}
}

func (suite *GeneratorTestSuite) TestArg() {
	subTests := []struct {
		name     string
		input    *goast.BlockStmt
		expected string
	}{
		{name: "add1", input: ast.Block(ast.Ret(ast.Bin(ast.Bin(1, "+", "arg"), "+", 2))), expected: "return (arg+3);"},
		{name: "add2", input: ast.Block(ast.Ret(ast.Bin(2, "+", ast.Bin(1, "+", "arg")))), expected: "return (arg+3);"},
		{name: "add zero", input: ast.Block(ast.Ret(ast.Bin(0, "+", "arg"))), expected: "return arg;"},
		{name: "add to mul", input: ast.Block(ast.Ret(ast.Bin("arg", "+", "arg"))), expected: "return (arg*2);"},
		{name: "add to mul and consts", input: ast.Block(ast.Ret(ast.Bin(1, "+", ast.Bin("arg", "+", ast.Bin(ast.Bin(2, "+", "arg"), "+", 3))))), expected: "return ((arg*2)+6);"},
		{name: "mul one", input: ast.Block(ast.Ret(ast.Bin(1, "*", "arg"))), expected: "return arg;"},
		{name: "div one", input: ast.Block(ast.Ret(ast.Bin("arg", "/", 1))), expected: "return arg;"},
		{name: "minus add", input: ast.Block(ast.Ret(ast.Un("-", ast.Bin(1, "-", "arg")))), expected: "return (arg-1);"},
		{name: "notnotnot", input: ast.Block(ast.Ret(ast.Un("!", ast.Un("!", ast.Un("!", "arg"))))), expected: "return (!arg);"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			retNode, err := NewGenerator(types.Bottom).Generate(test.input)
			suite.NoError(err)
			suite.Equal(test.expected, ToString(retNode))
		})
	}
}

func (suite *GeneratorTestSuite) TestConstArg() {
	retNode, err := NewGenerator(types.NewInt(2)).Generate(ast.Block(ast.Ret("arg")))
	suite.NoError(err)
	suite.Equal("return 2;", ToString(retNode))
}

func (suite *GeneratorTestSuite) TestBool() {
	subTests := []struct {
		name     string
		input    *goast.BlockStmt
		expected string
	}{
		{name: "eq true", input: ast.Block(ast.Ret(ast.Bin(3, "==", 3))), expected: "return 1;"},
		{name: "eq false", input: ast.Block(ast.Ret(ast.Bin(3, "==", 4))), expected: "return 0;"},
		{name: "neq true", input: ast.Block(ast.Ret(ast.Bin(3, "!=", 4))), expected: "return 1;"},
		{name: "neq false", input: ast.Block(ast.Ret(ast.Bin(3, "!=", 3))), expected: "return 0;"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			retNode, err := NewGenerator(types.Bottom).Generate(test.input)
			suite.NoError(err)
			suite.Equal(test.expected, ToString(retNode))
		})
	}
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
