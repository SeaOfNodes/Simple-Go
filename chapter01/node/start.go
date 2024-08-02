package node

// We only need one StartNode
var StartNode = newStartNode()

type startNode struct {
	baseNode
}

func newStartNode() *startNode {
	return initBaseNode(&startNode{})
}

func (s *startNode) IsControl() bool { return true }
