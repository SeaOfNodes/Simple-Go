package ir

import (
	"maps"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
	"github.com/pkg/errors"
)

type symbolTable map[string]int

const Control = "$ctrl"
const Arg0 = "arg"

type ScopeNode struct {
	Scopes []symbolTable
	baseNode
}

func NewScopeNode() *ScopeNode {
	s := initBaseNode(&ScopeNode{})
	s.typ = types.Bottom
	return s
}

func (s *ScopeNode) IsControl() bool      { return false }
func (s *ScopeNode) GraphicLabel() string { return s.label() }

func (s *ScopeNode) compute() (types.Type, error) { return types.Bottom, nil }
func (s *ScopeNode) idealize() (Node, error)      { return nil, nil }
func (s *ScopeNode) label() string                { return "Scope" }

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

func (s *ScopeNode) Clone() *ScopeNode {
	clone := NewScopeNode()
	clone.Scopes = make([]symbolTable, 0, len(s.Scopes))
	for _, table := range s.Scopes {
		clone.Scopes = append(clone.Scopes, (maps.Clone(table)))
	}
	for _, i := range Ins(s) {
		addIn(clone, i)
	}
	return clone
}

func (s *ScopeNode) dataNames() []string {
	names := make([]string, NumOfIns(s))
	for _, table := range s.Scopes {
		for name, i := range table {
			names[i] = name
		}
	}
	// Don't include $ctrl
	return names[1:]
}

func (s *ScopeNode) dataInputs() []Node { return Ins(s)[1:] }
func (s *ScopeNode) setDataInput(i int, n Node) error {
	return setIn(s, i+1, n)
}

func (s *ScopeNode) Merge(x *ScopeNode) (*RegionNode, error) {
	r, err := peepholeT(pin(NewRegionNode(s.Control(), x.Control())))
	if err != nil {
		return nil, err
	}
	names := s.dataNames()
	for i, d := range s.dataInputs() {
		d1 := x.dataInputs()[i]
		if d == d1 {
			continue
		}

		n, err := peephole(NewPhiNode(names[i], r, d, d1))
		if err != nil {
			return nil, err
		}
		err = s.setDataInput(i, n)
		if err != nil {
			return nil, err
		}
	}
	err = kill(x)
	if err != nil {
		return nil, err
	}
	return r, nil
}
