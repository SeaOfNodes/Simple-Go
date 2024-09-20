package ir

import (
	"go/ast"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type DivNode struct {
	expr ast.Expr
	binaryNode
}

func NewDivNode(lhs Node, rhs Node) *DivNode {
	return initBinaryNode(&DivNode{}, lhs, rhs)
}

func (d *DivNode) GraphicLabel() string { return "/" }
func (d *DivNode) label() string        { return "Div" }

func (d *DivNode) compute() (types.Type, error) {
	lType, ok := Type(d.Lhs()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}
	rType, ok := Type(d.Rhs()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}

	if lType.Constant() && rType.Constant() {
		if rType.Value == 0 {
			return nil, computeError(d.expr, "divide by zero")
		}
		return types.NewInt(lType.Value / rType.Value), nil
	}
	return types.Bottom, nil
}

func (d *DivNode) idealize() (Node, error) {
	if rType, ok := Type(d.Rhs()).(*types.Int); ok && rType.Value == 1 {
		return d.Lhs(), nil
	}

	return nil, nil
}

func (d *DivNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(d.Lhs(), sb)
	sb.WriteString("/")
	toString(d.Rhs(), sb)
	sb.WriteString(")")
}
