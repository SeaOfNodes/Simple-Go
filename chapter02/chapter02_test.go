package simple_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Chapter02TestSuite struct {
	suite.Suite
}

func (suite *Chapter02TestSuite) TestParseGrammer() {
	//    Node._disablePeephole = true; // disable peephole so we can observe full graph
	// Parser parser = new Parser("return 1+2*3+-5;");
	// ReturnNode ret = parser.parse();
	// assertEquals("return (1+((2*3)+(-5)));", ret.print());
	// GraphVisualizer gv = new GraphVisualizer();
	// System.out.println(gv.generateDotOutput(parser));
	// Node._disablePeephole = false;
	// TODO: This
}

func TestParser(t *testing.T) {
	suite.Run(t, new(Chapter02TestSuite))
}
