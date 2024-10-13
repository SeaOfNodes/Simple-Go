package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type SubNode struct {
	binaryNode
}

func NewSubNode(lhs Node, rhs Node) *SubNode {
	return initBinaryNode(&SubNode{}, lhs, rhs)
}

func (s *SubNode) GraphicLabel() string { return "-" }
func (s *SubNode) label() string        { return "Sub" }

func (s *SubNode) compute() (types.Type, error) {
	// x - x => 0
	if s.Lhs() == s.Rhs() {
		return types.NewInt(0), nil
	}

	lType, ok := s.Lhs().base().typ.(*types.Int)
	if !ok || !lType.Constant() {
		return types.Bottom, nil
	}
	rType, ok := s.Rhs().base().typ.(*types.Int)
	if !ok || !rType.Constant() {
		return types.Bottom, nil
	}
	return types.NewInt(lType.Value - rType.Value), nil
}

func (s *SubNode) idealize() (Node, error) {
	// 0 - x => -x
	if lType, ok := Type(s.Lhs()).(*types.Int); ok && lType.Value == 0 {
		return NewMinusNode(s.Rhs()), nil
	}
	// x - 0 => x
	if rType, ok := Type(s.Rhs()).(*types.Int); ok && rType.Value == 0 {
		return s.Lhs(), nil
	}

	return nil, nil
}

func (s *SubNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(s.Lhs(), sb)
	sb.WriteString("-")
	toString(s.Rhs(), sb)
	sb.WriteString(")")
}
