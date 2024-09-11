package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"
)

// We only need one StartNode
var StartNode = newStartNode()

type startNode struct {
	baseNode
}

func newStartNode() *startNode {
	return initBaseNode(&startNode{})
}

func (s *startNode) IsControl() bool      { return true }
func (s *startNode) GraphicLabel() string { return "Start" }

func (s *startNode) compute() (types.Type, error)         { return types.BottomType, nil }
func (s *startNode) label() string                        { return "Start" }
func (s *startNode) toStringInternal(sb *strings.Builder) { sb.WriteString(s.label()) }
