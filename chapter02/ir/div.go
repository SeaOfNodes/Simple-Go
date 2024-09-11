package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"
)

type DivNode struct {
	binaryNode
}

func NewDivNode(lhs Node, rhs Node) *DivNode {
	return initBinaryNode(&DivNode{}, lhs, rhs)
}

func (d *DivNode) GraphicLabel() string { return "/" }
func (d *DivNode) label() string        { return "Div" }

func (d *DivNode) compute() types.Type {
	lType, ok := Type(d.Lhs()).(*types.IntType)
	if !ok {
		return types.BottomType
	}
	rType, ok := Type(d.Rhs()).(*types.IntType)
	if !ok {
		return types.BottomType
	}

	if lType.Constant() && rType.Constant() {
		if rType.Value == 0 {
			return types.NewIntType(0)
		}
		return types.NewIntType(lType.Value / rType.Value)
	}
	return types.BottomType
}

func (d *DivNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(d.Lhs(), sb)
	sb.WriteString("/")
	toString(d.Rhs(), sb)
	sb.WriteString(")")
}
