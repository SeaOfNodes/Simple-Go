package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter03/ir/types"
	"github.com/pkg/errors"
)

type symbolTable map[string]int

type ScopeNode struct {
	Scopes []symbolTable
	baseNode
}

func NewScopeNode() *ScopeNode {
	s := initBaseNode(&ScopeNode{})
	s.typ = types.BottomType
	return s
}

func (s *ScopeNode) IsControl() bool      { return false }
func (s *ScopeNode) GraphicLabel() string { return s.label() }

func (s *ScopeNode) compute() types.Type { return types.BottomType }
func (s *ScopeNode) label() string       { return "Scope" }

func (s *ScopeNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString(s.label())
	for _, table := range s.Scopes {
		sb.WriteString("[")
		first := true
		for name := range table {
			if first {
				sb.WriteString(", ")
			}
			first = false

			sb.WriteString(name)
			sb.WriteString(":")
			n := In(s, table[name])
			toString(n, sb)
		}
		sb.WriteString("]")
	}
}

func (s *ScopeNode) Control() Node                 { return In(s, 0) }
func (s *ScopeNode) SetControl(control Node) error { return setIn(s, 0, control) }

func (s *ScopeNode) Define(name string, n Node) error {
	table := s.Scopes[len(s.Scopes)-1]
	if _, ok := table[name]; ok {
		return errors.Errorf("Cannot define a name that already exists: %s", name)
	}
	table[name] = NumOfIns(s)
	addIn(s, n)
	return nil
}

func (s *ScopeNode) Lookup(name string) (Node, bool) {
	i, ok := s.lookup(name)
	if !ok {
		return nil, false
	}
	return In(s, i), true
}

// Update returns true if the name exists in the symbol table
func (s *ScopeNode) Update(name string, n Node) (bool, error) {
	i, ok := s.lookup(name)
	if !ok {
		return false, nil
	}
	err := setIn(s, i, n)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (s *ScopeNode) lookup(name string) (int, bool) {
	for i := len(s.Scopes) - 1; i >= 0; i-- {
		if n, ok := s.Scopes[i][name]; ok {
			return n, true
		}
	}
	return 0, false
}

func (s *ScopeNode) Push() { s.Scopes = append(s.Scopes, symbolTable{}) }
func (s *ScopeNode) Pop() error {
	last := s.Scopes[len(s.Scopes)-1]
	for range last {
		err := removeLastIn(s)
		if err != nil {
			return err
		}
	}
	s.Scopes = s.Scopes[:len(s.Scopes)-1]
	return nil
}
