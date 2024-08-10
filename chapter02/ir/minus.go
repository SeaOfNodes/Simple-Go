package ir

import "github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"

type MinusNode struct {
	baseNode
}

func NewMinusNode(value Node) *MinusNode {
	return initBaseNode(&MinusNode{}, value)
}

func (m *MinusNode) IsControl() bool { return false }

func (m *MinusNode) compute() types.Type {
	typ, ok := m.Value().base().typ.(*types.IntType)
	if ok {
		if typ.Constant() {
			return types.NewIntType(-typ.Value)
		}
		return typ
	}

	return types.BottomType
}

func (m *MinusNode) label() string        { return "Minus" }
func (m *MinusNode) GraphicLabel() string { return "-" }

func (m *MinusNode) Value() Node { return m.ins[0] }
