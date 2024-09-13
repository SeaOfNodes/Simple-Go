package simple_test

import (
	"testing"

	simple "github.com/SeaOfNodes/Simple-Go/chapter03"
	"github.com/SeaOfNodes/Simple-Go/chapter03/ir"
	"github.com/SeaOfNodes/Simple-Go/chapter03/ir/types"
	"github.com/stretchr/testify/suite"
)

type SimpleTestSuite struct {
	suite.Suite
}

func (suite *SimpleTestSuite) TestValidPrograms() {
	subTests := []struct {
		name  string
		input string
		num   int
	}{
		{name: "One", input: "return 1;", num: 1},
		{name: "Zero", input: "return 0;", num: 0},
		{name: "MinusNumber", input: "return -2;", num: -2},
		{name: "EmptyBlock", input: "{ } return -2;", num: -2},
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, _, err := simple.Simple(test.input)
			suite.Require().NoError(err)
			suite.Equal(ir.StartNode, ret.Control())

			expr := ret.Expr()
			suite.IsType(&ir.ConstantNode{}, expr)
			suite.Equal(ir.StartNode, ir.In(expr, 0))
			typ := ir.Type(expr)
			suite.IsType(&types.IntType{}, typ)
			suite.Equal(test.num, typ.(*types.IntType).Value)
		})
	}
}

func (suite *SimpleTestSuite) TestInvalidPrograms() {
	subTests := []struct {
		name  string
		input string
		error string
	}{
		{name: "InvalidStatement", input: "ret", error: "Syntax error: expected assignment"},
		{name: "InvalidNumber", input: "return 0123;", error: "Syntax error: integer values cannot start with '0'"},
		{name: "MissingSemicolon", input: "return 123", error: "Syntax error: expected ; after expression"},
		{name: "MissingWhitespace", input: "return123;", error: "Syntax error: expected assignment"},
		{name: "ByteAfterSemicolon", input: "return 1;}", error: "Syntax error: expected a statement got }"},
		{name: "SelfAssign", input: "int a=a; return a;", error: "Compute error: unknown identifier"},
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, _, err := simple.Simple(test.input)
			suite.IsType(&simple.SourceError{}, err)
			suite.Contains(err.Error(), test.error)
			suite.Nil(ret)
		})
	}
}

func (suite *SimpleTestSuite) TestVars() {
	subTests := []struct {
		name   string
		input  string
		output string
	}{
		{name: "OneDecl", input: "int a = 1; return a;", output: "return 1;"},
		{name: "Add", input: "int a = 1; int b = 2; return a+b;", output: "return 3;"},
		{name: "Scope", input: "int a=1; int b=2; int c=0; { int b=3; c=a+b; } return c;", output: "return 4;"},
		{name: "Dist", input: "int x0=1; int y0=2; int x1=3; int y1=4; return (x0-x1)*(x0-x1) + (y0-y1)*(y0-y1);", output: "return 8;"},
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, _, err := simple.Simple(test.input)
			suite.NoError(err)
			suite.Equal(test.output, ir.ToString(ret))
		})
	}
}

func TestSimple(t *testing.T) {
	suite.Run(t, new(SimpleTestSuite))
}
