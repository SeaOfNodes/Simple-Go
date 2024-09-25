package ir

import (
	"fmt"
	"go/ast"
	"slices"
	"strconv"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
	"github.com/pkg/errors"
)

var DisablePeephole = false

var nodeID = 0

func computeError(n ast.Node, msg string) *ASTError {
	internal := errors.New("Compute error: " + msg)
	return &ASTError{error: internal, Pos: n.Pos()}
}

// Node is the interface every node type must implement. In order to avoid duplicate code, nodes should embed `baseNode`.
type Node interface {
	// IsControl indicates whether or not this node is part of the control flow graph
	IsControl() bool

	idealize() (Node, error)
	compute() (types.Type, error)
	label() string
	GraphicLabel() string
	toStringInternal(*strings.Builder)

	// Implemented by baseNode to get baseNode
	base() *baseNode
}

// MultiNode is a node that has a tuple result
type MultiNode interface {
	Node
	multinode()
}

type baseNode struct {
	ins    []Node
	outs   []Node
	id     int
	typ    types.Type
	pinned bool
}

// initBaseNode initializes the baseNode in the given node n. It returns n for convenience.
func initBaseNode[T Node](n T, ins ...Node) T {
	b := n.base()
	b.id = nodeID
	nodeID++
	b.ins = ins
	for _, in := range ins {
		if in != nil {
			addOut(in, n)
		}
	}
	return n
}

func UniqueName(n Node) string {
	id := strconv.Itoa(n.base().id)
	if _, ok := n.(*ConstantNode); ok {
		return "Con_" + id
	}
	if p, ok := n.(*ProjNode); ok {
		return fmt.Sprintf("%s:p%d", UniqueName(p.control()), p.i)
	}
	return n.label() + id
}

func In(n Node, i int) Node {
	return n.base().ins[i]
}

func NumOfIns(n Node) int {
	return len(n.base().ins)
}

func NumOfOuts(n Node) int {
	return len(n.base().outs)
}

func Unused(n Node) bool {
	return NumOfOuts(n) == 0 && !n.base().pinned
}

func Type(n Node) types.Type {
	return n.base().typ
}

func Ins(n Node) []Node {
	return n.base().ins
}

func Outs(n Node) []Node {
	return n.base().outs
}

func id(n Node) int {
	return n.base().id
}

func addIn(n Node, in Node) {
	n.base().ins = append(n.base().ins, in)
	if in != nil {
		addOut(in, n)
	}
}

func addOut(n Node, out Node) {
	n.base().outs = append(n.base().outs, out)
}

func removeOut(n Node, out Node) {
	i := slices.Index(n.base().outs, out)
	last := len(n.base().outs) - 1
	n.base().outs[i] = n.base().outs[last]
	n.base().outs = n.base().outs[:last]
}

func setIn(n Node, i int, in Node) error {
	old := In(n, i)
	if old == in {
		return nil
	}

	if old != nil {
		removeOut(old, n)
		if Unused(old) {
			err := kill(old)
			if err != nil {
				return err
			}
		}
	}

	n.base().ins[i] = in
	if in != nil {
		addOut(in, n)
	}
	return nil
}

func removeLastIn(n Node) error {
	lastIn := len(n.base().ins) - 1
	err := setIn(n, lastIn, nil)
	if err != nil {
		return err
	}
	n.base().ins = n.base().ins[:lastIn]
	return nil
}

func kill(n Node) error {
	if !Unused(n) {
		return errors.New("Cannot kill a node that is in use")
	}

	for i := range n.base().ins {
		err := setIn(n, i, nil)
		if err != nil {
			return err
		}
	}
	n.base().ins = []Node{}
	n.base().typ = nil

	if !dead(n) {
		return errors.Errorf("Node not dead after killing it: %v", n)
	}
	return nil
}

func dead(n Node) bool {
	return Unused(n) && len(n.base().ins) == 0 && n.base().typ == nil
}

func peephole(n Node) (Node, error) {
	typ, err := n.compute()
	if err != nil {
		return nil, err
	}
	n.base().typ = typ

	if DisablePeephole {
		return n, nil
	}

	var opt Node
	if _, ok := n.(*ConstantNode); !ok && Type(n).Constant() {
		opt = NewConstantNode(typ)
	} else {
		var err error
		opt, err = n.idealize()
		if err != nil {
			return nil, err
		}
	}

	if opt != nil {
		opt, err := peephole(opt)
		if err != nil {
			return nil, err
		}
		err = replace(n, opt)
		if err != nil {
			return nil, err
		}
		return opt, nil
	}
	return n, nil
}

func peepholeT[T Node](n T) (T, error) {
	p, err := peephole(n)
	if err != nil {
		return n, err
	}
	t, ok := p.(T)
	if !ok {
		return n, errors.Errorf("Expected node of type %T but got %T", n, p)
	}
	return t, nil
}

func pin[T Node](n T) T {
	n.base().pinned = true
	return n
}

func unpin(n Node) {
	n.base().pinned = false
}

func replace(old Node, new Node) error {
	if old != new && Unused(old) {
		pin(new)
		err := kill(old)
		if err != nil {
			return err
		}
		unpin(new)
	}
	return nil
}

func ToString(n Node) string {
	sb := &strings.Builder{}
	toString(n, sb)
	return sb.String()
}

func toString(n Node, sb *strings.Builder) {
	if n == nil {
		sb.WriteString("nil")
		return
	}

	if dead(n) {
		sb.WriteString(UniqueName(n))
		sb.WriteString(":DEAD")
		return
	}

	n.toStringInternal(sb)
}

func (b *baseNode) base() *baseNode {
	return b
}
