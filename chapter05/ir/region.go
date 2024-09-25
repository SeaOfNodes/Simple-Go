package ir

import (
	"strconv"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type RegionNode struct {
	baseNode
}

func NewRegionNode(inputs ...Node) *RegionNode {
	return initBaseNode(&RegionNode{}, inputs...)
}

func (r *RegionNode) label() string        { return "Region" }
func (r *RegionNode) GraphicLabel() string { return "Region" }
func (r *RegionNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("Region" + strconv.Itoa(id(r)))
}
func (r *RegionNode) IsControl() bool              { return true }
func (r *RegionNode) compute() (types.Type, error) { return types.Control, nil }
func (r *RegionNode) idealize() (Node, error)      { return nil, nil }
