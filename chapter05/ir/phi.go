package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type PhiNode struct {
	baseNode
	s string
}

func NewPhiNode(label string, region *RegionNode, inputs ...Node) *PhiNode {
	return initBaseNode(&PhiNode{s: label}, append([]Node{region}, inputs...)...)
}

func (p *PhiNode) region() Node                 { return In(p, 0) }
func (p *PhiNode) label() string                { return "Phi_" + p.s }
func (p *PhiNode) GraphicLabel() string         { return "&phi;_" + p.s }
func (p *PhiNode) IsControl() bool              { return true }
func (p *PhiNode) compute() (types.Type, error) { return types.Bottom, nil }

func (p *PhiNode) idealize() (Node, error) {
	if p.sameInputs() {
		return In(p, 1), nil
	}

	// TODO: Implement
	return nil, nil
}

func (p *PhiNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("Phi(")
	for i, n := range Ins(p) {
		if i != 0 {
			sb.WriteString(", ")
		}
		toString(n, sb)
	}
	sb.WriteString(")")
}

func (p *PhiNode) sameInputs() bool {
	f := In(p, 1)
	for _, n := range Ins(p)[2:] {
		if n != f {
			return false
		}
	}
	return true
}
