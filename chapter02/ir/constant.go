package ir

import (
	"strconv"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"
)

type ConstantNode struct {
	baseNode
}

func NewConstantNode(typ types.Type) *ConstantNode {
	n := initBaseNode(&ConstantNode{}, StartNode)
	n.typ = typ
	return n
}

func (c *ConstantNode) IsControl() bool { return false }

func (c *ConstantNode) compute() types.Type  { return c.typ }
func (c *ConstantNode) label() string        { return "#" + strconv.Itoa(c.value()) }
func (c *ConstantNode) GraphicLabel() string { return c.label() }

func (c *ConstantNode) value() int { return c.typ.(*types.IntType).Value }
