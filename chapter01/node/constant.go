package node

type ConstantNode struct {
	baseNode
	Value int
}

func NewConstantNode(value int) *ConstantNode {
	return initBaseNode(&ConstantNode{Value: value}, StartNode)
}

func (c *ConstantNode) IsControl() bool { return false }
