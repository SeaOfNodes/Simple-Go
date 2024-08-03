package ir

type ReturnNode struct {
	baseNode
}

func NewReturnNode(control Node, data Node) *ReturnNode {
	return initBaseNode(&ReturnNode{}, control, data)
}

func (r *ReturnNode) IsControl() bool { return true }

func (r *ReturnNode) Control() Node { return r.In(0) }
func (r *ReturnNode) Expr() Node    { return r.In(1) }
