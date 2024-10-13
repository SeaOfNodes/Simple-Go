package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

// We only need one StartNode
var StartNode *startNode

type startNode struct {
	args *types.Tuple
	baseNode
}

func newStartNode(args *types.Tuple) *startNode {
	s := initBaseNode(&startNode{args: args})
	s.typ = args
	return s
}

func (s *startNode) IsControl() bool      { return true }
func (s *startNode) GraphicLabel() string { return "Start" }

func (s *startNode) multinode()                           {}
func (s *startNode) idealize() (Node, error)              { return nil, nil }
func (s *startNode) compute() (types.Type, error)         { return types.Bottom, nil }
func (s *startNode) label() string                        { return "Start" }
func (s *startNode) toStringInternal(sb *strings.Builder) { sb.WriteString(s.label()) }
