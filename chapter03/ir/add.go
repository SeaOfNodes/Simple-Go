package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter03/ir/types"
)

type AddNode struct {
	binaryNode
}

func NewAddNode(lhs Node, rhs Node) *AddNode {
	return initBinaryNode(&AddNode{}, lhs, rhs)
}

func (a *AddNode) GraphicLabel() string { return "+" }
func (a *AddNode) label() string        { return "Add" }

func (a *AddNode) compute() (types.Type, error) {
	lType, ok := a.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType, nil
	}
	rType, ok := a.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType, nil
	}

	if lType.Constant() && rType.Constant() {
		return types.NewIntType(lType.Value + rType.Value), nil
	}
	return types.BottomType, nil
}

func (a *AddNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(a.Lhs(), sb)
	sb.WriteString("+")
	toString(a.Rhs(), sb)
	sb.WriteString(")")
}
