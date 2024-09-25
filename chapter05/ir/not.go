package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type NotNode struct {
	baseNode
}

func NewNotNode(value Node) *NotNode {
	return initBaseNode(&NotNode{}, value)
}

func (n *NotNode) value() Node     { return In(n, 0) }
func (n *NotNode) IsControl() bool { return false }

func (n *NotNode) idealize() (Node, error) {
	// Idealize !!!x => !x
	if n1, ok := n.value().(*NotNode); ok {
		if n2, ok := n1.value().(*NotNode); ok {
			return NewNotNode(n2.value()), nil
		}
	}
	return nil, nil
}

func (n *NotNode) compute() (types.Type, error) {
	t, ok := Type(n.value()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}
	if !t.Constant() {
		return t, nil
	}

	if t.Value == 0 {
		return types.NewInt(1), nil
	}
	return types.NewInt(0), nil
}

func (n *NotNode) label() string        { return "Not" }
func (n *NotNode) GraphicLabel() string { return "!" }
func (n *NotNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(!")
	toString(n.value(), sb)
	sb.WriteString(")")
}
