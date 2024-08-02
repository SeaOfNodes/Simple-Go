package simple_test

import (
	"testing"

	simple "github.com/SeaOfNodes/Simple-Go/chapter01"
	"github.com/SeaOfNodes/Simple-Go/chapter01/node"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func (suite *ParserTestSuite) TestValidPrograms() {
	subTests := []struct {
		name  string
		input string
		num   int
	}{
		{name: "One", input: "return 1;", num: 1},
		{name: "Zero", input: "return 0;", num: 0},
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, err := simple.Simple(test.input)
			suite.NoError(err)
			suite.Equal(node.StartNode, ret.Control())

			expr := ret.Expr()
			suite.IsType(&node.ConstantNode{}, expr)
			suite.Equal(node.StartNode, expr.In(0))
			suite.Equal(test.num, expr.(*node.ConstantNode).Value)
		})
	}
}

func (suite *ParserTestSuite) TestInvalidPrograms() {
	subTests := []struct {
		name  string
		input string
		error string
	}{
		{name: "InvalidStatement", input: "ret", error: "Syntax error: expected a statement got ret"},
		{name: "InvalidNumber", input: "return 0123;", error: "Syntax error: integer values cannot start with '0'"},
		{name: "MinusNumber", input: "return -2;", error: "Syntax error: not a number"},
		{name: "MissingSemicolon", input: "return 123", error: "Syntax error: expected ; after expression"},
		{name: "MissingWhitespace", input: "return123;", error: "Syntax error: expected a statement got return123"},
		{name: "ByteAfterSemicolon", input: "return 1;}", error: "Syntax error: unexpected }"},
	}
	for _, test := range subTests {
		suite.Run(test.name, func() {
			ret, err := simple.Simple(test.input)
			suite.EqualError(err, test.error)
			suite.Nil(ret)
		})
	}
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
