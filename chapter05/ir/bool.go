package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type BoolType string

const (
	EQ = BoolType("==")
	LT = BoolType("<")
	LE = BoolType("<=")
)

type BoolNode struct {
	binaryNode
	op BoolType
}

func NewBoolNode(lhs Node, op BoolType, rhs Node) *BoolNode {
	return initBinaryNode(&BoolNode{op: op}, lhs, rhs)
}

func (b *BoolNode) doOp(lhs int, rhs int) types.Type {
	val := false
	switch b.op {
	case EQ:
		val = lhs == rhs
	case LT:
		val = lhs < rhs
	case LE:
		val = lhs <= rhs
	}
	if val {
		return types.NewInt(1)
	}
	return types.NewInt(0)
}

func (b *BoolNode) compute() (types.Type, error) {
	lType, ok := Type(b.Lhs()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}
	rType, ok := Type(b.Rhs()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}

	if lType.Constant() && rType.Constant() {
		return b.doOp(lType.Value, rType.Value), nil
	}
	return lType.Meet(rType), nil
}

func (b *BoolNode) idealize() (Node, error) {
	if b.Lhs() == b.Rhs() {
		return NewConstantNode(b.doOp(3, 3)), nil
	}
	return nil, nil
}

func (b *BoolNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(b.Lhs(), sb)
	sb.WriteString(string(b.op))
	toString(b.Rhs(), sb)
	sb.WriteString(")")
}

func (b *BoolNode) GraphicLabel() string { return string(b.op) }
func (b *BoolNode) label() string {
	switch b.op {
	case EQ:
		return "eq"
	case LE:
		return "le"
	case LT:
		return "lt"
	}
	return ""
}
