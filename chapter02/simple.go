package simple

import (
	"github.com/SeaOfNodes/Simple-Go/chapter01/ir"
	"github.com/SeaOfNodes/Simple-Go/chapter01/parser"
)

func Simple(source string) (*ir.ReturnNode, error) {
	p := parser.NewParser(source)
	n, err := p.Parse()
	if err != nil {
		return nil, err
	}

	generator := ir.NewGenerator()
	return generator.Generate(n)
}
