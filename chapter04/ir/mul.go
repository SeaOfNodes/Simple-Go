package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type MulNode struct {
	binaryNode
}

func NewMulNode(lhs Node, rhs Node) *MulNode {
	return initBinaryNode(&MulNode{}, lhs, rhs)
}

func (m *MulNode) GraphicLabel() string { return "*" }
func (m *MulNode) label() string        { return "Mul" }

func (m *MulNode) compute() (types.Type, error) {
	rType, ok := m.Rhs().base().typ.(*types.Int)
	if !ok || !rType.Constant() {
		return types.Bottom, nil
	}
	// x*0=>0
	if rType.Value == 0 {
		return types.NewInt(0), nil
	}
	lType, ok := m.Lhs().base().typ.(*types.Int)
	if !ok || !lType.Constant() {
		return types.Bottom, nil
	}

	return types.NewInt(lType.Value * rType.Value), nil
}

func (m *MulNode) idealize() (Node, error) {
	if rType, ok := Type(m.Rhs()).(*types.Int); ok && rType.Value == 1 {
		return m.Lhs(), nil
	}

	if Type(m.Lhs()).Constant() && !Type(m.Rhs()).Constant() {
		m.swap()
		return m, nil
	}

	return nil, nil
}

func (m *MulNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(m.Lhs(), sb)
	sb.WriteString("*")
	toString(m.Rhs(), sb)
	sb.WriteString(")")
}
