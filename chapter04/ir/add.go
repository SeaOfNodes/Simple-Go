package ir

import (
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type AddNode struct {
	binaryNode
}

func NewAddNode(lhs Node, rhs Node) *AddNode {
	return initBinaryNode(&AddNode{}, lhs, rhs)
}

func (a *AddNode) GraphicLabel() string { return "+" }
func (a *AddNode) label() string        { return "Add" }

func (a *AddNode) compute() (types.Type, error) {
	lType, ok := Type(a.Lhs()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}
	rType, ok := Type(a.Rhs()).(*types.Int)
	if !ok {
		return types.Bottom, nil
	}

	if lType.Constant() && rType.Constant() {
		return types.NewInt(lType.Value + rType.Value), nil
	}
	return lType.Meet(rType), nil
}

func (a *AddNode) idealize() (Node, error) {
	if c, ok := Type(a.Rhs()).(*types.Int); ok && c.Value == 0 {
		return a.Lhs(), nil
	}

	if a.Lhs() == a.Rhs() {
		mul := NewMulNode(a.Lhs(), NewConstantNode(types.NewInt(2)))
		return peephole(mul)
	}

	/* Move all adds to lhs
	The possibilities are:
	* No adds: a + b -> nothing to do
	* lhs add: (a + b) + c -> nothing to do
	* rhs add: a + (b + c) -> swap so (a + b) + c
	* lhs&rhs add: (a + b) + (c + d) -> rebuild so ((a + b) + c) + d
	*/

	lAdd, lhsIsAdd := a.Lhs().(*AddNode)
	if rAdd, ok := a.Rhs().(*AddNode); ok {
		if lhsIsAdd {
			// We have (a + b) + (c + d)
			// lhs: (a + b) + c
			lhs, err := peephole(NewAddNode(a.Lhs(), rAdd.Lhs()))
			if err != nil {
				return nil, err
			}
			// new add is ((a + b) + c) + d
			return NewAddNode(lhs, rAdd.Rhs()), nil
		}
		// We have a + (b + c), swap so (a + b) + c
		a.swap()
		return a, nil
	}

	if !lhsIsAdd {
		// We have a + b
		if shouldSwapNonAdds(a.Lhs(), a.Rhs()) {
			a.swap()
			return a, nil
		}
		return nil, nil
	}

	if Type(lAdd.Rhs()).Constant() && Type(a.Rhs()).Constant() {
		// We have (v + c1) + c2 (c1 and c2 are constants)
		// rhs: c1 + c2 (folded)
		rhs, err := peephole(NewAddNode(lAdd.Rhs(), a.Rhs()))
		if err != nil {
			return nil, err
		}
		// new add is v + [c1 + c2]
		return NewAddNode(lAdd.Lhs(), rhs), nil
	}

	// Maybe change (a + b) + c to (a + c) + b
	if shouldSwapNonAdds(lAdd.Rhs(), a.Rhs()) {
		lhs, err := peephole(NewAddNode(lAdd.Lhs(), a.Rhs()))
		if err != nil {
			return nil, err
		}
		return NewAddNode(lhs, lAdd.Rhs()), nil
	}

	return nil, nil
}

func shouldSwapNonAdds(lhs Node, rhs Node) bool {
	return !Type(rhs).Constant() && (Type(lhs).Constant() || id(rhs) > id(lhs))
}

func (a *AddNode) toStringInternal(sb *strings.Builder) {
	sb.WriteString("(")
	toString(a.Lhs(), sb)
	sb.WriteString("+")
	toString(a.Rhs(), sb)
	sb.WriteString(")")
}
