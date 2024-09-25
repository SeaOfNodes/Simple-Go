package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type MinusNode struct {
	baseNode
}

func NewMinusNode(value Node) *MinusNode {
	return initBaseNode(&MinusNode{}, value)
}

func (m *MinusNode) IsControl() bool      { return false }
func (m *MinusNode) GraphicLabel() string { return "-" }
func (m *MinusNode) label() string        { return "Minus" }

func (m *MinusNode) compute() (types.Type, error) {
	typ, ok := m.Value().base().typ.(*types.Int)
	if ok {
		if typ.Constant() {
			return types.NewInt(-typ.Value), nil
		}
		return typ, nil
	}

	return types.Bottom, nil
}

func (m *MinusNode) idealize() (Node, error) {
	// -(x-y) => y-x
	if s, ok := m.Value().(*SubNode); ok {
		lhs, rhs := s.Rhs(), s.Lhs()
		return NewSubNode(lhs, rhs), nil
	}
	return nil, nil
}

func (m *MinusNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(-")
	toString(m.Value(), sb)
	sb.WriteString(")")
}

func (m *MinusNode) Value() Node { return m.ins[0] }
