package ir

type binaryNode struct {
	baseNode
}

type BinaryNode interface {
	Node

	Lhs() Node
	Rhs() Node

	binary() *binaryNode
}

// initBinaryNode initializes the binaryNode in the given node n. It returns n for convenience.
func initBinaryNode[T BinaryNode](n T, lhs Node, rhs Node) T {
	return initBaseNode(n, lhs, rhs)
}

func (b *binaryNode) binary() *binaryNode { return b }

func (b *binaryNode) IsControl() bool { return false }

func (b *binaryNode) Lhs() Node { return b.ins[0] }
func (b *binaryNode) Rhs() Node { return b.ins[1] }
