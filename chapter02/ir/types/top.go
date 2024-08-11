package types

import "strings"

type topType struct{}

func (t *topType) Simple() bool                 { return true }
func (t *topType) Constant() bool               { return false }
func (t *topType) ToString(sb *strings.Builder) { sb.WriteString("Top") }

var TopType = &topType{}
