package simple_test

import (
	"testing"

	simple "github.com/SeaOfNodes/Simple-Go/chapter02"
	"github.com/SeaOfNodes/Simple-Go/chapter02/ir"
	"github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"
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
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, err := simple.Simple(test.input)
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
		{name: "InvalidStatement", input: "ret", error: "Syntax error: expected a statement got ret"},
		{name: "InvalidNumber", input: "return 0123;", error: "Syntax error: integer values cannot start with '0'"},
		{name: "MissingSemicolon", input: "return 123", error: "Syntax error: expected ; after expression"},
		{name: "MissingWhitespace", input: "return123;", error: "Syntax error: expected a statement got return123"},
		{name: "ByteAfterSemicolon", input: "return 1;}", error: "Syntax error: unexpected }"},
		{name: "DivideByZero", input: "return 0/0;", error: "Compute error: divide by zero"},
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, err := simple.Simple(test.input)
			suite.IsType(&simple.SourceError{}, err)
			suite.Contains(err.Error(), test.error)
			suite.Nil(ret)
		})
	}
}

func TestSimple(t *testing.T) {
	suite.Run(t, new(SimpleTestSuite))
}
