package types

import (
	"strings"
)

type Tuple struct {
	Types []Type
}

func NewTuple(types ...Type) *Tuple {
	return &Tuple{Types: types}
}

func (t *Tuple) Simple() bool   { return false }
func (t *Tuple) Constant() bool { return false }

func (t *Tuple) ToString(sb *strings.Builder) {
	sb.WriteString("[ ")
	for i, typ := range t.Types {
		typ.ToString(sb)
		if i != len(t.Types)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ]")
}

func (t *Tuple) Meet(_ Type) Type {
	panic("not implemented") // TODO: Implement
}
