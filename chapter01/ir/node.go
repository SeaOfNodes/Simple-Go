package ir

var DisablePeephole = false

var nodeID = 0

// Node is the interface every node type must implement. In order to avoid duplicate code, nodes should embed `baseNode`.
type Node interface {
	// IsControl indicates whether or not this node is part of the control flow graph
	IsControl() bool

	// Implemented by baseNode to get baseNode
	base() *baseNode
}

type baseNode struct {
	ins  []Node
	outs []Node
	id   int
}

// initBaseNode initializes the baseNode in the given node n. It returns n for convenience.
func initBaseNode[T Node](n T, ins ...Node) T {
	b := n.base()
	b.id = nodeID
	nodeID++
	b.ins = ins
	for _, in := range ins {
		if in != nil {
			addOut(in, n)
		}
	}
	return n
}

func In(n Node, i int) Node {
	return n.base().ins[i]
}

func NumOfIns(n Node) int {
	return len(n.base().ins)
}

func NumOfOuts(n Node) int {
	return len(n.base().outs)
}

func Unused(n Node) bool {
	return NumOfOuts(n) == 0
}

func Ins(n Node) []Node {
	return n.base().ins
}

func Outs(n Node) []Node {
	return n.base().outs
}

func addOut(n Node, out Node) {
	n.base().outs = append(n.base().outs, out)
}

func (b *baseNode) base() *baseNode {
	return b
}
