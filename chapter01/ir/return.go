package ir

type ReturnNode struct {
	baseNode
}

func NewReturnNode(control Node, data Node) *ReturnNode {
	return initBaseNode(&ReturnNode{}, control, data)
}

func (r *ReturnNode) IsControl() bool { return true }

func (r *ReturnNode) Control() Node { return In(r, 0) }
func (r *ReturnNode) Expr() Node    { return In(r, 1) }
