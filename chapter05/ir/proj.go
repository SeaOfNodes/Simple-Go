package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type ProjNode struct {
	baseNode
	i int
	s string
}

func NewProjNode(control MultiNode, i int, label string) *ProjNode {
	return initBaseNode(&ProjNode{i: i, s: label}, control)
}

func (p *ProjNode) control() Node           { return In(p, 0) }
func (p *ProjNode) IsControl() bool         { return p.i == 0 }
func (p *ProjNode) idealize() (Node, error) { return nil, nil }

func (p *ProjNode) compute() (types.Type, error) {
	if t, ok := Type(p.control()).(*types.Tuple); ok {
		return t.Types[p.i], nil
	}
	return types.Bottom, nil
}

func (p *ProjNode) label() string                        { return p.s }
func (p *ProjNode) GraphicLabel() string                 { return p.s }
func (p *ProjNode) toStringInternal(sb *strings.Builder) { sb.WriteString(p.s) }
