package node

var nodeID = 0

// Node is the interface every node type must implement. In order to avoid duplicate code, nodes should embed `baseNode`.
type Node interface {
	// IsControl must be implemented by each node.
	IsControl() bool

	// All methods below are implemented by baseNode
	In(int) Node
	NumOfIns() int
	NumOfOuts() int
	Unused() bool

	base() *baseNode
	addOutput(Node)
}

type baseNode struct {
	ins  []Node
	outs []Node
	id   int
}

// NewNode initializes the baseNode in the given node n. It returns n for convenience.
func initBaseNode[T Node](n T, ins ...Node) T {
	b := n.base()
	b.id = nodeID
	nodeID++
	b.ins = ins
	for _, in := range ins {
		if in != nil {
			in.addOutput(n)
		}
	}
	return n
}

func (b *baseNode) base() *baseNode {
	return b
}

func (b *baseNode) addOutput(out Node) {
	b.outs = append(b.outs, out)
}

func (b *baseNode) In(i int) Node {
	return b.ins[i]
}

func (b *baseNode) NumOfIns() int {
	return len(b.ins)
}

func (b *baseNode) NumOfOuts() int {
	return len(b.outs)
}

func (b *baseNode) Unused() bool {
	return b.NumOfOuts() == 0
}
