package types

import (
	"strconv"
	"strings"
)

var IntTop = &Int{Value: 0, con: false}
var IntBottom = &Int{Value: 1, con: false}

type Int struct {
	Value int
	con   bool
}

func NewInt(value int) Type {
	return &Int{Value: value, con: true}
}

func (i *Int) Simple() bool                 { return false }
func (i *Int) Constant() bool               { return true }
func (i *Int) ToString(sb *strings.Builder) { sb.WriteString(strconv.Itoa(i.Value)) }
func (i *Int) Meet(t Type) Type {
	if i == t {
		return i
	}
	i0, ok := t.(*Int)
	if !ok {
		return Bottom
	}
	if i.Bottom() || i0.Bottom() {
		return IntBottom
	}
	if i.Top() {
		return i0
	}
	if i0.Top() {
		return i
	}
	if i.Value == i0.Value {
		return i
	}
	return IntBottom
}

func (i *Int) Top() bool    { return i == IntTop }
func (i *Int) Bottom() bool { return i == IntBottom }
