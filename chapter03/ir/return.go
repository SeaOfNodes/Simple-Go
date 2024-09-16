package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter03/ir/types"
)

type ReturnNode struct {
	baseNode
}

func NewReturnNode(control Node, data Node) *ReturnNode {
	return initBaseNode(&ReturnNode{}, control, data)
}

func (r *ReturnNode) IsControl() bool              { return true }
func (r *ReturnNode) GraphicLabel() string         { return "Return" }
func (r *ReturnNode) label() string                { return "Return" }
func (r *ReturnNode) compute() (types.Type, error) { return types.BottomType, nil }

func (r *ReturnNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("return ")
	toString(r.Expr(), sb)
	sb.WriteString(";")
}

func (r *ReturnNode) Control() Node { return In(r, 0) }
func (r *ReturnNode) Expr() Node    { return In(r, 1) }
