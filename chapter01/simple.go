package simple

import (
	"go/ast"
	"strconv"

	"github.com/SeaOfNodes/Simple-Go/chapter01/node"
	"github.com/SeaOfNodes/Simple-Go/chapter01/parser"
)

func Simple(source string) (*node.ReturnNode, error) {
	p := parser.NewParser(source)
	n, err := p.Parse()
	if err != nil {
		return nil, err
	}

	var retNode *node.ReturnNode
	ast.Inspect(n, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			res := ret.Results[0].(*ast.BasicLit)
			num, err := strconv.Atoi(res.Value)
			if err != nil {
				return false
			}
			expr := node.NewConstantNode(num)
			retNode = node.NewReturnNode(node.StartNode, expr)
			return false
		}
		return true
	})
	return retNode, nil
}
