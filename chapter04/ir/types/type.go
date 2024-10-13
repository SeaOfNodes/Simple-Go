package types

import (
	"strings"
)

type Type interface {
	Simple() bool
	Constant() bool
	ToString(*strings.Builder)
	Meet(Type) Type
}

type simple struct {
	s        string
	constant bool
}

func (s *simple) Simple() bool                 { return true }
func (s *simple) Constant() bool               { return s.constant }
func (s *simple) ToString(sb *strings.Builder) { sb.WriteString(s.s) }
func (s *simple) Meet(_ Type) Type             { return Bottom }

var Top = &simple{s: "Top", constant: true}
var Bottom = &simple{s: "Bottom", constant: false}
var Control = &simple{s: "Control", constant: false}
