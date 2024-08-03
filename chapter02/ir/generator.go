package ir

import (
	"go/ast"
	"strconv"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(n ast.Node) (*ReturnNode, error) {
	var retNode *ReturnNode
	ast.Inspect(n, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			res := ret.Results[0].(*ast.BasicLit)
			num, err := strconv.Atoi(res.Value)
			if err != nil {
				return false
			}
			expr := NewConstantNode(num)
			retNode = NewReturnNode(StartNode, expr)
			return false
		}
		return true
	})
	return retNode, nil
}
