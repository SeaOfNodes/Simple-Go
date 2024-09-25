package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type IfNode struct {
	baseNode
}

func NewIfNode(control Node, data Node) *IfNode {
	return initBaseNode(&IfNode{}, control, data)
}

func (i *IfNode) multinode()           {}
func (i *IfNode) control() Node        { return In(i, 0) }
func (i *IfNode) pred() Node           { return In(i, 1) }
func (i *IfNode) label() string        { return "If" }
func (i *IfNode) GraphicLabel() string { return "If" }
func (i *IfNode) IsControl() bool      { return true }
func (i *IfNode) compute() (types.Type, error) {
	return types.NewTuple(types.Control, types.Control), nil
}
func (i *IfNode) idealize() (Node, error) { return nil, nil }
func (i *IfNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("if (")
	toString(i.pred(), sb)
	sb.WriteString(")")
}
