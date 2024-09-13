package types

import (
	"strconv"
	"strings"
)

type IntType struct {
	Value int
}

func NewIntType(value int) Type {
	return &IntType{Value: value}
}

func (i *IntType) Simple() bool                 { return false }
func (i *IntType) Constant() bool               { return true }
func (i *IntType) ToString(sb *strings.Builder) { sb.WriteString(strconv.Itoa(i.Value)) }
