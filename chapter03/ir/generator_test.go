package ir

import (
	"testing"

	"github.com/SeaOfNodes/Simple-Go/chapter03/parser"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
}

func (suite *GeneratorTestSuite) TestPrint() {
	DisablePeephole = true
	p := parser.NewParser("return 1+2*3+-5;")
	n, err := p.Parse()
	suite.NoError(err)
	retNode, err := NewGenerator().Generate(n)
	suite.NoError(err)
	suite.Equal("return ((1+(2*3))+(-5));", ToString(retNode))
	DisablePeephole = false
}

func (suite *GeneratorTestSuite) TestPeephole() {
	subTests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "add", input: "return 1+2;", expected: "return 3;"},
		{name: "sub", input: "return 1-2;", expected: "return -1;"},
		{name: "mul", input: "return 2*3;", expected: "return 6;"},
		{name: "div", input: "return 6/2;", expected: "return 3;"},
		{name: "minus", input: "return 1--2;", expected: "return 3;"},
		{name: "arithmetic", input: "return 1+2*3+-5;", expected: "return 2;"},
	}

	for _, test := range subTests {
		suite.Run(test.name, func() {
			p := parser.NewParser(test.input)
			n, err := p.Parse()
			suite.NoError(err)
			retNode, err := NewGenerator().Generate(n)
			suite.NoError(err)
			suite.Equal(test.expected, ToString(retNode))
		})
	}
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
