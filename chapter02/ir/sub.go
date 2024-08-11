package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"
)

type SubNode struct {
	binaryNode
}

func NewSubNode(lhs Node, rhs Node) *SubNode {
	return initBinaryNode(&SubNode{}, lhs, rhs)
}

func (s *SubNode) GraphicLabel() string { return "-" }
func (s *SubNode) label() string        { return "Sub" }

func (s *SubNode) compute() types.Type {
	lType, ok := s.Lhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType
	}
	rType, ok := s.Rhs().base().typ.(*types.IntType)
	if !ok {
		return types.BottomType
	}

	if lType.Constant() && rType.Constant() {
		return types.NewIntType(lType.Value - rType.Value)
	}
	return types.BottomType
}

func (s *SubNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(s.Lhs(), sb)
	sb.WriteString("-")
	toString(s.Rhs(), sb)
	sb.WriteString(")")
}
