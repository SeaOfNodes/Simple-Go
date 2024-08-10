package ir

import "github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"

type MulNode struct {
	binaryNode
}

func NewMulNode(lhs Node, rhs Node) *MulNode {
	return initBinaryNode(&MulNode{}, lhs, rhs)
}

func (m *MulNode) compute() types.Type {
	lType, ok := m.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType
	}
	rType, ok := m.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType
	}

	if lType.Constant() && rType.Constant() {
		return types.NewIntType(lType.Value * rType.Value)
	}
	return types.BottomType
}

func (m *MulNode) label() string        { return "Mul" }
func (m *MulNode) GraphicLabel() string { return "*" }
